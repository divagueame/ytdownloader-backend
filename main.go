package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var (
	// addr      = flag.String("addr", ":8080", "http service address")
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Message struct {
	Url string `json:"url"`
}

type State int

const (
	Idle = iota
	Processing
	Done
	Error
)

type DownloadMsg struct {
	Id    string `json:"id"`
	State State  `json:"state"`
	Url   string `json:"url"`
}

var downloadManager = NewDownloadManager()

func main() {
	http.HandleFunc("/health", health)
	http.HandleFunc("/notify", notify)

	http.HandleFunc("/download", triggerDownload)
	http.HandleFunc("/download/", serveFile)
	http.ListenAndServe(":5000", nil)
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}

func triggerDownload(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		w.WriteHeader(http.StatusOK)
		return
	}

	session_id := setSessionCookie(w, r)

	var msg Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	downloadR := downloadManager.createDownload(msg.Url, session_id.Value)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	go downloadManager.downloadFile(downloadR.Id)
	jsonResponse, err := json.Marshal(downloadR)

	// w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Write([]byte(jsonResponse))
}

type Notification struct {
	Data string `json:"data"`
}

func notify(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	// w.Header().Set("Access-Control-Allow-Origin", "*")

	cookie, err := r.Cookie("chikiyt_session_id")
	if err != nil && err == http.ErrNoCookie {
		log.Print("No cookie present on notifiy")
		return
	}

	// <-r.Context().Done()

	downloadManager.mu.Lock()
	ch, exists := downloadManager.completed[cookie.Value]
	if !exists {
		ch = make(chan string, 100)
		downloadManager.completed[cookie.Value] = ch
	}
	downloadManager.mu.Unlock()

	for {
		select {
		case id := <-ch:
			download := downloadManager.getDownloadById(id)
			log.Printf("Owner Download done received for --> %s", download.Title)

			type DownloadNotification struct {
				ID    string `json:"id"`
				Title string `json:"title"`
				State int    `json:"state"`
			}

			type Notification struct {
				Data DownloadNotification `json:"data"`
			}

			notification := Notification{
				Data: DownloadNotification{
					ID:    id,
					Title: download.Title,
					State: int(download.State),
				},
			}

			notificationJSON, _ := json.Marshal(notification)
			fmt.Fprintf(w, "data: %s\n\n", notificationJSON)
			flusher.Flush()
		default:
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func serveFile(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		http.Error(w, "File Not found", http.StatusNotFound)
	}

	id := parts[2]
	download := downloadManager.getDownloadById(id)
	if download == nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

	http.ServeFile(w, r, download.FilePath)
}

package main

type song struct {
  Title string `json:"title"`
import (
	"context"

	"encoding/json"
	"fmt"
	"github.com/lrstanley/go-ytdlp"
	"log"
	// "io"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
	// "github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
)

var (
	// addr      = flag.String("addr", ":8080", "http service address")
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)
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

func main() {
	// go logDownloads()
	http.HandleFunc("/mp3/", serveFile)
	http.HandleFunc("/", setSessionCookieReq)
	http.HandleFunc("/ws", serveWs)

	http.HandleFunc("/download", triggerDownload)
	// go logDownloads()
	http.ListenAndServe(":5000", nil)
	// http.HandleFunc("/", serveHome)
}

}

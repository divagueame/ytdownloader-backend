package main

import (
	"context"
	"log"
	"sync"
	"github.com/google/uuid"
	"github.com/lrstanley/go-ytdlp"
)

type Download struct {
	Id       string     `json:"id"`
	State    State      `json:"state"`
	Channel  chan State `json:"-"`
	Owner    string     `json:"-"`
	FilePath string     `json:"-"`
	Title    string     `json:"title"`
	Url      string     `json:"url"`
}

type DownloadManager struct {
	downloads map[string]*Download
	completed map[string]chan string
	mu        sync.Mutex
}

func NewDownloadManager() *DownloadManager {
	return &DownloadManager{
		downloads: make(map[string]*Download),
		completed: make(map[string]chan string),
	}
}

func (dm *DownloadManager) createDownload(url string, session_id string) *Download {

	download := &Download{
		Id:      uuid.New().String(),
		State:   Processing,
		Channel: make(chan State),
		Owner:   session_id,
		Url:     url,
	}

	dm.mu.Lock()
	dm.downloads[download.Id] = download
	dm.mu.Unlock()
	return download
}

func (dm *DownloadManager) getDownloadById(id string) *Download {
	return dm.downloads[id]
}

func (dm *DownloadManager) getActiveOwnerDownloads(owner_id string) map[string]*Download {
	filtered := make(map[string]*Download)

	log.Printf("Geting active downloads for: %s", owner_id)
	for k, download := range dm.downloads {

		if download.Owner == owner_id && download.State != 2 {
			filtered[k] = download
		}
	}

	return filtered
}

func (dm *DownloadManager) removeDownloadFromQueue(key string) {
	dm.mu.Lock()
	delete(dm.downloads, key)
	defer dm.mu.Unlock()
}

func (dm *DownloadManager) downloadFile(id string) {
	log.Print("downloading....", id)

	ytdlp.MustInstall(context.TODO(), nil)

	downloadR := dm.downloads[id]
	dl := ytdlp.New().
		ExtractAudio().
		AudioFormat("mp3").
		PrintJSON().
		NoPlaylist().
		Progress().
		Paths("downloads/").
		Verbose().
		// FormatSort("res,ext:mp4:m4a").
		// RecodeVideo("mp4").
		Output(id)

	response, err := dl.Run(context.TODO(), downloadR.Url)

	if err != nil {
		log.Print("Error downlaoding1", err)
		// downloadR.Channel <- Error
		// panic(err)
	} else {

		title, _ := extractTitle(response.Stdout)
		downloadR.Title = title
		downloadR.State = Done
		downloadR.FilePath = downloadR.buildPathToFile()
		dm.mu.Lock()
		if ch, exists := dm.completed[downloadR.Owner]; exists {
			select {
			case ch <- id:
			default:
			}
		}
		dm.mu.Unlock()
	}
}

func (dm *DownloadManager) logDownloads() {
	dm.mu.Lock()
	log.Print("Current Downloads")
	for key := range dm.downloads {
		log.Printf("Current Downloads %s", key)
		// fmt.Printf("download key: %d\n", key)
	}

	dm.mu.Unlock()
}

func extractTitle(stdout string) (string, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(stdout), &result)
	if err != nil {
		jsonStart := strings.Index(stdout, "{")
		jsonEnd := strings.LastIndex(stdout, "}")
		if jsonStart != -1 && jsonEnd != -1 && jsonEnd > jsonStart {
			jsonPart := stdout[jsonStart : jsonEnd+1]
			err = json.Unmarshal([]byte(jsonPart), &result)
			if err != nil {
				return "", fmt.Errorf("error unmarshaling JSON part: %v", err)
			}
		} else {
			return "", fmt.Errorf("error unmarshaling JSON and couldn't find valid JSON part: %v", err)
		}
	}

	title, ok := result["title"].(string)
	if !ok {
		return "", fmt.Errorf("title not found or not a string")
	}

	return title, nil
}

func (download *Download) buildPathToFile() string {
	return "downloads/" + "/" + download.Id + ".mp3"



}


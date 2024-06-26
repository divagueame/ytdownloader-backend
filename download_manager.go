package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/lrstanley/go-ytdlp"
	"log"
	"sync"
)

type Download struct {
	Id       string
	State    State
	Channel  chan State `json:"-"`
	Owner    string
	FilePath string
	Url      string `json:"url"`
}

type DownloadManager struct {
	downloads map[string]*Download
	mu        sync.Mutex
}

func NewDownloadManager() *DownloadManager {
	return &DownloadManager{
		downloads: make(map[string]*Download),
	}
}

func (dm *DownloadManager) createDownload(url string, session_id string) *Download {

	download := &Download{
		Id:      uuid.New().String(),
		State:   Idle,
		Channel: make(chan State),
		Owner:   session_id,
		Url:     url,
	}

	dm.mu.Lock()
	dm.downloads[download.Id] = download
	dm.mu.Unlock()
	return download
}
func (dm *DownloadManager) logDownloads() {
	log.Print("Current Downloads")
	for key := range dm.downloads {
		log.Printf("Current Downloads %s", key)
		// fmt.Printf("download key: %d\n", key)
	}

}

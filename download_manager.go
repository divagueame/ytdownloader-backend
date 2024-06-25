package main

import (
	"log"
	"sync"
)

type DownloadManager struct {
	downloads map[string]chan bool
	mu        sync.Mutex
}

func NewDownloadManager() *DownloadManager {
	return &DownloadManager{
		downloads: make(map[string]chan bool),
	}
}

func (dm *DownloadManager) logDownloads() {
	log.Print("Current Downloads")
	for key := range dm.downloads {
		log.Printf("Current Downloads %s", key)
		// fmt.Printf("download key: %d\n", key)
	}

}

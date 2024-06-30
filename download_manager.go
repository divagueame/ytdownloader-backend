package main

import (
	"context"
	"log"
	"sync"
	"github.com/google/uuid"
	"github.com/lrstanley/go-ytdlp"
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
func (dm *DownloadManager) downloadFile(id string) {

	// dm.mu.Lock()
	//  defer dm.mu.Unlock()

	ytdlp.MustInstall(context.TODO(), nil)

	downloadR := dm.downloads[id]
	dl := ytdlp.New().
		ExtractAudio().
		AudioFormat("mp3").
		PrintJSON().
		NoPlaylist().
		Progress().
		Paths(fmt.Sprintf("downloads/%s", id)).
		// FormatSort("res,ext:mp4:m4a").
		// RecodeVideo("mp4").
		Output("%(extractor)s - %(title)s.%(ext)s")

	// response, err := dl.Run(context.TODO(), downloadR.Url)
	_, err := dl.Run(context.TODO(), downloadR.Url)
	if err != nil {

		// log.Println("Error donwloading:", err)
		downloadR.Channel <- Error
		// panic(err)
	} else {

		// log.Println("---> Downloading DONE!", response.Stdout)

		// outputFilePath := "output.log"
		// err = ioutil.WriteFile(outputFilePath, []byte(download.Stdout), 0644)
		// if err != nil {
		// fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
		// return
		// }
		log.Println("---> Done. Filepath:", downloadR.FilePath)
		downloadR.FilePath = "meow"
		downloadR.Channel <- Done
		log.Println("---> Done after. Filepath:", downloadR.FilePath)
	}
	close(downloadR.Channel)
}

func (dm *DownloadManager) logDownloads() {
	log.Print("Current Downloads")
	for key := range dm.downloads {
		log.Printf("Current Downloads %s", key)
		// fmt.Printf("download key: %d\n", key)
	}

}

package main

import (
	"context"

	"fmt"
	"github.com/lrstanley/go-ytdlp"
	"log"
)

func download(url string, done chan<- bool) {
	fmt.Println("---> Downloading")
	ytdlp.MustInstall(context.TODO(), nil)

	dl := ytdlp.New().
		ExtractAudio().
		AudioFormat("mp3").
		// FormatSort("res,ext:mp4:m4a").
		// RecodeVideo("mp4").
		Output("%(extractor)s - %(title)s.%(ext)s")

	// _, err := dl.Run(context.TODO(), "https://www.youtube.com/watch?v=3lr9STd2Ryk")
	download_res, err := dl.Run(context.TODO(), url)
	if err != nil {

		log.Println("Error donwloading:", err)
		done <- false
		// panic(err)
	} else {
		log.Println("---> Downloading DONE!", download_res)
		done <- true
		close(done)
	}
}

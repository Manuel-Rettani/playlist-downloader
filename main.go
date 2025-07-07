package main

import (
	"log"
	"yt-playlist-downloader/constants"
	"yt-playlist-downloader/service"
)

func main() {
	videoUrl := "https://www.youtube.com/watch?v=k-x1n5v3RvM"
	setupService := service.NewSetupService()

	log.Println("setting up environment")
	err := setupService.DeleteFolder(constants.TempFolder)
	if err != nil {
		log.Fatal(err)
	}
	err = setupService.CreateFolder(constants.TempFolder)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("setting up environment done")

	downloaderService := service.NewDownloader()
	_, err = downloaderService.DownloadVideo(videoUrl)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Cleaning temp files...")
	err = setupService.DeleteFolder(constants.TempFolder)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Execution completed")
}

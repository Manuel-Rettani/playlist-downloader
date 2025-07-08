package main

import (
	"log"
	"playlist-downloader/constants"
	"playlist-downloader/service"
	"playlist-downloader/utils"
)

func main() {
	//playlistId := "PLqlu7ZxfTBeiDDYv-a0NTwO6lvc4C3kBv"
	setupEnvironment()
	videoUrl := "https://www.youtube.com/watch?v=k-x1n5v3RvM"

	downloaderService := service.NewDownloader()
	_, err := downloaderService.DownloadVideo(videoUrl)
	if err != nil {
		log.Fatal(err)
	}

	teardownEnvironment()
}

func setupEnvironment() {
	log.Println("setting up environment")
	err := utils.DeleteFolder(constants.TempFolder)
	if err != nil {
		log.Fatal(err)
	}
	err = utils.CreateFolder(constants.TempFolder)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("setting up environment done")
}

func teardownEnvironment() {
	log.Println("Cleaning temp files...")
	err := utils.DeleteFolder(constants.TempFolder)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Execution completed")
}

package main

import (
	"flag"
	"log"
	"playlist-downloader/client"
	"playlist-downloader/conf"
	"playlist-downloader/constants"
	"playlist-downloader/service"
	"playlist-downloader/utils"
)

var configPath string

func main() {
	setupEnvironment()

	flag.StringVar(&configPath, "config", "config.yml", "config path")
	flag.Parse()
	config, err := conf.FromYaml(configPath)
	if err != nil {
		panic(err)
	}

	youtubeClient := client.NewYoutubeClient(config.Keys.Youtube)
	downloader := service.NewDownloader()
	_ = service.NewYoutubeProcessor(youtubeClient, downloader)

	videoUrl := "https://www.youtube.com/watch?v=k-x1n5v3RvM"

	downloaderService := service.NewDownloader()
	_, err = downloaderService.DownloadVideo(videoUrl)
	if err != nil {
		panic(err)
	}

	teardownEnvironment()
}

func setupEnvironment() {
	log.Println("setting up environment")
	err := utils.DeleteFolder(constants.TempFolder)
	if err != nil {
		panic(err)
	}
	err = utils.CreateFolder(constants.TempFolder)
	if err != nil {
		panic(err)
	}
	log.Println("setting up environment done")
}

func teardownEnvironment() {
	log.Println("Cleaning temp files...")
	err := utils.DeleteFolder(constants.TempFolder)
	if err != nil {
		panic(err)
	}
	log.Println("Execution completed")
}

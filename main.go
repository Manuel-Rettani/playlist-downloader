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
		teardownEnvironment(nil)
		panic(err)
	}

	youtubeClient := client.NewYoutubeClient(config.YoutubeKey, constants.YoutubeApiBaseUrl)
	downloader := service.NewDownloader(config.MaxRetries)
	s3Service := service.NewS3Service(config.Aws.Region, config.Aws.S3.Bucket, config.Aws.AccessKey, config.Aws.SecretKey)
	mailService := service.NewMailService(config.Email.SmtpServer, config.Email.SmtpPort, config.Email.Email, config.Email.AppPassword, config.RequesterEmail)
	youtubeProcessor := service.NewYoutubeProcessor(youtubeClient, downloader, config.ChunkSize, s3Service, mailService)
	zipName, err := youtubeProcessor.Process(config.PlayListId)
	if err != nil {
		teardownEnvironment(nil)
		panic(err)
	}

	teardownEnvironment(&zipName)
}

func setupEnvironment() {
	log.Println("setting up environment...")
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

func teardownEnvironment(zipName *string) {
	log.Println("Cleaning temp files...")
	err := utils.DeleteFolder(constants.TempFolder)
	if err != nil {
		panic(err)
	}
	if zipName != nil {
		err = utils.DeleteFolder(*zipName)
		if err != nil {
			panic(err)
		}
	}
	log.Println("Execution completed")
}

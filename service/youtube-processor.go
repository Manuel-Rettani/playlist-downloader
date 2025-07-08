package service

import "playlist-downloader/client"

type YoutubeProcessor struct {
	client     client.IYoutubeClient
	downloader IDownloader
}

func NewYoutubeProcessor(client client.IYoutubeClient, downloader IDownloader) *YoutubeProcessor {
	return &YoutubeProcessor{
		client:     client,
		downloader: downloader,
	}
}

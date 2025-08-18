package service

import (
	"fmt"
	"log"
	"playlist-downloader/client"
	"playlist-downloader/constants"
	"playlist-downloader/models"
)

type IYoutubeProcessor interface {
	Process(playlistId string) error
}

type YoutubeProcessor struct {
	client     client.IYoutubeClient
	downloader IDownloader
	chunkSize  int
}

func NewYoutubeProcessor(client client.IYoutubeClient, downloader IDownloader, chunkSize int) *YoutubeProcessor {
	return &YoutubeProcessor{
		client:     client,
		downloader: downloader,
		chunkSize:  chunkSize,
	}
}

func (y YoutubeProcessor) Process(playlistId string) error {

	playlist, err := y.client.FetchPlaylist(playlistId, y.chunkSize, nil)
	if err != nil {
		return err
	}
	playListLength := playlist.PageInfo.TotalResults
	log.Printf("found %d videos", playListLength)
	videoURLs := make([]string, 0, playListLength)
	fillLinkSlice(&videoURLs, playlist)

	nextPageToken := playlist.NextPageToken
	for nextPageToken != "" {
		playlist, err = y.client.FetchPlaylist(playlistId, y.chunkSize, &nextPageToken)
		if err != nil {
			return err
		}
		fillLinkSlice(&videoURLs, playlist)
		nextPageToken = playlist.NextPageToken
	}

	//_, err := y.downloader.DownloadVideo(videoUrl)

	return nil
}

func fillLinkSlice(videoURLs *[]string, playlist *models.Playlist) {
	for _, item := range playlist.Items {
		videoId := item.Snippet.ResourceId.VideoId
		if videoId != "" {
			url := fmt.Sprintf(constants.YoutubeVideoLinkFormat, videoId)
			*videoURLs = append(*videoURLs, url)
		}
	}
}

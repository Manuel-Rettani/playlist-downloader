package service

import (
	"fmt"
	"log"
	"playlist-downloader/client"
	"playlist-downloader/constants"
	"playlist-downloader/models"
	"playlist-downloader/utils"
)

type IYoutubeProcessor interface {
	Process(playlistId string) error
}

type YoutubeProcessor struct {
	client      client.IYoutubeClient
	downloader  IDownloader
	chunkSize   int
	s3Service   IS3Service
	mailService IMailService
}

func NewYoutubeProcessor(client client.IYoutubeClient, downloader IDownloader, chunkSize int, s3Service IS3Service, mailService IMailService) *YoutubeProcessor {
	return &YoutubeProcessor{
		client:      client,
		downloader:  downloader,
		chunkSize:   chunkSize,
		s3Service:   s3Service,
		mailService: mailService,
	}
}

func (y YoutubeProcessor) Process(playlistId string) (string, error) {
	playlist, err := y.client.FetchPlaylist(playlistId, y.chunkSize, nil)
	if err != nil {
		return "", err
	}
	playlistLength := playlist.PageInfo.TotalResults
	videoURLs := make([]string, 0, playlistLength)
	fillLinkSlice(&videoURLs, playlist)

	nextPageToken := playlist.NextPageToken
	for nextPageToken != "" {
		playlist, err = y.client.FetchPlaylist(playlistId, y.chunkSize, &nextPageToken)
		if err != nil {
			return "", err
		}
		fillLinkSlice(&videoURLs, playlist)
		nextPageToken = playlist.NextPageToken
	}
	playlistLength = len(videoURLs)
	log.Printf("found %d videos", playlistLength)

	var failed []string
	for i := 1; i < len(videoURLs); i++ {
		_, err := y.downloader.DownloadVideoWithRetry(videoURLs[i-1])
		if err != nil {
			log.Printf("Skipping %s: %v", videoURLs[i-1], err)
			failed = append(failed, videoURLs[i-1])
			continue
		}
		if i%5 == 0 || i == len(videoURLs) {
			log.Printf("downloaded %d/%d videos", i, playlistLength)
		}
	}

	if len(failed) > 0 {
		log.Printf("Failed to download %d videos:", len(failed))
		for _, f := range failed {
			log.Println("   ", f)
		}
	}

	log.Println("Creating zip file...")
	zipName := fmt.Sprintf("%s.zip", playlistId)
	err = utils.ZipFolder("temp", zipName)
	if err != nil {
		return "", err
	}

	log.Println("Uploading zip file to S3...")
	fileLink, err := y.s3Service.Upload(zipName)
	if err != nil {
		return "", err
	}

	playlistInfo, err := y.client.GetPlaylistInfo(playlistId)
	if err != nil {
		return "", err
	}

	err = y.mailService.SendMail(fileLink, playlistInfo.Items[0].Snippet.Title)
	if err != nil {
		return "", err
	}

	return zipName, nil
}

func fillLinkSlice(videoURLs *[]string, playlist *models.YoutubeResponse) {
	for _, item := range playlist.Items {
		videoId := item.Snippet.ResourceId.VideoId
		videoOwnerChannelId := item.Snippet.VideoOwnerChannelId
		if videoId != "" && videoOwnerChannelId != "" {
			url := fmt.Sprintf(constants.YoutubeVideoLinkFormat, videoId)
			*videoURLs = append(*videoURLs, url)
		}
	}
}

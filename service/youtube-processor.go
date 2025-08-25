package service

import (
	"fmt"
	"log"
	"playlist-downloader/client"
	"playlist-downloader/constants"
	"playlist-downloader/models"
	"playlist-downloader/utils"
	"sync"
	"sync/atomic"
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
	videoUrls := make([]string, 0, playlistLength)
	fillLinkSlice(&videoUrls, playlist)

	nextPageToken := playlist.NextPageToken
	for nextPageToken != "" {
		playlist, err = y.client.FetchPlaylist(playlistId, y.chunkSize, &nextPageToken)
		if err != nil {
			return "", err
		}
		fillLinkSlice(&videoUrls, playlist)
		nextPageToken = playlist.NextPageToken
	}
	playlistLength = len(videoUrls)
	log.Printf("found %d videos", playlistLength)

	y.downloadVideos(playlistLength, videoUrls)

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

func (y YoutubeProcessor) downloadVideos(playlistLength int, videoUrls []string) {
	log.Println("Starting download...")

	var downloadCounter int32
	var wg sync.WaitGroup
	urls := make(chan string, playlistLength)
	failed := make(chan string, playlistLength)
	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go y.downloadWorker(urls, failed, &wg, &downloadCounter, playlistLength)
	}
	for _, url := range videoUrls {
		urls <- url
	}
	close(urls)

	wg.Wait()
	close(failed)

	var failedList []string
	for f := range failed {
		failedList = append(failedList, f)
	}

	if len(failedList) > 0 {
		log.Printf("Failed to download %d videos:", len(failedList))
		for _, f := range failedList {
			log.Println("   ", f)
		}
	}
}

func (y YoutubeProcessor) downloadWorker(
	urls <-chan string,
	failed chan<- string,
	wg *sync.WaitGroup,
	downloadCounter *int32,
	playlistLength int,
) {
	defer wg.Done()
	for url := range urls {
		_, err := y.downloader.DownloadVideoWithRetry(url)
		done := atomic.AddInt32(downloadCounter, 1)
		if err != nil {
			failed <- url
		}

		if done%5 == 0 || int(done) == playlistLength {
			log.Printf("Downloaded %d/%d videos", done, playlistLength)
		}
	}
}

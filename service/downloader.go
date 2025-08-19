package service

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

type IDownloader interface {
	DownloadVideo(videoUrl string) (string, error)
	DownloadVideoWithRetry(videoUrl string) (string, error)
}

func NewDownloader(maxRetries int) *Downloader {
	return &Downloader{
		maxRetries: maxRetries,
	}
}

type Downloader struct {
	maxRetries int
}

func (d Downloader) DownloadVideo(videoUrl string) (string, error) {
	videoName, err := getVideoName(videoUrl)
	if err != nil {
		return "", err
	}

	downloadCmd := exec.Command(
		"yt-dlp",
		"-x", "--audio-format", "mp3",
		"-o", "temp/%(title)s.%(ext)s",
		videoUrl,
	)
	downloadOutput, err := downloadCmd.CombinedOutput()
	if err != nil {
		if _, statErr := os.Stat(videoName); statErr == nil {
			log.Printf("yt-dlp returned error, but file was created: %s", videoName)
		} else {
			return "", fmt.Errorf("yt-dlp failed: %v\nOutput:\n%s", err, string(downloadOutput))
		}
	}

	return videoName, nil
}

func (d Downloader) DownloadVideoWithRetry(videoUrl string) (string, error) {
	var err error
	var name string
	for i := 0; i < d.maxRetries; i++ {
		name, err = d.DownloadVideo(videoUrl)
		if err == nil {
			return name, nil
		}
		log.Printf("retry %d for %s due to error: %v", i+1, videoUrl, err)
		time.Sleep(time.Second * time.Duration(i+1))
	}
	return "", fmt.Errorf("all retries failed for %s: %w", videoUrl, err)
}

func getVideoName(videoUrl string) (string, error) {
	metadataCmd := exec.Command("yt-dlp", "--dump-json", videoUrl)
	output, err := metadataCmd.Output()
	if err != nil {
		return "", fmt.Errorf("yt-dlp failed: %s", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %s", err)
	}
	return result["title"].(string), nil
}

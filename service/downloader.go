package service

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type IDownloader interface {
	DownloadVideo(videoUrl string) (string, error)
}

func NewDownloader() Downloader {
	return Downloader{}
}

type Downloader struct {
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
		return "", fmt.Errorf("yt-dlp failed: %v\nOutput:\n%s", err, string(downloadOutput))
	}

	return videoName, nil
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

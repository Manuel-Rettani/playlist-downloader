package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"playlist-downloader/constants"
	"playlist-downloader/models"
	"time"
)

type IYoutubeClient interface {
	FetchPlaylist(playlistId string) (*models.Playlist, error)
}

type YoutubeClient struct {
	apikey string
}

func NewYoutubeClient(apikey string) *YoutubeClient {
	return &YoutubeClient{apikey: apikey}
}

func (*YoutubeClient) FetchPlaylist(playlistId string) (*models.Playlist, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	url := fmt.Sprintf("%s/%s", constants.YoutubeApiBaseUrl, playlistId)
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	var playlist models.Playlist
	if err := json.NewDecoder(resp.Body).Decode(&playlist); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return &playlist, nil
}

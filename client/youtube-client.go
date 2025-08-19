package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"playlist-downloader/models"
	"strconv"
	"time"
)

type IYoutubeClient interface {
	FetchPlaylist(playlistId string, pageSize int, pageToken *string) (*models.YoutubeResponse, error)
	GetPlaylistInfo(playlistId string) (*models.YoutubeResponse, error)
}

type YoutubeClient struct {
	apikey     string
	httpClient *http.Client
	baseUrl    string
}

func NewYoutubeClient(apikey string, baseUrl string) *YoutubeClient {
	return &YoutubeClient{
		apikey: apikey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseUrl: baseUrl,
	}
}

func (c *YoutubeClient) FetchPlaylist(playlistId string, pageSize int, pageToken *string) (*models.YoutubeResponse, error) {
	endpoint, _ := url.Parse(fmt.Sprintf("%s/playlistItems", c.baseUrl))
	query := endpoint.Query()
	query.Set("key", c.apikey)
	query.Set("part", "snippet")
	query.Set("playlistId", playlistId)
	query.Set("maxResults", strconv.Itoa(pageSize))
	if pageToken != nil {
		query.Set("pageToken", *pageToken)
	}
	endpoint.RawQuery = query.Encode()

	resp, err := c.httpClient.Get(endpoint.String())
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status: %s, body: %s", resp.Status, string(body))
	}

	var playlist models.YoutubeResponse
	if err := json.NewDecoder(resp.Body).Decode(&playlist); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return &playlist, nil
}

func (c *YoutubeClient) GetPlaylistInfo(playlistId string) (*models.YoutubeResponse, error) {
	endpoint, _ := url.Parse(fmt.Sprintf("%s/playlists", c.baseUrl))
	query := endpoint.Query()
	query.Set("key", c.apikey)
	query.Set("part", "id,snippet")
	query.Set("id", playlistId)
	endpoint.RawQuery = query.Encode()

	resp, err := c.httpClient.Get(endpoint.String())
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status: %s, body: %s", resp.Status, string(body))
	}

	var playlistInfo models.YoutubeResponse
	if err := json.NewDecoder(resp.Body).Decode(&playlistInfo); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return &playlistInfo, nil
}

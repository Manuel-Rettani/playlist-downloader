package models

import "encoding/json"

type Playlist struct {
	Kind          string         `json:"kind"`
	Etag          string         `json:"etag"`
	NextPageToken string         `json:"nextPageToken"`
	PrevPageToken string         `json:"prevPageToken"`
	Items         []PlaylistItem `json:"items"`
	PageInfo      PageInfo       `json:"pageInfo"`
}

type PlaylistItem struct {
	Kind    string  `json:"kind"`
	Etag    string  `json:"etag"`
	Id      string  `json:"id"`
	Snippet Snippet `json:"snippet"`
}

type Snippet struct {
	PublishedAt            string          `json:"publishedAt"`
	ChannelId              string          `json:"channelId"`
	Title                  string          `json:"title"`
	Description            string          `json:"description"`
	Thumbnails             json.RawMessage `json:"thumbnails"`
	ChannelTitle           string          `json:"channelTitle"`
	PlaylistId             string          `json:"playlistId"`
	Position               int             `json:"position"`
	ResourceId             ResourceId      `json:"resourceId"`
	VideoOwnerChannelTitle string          `json:"videoOwnerChannelTitle"`
	VideoOwnerChannelId    string          `json:"videoOwnerChannelId"`
}

type ResourceId struct {
	Kind    string `json:"kind"`
	VideoId string `json:"videoId"`
}

type PageInfo struct {
	TotalResults   int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}

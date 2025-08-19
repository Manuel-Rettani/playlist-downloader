package models

type YoutubeResponse struct {
	Kind          string   `json:"kind"`
	Etag          string   `json:"etag"`
	NextPageToken string   `json:"nextPageToken"`
	PrevPageToken string   `json:"prevPageToken"`
	Items         []Item   `json:"items"`
	PageInfo      PageInfo `json:"pageInfo"`
}

type Item struct {
	Kind    string  `json:"kind"`
	Etag    string  `json:"etag"`
	Id      string  `json:"id"`
	Snippet Snippet `json:"snippet"`
}

type Snippet struct {
	PublishedAt            string     `json:"publishedAt"`
	ChannelId              string     `json:"channelId"`
	Title                  string     `json:"title"`
	Description            string     `json:"description"`
	ChannelTitle           string     `json:"channelTitle"`
	PlaylistId             string     `json:"playlistId"`
	Position               int        `json:"position"`
	ResourceId             ResourceId `json:"resourceId"`
	VideoOwnerChannelTitle string     `json:"videoOwnerChannelTitle"`
	VideoOwnerChannelId    string     `json:"videoOwnerChannelId"`
	Localized              Localized  `json:"localized"`
}

type ResourceId struct {
	Kind    string `json:"kind"`
	VideoId string `json:"videoId"`
}

type PageInfo struct {
	TotalResults   int `json:"totalResults"`
	ResultsPerPage int `json:"resultsPerPage"`
}

type Localized struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

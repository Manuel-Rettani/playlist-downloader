package service

import (
	"github.com/stretchr/testify/require"
	"playlist-downloader/constants"
	"playlist-downloader/utils"
	"testing"
)

var downloaderService = NewDownloader(3)

func TestDownload(t *testing.T) {
	videoUrl := "https://www.youtube.com/watch?v=k-x1n5v3RvM"
	name, err := downloaderService.DownloadVideo(videoUrl)
	require.NoError(t, err)
	require.Equal(t, "Men at Work - Who Can it Be Now? (Lyrics)", name)
	utils.Teardown(constants.TempFolder)
}

func TestDownloadLinkNotFound(t *testing.T) {
	videoUrl := "https://www.youtube.com/watch?v=invalid-video-id"
	name, err := downloaderService.DownloadVideo(videoUrl)
	require.Error(t, err)
	require.Empty(t, name)
	utils.Teardown(constants.TempFolder)
}

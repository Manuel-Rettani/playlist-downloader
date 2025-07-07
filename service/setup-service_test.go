package service

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateFolder(t *testing.T) {
	err := setupService.CreateFolder("test")
	require.NoError(t, err)
	Teardown("test")
}

func TestCreateFolderAlreadyExists(t *testing.T) {
	err := setupService.CreateFolder("test")
	err = setupService.CreateFolder("test")
	require.Error(t, err)
	Teardown("test")
}

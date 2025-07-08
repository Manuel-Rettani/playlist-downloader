package utils

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateFolder(t *testing.T) {
	err := CreateFolder("test")
	require.NoError(t, err)
	Teardown("test")
}

func TestCreateFolderAlreadyExists(t *testing.T) {
	err := CreateFolder("test")
	err = CreateFolder("test")
	require.Error(t, err)
	Teardown("test")
}

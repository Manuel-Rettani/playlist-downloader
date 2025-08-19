package utils

import (
	"archive/zip"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateFolder(t *testing.T) {
	err := CreateFolder("test")
	require.NoError(t, err)
	defer Teardown("test")
}

func TestCreateFolderAlreadyExists(t *testing.T) {
	err := CreateFolder("test")
	err = CreateFolder("test")
	require.Error(t, err)
	defer Teardown("test")
}

func TestDeleteFolder(t *testing.T) {
	err := CreateFolder("test")
	err = DeleteFolder("test")
	require.NoError(t, err)
}

func TestDeleteFolderFolderNotExists(t *testing.T) {
	err := DeleteFolder("test")
	require.NoError(t, err)
}

func TestZipFolder(t *testing.T) {
	sourceDir := "test"
	err := CreateFolder(sourceDir)
	require.NoError(t, err)
	defer Teardown(sourceDir)

	files := []struct {
		name    string
		content string
	}{
		{"file1.txt", "hello"},
		{"file2.txt", "world"},
	}
	for _, f := range files {
		err := os.WriteFile(filepath.Join(sourceDir, f.name), []byte(f.content), 0644)
		require.NoError(t, err)
	}

	zipFile := "test.zip"
	defer os.Remove(zipFile)

	err = ZipFolder(sourceDir, zipFile)
	require.NoError(t, err)

	r, err := zip.OpenReader(zipFile)
	require.NoError(t, err)
	defer r.Close()

	require.Len(t, r.File, len(files))
	for _, f := range files {
		found := false
		for _, zf := range r.File {
			if zf.Name == f.name {
				rc, err := zf.Open()
				require.NoError(t, err)
				data, err := io.ReadAll(rc)
				require.NoError(t, err)
				_ = rc.Close()
				require.Equal(t, f.content, string(data))
				found = true
			}
		}
		require.True(t, found, "file %s not found in zip", f.name)
	}
}

func TestZipFolderNonExistingSource(t *testing.T) {
	zipFile := "test.zip"
	defer os.Remove(zipFile)

	err := ZipFolder("nonexistent", zipFile)
	require.Error(t, err)
}

func TestZipFolderEmptyFolder(t *testing.T) {
	sourceDir := "emptydir"
	err := CreateFolder(sourceDir)
	require.NoError(t, err)
	defer Teardown(sourceDir)

	zipFile := "empty.zip"
	defer os.Remove(zipFile)

	err = ZipFolder(sourceDir, zipFile)
	require.NoError(t, err)

	r, err := zip.OpenReader(zipFile)
	require.NoError(t, err)
	defer r.Close()

	require.Len(t, r.File, 0, "zip should not contain any files")
}

func TestGetFileName(t *testing.T) {
	fileName := GetFileName("README.md")
	require.Equal(t, "README", fileName)
}

func TestGetFileNameFileInsideFolder(t *testing.T) {
	fileName := GetFileName("test/testfile.zip")
	require.Equal(t, "testfile", fileName)
}

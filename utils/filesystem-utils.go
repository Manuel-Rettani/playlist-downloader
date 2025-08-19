package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CreateFolder(name string) error {
	err := os.Mkdir(name, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func DeleteFolder(name string) error {
	err := os.RemoveAll(name)
	if err != nil {
		return fmt.Errorf("failed to delete %s folder: %s", name, err)
	}
	return nil
}

func ZipFolder(source, destination string) error {
	zipFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer func(zipFile *os.File) {
		err := zipFile.Close()
		if err != nil {
			log.Printf("failed to close zip file: %s", err)
		}
	}(zipFile)

	archive := zip.NewWriter(zipFile)

	err = filepath.WalkDir(source, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Printf("failed to close file: %s", err)
			}
		}(file)

		f, err := archive.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(f, file)
		return err
	})
	if err != nil {
		return err
	}

	return archive.Close()
}

func GetFileName(filePath string) string {
	baseName := filepath.Base(filePath)
	ext := filepath.Ext(baseName)
	return strings.TrimSuffix(baseName, ext)
}

package utils

import (
	"fmt"
	"os"
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

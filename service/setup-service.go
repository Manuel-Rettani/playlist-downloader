package service

import (
	"fmt"
	"os"
)

type ISetupService interface {
	CreateFolder(name string) error
	DeleteFolder(name string) error
}

func NewSetupService() ISetupService {
	return &SetupService{}
}

type SetupService struct {
}

func (s *SetupService) CreateFolder(name string) error {
	err := os.Mkdir(name, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (s *SetupService) DeleteFolder(name string) error {
	err := os.RemoveAll(name)
	if err != nil {
		return fmt.Errorf("failed to delete %s folder: %s", name, err)
	}
	return nil
}

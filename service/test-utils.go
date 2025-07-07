package service

var setupService = NewSetupService()
var downloaderService = NewDownloader()

func Teardown(dirName string) {
	err := setupService.DeleteFolder(dirName)
	if err != nil {
		return
	}
}

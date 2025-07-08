package utils

func Teardown(dirName string) {
	err := DeleteFolder(dirName)
	if err != nil {
		return
	}
}

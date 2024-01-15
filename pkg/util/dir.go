package util

import (
	"os"
)

func WriteDataToFile(filename string, dataBytes []byte) {
	var err = os.WriteFile(filename, dataBytes, 0644)
	HandleError(err)
}

func ReadFile(filename string) []byte {
	file, err := os.ReadFile(filename)
	HandleError(err)
	return file
}
func CreateDirectory(dir_name string) {
	_ = os.Mkdir(dir_name, 0700)
}

func DoesDirectoryExist(dir_name string) bool {
	if _, err := os.Stat(dir_name); os.IsNotExist(err) {
		return false
	}
	return true
}
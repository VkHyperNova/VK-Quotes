package dir

import (
	"os"
	"vk-quotes/pkg/print"
	"vk-quotes/pkg/global"
)

func WriteDataToFile(filename string, dataBytes []byte) {
	var err = os.WriteFile(filename, dataBytes, 0644)
	print.HandleError(err)
}

func ReadFile(filename string) []byte {
	file, err := os.ReadFile(filename)
	print.HandleError(err)
	return file
}

func ValidateRequiredFiles() {
	if !DoesDirectoryExist(global.DatabasePath) {
		CreateDirectory("database")
		WriteDataToFile(global.DatabasePath, []byte("[]"))
		print.PrintRed("New Database Created!\n")
	} 
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
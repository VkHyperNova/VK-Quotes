package config

import (
	"path/filepath"
)

var ProgramVersion = "1.24"

/* All Paths */
var DefaultContent = `{"quotes": []}`
var	file = "quotes.json"
var BaseDB = "QUOTES"
var BaseLocal = "DATABASES"
var	BaseBackup = "/media/veikko/VK DATA/"

var LocalFile = filepath.Join(BaseLocal, BaseDB, file)
var BackupFile = filepath.Join(BaseBackup, BaseLocal, BaseDB, file)

var ReadCounter string
var Counter int
var MainQuoteID int
var UserInputs []string

var Messages []string
func AddMessage(message string) {
	Messages = append(Messages, message)
}

func FormatMessages() string {

	formattedString := ""

	for _, message := range Messages{
		formattedString += message + "\n------------------------------\n"
	}
	
	return formattedString
}

var RandomIDs []int
func DeleteUsedID(index int) {
	RandomIDs = append(RandomIDs[:index], RandomIDs[index+1:]...)
}



package config

var ProgramVersion = "1.24"

/* All Paths */
var SaveQuotesPath = "./QuotesDB/quotes.json"
var SaveSimilarPath = "./QuotesDB/similar.json"
var SaveFolderName = "QuotesDB"
var CopyFilePathLinux = "/media/veikko/VK DATA/DATABASES/QUOTESdb/quotes.json"
var CopyFilePathWindows = "D:\\DATABASES\\QUOTESdb\\quotes.json"

var ReadCounter string
var Counter int
var RandomIDs []int

var ID int
var Messages []string
var UserInputs []string

func FormatMessages() string {

	formattedString := ""
	for _, m := range Messages {
		formattedString += m + "\n"
	}

	return formattedString
}

func ResetReadCounter() {
	Counter = 0
	ReadCounter = ""
}

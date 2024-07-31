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

var MainQuoteID int
var Messages []string
var UserInputs []string

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
)

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

package config

import "strconv"

var ProgramVersion = "1.24"

/* All Paths */

const LocalPath = "./QUOTES/quotes.json"
var ReadCounter string
var Counter int
var MainQuoteID int
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


var Messages []string
func AddMessage(message string) {
	Messages = append(Messages, message)
}

func FormatMessages() string {

	formattedString := ""

	for nr, message := range Messages{
		formattedString += strconv.Itoa(nr+1) + ". " + message + "\n"
	}
	
	return formattedString
}

var RandomIDs []int
func DeleteUsedID(index int) {
	RandomIDs = append(RandomIDs[:index], RandomIDs[index+1:]...)
}

func DeleteAllRandomIDs() {
	RandomIDs = RandomIDs[:0]
}

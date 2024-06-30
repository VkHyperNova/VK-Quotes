package cmd

import (
	"strconv"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"

	"github.com/peterh/liner"
)

func Abort(a string) bool {
	if a == "q" {
		ErrorMsg = "<< previous action aborted by user. >>"
		return true
	}
	return false
}

func UserInput(Database *[]db.Quotes, id int) bool {

	util.PrintGray("\n(" + strconv.Itoa(len(*Database)) + ")\n")

	/* Scan Quote */
	input := ScanOrEditWithLiner("Quote", "")
	if Abort(input) {
		return false
	}
	if db.FindDublicates(input, Database) != -1 {
		ErrorMsg = "<< there are dublicates in database. >>"
		CurrentQuoteIndex = db.FindDublicates(input, Database)
		return false
	}
	Quote = input

	/* Scan Author */
	input = ScanOrEditWithLiner("Author", "")
	if Abort(input) {
		return false
	}
	Author = input
	
	/* Scan Language */
	input = ScanOrEditWithLiner("Language", "English")
	if Abort(input) {
		return false
	}
	Language = input

	return true
}

func EditUserInput(Database *[]db.Quotes, id int) {
	index := db.FindIndex(id, Database)
	Quote = ScanOrEditWithLiner("Quote", (*Database)[index].QUOTE)
	Author = ScanOrEditWithLiner("Author", (*Database)[index].AUTHOR)
	Language = ScanOrEditWithLiner("Language", (*Database)[index].LANGUAGE)
}

func ScanOrEditWithLiner(name string, editableString string) string {

	line := liner.NewLiner()
	defer line.Close()

	if editableString != "" {
		input, err := line.PromptWithSuggestion("   "+name+": ", editableString, -1)
		util.HandleError(err)
		return input
	} else {
		input, err := line.Prompt("   "+name+": ")
		util.HandleError(err)
		return input
	}
}

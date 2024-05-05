package cmd

import (
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func Abort(input string) bool {
	if input == "q" {
		ErrorMsg = "<< previous action aborted by user. >>"
		return true
	}
	return false
}

func GetQuote(Database *[]db.Quotes, id int) (string, bool) {

	if id > 0 {
		index := db.FindIndex(id, Database)
		util.PrintGreen("\n\"" + (*Database)[index].QUOTE + "\"\n")
	} else {
		util.PrintGreen("\n\"" + "Unknown" + "\"\n")
	}

	util.PrintCyan("Quote: ")
	input := util.ScanUserInput()

	if Abort(input) {
		return "", false
	}

	if db.FindDublicates(input, Database) != -1 {
		ErrorMsg = "<< there are dublicates in database. >>"
		CurrentQuoteIndex = db.FindDublicates(input, Database)
		return "", false
	}

	return input, true
}

func GetAuthor(Database *[]db.Quotes, id int) (string, bool) {

	if id > 0 {
		index := db.FindIndex(id, Database)
		util.PrintGreen("\n\"" + (*Database)[index].AUTHOR + "\"\n")
	} else {
		util.PrintGreen("\n\"" + "Unknown" + "\"\n")
	}

	util.PrintCyan("Author: ")
	input := util.ScanUserInput()

	if Abort(input) {
		return "", false
	}
	return input, true
}

func GetLanguage(Database *[]db.Quotes, id int) (string, bool) {

	if id > 0 {
		index := db.FindIndex(id, Database)
		util.PrintGreen("\n\"" + (*Database)[index].LANGUAGE + "\"\n")
	} else {
		util.PrintGreen("\n\"" + "English" + "\"\n")
	}

	util.PrintCyan("Language: ")
	input := util.ScanUserInput()

	if Abort(input) {
		return "", false
	}
	return input, true
}

func UserInput(Database *[]db.Quotes, id int) ([]string, bool) {

	quote, validation := GetQuote(Database, id)
	if !validation {
		return []string{""}, false
	}

	author, validation := GetAuthor(Database, id)
	if !validation {
		return []string{""}, false
	}
	
	language, validation := GetLanguage(Database, id)
	if !validation {
		return []string{""}, false
	}
	
	return []string{quote, author, language}, true
}
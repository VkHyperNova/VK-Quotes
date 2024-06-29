package cmd

import (
	"strconv"
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

func GetQuote(Database *[]db.Quotes) (string, bool) {

	input := util.ScanUserInputWithLiner("   Quote: ")

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

func GetAuthor(Database *[]db.Quotes) (string, bool) {

	input := util.ScanUserInputWithLiner("   Author: ")

	if Abort(input) {
		return "", false
	}
	return input, true
}

func GetLanguage(Database *[]db.Quotes) (string, bool) {

	input := util.ScanUserInputWithLiner("   Language: ")

	if Abort(input) {
		return "", false
	}
	return input, true
}

func UserInput(Database *[]db.Quotes, id int) ([]string, bool) {

	util.PrintGray("\n(" + strconv.Itoa(len(*Database)) + ")\n")

	quote, validation := GetQuote(Database)
	if !validation {
		return []string{""}, false
	}

	author, validation := GetAuthor(Database)
	if !validation {
		return []string{""}, false
	}

	language, validation := GetLanguage(Database)
	if !validation {
		return []string{""}, false
	}

	return []string{quote, author, language}, true
}

func EditUserInput(Database *[]db.Quotes, id int) []string {
	index := db.FindIndex(id, Database)
	quote := util.EditUserInputWithLiner("Quote: ", (*Database)[index].QUOTE)
	author := util.EditUserInputWithLiner("Author: ", (*Database)[index].AUTHOR)
	language := util.EditUserInputWithLiner("Language: ", (*Database)[index].LANGUAGE)
	return []string{quote, author, language}
}

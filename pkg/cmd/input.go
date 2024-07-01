package cmd

import (
	"strconv"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func UserInput(Database *[]db.Quotes, id int) bool {

	util.PrintGray("\n(" + strconv.Itoa(len(*Database)) + ")\n")

	/* Scan Quote */
	input := util.ScanOrEditWithLiner("Quote", "")
	if util.Abort(input) {
		IsMessage = "<< previous action aborted by user. >>"
		return false
	}
	if db.FindDublicates(input, Database) != -1 {
		IsMessage = "<< there are dublicates in database. >>"
		IsEditedQuote = db.FindDublicates(input, Database)
		return false
	}
	IsQuote = input

	/* Scan Author */
	input = util.ScanOrEditWithLiner("Author", "")
	if util.Abort(input) {
		IsMessage = "<< previous action aborted by user. >>"
		return false
	}
	IsAuthor = input

	/* Scan Language */
	input = util.ScanOrEditWithLiner("Language", "English")
	if util.Abort(input) {
		IsMessage = "<< previous action aborted by user. >>"
		return false
	}
	IsLanguage = input

	return true
}

func EditUserInput(Database *[]db.Quotes, id int) {
	index := db.FindIndex(id, Database)
	IsQuote = util.ScanOrEditWithLiner("Quote", (*Database)[index].QUOTE)
	IsAuthor = util.ScanOrEditWithLiner("Author", (*Database)[index].AUTHOR)
	IsLanguage = util.ScanOrEditWithLiner("Language", (*Database)[index].LANGUAGE)
}

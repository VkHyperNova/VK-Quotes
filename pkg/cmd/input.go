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
		ErrorMsg = "<< previous action aborted by user. >>"
		return false
	}
	if db.FindDublicates(input, Database) != -1 {
		ErrorMsg = "<< there are dublicates in database. >>"
		CurrentQuoteIndex = db.FindDublicates(input, Database)
		return false
	}
	Quote = input

	/* Scan Author */
	input = util.ScanOrEditWithLiner("Author", "")
	if util.Abort(input) {
		ErrorMsg = "<< previous action aborted by user. >>"
		return false
	}
	Author = input
	
	/* Scan Language */
	input = util.ScanOrEditWithLiner("Language", "English")
	if util.Abort(input) {
		ErrorMsg = "<< previous action aborted by user. >>"
		return false
	}
	Language = input

	return true
}

func EditUserInput(Database *[]db.Quotes, id int) {
	index := db.FindIndex(id, Database)
	Quote = util.ScanOrEditWithLiner("Quote", (*Database)[index].QUOTE)
	Author = util.ScanOrEditWithLiner("Author", (*Database)[index].AUTHOR)
	Language = util.ScanOrEditWithLiner("Language", (*Database)[index].LANGUAGE)
}
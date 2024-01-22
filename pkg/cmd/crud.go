package cmd

import (
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func Ask() (string, string, string) {

	quote := util.GetInput("Quote: ")

	if db.CheckDublicates(quote) != -1 {
		util.PrintRed("\n<< This quote is in the database >>\n")
		PrintQuote(db.CheckDublicates(quote))
		ReturnToCMD()
	}

	author := util.GetInput("Auhtor: ")
	language := util.GetInput("Language: ")

	return quote, author, language
}



func Add(quote, author, language, databasePath string) bool {

	uniqueID := 1

	if len(db.DATABASE) != 0 {
		uniqueID = db.DATABASE[len(db.DATABASE)-1].ID + 1 // Get Last item + 1
	}

	db.LastItemID = uniqueID

	NewQuote := db.Quotes{
		ID:       uniqueID,
		QUOTE:    util.FillEmptyInput(quote, "Unknown"),
		AUTHOR:   util.FillEmptyInput(author, "Unknown"),
		LANGUAGE: util.FillEmptyInput(language, "Unknown"),
		DATE:     util.GetFormattedDate(),
	}

	db.DATABASE = append(db.DATABASE, NewQuote)
	db.SaveDB(databasePath)

	return true
}

func Update(updateID int, updatedQuote, updatedAuthor, updatedLanguage, DatabasePath string) bool {

	db.LastItemID = updateID

	index := db.SearchIndexByID(updateID)

	if index == -1 {
		util.PrintRed("\nIndex out of range!\n")
		ReturnToCMD()
	}

	db.DATABASE[index].QUOTE = util.FillEmptyInput(updatedQuote, db.DATABASE[index].QUOTE)
	db.DATABASE[index].AUTHOR = util.FillEmptyInput(updatedAuthor, db.DATABASE[index].AUTHOR)
	db.DATABASE[index].LANGUAGE = util.FillEmptyInput(updatedLanguage, db.DATABASE[index].LANGUAGE)
	
	db.SaveDB(DatabasePath)

	return true
}

func Delete(deleteID int, DatabasePath string) bool {
	index := db.SearchIndexByID(deleteID)

	if index == -1 {
		util.PrintRed("\nIndex out of range!\n")
		ReturnToCMD()
	}

	db.DATABASE = append(db.DATABASE[:index], db.DATABASE[index+1:]...)
	db.SaveDB(DatabasePath)

	return true
}



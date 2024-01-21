package cmd

import (
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func Ask() (string, string, string) {
	Quote := util.GetInput("Quote: ")
	if db.CheckDublicates(Quote) != -1 {
		util.PrintRed("\n<< This quote is in the database >>\n")
		PrintQuote(db.CheckDublicates(Quote))
		util.PressAnyKey()
		util.ClearScreen()
		CMD()
	}

	Author := util.GetInput("Auhtor: ")
	if Author == "" {
		Author = "Unknown"
	}

	Language := util.GetInput("Language: ")
	if Language == "" {
		Language = "Unknown"
	}

	return Quote, Author, Language
}

func Add(quote, author, language, databasePath string) bool {

	uniqueID := 1

	if len(db.DATABASE) != 0 {
		uniqueID = db.DATABASE[len(db.DATABASE)-1].ID + 1 // Get Last item + 1
	}

	db.LastItemID = uniqueID

	NewQuote := db.Quotes{
		ID:       uniqueID,
		QUOTE:    quote,
		AUTHOR:   author,
		LANGUAGE: language,
		DATE:     util.GetFormattedDate(),
	}

	db.DATABASE = append(db.DATABASE, NewQuote)
	db.SaveDB(databasePath)

	return true
}

func Update(updatedID int, updatedQuote, updatedAuthor, updatedLanguage, DatabasePath string) bool {

	index := db.SearchIndexByID(updatedID)

	if index == -1 {
		util.PrintRed("\nIndex out of range!\n")
		return false
	}

	PrintQuote(index)
	db.LastItemID = updatedID

	if updatedQuote == "" {
		updatedQuote = db.DATABASE[index].QUOTE
	}

	if updatedAuthor == "" {
		updatedAuthor = db.DATABASE[index].AUTHOR
	}

	if updatedLanguage == "" {
		updatedLanguage = db.DATABASE[index].LANGUAGE
	}

	db.DATABASE[index].QUOTE = updatedQuote
	db.DATABASE[index].AUTHOR = updatedAuthor
	db.DATABASE[index].LANGUAGE = updatedLanguage

	db.SaveDB(DatabasePath)

	return true
}

func Delete(deleteID int, DatabasePath string) bool {

	index := db.SearchIndexByID(deleteID)

	if index == -1 {
		util.PrintRed("\nIndex out of range!\n")
		return false
	}

	PrintQuote(index)
	db.DATABASE = append(db.DATABASE[:index], db.DATABASE[index+1:]...)
	db.SaveDB(DatabasePath)

	return true
}

func ReturnToCMD() {
	util.PressAnyKey()
	util.ClearScreen()
	CMD()
}

package cmd

import (
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func GetQuoteDetails(Database *[]db.Quotes) (string, string, string) {

	quote := util.GetInput("Quote: ")
	if quote == "q" {
		CMD()
	}

	if db.CheckDublicates(quote, Database) != -1 {
		util.PrintRed("\n<< This quote is in the database >>\n")
		PrintQuote(db.CheckDublicates(quote, Database), Database)
		util.PressAnyKey()
		CMD()
	}

	author := util.GetInput("Auhtor: ")
	if author == "q" {
		CMD()
	}
	language := util.GetInput("Language: ")

	if language == "q" {
		CMD()
	}

	return quote, author, language
}

func Add(id int, quote, author, language string, Database *[]db.Quotes) bool {

	NewQuote := db.Quotes{
		ID:       id,
		QUOTE:    util.FillEmptyInput(quote, "Unknown"),
		AUTHOR:   util.FillEmptyInput(author, "Unknown"),
		LANGUAGE: util.FillEmptyInput(language, "Unknown"),
		DATE:     util.GetFormattedDate(),
	}

	*Database = append(*Database, NewQuote)

	return true
}

func Update(index int, updatedQuote, updatedAuthor, updatedLanguage, DatabasePath string, Database *[]db.Quotes) bool {

	(*Database)[index].QUOTE = util.FillEmptyInput(updatedQuote, (*Database)[index].QUOTE)
	(*Database)[index].AUTHOR = util.FillEmptyInput(updatedAuthor, (*Database)[index].AUTHOR)
	(*Database)[index].LANGUAGE = util.FillEmptyInput(updatedLanguage, (*Database)[index].LANGUAGE)

	return true
}

func Delete(index int, DatabasePath string, Database *[]db.Quotes) bool {

	(*Database) = append((*Database)[:index], (*Database)[index+1:]...)

	return true
}

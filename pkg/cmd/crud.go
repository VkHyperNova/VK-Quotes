package cmd

import (
	"strings"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func Add(id int, Database *[]db.Quotes) bool {

	NewQuote := db.Quotes{
		ID:       id,
		QUOTE:    GetInput("Quote: ", Database),
		AUTHOR:   GetInput("Author: ", Database),
		LANGUAGE: GetInput("Language: ", Database),
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

func Search(Database *[]db.Quotes, searchString string) {

	for key, value := range *Database {
		if strings.Contains(strings.ToUpper(value.AUTHOR), strings.ToUpper(searchString)) {
			PrintQuote(key, Database)
		}
	}
}
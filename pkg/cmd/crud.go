package cmd

import (
	"strings"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func Add(id int, Database *[]db.Quotes) bool {

	quoteDetails := []string{"Quote: ", "Author: ", "Language: "}
	var inputs []string

    for _, value := range quoteDetails {
        input, err := db.GetInput(value, Database)
		if err != nil {
			LastError = err.Error()
			util.ClearScreen()
			CMD()
		}
		inputs = append(inputs, input)	
    }

	NewQuote := db.Quotes{
		ID:       id,
		QUOTE:    util.FillEmptyInput(inputs[0], "Unknown"),
		AUTHOR:   util.FillEmptyInput(inputs[1], "Unknown"),
		LANGUAGE: util.FillEmptyInput(inputs[2], "Unknown"),
		DATE:     util.GetFormattedDate(),
	}

	*Database = append(*Database, NewQuote)

	return true
}

func Update(index int, DatabasePath string, Database *[]db.Quotes) bool {

	quoteDetails := []string{"Quote: ", "Author: ", "Language: "}
	var inputs []string

    for _, value := range quoteDetails {
        input, err := db.GetInput(value, Database)
		if err != nil {
			LastError = err.Error()
			util.ClearScreen()
			CMD()
		}
		inputs = append(inputs, input)	
    }

	(*Database)[index].QUOTE = util.FillEmptyInput(inputs[0], (*Database)[index].QUOTE)
	(*Database)[index].AUTHOR = util.FillEmptyInput(inputs[1], (*Database)[index].AUTHOR)
	(*Database)[index].LANGUAGE = util.FillEmptyInput(inputs[2], (*Database)[index].LANGUAGE)

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
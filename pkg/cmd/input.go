package cmd

import (
	"fmt"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func UserInput(quotes *db.Quotes) ([]string, bool) {

quote:

	var inputs []string

	/* Scan Quote */

	inputQuote := util.ScanOrEditWithLiner("Quote", "")
	if util.Abort(inputQuote) {
		IsMessage = "<< previous action aborted by user. >>"
		return []string{}, false
	}

	if quotes.FindDuplicates(inputQuote) != -1 {
		IsMessage = "<< there are dublicates in database. >>"
		UpdatedQuoteID = quotes.FindDuplicates(inputQuote)
		return []string{}, false
	}

	/* Scan Author */
author:
	inputAuthor := util.ScanOrEditWithLiner("Author", "")
	if util.MoveBack(inputAuthor) {
		goto quote
	}
	if util.Abort(inputAuthor) {
		IsMessage = "<< previous action aborted by user. >>"
		return []string{}, false
	}

	/* Scan Language */

	inputLanguage := util.ScanOrEditWithLiner("Language", "English")

	if util.MoveBack(inputLanguage) {
		IsMessage = " does something"
		goto author
	}

	if util.Abort(inputLanguage) {
		IsMessage = "<< previous action aborted by user. >>"
		return []string{}, false
	}

	inputs = append(inputs, util.FillEmptyInput(inputQuote, "Unknown"))
	inputs = append(inputs, util.FillEmptyInput(inputAuthor, "Unknown"))
	inputs = append(inputs, inputLanguage)

	return inputs, true
}

func EditUserInput(quotes *db.Quotes, id int) []string {

	index, err := quotes.FindIndexByID(id)
	if err != nil {
		fmt.Println(err)
	}

	var updatedInputs []string

	updatedQuote := util.ScanOrEditWithLiner("Quote", quotes.QUOTES[index].QUOTE)
	updatedAuthor := util.ScanOrEditWithLiner("Author", quotes.QUOTES[index].AUTHOR)
	updatedLanguage := util.ScanOrEditWithLiner("Language", quotes.QUOTES[index].LANGUAGE)

	updatedInputs = append(updatedInputs, util.FillEmptyInput(updatedQuote, "Unknown"))
	updatedInputs = append(updatedInputs, util.FillEmptyInput(updatedAuthor, "Unknown"))
	updatedInputs = append(updatedInputs, updatedLanguage)

	return updatedInputs
}

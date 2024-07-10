package cmd

import (
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func UserInput(quotes *db.Quotes) ([]string, bool) {

quote:

	var inputs []string

	/* Scan Quote */

	inputQuote := util.PromptWithSuggestion("Quote", "")
	if util.Abort(inputQuote) {
		PrintMessage = "<< previous action aborted by user. >>"
		return []string{}, false
	}

	if quotes.Duplicates(inputQuote) != -1 {
		PrintMessage = "<< there are dublicates in database. >>"
		PrintID = quotes.Duplicates(inputQuote)
		return []string{}, false
	}

	/* Scan Author */
author:
	inputAuthor := util.PromptWithSuggestion("Author", "")
	if util.MoveBack(inputAuthor) {
		goto quote
	}
	if util.Abort(inputAuthor) {
		PrintMessage = "<< previous action aborted by user. >>"
		return []string{}, false
	}

	/* Scan Language */

	inputLanguage := util.PromptWithSuggestion("Language", "English")

	if util.MoveBack(inputLanguage) {
		PrintMessage = "<< Moved Back! >>"
		goto author
	}

	if util.Abort(inputLanguage) {
		PrintMessage = "<< previous action aborted by user. >>"
		return []string{}, false
	}

	inputs = append(inputs, util.FillEmptyInput(inputQuote, "Unknown"))
	inputs = append(inputs, util.FillEmptyInput(inputAuthor, "Unknown"))
	inputs = append(inputs, inputLanguage)

	return inputs, true
}

func UpdateUserInput(quotes *db.Quotes, index int) []string {

	var updatedInputs []string

	updatedQuote := util.PromptWithSuggestion("Quote", quotes.QUOTES[index].QUOTE)
	updatedAuthor := util.PromptWithSuggestion("Author", quotes.QUOTES[index].AUTHOR)
	updatedLanguage := util.PromptWithSuggestion("Language", quotes.QUOTES[index].LANGUAGE)

	updatedInputs = append(updatedInputs, util.FillEmptyInput(updatedQuote, "Unknown"))
	updatedInputs = append(updatedInputs, util.FillEmptyInput(updatedAuthor, "Unknown"))
	updatedInputs = append(updatedInputs, updatedLanguage)

	return updatedInputs
}

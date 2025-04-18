package db

import (
	"vk-quotes/pkg/config"
	"vk-quotes/pkg/util"

	"github.com/peterh/liner"
)

type Quote struct {
	ID       int    `json:"id"`
	QUOTE    string `json:"quote"`
	AUTHOR   string `json:"author"`
	LANGUAGE string `json:"language"`
	DATE     string `json:"date"`
}

type Quotes struct {
	QUOTES []Quote `json:"quotes"` // Slice containing multiple Quote instances.
}

func (q *Quotes) UserInput(quote Quote) bool {

	config.UserInputs = []string{} // Clear previous user inputs

	if !q.GetQuote(quote) || !q.GetAuthor(quote) {
		return false
	}
	
	lang := util.AutoDetectLanguage(config.UserInputs[0])

	config.UserInputs = append(config.UserInputs, lang)

	return true
}

func (q *Quotes) GetQuote(quote Quote) bool {

	input := q.PromptWithSuggestion("Quote", quote.QUOTE)

	if !util.Quit(input) {
		return false
	}

	input = util.CapitalizeFirstLetter(input)

	input = util.EnsureSentenceEnd(input)


	if q.FindDuplicates(input, quote.ID) {
		return false
	}

	config.UserInputs = append(config.UserInputs, input)

	return true
}

func (q *Quotes) GetAuthor(quote Quote) bool {

	input := q.PromptWithSuggestion("Author", quote.AUTHOR)

	if !util.Quit(input) {
		return false
	}

	if input == "" {
		input = "Unknown"
	}

	config.UserInputs = append(config.UserInputs, input)

	return true
}

func (q *Quotes) PromptWithSuggestion(name string, edit string) string {

	line := liner.NewLiner()
	defer line.Close()

	input, err := line.PromptWithSuggestion("   "+name+": ", edit, -1)
	if err != nil {
		panic(err)
	}

	return input
}

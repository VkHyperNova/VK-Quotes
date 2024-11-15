package db

import (
	"fmt"
	"vk-quotes/pkg/config"
	"vk-quotes/pkg/util"
)

func (q *Quotes) SetToDefaultQuote() {

	index := len(q.QUOTES) - 1

	if index > 0 {
		config.MainQuoteID = q.QUOTES[index].ID
	}

}

func (q *Quotes) Find() {

	fmt.Print("Search: ")

	var searchQuote string
	fmt.Scanln(&searchQuote)

	q.PrintQuote(searchQuote)

	util.PressAnyKey()
}

func (q *Quotes) FindQuoteByQuote(searchQuote string) Quote {

	var foundQuote Quote

	for _, quote := range q.QUOTES {
		if quote.QUOTE == searchQuote {
			foundQuote = quote
		}
	}
	return foundQuote
}

func (q *Quotes) FindQuoteByID(id int) Quote {

	var foundQuote Quote

	for _, quote := range q.QUOTES {
		if quote.ID == id {
			foundQuote = quote
		}
	}
	return foundQuote
}

func (q *Quotes) FindDuplicates(searchQuote string) bool {

	if searchQuote == "" || searchQuote == "Unknown" {
		return false
	}

	for _, quote := range q.QUOTES {
		if quote.QUOTE == searchQuote {
			if quote.ID != config.MainQuoteID {
				message := config.Red + "There are dublicates in database" + config.Reset
				config.AddMessage(message)
				config.MainQuoteID = quote.ID
				return true
			}
		}
	}
	return false
}

package db

import "vk-quotes/pkg/config"

func (q *Quotes) ResetMainQuote() {

	index := len(q.QUOTES) - 1

	config.MainQuoteID = q.QUOTES[index].ID
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

package db

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
	fmt.Print("Find: ")

	// Read user input
	reader := bufio.NewReader(os.Stdin)
	searchQuote, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		util.PressAnyKey()
		return
	}

	// Clean and process the input
	searchQuote = strings.TrimSpace(searchQuote)
	searchQuote = util.EnsureSentenceEnd(searchQuote)

	// Search for the quote
	foundQuote, found := q.FindQuoteByQuote(searchQuote)
	if !found {
		fmt.Println("Quote not found.")
		util.PressAnyKey()
		return
	}

	// Print the found quote
	fmt.Println(FormatQuote(foundQuote))
	util.PressAnyKey()
}

func (q *Quotes) FindQuoteByQuote(searchQuote string) (Quote, bool) {
	for _, quote := range q.QUOTES {
		if strings.EqualFold(quote.QUOTE, searchQuote) {
			return quote, true
		}
	}
	return Quote{}, false
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

func (q *Quotes) FindDuplicates(searchQuote string, excludeID int) bool {

	if searchQuote == "" || searchQuote == "Unknown" {
		return false
	}

	for _, quote := range q.QUOTES {
		if  searchQuote == quote.QUOTE {
			if excludeID != quote.ID {
				message := config.Red + "There are dublicates in database" + config.Reset
				config.AddMessage(message)
				config.MainQuoteID = quote.ID
				return true
			}
		}
	}

	return false
}

package db

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"vk-quotes/pkg/config"
	"vk-quotes/pkg/util"
)

func (q *Quotes) PrintQuote(command string) {

	for _, quote := range q.QUOTES {

		isID, _ := strconv.Atoi(command)
		if quote.ID == isID {
			fmt.Println(q.FormatQuote(quote))
			return
		}

		normalizedAuthor := strings.ToLower(quote.AUTHOR)
		normalizedQuote := strings.ToLower(quote.QUOTE)

		if strings.Contains(normalizedAuthor, command) || strings.Contains(normalizedQuote, command) {
			fmt.Println(q.FormatQuote(quote))
		}

	}
}

func (q *Quotes) PrintAllQuotes() {

	util.ClearScreen()

	for _, quote := range q.QUOTES {
		fmt.Print(q.FormatQuote(quote))
	}

	util.PressAnyKey()
}

func (q *Quotes) FormatQuote(quote Quote) string {

	var (
		quoteBuffer    bytes.Buffer
		formattedQuote string
	)

	stringFormat := `` + "\n" + config.Cyan + `%d. ` + "\"" + config.Reset + `%s` + `` + config.Cyan + "\"" +
		"\n" + strings.Repeat(" ", 50) + `By %s (%s %s)` + config.Reset + "\n" + ``

	formattedQuote = fmt.Sprintf(
		stringFormat,
		quote.ID,
		quote.QUOTE,
		quote.AUTHOR,
		quote.DATE,
		quote.LANGUAGE)

	quoteBuffer.WriteString(formattedQuote)

	return quoteBuffer.String()
}
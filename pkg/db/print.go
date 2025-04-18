package db

import (
	"bytes"
	"fmt"
	"strings"
	"vk-quotes/pkg/config"
	"vk-quotes/pkg/util"
)

func (q *Quotes) PrintCLI() {

	util.ClearScreen()

	if config.MainQuoteID <= 0 {
		q.SetToDefaultQuote()
	}

	stringFormat := `` +
		config.Cyan + "VK-Quotes" + config.Reset + " %s" + "\n\n" +
		"%s" +
		config.Cyan + `%s` + config.Reset +
		"%s" +
		config.Yellow + `%s` + config.Reset +
		``


	lastQuote := q.FindQuoteByID(config.MainQuoteID)
	formattedLastQuote := FormatQuote(lastQuote)

	messages := config.FormatMessages()

	commands := "\nAdd Update Delete Find Read Showall Stats SimilarQuotes ResetIDs Quit\n"

	cli := fmt.Sprintf(stringFormat, config.ProgramVersion, messages, config.ReadCounter, formattedLastQuote, commands)

	fmt.Print(cli)
}

func (q *Quotes) PrintAllQuotes() {

	util.ClearScreen()

	for _, quote := range q.QUOTES {
		fmt.Print(FormatQuote(quote))
	}

	util.PressAnyKey()
}

func FormatQuote(quote Quote) string {

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

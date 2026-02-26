package db

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"vk-quotes/pkg/color"
	"vk-quotes/pkg/config"
	"vk-quotes/pkg/util"
)

func (q *Quotes) PrintCLI() {

	util.ClearScreen()

	if config.MainQuoteID <= 0 {
		q.SetToDefaultQuote()
	}

	nowPlaying := "Now playing: Flute.mp3"

	stringFormat := `` +
		color.Cyan + "VK-Quotes" + color.Reset + " %s" + "\n" + // Program Name
		color.Purple + "%s" + color.Reset + "\n" + // Now Playing
		"%s" + // Messages
		color.Cyan + `%s` + color.Reset + // Read Counter
		"%s" + // Last Quote
		color.Yellow + `%s` + color.Reset + // Commands
		``

	lastQuote, _ := q.FindQuoteByID(config.MainQuoteID)

	formattedLastQuote := FormatQuote(lastQuote)

	messages := config.FormatMessages()

	commands := "\nRandom Add Update Delete Find Read Showall Stats SimilarQuotes ResetIDs Quit\n"

	cli := fmt.Sprintf(stringFormat, config.ProgramVersion, nowPlaying, messages, config.ReadCounter, formattedLastQuote, commands)

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

	stringFormat := `` + "\n" + color.Cyan + `%d. ` + "\"" + color.Reset + `%s` + `` + color.Cyan + "\"" +
		"\n" + strings.Repeat(" ", 50) + `By %s (%s %s)` + color.Reset + "\n" + ``

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

const DefaultQuoteAmount = 10

func (q *Quotes) PrintRandomQuotes(amount int) {

	util.ClearScreen()

	if amount <= 0 {
		amount = DefaultQuoteAmount
	}

	if len(q.QUOTES) == 0 {
		fmt.Print(FormatQuote(DefaultQuote))
		util.PressAnyKey()
		return
	}

	for i := 0; i < amount; i++ {
		fmt.Printf("Number %d\n", i+1)
		idx := rand.Intn(len(q.QUOTES))
		quote := q.QUOTES[idx] 
		fmt.Print(FormatQuote(quote))
	}

	util.PressAnyKey()
}

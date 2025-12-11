package db

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"vk-quotes/pkg/config"
	"vk-quotes/pkg/util"
)

func (q *Quotes) PrintCLI() {

	util.ClearScreen()

	util.IsVKDataMounted()
	
	if config.MainQuoteID <= 0 {
		q.SetToDefaultQuote()
	}

	nowPlaying := "Now playing: Flute.mp3"

	stringFormat := `` +
		config.Cyan + "VK-Quotes" + config.Reset + " %s" + "\n" + // Program Name
		config.Purple + "%s" + config.Reset + "\n" + // Now Playing
		"%s" + // Messages
		config.Cyan + `%s` + config.Reset + // Read Counter
		"%s" + // Last Quote
		config.Yellow + `%s` + config.Reset + // Commands
		``

	lastQuote, _ := q.FindQuoteByID(config.MainQuoteID)
	formattedLastQuote := FormatQuote(lastQuote)

	messages := config.FormatMessages()

	commands := "\ndayly Add Update Delete Find Read Showall Stats SimilarQuotes ResetIDs Quit\n"

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

func (q *Quotes) PrintDaylyTen() {
	util.ClearScreen()

	maxNumber := len(q.QUOTES)

	numbers := make([]int, 10)
	for i := 0; i < 10; i++ {
		numbers[i] = rand.Intn(maxNumber) + 1 // random number 1–5000
	}

	fmt.Println(numbers)

	for id, quote := range q.QUOTES {
		if util.Contains(numbers, id) {
			fmt.Print(FormatQuote(quote))
		}
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



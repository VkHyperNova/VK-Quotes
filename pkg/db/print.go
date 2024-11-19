package db

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strconv"
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

	quote := q.FindQuoteByID(config.MainQuoteID)
	formattedQuote := q.FormatQuote(quote)

	messages := config.FormatMessages()

	commands := "\nAdd Update Delete Search Read Showall Stats SimilarQuotes ResetIDs Quit\n"

	cli := fmt.Sprintf(stringFormat, config.ProgramVersion, messages, config.ReadCounter, formattedQuote, commands)

	fmt.Print(cli)
}

func (q *Quotes) PrintQuote(command string) {

	reader := bufio.NewReader(os.Stdin)
	
	for _, quote := range q.QUOTES {

		// Find by ID
		isID, _ := strconv.Atoi(command)
		if quote.ID == isID {
			fmt.Println(q.FormatQuote(quote))
			return
		}

		// Find by Quote or Author
		Author := strings.ToLower(quote.AUTHOR)
		Quote := strings.ToLower(quote.QUOTE)

		findAuthor := strings.Contains(Author, command)
		findQuote := strings.Contains(Quote, command)

		
		if findAuthor || findQuote {
			fmt.Println(q.FormatQuote(quote))

			fmt.Print("Press Enter to continue or type 'q' to quit: ")

			input, _ := reader.ReadString('\n') 

			input = strings.TrimSpace(input)    

			if strings.ToLower(input) == "q" { 
				fmt.Println("Exiting...")
				break
			}
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

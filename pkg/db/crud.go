package db

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
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

func (q *Quotes) Add() bool {

	newQuote := Quote{}
	newQuote.ID = q.GenerateUniqueID()

	newQuote.QUOTE = util.PromptWithSuggestion("Quote", "")
	if !util.Quit(newQuote.QUOTE) {
		return false
	}

	newQuote.QUOTE = util.CapitalizeFirstLetter(newQuote.QUOTE)
	newQuote.QUOTE = util.EnsureSentenceEnd(newQuote.QUOTE)
	if q.FindDuplicates(newQuote.QUOTE, newQuote.ID) {
		return false
	}

	newQuote.AUTHOR = util.PromptWithSuggestion("Author", "Unknown")
	newQuote.LANGUAGE = util.AutoDetectLanguage(newQuote.QUOTE)
	newQuote.DATE = time.Now().Format("02.01.2006")

	q.QUOTES = append(q.QUOTES, newQuote)

	message := config.Green + strconv.Itoa(newQuote.ID) + ". " + newQuote.QUOTE + "\n\t" + newQuote.AUTHOR + " (" + newQuote.LANGUAGE + " " + newQuote.DATE+ ")" + config.Reset

	q.SaveToFile(message)

	config.MainQuoteID = newQuote.ID

	return true
}

func (q *Quotes) Update(updateID int) bool {

	updateQuote, exists := q.FindQuoteByID(updateID)

	if !exists {
		config.AddMessage(config.Red + "Quote ID out of range! Range from 1 to " + strconv.Itoa(len(q.QUOTES)) + config.Reset)
		return false
	}

	updateQuote.QUOTE = util.PromptWithSuggestion("Quote", updateQuote.QUOTE)
	updateQuote.AUTHOR = util.PromptWithSuggestion("Author", updateQuote.AUTHOR)

	for i, quote := range q.QUOTES {
		if quote.ID == updateID {
			q.QUOTES[i] = updateQuote

		}
	}

	message := config.Yellow + strconv.Itoa(updateID) + " updated" + config.Reset

	q.SaveToFile(message)

	config.MainQuoteID = updateID

	return true
}

func (q *Quotes) Delete(deleteID int) bool {

	index := -1
	for i, quote := range q.QUOTES {
		if quote.ID == deleteID {
			index = i
		}
	}

	if index == -1 {
		message := config.Red + "Index Not Found" + config.Reset
		config.AddMessage(message)
		return false
	}

	for _, quote := range q.QUOTES {
		if quote.ID == deleteID {
			fmt.Println(FormatQuote(quote))
		}
	}

	line := liner.NewLiner()

	defer line.Close()

	confirm, _ := line.Prompt("(y/n)")

	if confirm != "y" {
		message := config.Red + "Delete Canceled" + config.Reset
		config.AddMessage(message)
		return false
	}

	q.QUOTES = append(q.QUOTES[:index], q.QUOTES[index+1:]...)

	message := config.Red + strconv.Itoa(deleteID) + " deleted" + config.Reset

	q.SaveToFile(message)

	q.SetToDefaultQuote()

	return true
}

func (q *Quotes) Read() {

	// Append All Quotes IDs
	for _, quote := range q.QUOTES {
		if !util.ArrayContainsInt(config.RandomIDs, quote.ID) {
			config.RandomIDs = append(config.RandomIDs, quote.ID)
		}
	}

	for len(config.RandomIDs) != 0 {

		config.Counter += 1

		source := rand.NewSource(time.Now().UnixNano())
		r := rand.New(source)
		randomIndex := r.Intn(len(config.RandomIDs))

		config.MainQuoteID = config.RandomIDs[randomIndex]
		config.DeleteUsedID(randomIndex)

		count := config.Counter
		size := len(q.QUOTES)
		percentage := float64(count) / float64(size) * 100
		config.ReadCounter = fmt.Sprintf("<< Reading [%d] %.0f%% >>", count, percentage)

		q.PrintCLI()

		var quit string
		fmt.Scanln(&quit)

		if quit == "q" {
			break
		}
	}

	message := config.Yellow + "Reading Done" + config.Reset
	config.AddMessage(message)

	config.RandomIDs = config.RandomIDs[:0]
	q.SetToDefaultQuote()
	config.Counter = 0
	config.ReadCounter = ""
}

func (q *Quotes) GenerateUniqueID() int {

	maxID := 0

	for _, quote := range q.QUOTES {
		if quote.ID > maxID {
			maxID = quote.ID
		}
	}

	return maxID + 1
}

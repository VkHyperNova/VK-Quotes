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

func (q *Quotes) Add() bool {

	maxID := 0

	for _, quote := range q.QUOTES {
		if quote.ID > maxID {
			maxID = quote.ID
		}
	}

	newID := maxID + 1

	quote := Quote{
		ID:       newID,
		QUOTE:    config.UserInputs[0],
		AUTHOR:   config.UserInputs[1],
		LANGUAGE: config.UserInputs[2],
		DATE:     time.Now().Format("02.01.2006")}

	q.QUOTES = append(q.QUOTES, quote)

	message := config.Green + strconv.Itoa(newID) + " added" + config.Reset

	q.SaveToFile(message)

	config.MainQuoteID = newID

	return true
}

func (q *Quotes) Update(updateID int) bool {

	updatedQuote := Quote{
		ID:       updateID,
		QUOTE:    config.UserInputs[0],
		AUTHOR:   config.UserInputs[1],
		LANGUAGE: config.UserInputs[2],
		DATE:     time.Now().Format("02.01.2006"),
	}

	for i, quote := range q.QUOTES {
		if quote.ID == updateID {
			q.QUOTES[i] = updatedQuote

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

	q.PrintQuote(strconv.Itoa(deleteID))

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

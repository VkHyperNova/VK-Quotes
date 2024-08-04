package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"vk-quotes/pkg/config"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"

	"math/rand"

	"github.com/peterh/liner"
)

func CommandLine(quotes *db.Quotes) {

	cli(quotes)

	input, inputID := util.CommandPrompt("> ")

	for {
		switch input {

		case "add", "a":
			if quotes.UserInput(inputID) {
				add(quotes)
			}
			CommandLine(quotes)

		case "update", "u":
			if quotes.UserInput(inputID) {
				update(quotes, inputID)
			}
			CommandLine(quotes)

		case "delete", "d":
			delete(quotes, inputID)
			CommandLine(quotes)

		case "showall", "s":
			quotes.PrintAllQuotes()
			CommandLine(quotes)

		case "stats":
			printStatistics(quotes)
			CommandLine(quotes)

		case "resetids":
			quotes.ResetIDs(quotes)
			CommandLine(quotes)

		case "read", "r":
			read(quotes)
			CommandLine(quotes)

		case "similarquotes", "sim":
			db.FindSimilarQuotes(quotes)
			CommandLine(quotes)

		case "q", "quit":
			util.ClearScreen()
			os.Exit(0)

		default:
			if input != "" {
				quotes.PrintQuote(input)
				util.PressAnyKey()
			}
			CommandLine(quotes)
		}
	}
}

func cli(quotes *db.Quotes) {

	util.ClearScreen() 

	if config.MainQuoteID <= 0 { 
		config.MainQuoteID = quotes.LastID() 
	}

	stringFormat := `` + 
		config.Cyan + "VK-Quotes" + config.Reset + " %s" + "\n\n" +
		"%s" +
		config.Cyan + `%s` + config.Reset +
		"%s" +
		config.Yellow + `%s` + config.Reset +
		``

	quote := quotes.FindByID(config.MainQuoteID)
	formattedQuote := quotes.FormatQuote(quote)

	commands := "\nAdd Update Delete Read Showall Stats SimilarQuotes ResetIDs Quit\n" 

	cli := fmt.Sprintf(stringFormat, config.ProgramVersion, util.FormatMessages(), config.ReadCounter, formattedQuote, commands) 

	fmt.Print(cli) 
}

func add(quotes *db.Quotes) bool {

	newID := quotes.CreateId()

	quotes.AppendQuote(db.Quote{
		ID:       newID,
		QUOTE:    config.UserInputs[0],
		AUTHOR:   config.UserInputs[1],
		LANGUAGE: config.UserInputs[2],
		DATE:     time.Now().Format("02.01.2006")})

	quotes.SaveToFile()

	config.MainQuoteID = newID

	return true
}

func update(quotes *db.Quotes, updateID int) bool {

	quotes.Update(db.Quote{
		ID:       updateID,
		QUOTE:    config.UserInputs[0],
		AUTHOR:   config.UserInputs[1],
		LANGUAGE: config.UserInputs[2],
		DATE:     time.Now().Format("02.01.2006")})

	quotes.SaveToFile()

	config.MainQuoteID = updateID

	return true
}

func delete(quotes *db.Quotes, deleteID int) bool {

	index := quotes.IndexOf(deleteID)

	if index == -1 {
		config.Messages = append(config.Messages, config.Red+"<< Index Not Found >>"+config.Reset)
		return false
	}

	quotes.PrintQuote(strconv.Itoa(deleteID))

	line := liner.NewLiner()

	defer line.Close()

	confirm, _ := line.Prompt("(y/n)")

	if confirm != "y" {
		config.Messages = append(config.Messages, config.Red+"<< Delete Canceled >>"+config.Reset)
		return false
	}

	quotes.Remove(index)

	quotes.SaveToFile()

	config.MainQuoteID = quotes.LastID()

	return true
}

func read(quotes *db.Quotes) {

	quotes.AppendRandomIDs()

	for len(config.RandomIDs) != 0 {

		config.Counter += 1

		source := rand.NewSource(time.Now().UnixNano())
		r := rand.New(source)
		randomIndex := r.Intn(len(config.RandomIDs))

		config.MainQuoteID = config.RandomIDs[randomIndex]
		config.RandomIDs = append(config.RandomIDs[:randomIndex], config.RandomIDs[randomIndex+1:]...)

		count := config.Counter
		size := quotes.Size()
		percentage := float64(count) / float64(size) * 100
		config.ReadCounter = fmt.Sprintf("<< Reading [%d] %.0f%% >>", count, percentage)

		cli(quotes)

		var quit string
		fmt.Scanln(&quit)

		if quit == "q" {
			break
		}
	}

	config.RandomIDs = config.RandomIDs[:0]
	config.Messages = append(config.Messages, config.Yellow+"<< Reading Done! >>"+config.Reset)
	config.MainQuoteID = quotes.LastID()
	config.Counter = 0
	config.ReadCounter = ""
	CommandLine(quotes)
}

func printStatistics(quotes *db.Quotes) {

	util.ClearScreen()

	format := "%s%s%s"

	name := config.Cyan + "Statistics: " + config.Reset

	stats := fmt.Sprintf(format, name, quotes.TopAuthors(), quotes.TopLanguages())

	fmt.Println(stats)

	util.PressAnyKey()
}

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

	messages := config.FormatMessages()

	commands := "\nAdd Update Delete Read Showall Stats SimilarQuotes ResetIDs Quit\n"

	cli := fmt.Sprintf(stringFormat, config.ProgramVersion, messages, config.ReadCounter, formattedQuote, commands)

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

	message := config.Green + strconv.Itoa(newID) + " added" + config.Reset

	quotes.SaveToFile(message)

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

	message := config.Yellow + strconv.Itoa(updateID) + " updated" + config.Reset

	quotes.SaveToFile(message)

	config.MainQuoteID = updateID

	return true
}

func delete(quotes *db.Quotes, deleteID int) bool {

	index := quotes.IndexOf(deleteID)
	

	if index == -1 {
		message := config.Red + "Index Not Found" + config.Reset
		config.AddMessage(message)
		return false
	}

	quotes.PrintQuote(strconv.Itoa(deleteID))

	line := liner.NewLiner()

	defer line.Close()

	confirm, _ := line.Prompt("(y/n)")

	if confirm != "y" {
		message := config.Red + "Delete Canceled" + config.Reset
		config.AddMessage(message)
		return false
	}

	quotes.Remove(index)

	message := config.Red + strconv.Itoa(deleteID) + " deleted" + config.Reset

	quotes.SaveToFile(message)

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
		config.DeleteUsedID(randomIndex)

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

	message := config.Yellow + "Reading Done" + config.Reset
	config.AddMessage(message)

	config.DeleteAllRandomIDs()
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

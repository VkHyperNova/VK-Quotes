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

func CMD(quotes *db.Quotes) {

	CLI(quotes)

	command, commandID := util.CommandPrompt("> ")

	for {
		switch command {
		case "add", "a":
			validation := quotes.UserInput(commandID)
			if validation {
				Add(quotes)
			}
			CMD(quotes)
		case "update", "u":
			validation := quotes.UserInput(commandID)
			if validation {
				Update(quotes, commandID)
			}
			CMD(quotes)
		case "delete", "d":
			Delete(quotes, commandID)
			CMD(quotes)
		case "showall", "s":
			quotes.PrintAllQuotes()
			util.PressAnyKey()
			CMD(quotes)
		case "stats":
			printStats(quotes)
			util.PressAnyKey()
			CMD(quotes)
		case "resetids":
			quotes.ResetIDs(quotes)
			util.PressAnyKey()
			CMD(quotes)
		case "read", "r":
			Read(quotes)
			CMD(quotes)
		case "similar":
			db.RunTaskWithProgress(quotes)
			util.PressAnyKey()
			CMD(quotes)
		case "q":
			util.ClearScreen()
			os.Exit(0)
		default:
			foundQuote := quotes.Search(command)
			formattedQuote := quotes.FormatQuote(foundQuote)
			fmt.Print(formattedQuote)
			CMD(quotes)
		}
	}

}

func CLI(quotes *db.Quotes) {

	util.ClearScreen()

	// Reset the ID if it's set to 0 or -1.
	if config.ID <= 0 {
		config.ID = quotes.LastID()
	}

	format := "%s %s %s %s"
	version := util.Cyan + "VK-Quotes" + " " + config.ProgramVersion + util.Reset
	message := util.Yellow + "\n\n" + config.Message + "\n" + util.Reset
	foundQuote := quotes.Search(strconv.Itoa(config.ID))
	formattedQuote := quotes.FormatQuote(foundQuote)
	commands := util.Yellow + "\n" + "add, update, delete, read, showall, stats, similar, resetids, quit" + "\n" + util.Reset

	cli := fmt.Sprintf(format, version, message, formattedQuote, commands)

	fmt.Print(cli)
}

func Add(quotes *db.Quotes) bool {

	newID := quotes.CreateId()

	quotes.AppendQuote(db.Quote{
		ID:       newID,
		QUOTE:    config.UserInputs[0],
		AUTHOR:   config.UserInputs[1],
		LANGUAGE: config.UserInputs[2],
		DATE:     time.Now().Format("02.01.2006")})

	quotes.SaveToFile(config.SaveQuotesPath)
	config.ID = newID
	config.Message = fmt.Sprintf("<< %d Quote Added! >>", newID)

	return true
}

func Update(quotes *db.Quotes, updateID int) bool {

	quotes.Update(db.Quote{
		ID:       updateID,
		QUOTE:    config.UserInputs[0],
		AUTHOR:   config.UserInputs[1],
		LANGUAGE: config.UserInputs[2],
		DATE:     time.Now().Format("02.01.2006")})

	quotes.SaveToFile(config.SaveQuotesPath)

	config.ID = updateID

	config.Message = fmt.Sprintf("<< %d Quote Updated! >>", updateID)

	return true
}

func Delete(quotes *db.Quotes, deleteID int) bool {

	index := quotes.IndexOf(deleteID)

	if index == -1 {
		config.Message = "Index Not Found!"
		return false
	}

	quote := quotes.Search(strconv.Itoa(deleteID))

	formattedQuote := quotes.FormatQuote(quote)

	fmt.Println(formattedQuote)

	line := liner.NewLiner()

	defer line.Close()

	confirm, err := line.Prompt("(y/n)")

	if err != nil {
		config.Message = err.Error()
		return false
	}

	if confirm != "y" {
		return false
	}

	quotes.Remove(index)

	quotes.SaveToFile(config.SaveQuotesPath)

	config.ID = quotes.LastID()

	config.Message = fmt.Sprintf("<< %d Quote Deleted! >>", deleteID)

	return true
}

func Read(quotes *db.Quotes) {

	quotes.AppendRandomIDs()

	for len(config.RandomIDs) != 0 {

		config.ReadCounter += 1

		source := rand.NewSource(time.Now().UnixNano())

		r := rand.New(source)

		randomIndex := r.Intn(len(config.RandomIDs))

		config.ID = config.RandomIDs[randomIndex]

		config.RandomIDs = append(config.RandomIDs[:randomIndex], config.RandomIDs[randomIndex+1:]...)

		count := config.ReadCounter
		size := quotes.Size()
		percentage := float64(count) / float64(size) * 100
		config.Message = fmt.Sprintf("<< Reading [%d] %.0f%% >>", count, percentage)

		CLI(quotes)

		var quit string
		fmt.Scanln(&quit)

		if quit == "q" {
			config.Message = "<< Reading Mode Off >>"

			config.ReadCounter = 0

			config.ID = quotes.LastID()

			if len(config.RandomIDs) > 0 {
				config.RandomIDs = config.RandomIDs[:0]
			}

			CMD(quotes)
		}
	}

	config.Message = "<< Read 100% >>"

	config.ReadCounter = 0
}

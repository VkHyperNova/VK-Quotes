package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"

	"github.com/peterh/liner"
	"math/rand"
)

func CMD(quotes *db.Quotes, settings *util.Settings) {

	CLI(quotes, settings)

	command, commandID := util.CommandPrompt(settings, "> ")

	for {
		switch command {
		case "add", "a":
			validation := quotes.UserInput(settings, commandID)
			if validation {
				Add(quotes, settings)
			}
			CMD(quotes, settings)
		case "update", "u":
			validation := quotes.UserInput(settings, commandID)
			if validation {
				Update(quotes, settings, commandID)
			}
			CMD(quotes, settings)
		case "delete", "d":
			Delete(quotes, settings, commandID)
			CMD(quotes, settings)
		case "showall", "s":
			quotes.PrintAllQuotes()
			util.PressAnyKey()
			CMD(quotes, settings)
		case "stats":
			printStats(quotes)
			util.PressAnyKey()
			CMD(quotes, settings)
		case "rearrange":
			quotes.ReArrangeIDs(settings)
			util.PressAnyKey()
			CMD(quotes, settings)
		case "read", "r":
			Read(quotes, settings)
			CMD(quotes, settings)
		case "similar":
			db.RunTaskWithProgress(quotes, settings)
			util.PressAnyKey()
			CMD(quotes, settings)
		case "q":
			util.ClearScreen()
			os.Exit(0)
		default:
			foundQuote := quotes.Search(command)
			formattedQuote := quotes.FormatQuote(foundQuote)
			fmt.Print(formattedQuote)
			CMD(quotes, settings)
		}
	}

}

func CLI(quotes *db.Quotes, settings *util.Settings) {

	util.ClearScreen()

	// Reset the ID if it's set to 0 or -1.
	if settings.ID == 0 || settings.ID == -1 {
		quotes.ResetID(settings)
	}

	format := "%s %s %s %s %s"
	version := util.Cyan + "VK-Quotes" + " " + settings.Version + util.Reset
	message := util.Yellow + "\n\n" + settings.Message + "\n" + util.Reset
	counter := ""
	foundQuote := quotes.Search(strconv.Itoa(settings.ID))
	formattedQuote := quotes.FormatQuote(foundQuote)
	commands := util.Yellow + "\n" + "add, update, delete, read, showall, stats, similar, reaarange, quit" + "\n" + util.Reset

	cli := fmt.Sprintf(format, version, message, counter, formattedQuote, commands)

	fmt.Print(cli)
}

func Add(quotes *db.Quotes, settings *util.Settings) bool {
	newID := quotes.CreateId()
	quotes.AppendQuote(db.Quote{ID: newID, QUOTE: settings.UserInputs[0], AUTHOR: settings.UserInputs[1], LANGUAGE: settings.UserInputs[2], DATE: time.Now().Format("02.01.2006")})
	quotes.SaveToFile(settings)
	settings.ID = newID
	settings.Message = fmt.Sprintf("<< %d Quote Added! >>", newID)

	return true
}

func Update(quotes *db.Quotes, settings *util.Settings, updateID int) bool {

	quotes.Update(db.Quote{
		ID:       updateID,
		QUOTE:    settings.UserInputs[0],
		AUTHOR:   settings.UserInputs[1],
		LANGUAGE: settings.UserInputs[2],
		DATE:     time.Now().Format("02.01.2006")})

	quotes.SaveToFile(settings)

	settings.ID = updateID

	settings.Message = fmt.Sprintf("<< %d Quote Updated! >>", updateID)

	return true
}

func Delete(quotes *db.Quotes, settings *util.Settings, deleteID int) bool {
	/*
		Delete removes a quote from the database if it exists, with user confirmation.
		It returns true if the quote was successfully deleted, otherwise false.
	*/

	// Check if the quote ID exists in the database

	index := quotes.IndexOf(deleteID)

	if index == -1 {
		settings.Message = "Index Not Found!"
		return false
	}

	// Convert ID to String and Search for the quote

	quote := quotes.Search(strconv.Itoa(deleteID))

	// Print the quote to the console for user verification

	fmt.Println(quotes.FormatQuote(quote))

	// Initialize a new line reader for user input

	line := liner.NewLiner()

	// Ensure the line reader is closed to free resources

	defer line.Close()

	// Prompt the user for confirmation to delete the quote

	confirm, err := line.Prompt("(y/n)")

	// Handle any errors from the prompt

	if err != nil {
		settings.Message = err.Error()
		return false
	}

	// If the user did not confirm with 'y', abort the deletion

	if confirm != "y" {
		return false
	}

	// Remove the quote from the database

	quotes.Remove(index)

	// Save the updated quotes database to file

	quotes.SaveToFile(settings)

	// Reset IDs after deletion

	quotes.ResetID(settings)

	// Set a success message in settings

	settings.Message = fmt.Sprintf("<< %d Quote Deleted! >>", deleteID)

	// Indicate successful deletion

	return true
}

/* Random IDs bug */
func Read(quotes *db.Quotes, settings *util.Settings) {

	quotes.AppendRandomIDs(settings)

	for len(settings.RandomIDs) != 0 {

		source := rand.NewSource(time.Now().UnixNano())

		r := rand.New(source)

		randomIndex := r.Intn(len(settings.RandomIDs))

	
		// asd = append(asd, settings.RandomIDs[randomIndex]) //

		settings.ID = settings.RandomIDs[randomIndex]

		settings.RandomIDs = append(settings.RandomIDs[:randomIndex], settings.RandomIDs[randomIndex+1:]...)


		count := settings.ReadCounter
		size := quotes.Size()
		percentage := float64(count) / float64(size) * 100
		settings.Message = fmt.Sprintf("Reading:[%d] [%d] %.0f%%", count, percentage)

		CLI(quotes, settings)

		settings.ReadCounter += 1

		var quit string
		fmt.Scanln(&quit)

		if quit == "q" {
			settings.Message = "<< Reading Mode Off >>"

			settings.ReadCounter = 0

			quotes.ResetID(settings)

			if len(settings.RandomIDs) > 0 {
				settings.RandomIDs = settings.RandomIDs[:0]
			}

			CMD(quotes, settings)
		}
	}

	settings.Message = "<< Read 100% >>"

	settings.ReadCounter = 0
}

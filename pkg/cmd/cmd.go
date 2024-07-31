package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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

	command = strings.ToLower(command)

	for {
		switch  command {
		case "add", "a":
			validation := quotes.UserInput(commandID)
			if validation {
				add(quotes)
			}
			CMD(quotes)
		case "update", "u":
			validation := quotes.UserInput(commandID)
			if validation {
				update(quotes, commandID)
			}
			CMD(quotes)
		case "delete", "d":
			delete(quotes, commandID)
			CMD(quotes)
		case "showall", "s":
			quotes.PrintAllQuotes()
			util.PressAnyKey()
			CMD(quotes)
		case "stats":
			printStatistics(quotes)
			util.PressAnyKey()
			CMD(quotes)
		case "resetids":
			quotes.ResetIDs(quotes)
			CMD(quotes)
		case "read", "r":
			read(quotes)
			CMD(quotes)
		case "similar", "similarquotes":
			db.RunTaskWithProgress(quotes)
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

	format := "%s %s %s %s %s"
	version := util.Cyan + "VK-Quotes" + " " + config.ProgramVersion + util.Reset + "\n\n"
	message := config.FormatMessages()
	readCounter := util.Cyan + config.ReadCounter + util.Reset
	foundQuote := quotes.Search(strconv.Itoa(config.ID))
	formattedQuote := quotes.FormatQuote(foundQuote)
	commands := util.Yellow + "\n" + "Add Update Delete Read Showall Stats SimilarQuotes ResetIDs Quit" + "\n" + util.Reset

	cli := fmt.Sprintf(format, version, message, readCounter, formattedQuote, commands)

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

	quotes.SaveToFile(config.SaveQuotesPath)
	config.ID = newID
	config.Messages = append(config.Messages, fmt.Sprintf("<< %d Quote Added! >>", newID))

	return true
}

func update(quotes *db.Quotes, updateID int) bool {

	quotes.Update(db.Quote{
		ID:       updateID,
		QUOTE:    config.UserInputs[0],
		AUTHOR:   config.UserInputs[1],
		LANGUAGE: config.UserInputs[2],
		DATE:     time.Now().Format("02.01.2006")})

	quotes.SaveToFile(config.SaveQuotesPath)

	config.ID = updateID

	config.Messages = append(config.Messages, fmt.Sprintf("<< %d Quote Updated! >>", updateID))

	return true
}

func delete(quotes *db.Quotes, deleteID int) bool {

	index := quotes.IndexOf(deleteID)

	if index == -1 {
		config.Messages = append(config.Messages, "Index Not Found!")
		return false
	}

	quote := quotes.Search(strconv.Itoa(deleteID))

	formattedQuote := quotes.FormatQuote(quote)

	fmt.Println(formattedQuote)

	line := liner.NewLiner()

	defer line.Close()

	confirm, err := line.Prompt("(y/n)")

	if err != nil {
		config.Messages = append(config.Messages, err.Error())
		return false
	}

	if confirm != "y" {
		config.Messages = append(config.Messages, "<< Delete Canceled! >>")
		return false
	}

	quotes.Remove(index)

	quotes.SaveToFile(config.SaveQuotesPath)

	config.ID = quotes.LastID()

	config.Messages = append(config.Messages, fmt.Sprintf("<< %d Quote Deleted! >>", deleteID))

	return true
}

func read(quotes *db.Quotes) {

    // Append a set of random IDs to the quotes object
	quotes.AppendRandomIDs()

	// Loop until all random IDs are processed
	for len(config.RandomIDs) != 0 {

		// Increment the counter for each iteration
		config.Counter += 1

		// Create a new random source based on the current time in nanoseconds
		source := rand.NewSource(time.Now().UnixNano())

		// Initialize a new random number generator with the source
		r := rand.New(source)

		// Generate a random index based on the length of RandomIDs
		randomIndex := r.Intn(len(config.RandomIDs))

		// Select a random ID from RandomIDs using the random index
		config.ID = config.RandomIDs[randomIndex]

		// Remove the selected ID from RandomIDs
		config.RandomIDs = append(config.RandomIDs[:randomIndex], config.RandomIDs[randomIndex+1:]...)

		// Get the current count of processed items
		count := config.Counter

		// Get the total size of the quotes
		size := quotes.Size()
		
		// Calculate the percentage of completed readings
		percentage := float64(count) / float64(size) * 100

		// Update the read counter message
		config.ReadCounter = fmt.Sprintf("<< Reading [%d] %.0f%% >>", count, percentage)

		// Call the CLI function to process/display the quotes
		CLI(quotes)

		// Wait for user input to continue or quit
		var quit string
		fmt.Scanln(&quit)

		// If the user inputs "q", exit the reading mode
		if quit == "q" {
			// Append a message indicating the reading mode is off
			config.Messages = append(config.Messages, "<< Reading Mode Off >>")

			// Reset the read counter
			config.ResetReadCounter()

			// Set the current ID to the last ID in quotes
			config.ID = quotes.LastID() 

			// Clear the remaining RandomIDs if there are any left
			if len(config.RandomIDs) > 0 {
				config.RandomIDs = config.RandomIDs[:0]
			}

			// Call the CMD function to process the command
			CMD(quotes)
		}
	}

	// Append a message indicating 100% of readings are done
	config.Messages = append(config.Messages, "<< Read 100% >>")

	// Reset the read counter after finishing
	config.ResetReadCounter()
}

func printStatistics(quotes *db.Quotes) {

	util.ClearScreen()

	format := "%s%s%s"

	name := util.Cyan + "Statistics: " + util.Reset

	stats := fmt.Sprintf(format, name, quotes.TopAuthors(), quotes.TopLanguages())
	
	fmt.Println(stats)
}

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

	util.ClearScreen() // Clear the screen for a fresh display.

	if config.MainQuoteID <= 0 { // Check if the main quote ID is not set (<= 0)
		config.MainQuoteID = quotes.LastID() // If so, set it to the ID of the last quote.
	}

	stringFormat := `` + // Define the string format for displaying the CLI information.
		config.Cyan + "VK-Quotes" + config.Reset + " %s" + "\n\n" +
		"%s" +
		config.Cyan + `%s` + config.Reset +
		"%s" +
		config.Yellow + `%s` + config.Reset +
		``

	quote := quotes.FindByID(config.MainQuoteID)
	formattedQuote := quotes.FormatQuote(quote)

	commands := "\nAdd Update Delete Read Showall Stats SimilarQuotes ResetIDs Quit\n" // Define the commands available to the user with color formatting.

	cli := fmt.Sprintf(stringFormat, config.ProgramVersion, config.FormatMessages(), config.ReadCounter, formattedQuote, commands) // Combine all the parts into the final CLI string.

	fmt.Print(cli) // Print the assembled CLI string.
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
		config.Messages = append(config.Messages, config.Red + "<< Index Not Found >>" + config.Reset)
		return false
	}

	quotes.PrintQuote(strconv.Itoa(deleteID))

	line := liner.NewLiner()

	defer line.Close()

	confirm, _ := line.Prompt("(y/n)")

	if confirm != "y" {
		config.Messages = append(config.Messages, config.Red + "<< Delete Canceled >>" + config.Reset)
		return false
	}

	quotes.Remove(index)

	quotes.SaveToFile()

	config.MainQuoteID = quotes.LastID()

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
		config.MainQuoteID = config.RandomIDs[randomIndex]

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
		cli(quotes)

		// Wait for user input to continue or quit
		var quit string
		fmt.Scanln(&quit)

		// If the user inputs "q", exit the reading mode
		if quit == "q" {
			// Append a message indicating the reading mode is off
			config.Messages = append(config.Messages, config.Red + "<< Reading Mode Off >>" + config.Reset)

			// Reset the read counter
			config.ResetReadCounter()

			// Set the current ID to the last ID in quotes
			config.MainQuoteID = quotes.LastID()

			// Clear the remaining RandomIDs if there are any left
			if len(config.RandomIDs) > 0 {
				config.RandomIDs = config.RandomIDs[:0]
			}

			// Call the CMD function to process the command
			CommandLine(quotes)
		}
	}

	// Append a message indicating 100% of readings are done
	config.Messages = append(config.Messages, config.Green + "<< Read 100% >>" + config.Reset)

	// Reset the read counter after finishing
	config.ResetReadCounter()
}

func printStatistics(quotes *db.Quotes) {

	util.ClearScreen()

	format := "%s%s%s"

	name := config.Cyan + "Statistics: " + config.Reset

	stats := fmt.Sprintf(format, name, quotes.TopAuthors(), quotes.TopLanguages())

	fmt.Println(stats)

	util.PressAnyKey()
}

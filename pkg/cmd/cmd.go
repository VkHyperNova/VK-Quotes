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

// CommandLine is the main command processing function that handles various commands
// entered by the user to manipulate quotes in the database.
func CommandLine(quotes *db.Quotes) {

	cli(quotes) // Initial display of the command line interface.

	command, commandID := util.CommandPrompt("> ") // Retrieve the user's command input and corresponding ID if applicable.

	command = strings.ToLower(command) // Convert command to lowercase for uniformity.

	// Enter an infinite loop to continuously process user commands.

	for {
		switch command {

		case "add", "a": // Add a new quote to the database.
			validation := quotes.UserInput(commandID) // Validate the user input for the add command.
			if validation {
				add(quotes) // Call the add function if validation passes.
			}
			CommandLine(quotes) // Restart command processing.

		case "update", "u": // Update a Quote in the database.
			validation := quotes.UserInput(commandID) // Validate the user input for the update command.
			if validation {
				update(quotes, commandID) // Call the update function if validation passes.
			}
			CommandLine(quotes) // Restart command processing.

		case "delete", "d": // Delete a Quote from the database.
			delete(quotes, commandID) // Call the delete function.
			CommandLine(quotes)       // Restart command processing.

		case "showall", "s": // Print all Quotes from the database.
			util.ClearScreen()      // Clear the screen in terminal.
			quotes.PrintAllQuotes() // Display all quotes.
			util.PressAnyKey()      // Wait for the user to press a key.
			CommandLine(quotes)     // Restart command processing.

		case "stats": // Print Statistics for top Authors and top Languages.
			util.ClearScreen()      // Clear the screen in terminal.
			printStatistics(quotes) // Print statistics about the quotes.
			util.PressAnyKey()      // Wait for the user to press a key.
			CommandLine(quotes)     // Restart command processing.

		case "resetids": // Reset all IDs in the database from 1 to n.
			quotes.ResetIDs(quotes) // Reset the IDs of the quotes.
			CommandLine(quotes)     // Restart command processing.

		case "read", "r": // Read all Quotes in random order.
			read(quotes)        // Call the read function.
			CommandLine(quotes) // Restart command processing.

		case "similarquotes", "sim": // Find similar Quotes in the database and save them to file.
			db.RunTaskWithProgress(quotes) // Run a task to find similar quotes.
			CommandLine(quotes)            // Restart command processing.

		case "q", "quit": // Exit the program.
			util.ClearScreen()
			os.Exit(0) 

		default:
			foundQuote := quotes.Search(command)             // Search for a quote based on the command input.
			formattedQuote := quotes.FormatQuote(foundQuote) // Format the found quote.
			fmt.Print(formattedQuote)                        // Print the formatted quote.
			CommandLine(quotes)                              // Restart command processing.
		}
	}
}

// cli initializes and displays the command-line interface (CLI) for the quotes application.
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

	foundQuote := quotes.Search(strconv.Itoa(config.MainQuoteID)) // Search for the quote with the current main quote ID.

	formattedQuote := quotes.FormatQuote(foundQuote) // Format the found Quote.

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

	quotes.SaveToFile()

	config.MainQuoteID = updateID

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

	quotes.SaveToFile()

	config.MainQuoteID = quotes.LastID()

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
			config.Messages = append(config.Messages, "<< Reading Mode Off >>")

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
	config.Messages = append(config.Messages, "<< Read 100% >>")

	// Reset the read counter after finishing
	config.ResetReadCounter()
}

func printStatistics(quotes *db.Quotes) {

	format := "%s%s%s"

	name := config.Cyan + "Statistics: " + config.Reset

	stats := fmt.Sprintf(format, name, quotes.TopAuthors(), quotes.TopLanguages())

	fmt.Println(stats)
}

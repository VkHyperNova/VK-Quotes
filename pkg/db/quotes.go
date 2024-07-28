package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"vk-quotes/pkg/util"

	"github.com/peterh/liner"
)

type Quote struct {
	ID       int    `json:"id"`       // Unique identifier for the quote.
	QUOTE    string `json:"quote"`    // The text of the quote.
	AUTHOR   string `json:"author"`   // The author of the quote.
	LANGUAGE string `json:"language"` // The language in which the quote is written.
	DATE     string `json:"date"`     // The date when the quote was made or published.
}

type Quotes struct {
	QUOTES []Quote `json:"quotes"` // Slice containing multiple Quote instances.
}

func (q *Quotes) Add(quote Quote) {
	// Append the provided quote to the QUOTES slice.
	q.QUOTES = append(q.QUOTES, quote)
}

func (q *Quotes) ReadFromFile(settings *util.Settings) error {
	// Define the path to the JSON file where quotes are stored.
	path := settings.SaveQuotesPath
	folder := settings.SaveFolderPath

	// Check if the file exists at the specified path.
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// If the file does not exist, create the "database" directory.
		_ = os.Mkdir(folder, 0700)

		// Create the JSON file with an initial empty quotes array.
		err = os.WriteFile(path, []byte(`{"quotes": []}`), 0644)
		if err != nil {
			// Print any error that occurs during file creation.
			fmt.Println(err)
		}
		// Print a message indicating that a new database file has been created.
		util.PrintRed("New Database Created!\n")
	}

	// Open the file for reading.
	file, err := os.Open(path)
	if err != nil {
		// Return an error if the file could not be opened.
		return err
	}
	// Ensure the file is closed after reading.
	defer file.Close()

	// Read the contents of the file into a byte slice.
	byteValue, err := io.ReadAll(file)
	if err != nil {
		// Return an error if reading the file fails.
		return err
	}

	// Unmarshal the JSON byte slice into the Quotes struct.
	err = json.Unmarshal(byteValue, q)
	if err != nil {
		// Return an error if JSON unmarshalling fails.
		return err
	}

	// Return nil to indicate that the operation was successful.
	return nil
}

func (q *Quotes) SaveToFile(settings *util.Settings) error {
	byteValue, err := json.MarshalIndent(q, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(settings.SaveQuotesPath, byteValue, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (q *Quotes) Update(updatedQuote Quote) error {
	for i, quote := range q.QUOTES {
		if quote.ID == updatedQuote.ID {
			q.QUOTES[i] = updatedQuote
			return nil
		}
	}
	return errors.New("quote not found")
}

func (q *Quotes) Remove(index int) {
	q.QUOTES = append(q.QUOTES[:index], q.QUOTES[index+1:]...)
}

func (q *Quotes) PrintQuotes() {

	util.ClearScreen()
	for _, quote := range q.QUOTES {
		fmt.Print(q.Quote(quote.ID))
	}
}

func (q *Quotes) Quote(id int) string {

	var (
		quoteBuffer    bytes.Buffer
		formattedQuote string
	)

	stringFormat := `` + "\n" + util.Cyan + `%s. ` + "\"" + util.Reset + `%s` + `` + util.Cyan + "\"" +
		"\n" + strings.Repeat(" ", 50) + `By %s (%s %s)` + util.Reset + ``

	for _, quote := range q.QUOTES {
		if quote.ID == id {
			formattedQuote = fmt.Sprintf(
				stringFormat,
				strconv.Itoa((quote.ID)),
				quote.QUOTE,
				quote.AUTHOR,
				quote.DATE,
				quote.LANGUAGE)

		}
	}
	// Write the formatted quote into the buffer.
	quoteBuffer.WriteString(formattedQuote)

	// Return the formatted quote string.
	return quoteBuffer.String()
}

func (q *Quotes) Size() int {
	return len(q.QUOTES)
}

func (q *Quotes) Duplicates(searchQuote string, settings *util.Settings) bool {

	if searchQuote == "" || searchQuote == "Unknown" {
		return false
	}

	for _, quote := range q.QUOTES {
		if quote.QUOTE == searchQuote {
			if quote.ID != settings.ID {
				settings.Message = "<< there are dublicates in database. >>"
				settings.ID = quote.ID
				return true
			}
		}
	}
	return false
}

func (q *Quotes) IdOf(index int) (int, error) {
	if index < 0 || index >= len(q.QUOTES) {
		return 0, errors.New("index out of bounds")
	}
	return q.QUOTES[index].ID, nil
}

func (q *Quotes) IndexOf(id int) int {

	for i, quote := range q.QUOTES {
		if quote.ID == id {
			return i
		}
	}
	return -1
}

func (q *Quotes) CreateId() int {
	maxID := 0
	for _, quote := range q.QUOTES {
		if quote.ID > maxID {
			maxID = quote.ID
		}
	}
	return maxID + 1
}

func (q *Quotes) Find(input string) {

	util.ClearScreen()

	for _, quote := range q.QUOTES {

		// Find by Author
		if strings.Contains(strings.ToUpper(quote.AUTHOR), strings.ToUpper(input)) {
			fmt.Print(q.Quote(quote.ID))
		}

		// Find by Quote
		if strings.Contains(strings.ToUpper(quote.QUOTE), strings.ToUpper(input)) {
			fmt.Print(q.Quote(quote.ID))
		}

		// Find by ID
		if quote.ID == util.StringToInt(input) {
			fmt.Print(q.Quote(quote.ID))
		}

	}
}

func (q *Quotes) AppendRandomIDs(settings *util.Settings) {
	for _, quote := range q.QUOTES {
		if !util.ArrayContainsInt(settings.RandomIDs, quote.ID) {
			settings.RandomIDs = append(settings.RandomIDs, quote.ID)
		}
	}
}

func (q *Quotes) ResetID(settings *util.Settings) error {

	// Calculate the index of the last quote. `Size()` returns the total number of quotes,
	// so subtracting 1 gives the index of the last quote.
	index := q.Size() - 1

	// Retrieve the ID of the quote at the calculated index.
	lastId, err := q.IdOf(index)

	// Check if an error occurred while retrieving the ID.
	if err != nil {
		return err
	}

	// Set the retrieved last quote's ID into the provided settings.
	settings.ID = lastId

	// Return nil to indicate that the operation was successful and the settings were updated correctly.
	return nil
}

func (q *Quotes) PromptWithSuggestion(name string, edit string, settings *util.Settings) bool {

	line := liner.NewLiner()
	defer line.Close()

	input, err := line.PromptWithSuggestion("   "+name+": ", edit, -1)
	if err != nil {
		fmt.Println("Error reading input: ", err)
		return false
	}

	if input == "q" {
		settings.Message = "<< previous action aborted by user. >>"
		return false
	}

	if name == "Quote" && q.Duplicates(input, settings) {
		return false
	}

	settings.UserInputs = append(settings.UserInputs, util.FillEmptyInput(input, "Unknown"))

	return true
}

func (q *Quotes) UserInput(settings *util.Settings) bool {

	// empty the old input before getting new values
	if len(settings.UserInputs) > 0 {
		settings.UserInputs = settings.UserInputs[:0]
	}

	type questionPairs struct {
		First  string
		Second string
	}

	// pairs for adding
	questions := [3]questionPairs{{"Quote", ""}, {"Author", ""}, {"Language", "English"}}

	// pairs for updating
	if settings.ID > 0 {
		index := q.IndexOf(settings.ID)
		if index == -1 {
			return false
		}
		questions = [3]questionPairs{{"Quote", q.QUOTES[index].QUOTE}, {"Author", q.QUOTES[index].AUTHOR}, {"Language", q.QUOTES[index].LANGUAGE}}
	}

	// prompt all three questions
	for _, question := range questions {
		validation := q.PromptWithSuggestion(question.First, question.Second, settings)
		if !validation {
			return false
		}
	}
	return true
}

func (q *Quotes) DetailsOf(searchQuote string) (int, string, string) {
	for _, value := range q.QUOTES {
		if value.QUOTE == searchQuote {
			return value.ID, value.QUOTE, value.AUTHOR
		}
	}
	return -1, "", ""
}

func (q *Quotes) GetAllQuotes() []string {
	var sentences []string
	for _, value := range q.QUOTES {
		sentences = append(sentences, value.QUOTE)
	}
	return sentences
}

func (q *Quotes) ReArrangeIDs(settings *util.Settings) {

	for key := range q.QUOTES {
		q.QUOTES[key].ID = key + 1
	}

	q.SaveToFile(settings)

	fmt.Println("Rearraging IDs Done!")
}

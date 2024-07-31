package db

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"vk-quotes/pkg/config"
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

func (q *Quotes) AppendQuote(quote Quote) {
	q.QUOTES = append(q.QUOTES, quote)
}

func (q *Quotes) ReadFromFile(path string, folder string) error {

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

func (q *Quotes) SaveToFile(path string) error {

	byteValue, err := json.MarshalIndent(q, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, byteValue, 0644)
	if err != nil {
		config.Messages = append(config.Messages, err.Error())
		return err
	}

	copyPath := config.CopyFilePathLinux

	if runtime.GOOS == "windows" {
		copyPath = config.CopyFilePathWindows
	}

	err = os.WriteFile(copyPath, byteValue, 0644)
	if err != nil {
		config.Messages = append(config.Messages, "Copy to hdd error: " +  err.Error())
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

func (q *Quotes) PrintAllQuotes() {

	util.ClearScreen()
	for _, quote := range q.QUOTES {
		fmt.Print(q.FormatQuote(quote))
	}
}

func (q *Quotes) FormatQuote(quote Quote) string {

	var (
		quoteBuffer    bytes.Buffer
		formattedQuote string
	)

	stringFormat := `` + "\n" + util.Cyan + `%d. ` + "\"" + util.Reset + `%s` + `` + util.Cyan + "\"" +
		"\n" + strings.Repeat(" ", 50) + `By %s (%s %s)` + util.Reset + "\n" + ``

	formattedQuote = fmt.Sprintf(
		stringFormat,
		quote.ID,
		quote.QUOTE,
		quote.AUTHOR,
		quote.DATE,
		quote.LANGUAGE)

	// Write the formatted quote into the buffer.
	quoteBuffer.WriteString(formattedQuote)

	// Return the formatted quote string.
	return quoteBuffer.String()
}

func (q *Quotes) Size() int {
	return len(q.QUOTES)
}

func (q *Quotes) Duplicates(searchQuote string) bool {

	if searchQuote == "" || searchQuote == "Unknown" {
		return false
	}

	for _, quote := range q.QUOTES {
		if quote.QUOTE == searchQuote {
			if quote.ID != config.ID {
				config.Messages = append(config.Messages, "<< there are dublicates in database. >>")
				config.ID = quote.ID
				return true
			}
		}
	}
	return false
}

func (q *Quotes) FindIdByIndex(index int) (int, error) {
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

func (q *Quotes) Search(command string) Quote {

	var foundQuote Quote

	normalizedCommand := strings.ToUpper(command)

	for _, quote := range q.QUOTES {

		isID, _ := strconv.Atoi(command)

		if quote.ID == isID {
			foundQuote = quote
			break
		}

		normalizedAuthor := strings.ToUpper(quote.AUTHOR)

		if strings.Contains(normalizedAuthor, normalizedCommand) {
			foundQuote = quote
		}

		normalizedQuote := strings.ToUpper(quote.QUOTE)

		if strings.Contains(normalizedQuote, normalizedCommand) {
			foundQuote = quote
		}
	}

	return foundQuote
}

func (q *Quotes) AppendRandomIDs() {
	for _, quote := range q.QUOTES {
		if !util.ArrayContainsInt(config.RandomIDs, quote.ID) {
			config.RandomIDs = append(config.RandomIDs, quote.ID)
		}
	}
}

func (q *Quotes) LastID() int {

	index := q.Size() - 1

	lastId, err := q.FindIdByIndex(index)

	if err != nil {
		return 0
	}

	return lastId
}

func (q *Quotes) PromptWithSuggestion(name string, edit string) bool {

	line := liner.NewLiner()
	defer line.Close()

	input, err := line.PromptWithSuggestion("   "+name+": ", edit, -1)
	if err != nil {
		fmt.Println("Error reading input: ", err)
		return false
	}

	if input == "q" {
		config.Messages = append(config.Messages, "<< previous action aborted by user. >>")
		return false
	}

	if name == "Quote" && q.Duplicates(input) {
		return false
	}

	config.UserInputs = append(config.UserInputs, util.FillEmptyInput(input, "Unknown"))

	return true
}

func (q *Quotes) UserInput(id int) bool {

	// empty the old input before getting new values
	if len(config.UserInputs) > 0 {
		config.UserInputs = config.UserInputs[:0]
	}

	type questionPairs struct {
		First  string
		Second string
	}

	// pairs for adding
	questions := [3]questionPairs{{"Quote", ""}, {"Author", ""}, {"Language", "English"}}

	// pairs for updating
	if id > 0 {
		index := q.IndexOf(id)
		if index == -1 {
			return false
		}
		questions = [3]questionPairs{{"Quote", q.QUOTES[index].QUOTE}, {"Author", q.QUOTES[index].AUTHOR}, {"Language", q.QUOTES[index].LANGUAGE}}
	}

	// prompt all three questions
	for _, question := range questions {
		validation := q.PromptWithSuggestion(question.First, question.Second)
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

func (q *Quotes) ResetIDs(quotes *Quotes) {

	for key := range q.QUOTES {
		q.QUOTES[key].ID = key + 1
	}

	q.SaveToFile(config.SaveQuotesPath)

	config.ID = q.LastID()

	fmt.Println("Reset IDs Done!")
}

func (q *Quotes) TopAuthors() string {

	/* Get All Author Names */
	var authors []string
	for _, value := range q.QUOTES {
		if !util.ArrayContainsString(authors, value.AUTHOR) {
			authors = append(authors, value.AUTHOR)
		}
	}

	/* Count Authors By Name */
	authorsMap := make(map[string]int)
	for _, name := range authors {
		for _, value := range q.QUOTES {
			if value.AUTHOR == name {
				authorsMap[name] += 1
			}
		}
	}

	/* Make Pairs */
	type pair struct {
		name  string
		count int
	}

	var pairs []pair
	for key, value := range authorsMap {
		pairs = append(pairs, pair{key, value})
	}

	/* Sort by count */
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})

	// Make a string
	authorsString := ""
	for i := 0; i < len(pairs) && i < 10; i++ {
		authorsString += "\n" + strconv.Itoa(pairs[i].count) + " " + util.Cyan + pairs[i].name + util.Reset
	}

	return authorsString
}

func (q *Quotes) TopLanguages() string {

	languages := []string{"English", "Russian"}

	// Count languages
	languagesMap := make(map[string]int)
	for _, name := range languages {
		for _, value := range q.QUOTES {
			if value.LANGUAGE == name {
				languagesMap[name] += 1
			}
		}
	}

	// Make a string
	languagesString := ""
	for name, num := range languagesMap {
		languagesString += "\n" + strconv.Itoa(num) + " " + util.Yellow + name + util.Reset
	}
	return languagesString
}

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
	"time"
	"vk-quotes/pkg/config"
	"vk-quotes/pkg/util"

	"github.com/peterh/liner"
)

type Quote struct {
	ID       int    `json:"id"`
	QUOTE    string `json:"quote"`
	AUTHOR   string `json:"author"`
	LANGUAGE string `json:"language"`
	DATE     string `json:"date"`
}

type Quotes struct {
	QUOTES []Quote `json:"quotes"` // Slice containing multiple Quote instances.
}

func (q *Quotes) UserInput(id int) bool {

	config.UserInputs = []string{} // Clear previous user inputs

	// Define question pairs with default values
	type questionPairs struct {
		Prompt  string
		Default string
	}
	questions := []questionPairs{
		{"Quote", ""},
		{"Author", ""},
		{"Language", "English"},
	}

	// If updating, populate questions with existing values
	if id > 0 {
		index := q.IndexOf(id)
		if index == -1 {
			return false
		}
		quote := q.QUOTES[index]
		questions = []questionPairs{
			{"Quote", quote.QUOTE},
			{"Author", quote.AUTHOR},
			{"Language", quote.LANGUAGE},
		}
	}

	// Prompt user for input
	for _, question := range questions {
		if !q.PromptWithSuggestion(question.Prompt, question.Default) {
			return false
		}
	}
	return true
}

func (q *Quotes) AppendQuote(quote Quote) {
	q.QUOTES = append(q.QUOTES, quote)
}

func (q *Quotes) ReadFromFile() {

	path := config.LocalPath

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(byteValue, q)
	if err != nil {
		panic(err)
	}
}

func (q *Quotes) SaveToFile(message string) {

	byteValue, err := json.MarshalIndent(q, "", "  ")
	if err != nil {
		panic(err)
	}

	path := config.LocalPath

	err = os.WriteFile(path, byteValue, 0644)
	if err != nil {
		panic(err)
	}

	config.AddMessage(message)
	
	// Backup
	q.Backup(byteValue)
}

func (q *Quotes) Backup(byteValue []byte) {

	currentTime := time.Now()
	layout := "(02.01.2006_15-04-05)"

	backupPath := "/media/veikko/VK DATA/DATABASES/QUOTES/" + strconv.Itoa(q.Size()) + ". quotes " + currentTime.Format(layout) + ".json"

	if runtime.GOOS == "windows" {
		backupPath = "D:\\DATABASES\\QUOTES\\" + strconv.Itoa(q.Size()) + ".json"
	}

	err := os.WriteFile(backupPath, byteValue, 0644)
	if err != nil {
		message := config.Red + "<< No Backup >>" + config.Reset
		config.AddMessage(message)
		config.AddMessage(err.Error())
		return
	}
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

	util.PressAnyKey()
}

func (q *Quotes) FormatQuote(quote Quote) string {

	var (
		quoteBuffer    bytes.Buffer
		formattedQuote string
	)

	stringFormat := `` + "\n" + config.Cyan + `%d. ` + "\"" + config.Reset + `%s` + `` + config.Cyan + "\"" +
		"\n" + strings.Repeat(" ", 50) + `By %s (%s %s)` + config.Reset + "\n" + ``

	formattedQuote = fmt.Sprintf(
		stringFormat,
		quote.ID,
		quote.QUOTE,
		quote.AUTHOR,
		quote.DATE,
		quote.LANGUAGE)

	quoteBuffer.WriteString(formattedQuote)

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
			if quote.ID != config.MainQuoteID {
				message := config.Red + "There are dublicates in database" + config.Reset
				config.AddMessage(message)
				config.MainQuoteID = quote.ID
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

func (q *Quotes) PrintQuote(command string) {

	for _, quote := range q.QUOTES {

		isID, _ := strconv.Atoi(command)
		if quote.ID == isID {
			fmt.Println(q.FormatQuote(quote))
			return
		}

		normalizedAuthor := strings.ToLower(quote.AUTHOR)
		normalizedQuote := strings.ToLower(quote.QUOTE)

		if strings.Contains(normalizedAuthor, command) || strings.Contains(normalizedQuote, command) {
			fmt.Println(q.FormatQuote(quote))
		}

	}
}

func (q *Quotes) FindByID(id int) Quote {

	var foundQuote Quote

	for _, quote := range q.QUOTES {
		if quote.ID == id {
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
		message := config.Red + "Previous action aborted by user" + config.Reset
		config.AddMessage(message)
		return false
	}

	if name == "Quote" && q.Duplicates(input) {
		return false
	}

	if name == "Author" && input == "" {
		input = "Unknown"
	}

	if name == "Language" && input == "" {
		input = "Russian"
	}

	config.UserInputs = append(config.UserInputs, input)

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

	q.SaveToFile("<< All ID's Reset >>")

	config.MainQuoteID = q.LastID()
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
		authorsString += "\n" + strconv.Itoa(pairs[i].count) + " " + config.Cyan + pairs[i].name + config.Reset
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
		languagesString += "\n" + strconv.Itoa(num) + " " + config.Yellow + name + config.Reset
	}
	return languagesString
}

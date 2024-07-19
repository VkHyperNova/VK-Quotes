package db

import (
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
	ID       int    `json:"id"`
	QUOTE    string `json:"quote"`
	AUTHOR   string `json:"author"`
	LANGUAGE string `json:"language"`
	DATE     string `json:"date"`
}

type Quotes struct {
	QUOTES []Quote `json:"quotes"`
}

func (q *Quotes) Add(quote Quote) {
	q.QUOTES = append(q.QUOTES, quote)
}

func (q *Quotes) ReadFromFile() error {

	path := "./database/quotes.json"

	if _, err := os.Stat(path); os.IsNotExist(err) {
		_ = os.Mkdir("database", 0700)
		err = os.WriteFile(path, []byte(`{"quotes": []}`), 0644)
		if err != nil {
			fmt.Println(err)
		}
		util.PrintRed("New Database Created!\n")
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, q)
	if err != nil {
		return err
	}

	return nil
}

func (q *Quotes) SaveToFile() error {
	byteValue, err := json.MarshalIndent(q, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile("./database/quotes.json", byteValue, 0644)
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

func (q *Quotes) Delete(settings *util.Settings) error {

	index := q.FindIndex(settings.ID)
	if index == -1 {
		settings.Message = fmt.Sprintf("<< %d Index Not Found! >>", settings.ID)
	}

	q.QUOTES = append(q.QUOTES[:index], q.QUOTES[index+1:]...)

	return nil
}

func (q *Quotes) PrintQuotes() {
	for _, quote := range q.QUOTES {
		q.PrintQuote(quote.ID)
	}
}

func (q *Quotes) PrintQuote(id int) error {
	for _, quote := range q.QUOTES {
		if quote.ID == id {
			util.PrintCyan("\n\n" + strconv.Itoa((quote.ID)) + ". \"")
			util.PrintGray(quote.QUOTE)
			util.PrintCyan("\"\n" + strings.Repeat(" ", 50) + " By " + quote.AUTHOR + " (" + quote.DATE + " " + quote.LANGUAGE + ")\n")
			return nil
		}
	}
	return errors.New("quote not found")
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

func (q *Quotes) FindId(index int) (int, error) {
	if index < 0 || index >= len(q.QUOTES) {
		return 0, errors.New("index out of bounds")
	}
	return q.QUOTES[index].ID, nil
}

func (q *Quotes) FindIndex(id int) int {

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

func (q *Quotes) PrintByAuthor(author string) {
	for _, quote := range q.QUOTES {
		if strings.Contains(strings.ToUpper(quote.AUTHOR), strings.ToUpper(author)) {
			q.PrintQuote(quote.ID)
		}
	}
}

func (q *Quotes) FindIDs(settings *util.Settings) {
	for _, quote := range q.QUOTES {
		if !util.ArrayContainsInt(settings.RandomIDs, quote.ID) {
			settings.RandomIDs = append(settings.RandomIDs, quote.ID)
		}
	}
}

func (q *Quotes) SetToLastID(settings *util.Settings) error {
	index := q.Size() - 1

	lastId, err := q.FindId(index)

	if err != nil {
		return err
	}

	settings.ID = lastId

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

	if len(settings.UserInputs) > 0 {
		settings.UserInputs = settings.UserInputs[:0]
	}

	type Pairs struct {
		First  string
		Second string
	}

	questions := [3]Pairs{{"Quote", ""}, {"Author", ""}, {"Language", "English"}}

	/* If Update */
	if settings.ID > 0 {
		index := q.FindIndex(settings.ID)
		if index == -1 {
			return false
		}
		questions = [3]Pairs{{"Quote", q.QUOTES[index].QUOTE}, {"Author", q.QUOTES[index].AUTHOR}, {"Language", q.QUOTES[index].LANGUAGE}}
	}

	for _, question := range questions {
		validation := q.PromptWithSuggestion(question.First, question.Second, settings)
		if !validation {
			return false
		}
	}
	return true
}

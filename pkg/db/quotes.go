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

func (q *Quotes) ReadFromFile(filename string) error {
	file, err := os.Open(filename)
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

func (q *Quotes) SaveToFile(filename string) error {
	byteValue, err := json.MarshalIndent(q, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, byteValue, 0644)
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

func (q *Quotes) Delete(id int) error {

	index := q.FindIndex(id)

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

func (q *Quotes) Duplicates(searchQuote string) int {

	if searchQuote == "" || searchQuote == "Unknown" {
		return -1
	}
	for _, quote := range q.QUOTES {
		if quote.QUOTE == searchQuote {
			return quote.ID
		}
	}
	return -1
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
	fmt.Println("quote not found")
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

func (q *Quotes) FindByAuthor(author string) {
	for _, quote := range q.QUOTES {
		if strings.Contains(strings.ToUpper(quote.AUTHOR), strings.ToUpper(author)) {
			q.PrintQuote(quote.ID)
		}
	}
}

func (q *Quotes) FindIds() {
	for _, quote := range q.QUOTES {
		if !util.ArrayContainsInt(IDs, quote.ID) {
			IDs = append(IDs, quote.ID)
		}
	}
}

func (q *Quotes) GetLastId() int {
	index := q.Size() - 1

	lastId, err := q.FindId(index)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	return lastId
}

package db

import "errors"

func (q *Quotes) AddQuote(quote Quote) {
	q.QUOTES = append(q.QUOTES, quote)
}

func (q *Quotes) UpdateQuote(updatedQuote Quote) error {
	for i, quote := range q.QUOTES {
		if quote.ID == updatedQuote.ID {
			q.QUOTES[i] = updatedQuote
			return nil
		}
	}
	return errors.New("quote not found")
}

func (q *Quotes) DeleteQuote(index int) {
	q.QUOTES = append(q.QUOTES[:index], q.QUOTES[index+1:]...)
}
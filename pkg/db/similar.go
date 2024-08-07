/*
Package db provides functionality for handling and processing a database of quotes. It includes the ability to
identify similar quotes using TF-IDF vectors and cosine similarity, save the results to a file, and display progress
while running long-running tasks.

The package includes the following main types and functions:

Types:
  - SimilarQuotePairs: Represents a pair of similar quotes with their respective details, including IDs, quotes, and authors.
  - SimilarQuotes: Holds a collection of SimilarQuotePairs and provides methods for adding pairs and saving to a file.

Functions:
  - (s *SimilarQuotes) Add(quote SimilarQuotePairs): Appends a new SimilarQuotePairs instance to the SimilarQuotes slice.
  - (q *SimilarQuotes) SaveToFile() error: Marshals the SimilarQuotes instance to JSON format and saves it to a file.
  - waitGroupDone(wg *sync.WaitGroup) <-chan struct{}: Returns a channel that is closed when the provided WaitGroup is done.
  - processSimilarQuotes(quotes *Quotes, similar *SimilarQuotes): Finds and processes similar quotes, adding the pairs to SimilarQuotes and saving the results.
  - calculateTFIDF(sentences []string) []map[string]float64: Calculates the TF-IDF vectors for a given set of sentences.
  - cosineSimilarity(vec1, vec2 map[string]float64) float64: Calculates the cosine similarity between two TF-IDF vectors.
  - findSimilarSentences(sentences []string, tfidfVectors []map[string]float64, threshold float64) [][2]string: Identifies pairs of sentences that are similar based on their TF-IDF vectors.
  - RunTaskWithProgress(quotes *Quotes): Runs a long-running task with a progress bar to indicate progress, processing quotes and updating the progress bar until the task is complete.

Example Usage:

To use this package, you would typically follow these steps:

1. Load your quotes into a Quotes instance.
2. Call the RunTaskWithProgress function to process similar quotes with a progress bar.

	// Create an instance of Quotes and load quotes from a file.
	quotes := &Quotes{}
	err := quotes.ReadFromFile("path/to/quotes.json")
	if err != nil {
	  log.Fatalf("Failed to read quotes: %v", err)
	}

	// Run the long-running task with a progress bar.
	RunTaskWithProgress(quotes)

The package handles reading quotes, processing them to find similar pairs based on TF-IDF and cosine similarity,
and saving the results to a file. It also provides a progress bar to give feedback on the task's progress.
*/
package db

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"sync"
	"time"
	"vk-quotes/pkg/config"

	"github.com/cheggaaa/pb/v3"
	"github.com/jdkato/prose/v2"
)

type SimilarQuotePairs struct {
	FirstID      int    `json:"firstid"`
	FirstQuote   string `json:"firstquote"`
	FirstAuthor  string `json:"firstauthor"`
	SecondID     int    `json:"secondid"`
	SecondQuote  string `json:"secondquote"`
	SecondAuthor string `json:"secondauthor"`
}

type SimilarQuotes struct {
	SimilarQuotes []SimilarQuotePairs `json:"similarquotes"`
}

func (s *SimilarQuotes) Add(quote SimilarQuotePairs) {
	// Append the provided quote to the SimilarQuotes slice.
	s.SimilarQuotes = append(s.SimilarQuotes, quote)
}

func (s *SimilarQuotes) SaveToFile() {

	byteValue, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err)
	}

	path := config.SimilarPath

	err = os.WriteFile(path, byteValue, 0644)
	if err != nil {
		panic(err)
	}
}

func waitGroupDone(wg *sync.WaitGroup) <-chan struct{} {
	ch := make(chan struct{})

	go func() {
		wg.Wait()

		close(ch)
	}()

	return ch
}

func processSimilarQuotes(quotes *Quotes, similar *SimilarQuotes) {

	sentences := quotes.GetAllQuotes()

	tfidfVectors := calculateTFIDF(sentences)

	threshold := 0.5
	similarSentences := findSimilarSentences(sentences, tfidfVectors, threshold)

	for _, pair := range similarSentences {

		firstID, firstQuote, FirstAuthor := quotes.DetailsOf(pair[0])

		secondID, secondQuote, secondAuthor := quotes.DetailsOf(pair[1])

		SimilarQuotePairs := SimilarQuotePairs{
			FirstID:      firstID,
			FirstQuote:   firstQuote,
			FirstAuthor:  FirstAuthor,
			SecondID:     secondID,
			SecondQuote:  secondQuote,
			SecondAuthor: secondAuthor,
		}

		similar.Add(SimilarQuotePairs)
	}

	similar.SaveToFile()

	message := config.Green+"Find Similar Quotes Process Done"+config.Reset
	config.AddMessage(message)
}

func calculateTFIDF(sentences []string) []map[string]float64 {
	var tfidfVectors []map[string]float64

	documentCount := len(sentences)

	wordDocCount := make(map[string]int)

	tokenizedSentences := make([][]string, documentCount)
	for i, sentence := range sentences {
		doc, err := prose.NewDocument(sentence)
		if err != nil {
			log.Fatalf("Failed to tokenize sentence: %v", err)
		}

		tokens := []string{}
		for _, token := range doc.Tokens() {
			if token.Tag != "PUNCT" {
				tokens = append(tokens, strings.ToLower(token.Text))

				wordDocCount[strings.ToLower(token.Text)]++
			}
		}

		tokenizedSentences[i] = tokens
	}

	for _, tokens := range tokenizedSentences {
		tfidf := make(map[string]float64)

		wordCount := len(tokens)

		wordFreq := make(map[string]int)
		for _, token := range tokens {
			wordFreq[token]++
		}

		for word, count := range wordFreq {
			tf := float64(count) / float64(wordCount)

			idf := math.Log(float64(documentCount) / float64(1+wordDocCount[word]))

			tfidf[word] = tf * idf
		}

		tfidfVectors = append(tfidfVectors, tfidf)
	}

	return tfidfVectors
}

func cosineSimilarity(vec1, vec2 map[string]float64) float64 {
	var dotProduct, magnitude1, magnitude2 float64

	for word, val1 := range vec1 {
		val2, exists := vec2[word]
		if exists {
			dotProduct += val1 * val2
		}
		magnitude1 += val1 * val1
	}

	for _, val2 := range vec2 {
		magnitude2 += val2 * val2
	}

	if magnitude1 == 0 || magnitude2 == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(magnitude1) * math.Sqrt(magnitude2))
}

func findSimilarSentences(sentences []string, tfidfVectors []map[string]float64, threshold float64) [][2]string {
	var similarPairs [][2]string

	for i, vec1 := range tfidfVectors {
		for j, vec2 := range tfidfVectors {
			if i != j {
				similarity := cosineSimilarity(vec1, vec2)
				if similarity >= threshold {
					similarPairs = append(similarPairs, [2]string{sentences[i], sentences[j]})
				}
			}
		}
	}
	return similarPairs
}

func FindSimilarQuotes(quotes *Quotes) {

	similar := SimilarQuotes{}

	bar := pb.New(0).Start()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		processSimilarQuotes(quotes, &similar)
	}()

	go func() {
		for {
			select {
			case <-time.After(100 * time.Millisecond):
				bar.Add(1)
			case <-waitGroupDone(&wg):
				bar.Finish()
				fmt.Println("Task completed!")
				return
			}
		}
	}()

	wg.Wait()
}

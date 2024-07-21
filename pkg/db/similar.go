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
	"vk-quotes/pkg/util"

	"github.com/cheggaaa/pb/v3"
	"github.com/jdkato/prose/v2"
)

// SimilarQuotePairs represents a pair of similar quotes with their respective details.
type SimilarQuotePairs struct {
	FirstID      int    `json:"firstid"`
	FirstQuote   string `json:"firstquote"`
	FirstAuthor  string `json:"firstauthor"`
	SecondID     int    `json:"secondid"`
	SecondQuote  string `json:"secondquote"`
	SecondAuthor string `json:"secondauthor"`
}

// SimilarQuotes holds a collection of SimilarQuotePairs.
type SimilarQuotes struct {
	SimilarQuotes []SimilarQuotePairs `json:"similarquotes"`
}

// Add appends a new SimilarQuotePairs instance to the SimilarQuotes slice.
func (s *SimilarQuotes) Add(quote SimilarQuotePairs) {
	// Append the provided quote to the SimilarQuotes slice.
	s.SimilarQuotes = append(s.SimilarQuotes, quote)
}

// SaveToFile marshals the SimilarQuotes instance to JSON format and saves it to a file.
// It returns an error if any step in the process fails.
func (q *SimilarQuotes) SaveToFile() error {
	// Marshal the SimilarQuotes instance to a JSON byte slice with indentation for readability.
	byteValue, err := json.MarshalIndent(q, "", "  ")
	if err != nil {
		// Return the error if marshalling fails.
		return err
	}

	// Write the JSON byte slice to a file named "similarquotes.json" with read/write permissions for the owner, and read permissions for others.
	err = os.WriteFile("./similarquotes.json", byteValue, 0644)
	if err != nil {
		// Return the error if writing to the file fails.
		return err
	}

	// Return nil to indicate success.
	return nil
}

// waitGroupDone returns a channel that is closed when the provided sync.WaitGroup is done.
// This allows you to use the channel to wait for the completion of the WaitGroup in a select statement.
func waitGroupDone(wg *sync.WaitGroup) <-chan struct{} {
	// Create a channel that will be returned.
	ch := make(chan struct{})

	// Start a goroutine to wait for the WaitGroup to be done.
	go func() {
		// Wait for the WaitGroup to finish.
		wg.Wait()

		// Close the channel to signal that the WaitGroup is done.
		close(ch)
	}()

	// Return the channel to the caller.
	return ch
}

// processSimilarQuotes finds and processes similar quotes from the provided quotes database,
// then adds the similar quote pairs to the SimilarQuotes collection and saves the results.
func processSimilarQuotes(quotes *Quotes, similar *SimilarQuotes) {

	// Step 1: Extract all quote strings from the quotes database into a slice.
	sentences := quotes.GetAllQuotes()

	// Step 2: Calculate the TF-IDF vectors for each quote sentence.
	tfidfVectors := calculateTFIDF(sentences)

	// Step 3: Define a threshold for similarity and find pairs of similar sentences.
	threshold := 0.5
	similarSentences := findSimilarSentences(sentences, tfidfVectors, threshold)

	// Step 4: Process each pair of similar sentences.
	for _, pair := range similarSentences {

		// Extract information for the first quote in the pair
		firstID, firstQuote, FirstAuthor := quotes.DetailsOf(pair[0])

		// Extract information for the second quote	in the pair
		secondID, secondQuote, secondAuthor := quotes.DetailsOf(pair[1])

		// Create an instance of SimilarQuotePairs with the details of both quotes.
		SimilarQuotePairs := SimilarQuotePairs{
			FirstID:      firstID,
			FirstQuote:   firstQuote,
			FirstAuthor:  FirstAuthor,
			SecondID:     secondID,
			SecondQuote:  secondQuote,
			SecondAuthor: secondAuthor,
		}

		// Add the created SimilarQuotePairs instance to the collection of similar quotes.
		similar.Add(SimilarQuotePairs)
	}

	// Step 5: Save the similar quotes data to a file.
	similar.SaveToFile()

	// Print a confirmation message to indicate the process is complete.
	util.PrintGreen("Process Done!:\n\n")
}

// calculateTFIDF calculates the TF-IDF vectors for a given set of sentences.
// TF-IDF (Term Frequency-Inverse Document Frequency) is a statistical measure used to evaluate the importance of a word in a document relative to a collection of documents.
func calculateTFIDF(sentences []string) []map[string]float64 {
	// Initialize a slice to hold the TF-IDF vectors for each sentence.
	var tfidfVectors []map[string]float64

	// Get the total number of sentences/documents.
	documentCount := len(sentences)

	// Map to count the number of documents in which each word appears.
	wordDocCount := make(map[string]int)

	// Tokenize each sentence and count word occurrences across all documents.
	tokenizedSentences := make([][]string, documentCount)
	for i, sentence := range sentences {
		// Tokenize the sentence using prose library.
		doc, err := prose.NewDocument(sentence)
		if err != nil {
			log.Fatalf("Failed to tokenize sentence: %v", err)
		}

		// Slice to hold tokens for the current sentence.
		tokens := []string{}
		for _, token := range doc.Tokens() {
			// Ignore punctuation tokens.
			if token.Tag != "PUNCT" {
				// Convert token to lowercase and add it to the tokens slice.
				tokens = append(tokens, strings.ToLower(token.Text))

				// Increment the document count for the token.
				wordDocCount[strings.ToLower(token.Text)]++
			}
		}

		// Store the tokens for the current sentence.
		tokenizedSentences[i] = tokens
	}

	// Calculate the TF-IDF for each sentence.
	for _, tokens := range tokenizedSentences {
		// Map to hold the TF-IDF values for the current sentence.
		tfidf := make(map[string]float64)

		// Get the total number of words in the current sentence.
		wordCount := len(tokens)

		// Map to count the frequency of each word in the current sentence.
		wordFreq := make(map[string]int)
		for _, token := range tokens {
			wordFreq[token]++
		}

		// Calculate the TF-IDF value for each word in the sentence.
		for word, count := range wordFreq {
			// Term Frequency (TF) calculation.
			tf := float64(count) / float64(wordCount)

			// Inverse Document Frequency (IDF) calculation.
			idf := math.Log(float64(documentCount) / float64(1+wordDocCount[word]))

			// TF-IDF calculation.
			tfidf[word] = tf * idf
		}

		// Add the TF-IDF map for the current sentence to the tfidfVectors slice.
		tfidfVectors = append(tfidfVectors, tfidf)
	}

	// Return the slice of TF-IDF vectors.
	return tfidfVectors
}

// cosineSimilarity calculates the cosine similarity between two TF-IDF vectors.
// Cosine similarity is a measure of similarity between two non-zero vectors
// that calculates the cosine of the angle between them.
func cosineSimilarity(vec1, vec2 map[string]float64) float64 {
	var dotProduct, magnitude1, magnitude2 float64

	// Calculate the dot product and the magnitude of the first vector.
	for word, val1 := range vec1 {
		// Check if the word also exists in the second vector.
		val2, exists := vec2[word]
		if exists {
			// If the word exists in both vectors, add to the dot product.
			dotProduct += val1 * val2
		}
		// Accumulate the square of the value for the first vector's magnitude.
		magnitude1 += val1 * val1
	}

	// Calculate the magnitude of the second vector.
	for _, val2 := range vec2 {
		// Accumulate the square of the value for the second vector's magnitude.
		magnitude2 += val2 * val2
	}

	// Check if either of the magnitudes is zero, which would make cosine similarity undefined.
	if magnitude1 == 0 || magnitude2 == 0 {
		// Return 0 as the similarity in case of zero magnitude to avoid division by zero.
		return 0
	}

	// Calculate and return the cosine similarity.
	return dotProduct / (math.Sqrt(magnitude1) * math.Sqrt(magnitude2))
}

// findSimilarSentences identifies pairs of sentences that are similar based on their TF-IDF vectors.
// It returns a slice of sentence pairs where the cosine similarity between their TF-IDF vectors is above the specified threshold.
func findSimilarSentences(sentences []string, tfidfVectors []map[string]float64, threshold float64) [][2]string {
	// Slice to hold pairs of similar sentences.
	var similarPairs [][2]string

	// Loop through each TF-IDF vector.
	for i, vec1 := range tfidfVectors {
		// Compare each TF-IDF vector with every other TF-IDF vector.
		for j, vec2 := range tfidfVectors {
			// Ensure that a sentence is not compared with itself.
			if i != j {
				// Calculate the cosine similarity between the two TF-IDF vectors.
				similarity := cosineSimilarity(vec1, vec2)
				// If the similarity is above or equal to the threshold, consider the sentences as similar.
				if similarity >= threshold {
					// Append the pair of similar sentences to the similarPairs slice.
					similarPairs = append(similarPairs, [2]string{sentences[i], sentences[j]})
				}
			}
		}
	}

	// Return the slice of similar sentence pairs.
	return similarPairs
}

// RunTaskWithProgress runs a long-running task with a progress bar to indicate progress.
// It processes quotes and updates the progress bar until the task is complete.
func RunTaskWithProgress(quotes *Quotes) {

	// Create an instance of SimilarQuotes to hold similar quotes.
	similar := SimilarQuotes{}

	// Create a new progress bar with an unknown total.
	bar := pb.New(0).Start()

	// Create a WaitGroup to signal the completion of the long-running task.
	var wg sync.WaitGroup
	wg.Add(1)

	// Run the long-running task in a separate goroutine.
	go func() {
		// Defer the Done call to ensure the WaitGroup counter is decremented when the task completes.
		defer wg.Done()
		// Process similar quotes, a long-running task.
		processSimilarQuotes(quotes, &similar)
	}()

	// Run the progress bar updates in a separate goroutine.
	go func() {
		for {
			select {
			// Update the progress bar every 100 milliseconds.
			case <-time.After(100 * time.Millisecond):
				bar.Add(1)
			// When the WaitGroup is done, finish the progress bar and exit the goroutine.
			case <-waitGroupDone(&wg):
				bar.Finish()
				fmt.Println("Task completed!")
				return
			}
		}
	}()

	// Wait for the long-running task to finish.
	wg.Wait()
}
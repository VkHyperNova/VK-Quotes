
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

func (s *SimilarQuotes) SaveToFile() {

	byteValue, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err)
	}

	path := "./QUOTES/similiar.json"

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

	var sentences []string
	for _, value := range quotes.QUOTES {
		sentences = append(sentences, value.QUOTE)
	}

	tfidfVectors := calculateTFIDF(sentences)

	threshold := 0.5
	similarSentences := findSimilarSentences(sentences, tfidfVectors, threshold)

	for _, pair := range similarSentences {

		firstQuote := quotes.FindQuoteByQuote(pair[0])
		secondQuote := quotes.FindQuoteByQuote(pair[1])

		SimilarQuotePairs := SimilarQuotePairs{
			FirstID:      firstQuote.ID,
			FirstQuote:   firstQuote.QUOTE,
			FirstAuthor:  firstQuote.AUTHOR,
			SecondID:     secondQuote.ID,
			SecondQuote:  secondQuote.QUOTE,
			SecondAuthor: secondQuote.AUTHOR,
		}

		similar.SimilarQuotes = append(similar.SimilarQuotes, SimilarQuotePairs)
	}

	similar.SaveToFile()

	message := config.Green + "Find Similar Quotes Process Done" + config.Reset
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


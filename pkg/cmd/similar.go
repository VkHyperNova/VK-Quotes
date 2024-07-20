package cmd

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"

	"github.com/jdkato/prose/v2"
)

func similarities(quotes *db.Quotes) {
	var sentences []string
	for _, value := range quotes.QUOTES {
		sentences = append(sentences, value.QUOTE)
	}

	tfidfVectors := calculateTFIDF(sentences)

	threshold := 0.5
	similarSentences := findSimilarSentences(sentences, tfidfVectors, threshold)

	util.PrintGreen("Similarities:\n\n")

	for _, pair := range similarSentences {
		printSimilarQuotes(quotes, pair[0])
		printSimilarQuotes(quotes, pair[1])
		fmt.Println()
	}
	util.PressAnyKey()
}

func printSimilarQuotes(quotes *db.Quotes, searchQuote string) {
	for _, value := range quotes.QUOTES {
		if value.QUOTE == searchQuote {
			util.PrintRed(strconv.Itoa(value.ID))
			util.PrintRed(". " + searchQuote)
		}
	}
}

func calculateTFIDF(sentences []string) []map[string]float64 {
	var tfidfVectors []map[string]float64
	documentCount := len(sentences)
	wordDocCount := make(map[string]int)

	// Tokenize sentences and count word occurrences
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
	// Calculate TF-IDF for each sentence
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

package db

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"vk-quotes/pkg/color"
	"vk-quotes/pkg/config"

	"vk-quotes/pkg/util"

	"github.com/cheggaaa/pb/v3"
	"github.com/jdkato/prose/v2"
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

/* Main */

func (q *Quotes) Add() error {

	newQuote, err := q.promptEntry(Quote{})
	if err != nil {
		return err
	}

	q.QUOTES = append(q.QUOTES, newQuote)

	q.save()

	return nil
}

func (q *Quotes) Update(id int) error {

	index, err := q.indexOf(id)
	if err != nil {
		return err
	}

	updated, err := q.promptEntry((q.QUOTES)[index])
	if err != nil {
		return err
	}

	(q.QUOTES)[index] = updated

	q.save()

	return nil

}

func (q *Quotes) Delete(id int) error {

	index, err := q.indexOf(id)
	if err != nil {
		return err
	}

	fmt.Println(FormatQuote(q.QUOTES[index]))

	confirm, err := util.PromptWithSuggestion("(y/n): ", "n")
	if err != nil {
		return err
	}
	if confirm != "y" && confirm != "yes" {
		return fmt.Errorf("Aborted")
	}

	q.QUOTES = append(q.QUOTES[:index], q.QUOTES[index+1:]...)

	// Reset IDs
	for key := range q.QUOTES {
		q.QUOTES[key].ID = key + 1
	}

	return q.save()
}

func (q *Quotes) Read() {

	// Append All Quotes IDs
	for _, quote := range q.QUOTES {
		if !util.ArrayContainsInt(config.RandomIDs, quote.ID) {
			config.RandomIDs = append(config.RandomIDs, quote.ID)
		}
	}

	for len(config.RandomIDs) != 0 {

		config.Counter += 1

		source := rand.NewSource(time.Now().UnixNano())
		r := rand.New(source)
		randomIndex := r.Intn(len(config.RandomIDs))

		config.MainQuoteID = config.RandomIDs[randomIndex]
		config.DeleteUsedID(randomIndex)

		count := config.Counter
		size := len(q.QUOTES)
		percentage := float64(count) / float64(size) * 100
		config.ReadCounter = fmt.Sprintf("<< Reading [%d] %.0f%% >>", count, percentage)

		q.PrintCLI()

		var quit string
		fmt.Scanln(&quit)

		if quit == "q" {
			break
		}
	}

	message := color.Yellow + "Reading Done" + color.Reset
	config.AddMessage(message)

	config.RandomIDs = config.RandomIDs[:0]
	q.SetToDefaultQuote()
	config.Counter = 0
	config.ReadCounter = ""
}

func (q *Quotes) Export() error {

	input, err := util.PromptWithSuggestion("Export db to d drive? (y/n) ", "n")
	if err != nil {
		return err
	}

	if input == "y" || input == "yes" {

		if err := util.InitBackupStorage(); err != nil {
			return err
		}

		if err := q.LoadFromFile(config.LocalFile); err != nil {
			return fmt.Errorf("load from file: %w", err)
		}

		finance, err := json.MarshalIndent(q, "", "  ")
		if err != nil {
			return err
		}

		if err := os.WriteFile(config.BackupFile, finance, 0644); err != nil {
			return err
		}

		fmt.Printf("Database exported to %s\nPress Enter!", config.BackupFile)
		return nil
	}

	fmt.Println("Export canceled!")
	return nil
}

func (q *Quotes) Import() error {

	input, err := util.PromptWithSuggestion("Import db from d drive? (y/n) ", "n")
	if err != nil {
		return err
	}

	if input == "y" || input == "yes" {

		if err := util.InitBackupStorage(); err != nil {
			return err
		}

		if err := q.LoadFromFile(config.BackupFile); err != nil {
			return fmt.Errorf("load from file: %w", err)
		}

		finance, err := json.MarshalIndent(q, "", "  ")
		if err != nil {
			return err
		}

		if err := os.WriteFile(config.LocalFile, finance, 0644); err != nil {
			return err
		}

		fmt.Printf("Database imported from %s\nPress Enter!", config.BackupFile)
		return nil
	}

	fmt.Println("Import canceled!")

	return nil
}

func (q *Quotes) newID() int {

	maxID := 0

	// maxID := q.QUOTES[len(q.QUOTES) - 1].ID

	for _, quote := range q.QUOTES {
		if quote.ID > maxID {
			maxID = quote.ID
		}
	}

	return maxID + 1
}

func (q *Quotes) LoadFromFile(path string) error {

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

/* Print */

func (q *Quotes) PrintCLI() {

	util.ClearScreen()

	if config.MainQuoteID <= 0 {
		q.SetToDefaultQuote()
	}

	nowPlaying := "Now playing: Flute.mp3"

	stringFormat := `` +
		color.Cyan + "VK-Quotes" + color.Reset + " %s" + "\n" + // Program Name
		color.Purple + "%s" + color.Reset + "\n" + // Now Playing
		"%s" + // Messages
		color.Cyan + `%s` + color.Reset + // Read Counter
		"%s" + // Last Quote
		color.Yellow + `%s` + color.Reset + // Commands
		``

	formattedLastQuote := FormatQuote(q.QUOTES[len(q.QUOTES) - 1])

	messages := config.FormatMessages()

	commands := "\n< add, update, delete, random, find, read, history, unmount, export, import, stats, findsimilar, quit\n"

	cli := fmt.Sprintf(stringFormat, config.ProgramVersion, nowPlaying, messages, config.ReadCounter, formattedLastQuote, commands)

	fmt.Print(cli)
}

func (q *Quotes) History() {

	util.ClearScreen()

	for _, quote := range q.QUOTES {
		fmt.Print(FormatQuote(quote))
	}

	util.PressAnyKey()
}

func FormatQuote(quote Quote) string {

	var (
		quoteBuffer    bytes.Buffer
		formattedQuote string
	)

	stringFormat := `` + "\n" + color.Cyan + `%d. ` + "\"" + color.Reset + `%s` + `` + color.Cyan + "\"" +
		"\n" + strings.Repeat(" ", 50) + `By %s (%s %s)` + color.Reset + "\n" + ``

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

/* Stats */

func (q *Quotes) PrintStatistics() {

	util.ClearScreen()

	format := "%s%s%s"

	name := color.Cyan + "Statistics: " + color.Reset

	stats := fmt.Sprintf(format, name, q.TopAuthors(), q.TopLanguages())

	fmt.Println(stats)

	util.PressAnyKey()
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
		authorsString += "\n" + strconv.Itoa(pairs[i].count) + " " + color.Cyan + pairs[i].name + color.Reset
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
		languagesString += "\n" + strconv.Itoa(num) + " " + color.Yellow + name + color.Reset
	}
	return languagesString
}

/* Find */

func (q *Quotes) SetToDefaultQuote() {

	index := len(q.QUOTES) - 1

	if index > 0 {
		config.MainQuoteID = q.QUOTES[index].ID
	}

}

func (q *Quotes) Find() bool {
	fmt.Print("Find: ")

	// Read user input
	reader := bufio.NewReader(os.Stdin)
	searchQuote, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		util.PressAnyKey()
		return false
	}

	// Clean and process the input
	searchQuote = strings.TrimSpace(searchQuote)
	searchQuote = util.EnsureSentenceEnd(searchQuote)

	// Search for the quote
	foundQuote, found := q.FindQuoteByQuote(searchQuote)
	if !found {
		config.AddMessage(color.Red + "Quote not found." + color.Reset)
		return false
	}

	// Print the found quote
	config.MainQuoteID = foundQuote.ID
	return true
}

func (q *Quotes) FindQuoteByQuote(searchQuote string) (Quote, bool) {
	for _, quote := range q.QUOTES {
		if strings.EqualFold(quote.QUOTE, searchQuote) {
			return quote, true
		}
	}
	return Quote{}, false
}

func (q *Quotes) FindQuoteByID(id int) (Quote, bool) {
	for _, quote := range q.QUOTES {
		if quote.ID == id {
			return quote, true // found
		}
	}
	return Quote{}, false // not found
}

func (q *Quotes) FindDuplicates(searchQuote string, excludeID int) bool {

	if searchQuote == "" || searchQuote == "Unknown" {
		return false
	}

	for _, quote := range q.QUOTES {
		if searchQuote == quote.QUOTE {
			if excludeID != quote.ID {
				return true
			}
		}
	}

	return false
}

func (q *Quotes) PrintRandomQuotes(amount int) {

	util.ClearScreen()

	if amount <= 0 {
		amount = 10
	}

	var DefaultQuote = Quote{
		ID:       0,
		QUOTE:    "Nothing Smart to Say",
		AUTHOR:   "Unknown",
		LANGUAGE: "Golang",
		DATE:     "00.00.00",
	}

	if len(q.QUOTES) == 0 {
		fmt.Print(FormatQuote(DefaultQuote))
		util.PressAnyKey()
		return
	}

	for i := 0; i < amount; i++ {
		fmt.Printf("Number %d\n", i+1)
		idx := rand.Intn(len(q.QUOTES))
		quote := q.QUOTES[idx]
		fmt.Print(FormatQuote(quote))
	}

	util.PressAnyKey()
}

/* Find Similar Quotes */

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

		firstQuote, _ := quotes.FindQuoteByQuote(pair[0])
		secondQuote, _ := quotes.FindQuoteByQuote(pair[1])

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

	message := color.Green + "Find Similar Quotes Process Done" + color.Reset
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

/* Other */

func (q *Quotes) save() error {

	copySlice := make([]Quote, len(q.QUOTES))
	copy(copySlice, q.QUOTES)
	copyQuotes := Quotes{QUOTES: copySlice}

	byteValue, err := json.MarshalIndent(copyQuotes, "", "  ")
	if err != nil {
		return err
	}

	// Save Local
	if err := os.WriteFile(config.LocalFile, byteValue, 0644); err != nil {
		return err
	}

	fmt.Println(color.Green + "Local save!" + color.Reset)

	var date = time.Now().Format("02.01.2006")
	backupWithDate := filepath.Join(config.BaseBackup, config.BaseLocal, config.BaseDB, strconv.Itoa(len(q.QUOTES))+". quotes "+date+".json")

	// Save Backup
	if err := util.InitBackupStorage(); err != nil {
		fmt.Println(color.Yellow + "Backup init failed: " + err.Error() + color.Reset)
		return nil // or return err, depending on your needs
	}

	if err := os.WriteFile(config.BackupFile, byteValue, 0644); err != nil {
		fmt.Println(color.Yellow + "Backup write failed: " + err.Error() + color.Reset)
		return nil // same decision here
	}

	if err := os.WriteFile(backupWithDate, byteValue, 0644); err != nil {
		fmt.Println(color.Yellow + "Backup write failed: " + err.Error() + color.Reset)
		return nil // same decision here
	}

	fmt.Println(color.Green + "Backup save!" + color.Reset)

	return nil
}

func (q *Quotes) promptEntry(suggestions Quote) (Quote, error) {
	prompts := []struct {
		label  string
		target *string
	}{
		{"Quote: ", &suggestions.QUOTE},
		{"Author: ", &suggestions.AUTHOR},
	}

	suggestions.ID = q.newID()

	for _, p := range prompts {
		input, err := util.PromptWithSuggestion(p.label, *p.target)
		if err != nil {
			return Quote{}, err
		}

		if input == "q" {
			return Quote{}, errors.New("Aborted")
		}

		if p.label == "Quote: " {
			if q.FindDuplicates(input, suggestions.ID) {
				return Quote{}, errors.New("Dublicates found")
			}
		}

		*p.target = input
	}

	if suggestions.QUOTE == "" {
		suggestions.QUOTE = "Unknown"
	}

	if suggestions.AUTHOR == "" {
		suggestions.AUTHOR = "Unknown"
	}

	suggestions.QUOTE = util.CapitalizeFirstLetter(suggestions.QUOTE)
	suggestions.QUOTE = util.EnsureSentenceEnd(suggestions.QUOTE)
	suggestions.LANGUAGE = util.AutoDetectLanguage(suggestions.QUOTE)
	suggestions.DATE = time.Now().Format("15:04 (02.01.2006)")

	return suggestions, nil
}

func (q *Quotes) indexOf(id int) (int, error) {
	if id <= 0 {
		return -1, fmt.Errorf("invalid ID: %d", id)
	}
	for i, w := range q.QUOTES {
		if w.ID == id {
			return i, nil
		}
	}
	return -1, fmt.Errorf("item with ID %d not found", id)
}

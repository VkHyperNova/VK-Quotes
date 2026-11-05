package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/peterh/liner"
)

type Quote struct {
	ID               int    `json:"id"`
	QUOTE            string `json:"quote"`
	AUTHOR           string `json:"author"`
	LANGUAGE         string `json:"language"`
	DATE             string `json:"date"`
	WRITTEN_ON_PAPER bool   `json:"written_on_paper"`
}

type Quotes struct {
	QUOTES []Quote `json:"quotes"` // Slice containing multiple Quote instances.
}

/* Main */

func (q *Quotes) Add() error {

	go func() {
		if err := PlayMP3(); err != nil {
			log.Printf("MP3 Playback Error: %v", err)
		}
	}()

	newQuote, err := q.promptEntry(Quote{}, true)
	if err != nil {
		return err
	}

	if matches := q.findSimilar(newQuote.QUOTE, dupThreshold); len(matches) > 0 {
		fmt.Printf("\n%sSimilar quote(s) already exist:%s\n", Gray, Reset)
		for _, m := range matches {
			label := fmt.Sprintf("%.0f%%", m.score*100)
			if m.score >= 1 {
				label = "exact"
			}
			fmt.Printf("  %s#%d (%s)%s %s%s%s - %s\n",
				Gray, m.quote.ID, label, Reset,
				Yellow, m.quote.QUOTE, Reset, m.quote.AUTHOR)
		}
		fmt.Println()

		answer, err := prompt("Add anyway? (y/N):", "")
		if err != nil {
			return err
		}
		if a := strings.ToLower(answer); a != "y" && a != "yes" {
			fmt.Println(Gray + "Not added." + Reset)
			return nil
		}
	}

	q.QUOTES = append(q.QUOTES, newQuote)

	return q.Save()
}

func (q *Quotes) Update() error {

	if len(os.Args) < 3 {
		return fmt.Errorf("usage: update <id>")
	}

	id, err := parseID(os.Args[2])
	if err != nil {
		return err
	}

	index, err := q.indexOf(id)
	if err != nil {
		return err
	}

	updated, err := q.promptEntry(q.QUOTES[index], false)
	if err != nil {
		return err
	}

	q.QUOTES[index] = updated

	return q.Save()
}

func (q *Quotes) Delete() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: delete <id>")
	}

	id, err := parseID(os.Args[2])
	if err != nil {
		return err
	}

	_, err = q.deleteByID(id)
	return err
}

func (q *Quotes) Search() error {

	query := strings.TrimSpace(strings.Join(os.Args[2:], " "))

	if query == "" {
		var err error
		query, err = promptRequired("Search:", "")
		if err != nil {
			return err
		}
	}

	needle := strings.ToLower(query)
	found := 0

	fmt.Println()
	for _, quote := range q.QUOTES {
		if !quote.matches(needle) {
			continue
		}
		found++
		fmt.Printf("%s#%d%s %s%s%s - %s %s(%s, %s)%s\n",
			Gray, quote.ID, Reset,
			Yellow, quote.QUOTE, Reset,
			quote.AUTHOR,
			Gray, quote.LANGUAGE, quote.DATE, Reset)
	}

	if found == 0 {
		fmt.Println(Gray + "No quotes found." + Reset)
		return nil
	}

	fmt.Printf("\n%s%d found.%s\n", Gray, found, Reset)
	return nil
}

func (q *Quotes) Read() error {

	if len(q.QUOTES) == 0 {
		fmt.Println(Gray + "No quotes yet." + Reset)
		return nil
	}

	onlyUnwritten := false
	startIdx := 0

	// Optional args: read [id] [unwritten]
	if len(os.Args) > 2 {
		for _, arg := range os.Args[2:] {
			switch strings.ToLower(arg) {
			case "u", "unwritten":
				onlyUnwritten = true
			default:
				id, err := parseID(arg)
				if err != nil {
					return fmt.Errorf("usage: read [id] [unwritten]")
				}
				if startIdx, err = q.indexOf(id); err != nil {
					return err
				}
			}
		}
	}

	// view holds indexes into q.QUOTES that are currently being read.
	buildView := func() []int {
		var v []int
		for idx, quote := range q.QUOTES {
			if !onlyUnwritten || !quote.WRITTEN_ON_PAPER {
				v = append(v, idx)
			}
		}
		return v
	}

	// firstAtOrAfter returns the position in view of the first quote whose
	// index is >= idx (or the last position if there is none).
	firstAtOrAfter := func(view []int, idx int) int {
		for pos, v := range view {
			if v >= idx {
				return pos
			}
		}
		return len(view) - 1
	}

	view := buildView()
	if len(view) == 0 {
		fmt.Println(Gray + "All quotes are already written on paper." + Reset)
		return nil
	}
	i := firstAtOrAfter(view, startIdx)

	msg := "" // one-shot message, shown after the screen is redrawn
	for {
		clearScreen()
		quote := q.QUOTES[view[i]]

		mode := ""
		if onlyUnwritten {
			mode = " - unwritten only"
		}
		fmt.Printf("\n%s#%d (%d/%d%s)%s\n", Gray, quote.ID, i+1, len(view), mode, Reset)
		fmt.Printf("%s%s%s\n", Yellow, quote.QUOTE, Reset)
		fmt.Printf("- %s %s(%s, %s)%s\n", quote.AUTHOR, Gray, quote.LANGUAGE, quote.DATE, Reset)
		if quote.WRITTEN_ON_PAPER {
			fmt.Printf("%s[written on paper]%s\n", Cyan, Reset)
		}
		fmt.Println()
		if msg != "" {
			fmt.Println(msg)
			msg = ""
		}

		answer, err := prompt("[Enter] next, [p] prev, [r] random, [w] written on paper, [f] filter, [q] quit:", "")
		if errors.Is(err, errCancelled) {
			return nil // Ctrl+C / Ctrl+D is a normal way to stop reading
		}
		if err != nil {
			return err
		}

		switch strings.ToLower(answer) {
		case "q", "quit", "exit":
			return nil

		case "p", "prev":
			if i > 0 {
				i--
			} else {
				msg = Gray + "   This is the first quote." + Reset
			}

		case "r", "random":
			i = rand.Intn(len(view))

		case "w":
			// Toggle written_on_paper and save.
			idx := view[i]
			q.QUOTES[idx].WRITTEN_ON_PAPER = !q.QUOTES[idx].WRITTEN_ON_PAPER
			if err := q.Save(); err != nil {
				return err
			}

			if onlyUnwritten {
				// The quote just left the list, so rebuild it.
				view = buildView()
				if len(view) == 0 {
					fmt.Println(Gray + "All quotes are written on paper." + Reset)
					return nil
				}
				if i >= len(view) {
					i = len(view) - 1
				}
				// otherwise i now points at the next unwritten quote
			}

		case "f", "filter":
			cur := view[i]
			onlyUnwritten = !onlyUnwritten

			newView := buildView()
			if len(newView) == 0 {
				onlyUnwritten = false // nothing to show, stay in the full list
				msg = Gray + "   All quotes are already written on paper." + Reset
				continue
			}
			view = newView
			i = firstAtOrAfter(view, cur)

		default: // Enter or anything else: next
			if i < len(view)-1 {
				i++
			} else {
				fmt.Println(Gray + "   That was the last quote." + Reset)
				return nil
			}
		}
	}
}

func (q *Quotes) Duplicates() error {

	threshold := dupThreshold

	// Optional threshold: duplicates 0.9 (or 90)
	if len(os.Args) >= 3 {
		t, err := strconv.ParseFloat(os.Args[2], 64)
		if err != nil || t <= 0 || t > 100 {
			return fmt.Errorf("usage: duplicates [threshold 0-1]")
		}
		if t > 1 {
			t /= 100
		}
		threshold = t
	}

	fmt.Println(Gray + "Scanning..." + Reset)

	norm := make([]string, len(q.QUOTES))
	for i, quote := range q.QUOTES {
		norm[i] = normalize(quote.QUOTE)
	}

	var pairs []dupPair
	for i := 0; i < len(norm); i++ {
		for j := i + 1; j < len(norm); j++ {
			if s := similarity(norm[i], norm[j], threshold); s >= threshold {
				pairs = append(pairs, dupPair{q.QUOTES[i], q.QUOTES[j], s})
			}
		}
	}

	if len(pairs) == 0 {
		fmt.Println(Gray + "No duplicates or similar quotes found." + Reset)
		return nil
	}

	sort.Slice(pairs, func(x, y int) bool { return pairs[x].score > pairs[y].score })

	// Show one pair per screen.
	msg := "" // one-shot message, shown after the screen is redrawn
	for i := 0; ; {
		p := pairs[i]

		label := fmt.Sprintf("%.0f%% similar", p.score*100)
		if p.score >= 1 {
			label = "exact duplicate"
		}

		clearScreen()
		fmt.Printf("%sPair %d/%d - %s%s\n", Cyan, i+1, len(pairs), label, Reset)
		fmt.Print(FormatQuote(p.a))
		fmt.Print(FormatQuote(p.b))
		fmt.Printf("\n%sUse 'd <id>' to delete one of them.%s\n", Gray, Reset)
		if msg != "" {
			fmt.Println(msg)
			msg = ""
		}
		fmt.Println()

		answer, err := prompt("[Enter] next, [p] prev, [d <id>] delete, [q] quit:", "")
		if errors.Is(err, errCancelled) {
			return nil // Ctrl+C / Ctrl+D is a normal way to stop
		}
		if err != nil {
			return err
		}

		fields := strings.Fields(strings.ToLower(answer))
		cmd := ""
		if len(fields) > 0 {
			cmd = fields[0]
		}

		switch cmd {
		case "q", "quit", "exit":
			return nil

		case "p", "prev":
			if i > 0 {
				i--
			}

		case "d", "delete":
			if len(fields) < 2 {
				msg = Yellow + "usage: d <id>" + Reset
				continue
			}
			id, err := parseID(fields[1])
			if err != nil {
				msg = Yellow + err.Error() + Reset
				continue
			}
			if id != p.a.ID && id != p.b.ID {
				msg = Yellow + "That ID is not in this pair." + Reset
				continue
			}

			deleted, err := q.deleteByID(id)
			if errors.Is(err, errCancelled) {
				return nil
			}
			if err != nil {
				msg = Yellow + err.Error() + Reset
				continue
			}
			if !deleted {
				continue
			}

			// Drop every pair that contained the deleted quote.
			kept := pairs[:0]
			for _, pr := range pairs {
				if pr.a.ID != id && pr.b.ID != id {
					kept = append(kept, pr)
				}
			}
			pairs = kept

			if len(pairs) == 0 {
				fmt.Println(Gray + "No more duplicates." + Reset)
				return nil
			}
			if i >= len(pairs) {
				i = len(pairs) - 1
			}
			// otherwise i now points at the next pair, which shifted into place

		default: // Enter or anything else: next
			if i < len(pairs)-1 {
				i++
			} else {
				fmt.Println(Gray + "   That was the last pair." + Reset)
				return nil
			}
		}
	}
}

func (q *Quotes) ResetIDs() error {
	// Reset IDs
	for key := range q.QUOTES {
		q.QUOTES[key].ID = key + 1
	}

	return q.Save()
}

/* Helper Functions */

func (q *Quotes) deleteByID(id int) (bool, error) {
	index, err := q.indexOf(id)
	if err != nil {
		return false, err
	}

	quote := q.QUOTES[index]
	fmt.Printf("\nDelete %s#%d%s %s%s%s - %s?\n",
		Gray, quote.ID, Reset, Yellow, quote.QUOTE, Reset, quote.AUTHOR)

	answer, err := prompt("(y/N):", "")
	if err != nil {
		return false, err
	}
	if a := strings.ToLower(answer); a != "y" && a != "yes" {
		fmt.Println(Gray + "Not deleted." + Reset)
		return false, nil
	}

	q.QUOTES = append(q.QUOTES[:index], q.QUOTES[index+1:]...)
	return true, q.Save()
}

func (q *Quotes) promptEntry(existing Quote, isNew bool) (Quote, error) {
	result := existing

	if isNew {
		result.ID = q.nextID()
	}

	var err error

	if result.QUOTE, err = promptRequired("Quote:", existing.QUOTE); err != nil {
		return Quote{}, err
	}
	if result.AUTHOR, err = promptRequired("Author:", existing.AUTHOR); err != nil {
		return Quote{}, err
	}
	if result.LANGUAGE, err = prompt("Language:", existing.LANGUAGE); err != nil {
		return Quote{}, err
	}

	defDate := existing.DATE
	if isNew && defDate == "" {
		defDate = time.Now().Format("02.01.2006")
	}
	if result.DATE, err = promptDate("Date (DD.MM.YYYY):", defDate); err != nil {
		return Quote{}, err
	}

	if result.WRITTEN_ON_PAPER, err = promptBool("Written on paper (y/n):", existing.WRITTEN_ON_PAPER); err != nil {
		return Quote{}, err
	}

	return result, nil
}

func promptRequired(label, suggestion string) (string, error) {
	for {
		v, err := prompt(label, suggestion)
		if err != nil {
			return "", err
		}
		if v != "" {
			return v, nil
		}
		fmt.Println(Gray + "   This field is required." + Reset)
	}
}

func promptDate(label, suggestion string) (string, error) {
	for {
		v, err := prompt(label, suggestion)
		if err != nil {
			return "", err
		}
		if v == "" {
			return "", nil
		}
		if _, perr := time.Parse("02.01.2006", v); perr != nil {
			fmt.Println(Gray + "   Invalid date, use DD.MM.YYYY." + Reset)
			suggestion = v // keep what they typed so they can fix it
			continue
		}
		return v, nil
	}
}

func promptBool(label string, def bool) (bool, error) {
	suggestion := "n"
	if def {
		suggestion = "y"
	}
	for {
		v, err := prompt(label, suggestion)
		if err != nil {
			return false, err
		}
		switch strings.ToLower(strings.TrimSpace(v)) {
		case "":
			return def, nil
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		}
		fmt.Println(Gray + "   Please answer y or n." + Reset)
		suggestion = v // keep what they typed so they can fix it
	}
}

func (q *Quotes) nextID() int {
	maxID := 0
	for _, quote := range q.QUOTES {
		if quote.ID > maxID {
			maxID = quote.ID
		}
	}
	return maxID + 1
}

func (q *Quotes) indexOf(id int) (int, error) {
	for i, quote := range q.QUOTES {
		if quote.ID == id {
			return i, nil
		}
	}
	return -1, fmt.Errorf("item with ID %d not found", id)
}

func parseID(s string) (int, error) {
	id, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid ID: %q", s)
	}
	return id, nil
}

func (qt Quote) matches(needle string) bool {
	fields := []string{qt.QUOTE, qt.AUTHOR, qt.LANGUAGE, qt.DATE}
	for _, f := range fields {
		if strings.Contains(strings.ToLower(f), needle) {
			return true
		}
	}
	return false
}

/* Prompt helpers */

var errCancelled = errors.New("cancelled")

var line *liner.State

var prompt = func(label, suggestion string) (string, error) {
	if line == nil {
		line = liner.NewLiner()
		line.SetCtrlCAborts(true)
	}
	input, err := line.PromptWithSuggestion("   "+label+" ", suggestion, -1)
	if errors.Is(err, liner.ErrPromptAborted) || errors.Is(err, io.EOF) {
		return "", errCancelled // Ctrl+C / Ctrl+D
	}
	return strings.TrimSpace(input), err
}

func closePrompt() {
	if line != nil {
		line.Close()
		line = nil
	}
}

/* Dublicates */

const dupThreshold = 0.85

type dupPair struct {
	a, b  Quote
	score float64
}

type similarQuote struct {
	quote Quote
	score float64
}

func normalize(s string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
		case unicode.IsSpace(r):
			b.WriteRune(' ')
		}
	}
	return strings.Join(strings.Fields(b.String()), " ")
}

func similarity(a, b string, threshold float64) float64 {
	if a == "" || b == "" {
		return 0
	}
	if a == b {
		return 1
	}

	ra, rb := []rune(a), []rune(b)
	if len(ra) > len(rb) {
		ra, rb = rb, ra
		a, b = b, a
	}

	// One quote is a longer version of the other (e.g. truncated copy).
	if len(ra) >= 20 && strings.Contains(b, a) {
		return 0.95
	}

	// Edit distance is at least the length difference, so the best possible
	// score is len(short)/len(long). Skip if that can't reach the threshold.
	if float64(len(ra))/float64(len(rb)) < threshold {
		return 0
	}

	maxDist := int((1 - threshold) * float64(len(rb)))
	d := levenshteinMax(ra, rb, maxDist)
	if d > maxDist {
		return 0
	}
	return 1 - float64(d)/float64(len(rb))
}

func levenshteinMax(a, b []rune, maxDist int) int {
	if len(b)-len(a) > maxDist {
		return maxDist + 1
	}

	prev := make([]int, len(b)+1)
	curr := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}

	for i := 1; i <= len(a); i++ {
		curr[0] = i
		rowMin := i

		lo := max(1, i-maxDist)
		hi := min(len(b), i+maxDist)
		if lo > 1 {
			curr[lo-1] = maxDist + 1
		}

		for j := lo; j <= hi; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			v := min(prev[j]+1, curr[j-1]+1, prev[j-1]+cost)
			curr[j] = v
			if v < rowMin {
				rowMin = v
			}
		}
		if hi < len(b) {
			curr[hi+1] = maxDist + 1
		}

		if rowMin > maxDist {
			return maxDist + 1
		}
		prev, curr = curr, prev
	}
	return prev[len(b)]
}

func (q *Quotes) findSimilar(text string, threshold float64) []similarQuote {
	n := normalize(text)

	var out []similarQuote
	for _, existing := range q.QUOTES {
		if s := similarity(n, normalize(existing.QUOTE), threshold); s >= threshold {
			out = append(out, similarQuote{existing, s})
		}
	}

	sort.Slice(out, func(i, j int) bool { return out[i].score > out[j].score })
	return out
}

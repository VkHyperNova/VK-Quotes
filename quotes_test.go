package main

import (
	"bytes"
	"errors"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

/* Test helpers */

// enter simulates pressing Enter on a pre-filled liner prompt: the
// suggestion is returned unchanged.
const enter = "\x00enter"

type promptCall struct{ label, suggestion string }

// fakePrompt replaces the package-level prompt with a scripted one. Each
// call consumes the next answer. When the answers run out it behaves like
// Ctrl+C and returns errCancelled. It returns the recorded calls.
func fakePrompt(t *testing.T, answers ...string) *[]promptCall {
	t.Helper()
	old := prompt
	calls := &[]promptCall{}
	prompt = func(label, suggestion string) (string, error) {
		n := len(*calls)
		*calls = append(*calls, promptCall{label, suggestion})
		if n >= len(answers) {
			return "", errCancelled
		}
		if answers[n] == enter {
			return suggestion, nil
		}
		return answers[n], nil
	}
	t.Cleanup(func() { prompt = old })
	return calls
}

// withArgs sets os.Args to ["quotes", args...] for the test.
func withArgs(t *testing.T, args ...string) {
	t.Helper()
	old := os.Args
	os.Args = append([]string{"quotes"}, args...)
	t.Cleanup(func() { os.Args = old })
}

// captureStdout returns everything fn prints to stdout.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	done := make(chan string)
	go func() {
		var b bytes.Buffer
		io.Copy(&b, r)
		done <- b.String()
	}()

	defer func() { os.Stdout = old }()
	fn()
	w.Close()
	return <-done
}

// setupTempStorage points the storage globals at a temp dir with a fake
// backup drive, and restores them afterwards. Returns the local file path.
func setupTempStorage(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	oldLocal, oldBackup, oldMount, oldNotices := LocalFile, BackupFile, mountBackup, notices
	LocalFile = filepath.Join(dir, "local", "quotes.json")
	BackupFile = filepath.Join(dir, "backup", "quotes.json")
	mountBackup = func() error { return nil }
	notices = nil

	t.Cleanup(func() {
		LocalFile, BackupFile, mountBackup, notices = oldLocal, oldBackup, oldMount, oldNotices
	})
	return LocalFile
}

func testQuotes() *Quotes {
	return &Quotes{QUOTES: []Quote{
		{ID: 1, QUOTE: "Be yourself; everyone else is already taken", AUTHOR: "Oscar Wilde", LANGUAGE: "English", DATE: "2026-01-02"},
		{ID: 2, QUOTE: "The only thing we have to fear is fear itself", AUTHOR: "F. D. Roosevelt", LANGUAGE: "English", DATE: "2026-02-03"},
		{ID: 3, QUOTE: "Tark nagu öö", AUTHOR: "Anonymous", LANGUAGE: "Estonian", DATE: ""},
	}}
}

func loadSaved(t *testing.T, path string) Quotes {
	t.Helper()
	var got Quotes
	if err := got.Load(path); err != nil {
		t.Fatalf("load saved file: %v", err)
	}
	return got
}

/* normalize */

func TestNormalize(t *testing.T) {
	tests := []struct{ in, want string }{
		{"Hello, World!", "hello world"},
		{"  Be   yourself.\n", "be yourself"},
		{"Don't", "dont"},
		{"Tark nagu ÖÖ", "tark nagu öö"},
		{"tab\tand\nnewline", "tab and newline"},
		{"Route 66", "route 66"},
		{"!!!", ""},
		{"", ""},
	}
	for _, tt := range tests {
		if got := normalize(tt.in); got != tt.want {
			t.Errorf("normalize(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

/* levenshteinMax */

func refLevenshtein(a, b []rune) int {
	prev := make([]int, len(b)+1)
	curr := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		curr[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			curr[j] = min(prev[j]+1, curr[j-1]+1, prev[j-1]+cost)
		}
		prev, curr = curr, prev
	}
	return prev[len(b)]
}

func TestLevenshteinMax_KnownValues(t *testing.T) {
	tests := []struct {
		a, b    string
		maxDist int
		want    int // exact distance if <= maxDist
	}{
		{"", "", 0, 0},
		{"abc", "abc", 0, 0},
		{"abc", "abd", 1, 1},
		{"kitten", "sitting", 3, 3},
		{"", "abc", 3, 3},
		{"tark öö", "tark oo", 2, 2},
	}
	for _, tt := range tests {
		got := levenshteinMax([]rune(tt.a), []rune(tt.b), tt.maxDist)
		if got != tt.want {
			t.Errorf("levenshteinMax(%q, %q, %d) = %d, want %d", tt.a, tt.b, tt.maxDist, got, tt.want)
		}
	}
}

func TestLevenshteinMax_ExceedingLimitReturnsMorethanMax(t *testing.T) {
	cases := []struct {
		a, b    string
		maxDist int
	}{
		{"kitten", "sitting", 2},
		{"abc", "xyz", 1},
		{"a", "abcdef", 2}, // length difference alone exceeds the limit
	}
	for _, c := range cases {
		got := levenshteinMax([]rune(c.a), []rune(c.b), c.maxDist)
		if got <= c.maxDist {
			t.Errorf("levenshteinMax(%q, %q, %d) = %d, want > %d", c.a, c.b, c.maxDist, got, c.maxDist)
		}
	}
}

func TestLevenshteinMax_MatchesReference(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	alphabet := []rune("abcö ")
	randRunes := func(n int) []rune {
		out := make([]rune, n)
		for i := range out {
			out[i] = alphabet[r.Intn(len(alphabet))]
		}
		return out
	}

	for i := 0; i < 3000; i++ {
		a, b := randRunes(r.Intn(16)), randRunes(r.Intn(16))
		if len(a) > len(b) { // the function requires a to be the shorter one
			a, b = b, a
		}
		maxDist := r.Intn(7)

		want := refLevenshtein(a, b)
		got := levenshteinMax(a, b, maxDist)

		if want <= maxDist {
			if got != want {
				t.Fatalf("a=%q b=%q max=%d: got %d, want %d", string(a), string(b), maxDist, got, want)
			}
		} else if got <= maxDist {
			t.Fatalf("a=%q b=%q max=%d: got %d but real distance %d exceeds the limit",
				string(a), string(b), maxDist, got, want)
		}
	}
}

/* similarity */

func TestSimilarity_EmptyIsZero(t *testing.T) {
	if s := similarity("", "abc", dupThreshold); s != 0 {
		t.Errorf("got %v, want 0", s)
	}
	if s := similarity("abc", "", dupThreshold); s != 0 {
		t.Errorf("got %v, want 0", s)
	}
	if s := similarity("", "", dupThreshold); s != 0 {
		t.Errorf("got %v, want 0", s)
	}
}

func TestSimilarity_Identical(t *testing.T) {
	if s := similarity("be yourself", "be yourself", dupThreshold); s != 1 {
		t.Errorf("got %v, want 1", s)
	}
}

func TestSimilarity_LongSubstringScores095(t *testing.T) {
	short := "the quick brown fox jumps"
	long := short + " over the lazy dog"
	if s := similarity(short, long, dupThreshold); s != 0.95 {
		t.Errorf("got %v, want 0.95", s)
	}
	if s := similarity(long, short, dupThreshold); s != 0.95 {
		t.Errorf("order must not matter, got %v", s)
	}
}

func TestSimilarity_ShortSubstringIsNotADuplicate(t *testing.T) {
	// Under 20 characters the substring shortcut does not apply, and the
	// length ratio rules the pair out.
	if s := similarity("be yourself", "be yourself always and forever", dupThreshold); s != 0 {
		t.Errorf("got %v, want 0", s)
	}
}

func TestSimilarity_OneTypo(t *testing.T) {
	a := "the only thing we have to fear is fear itself"
	b := "the only thing we have to fear is fear itselg"
	s := similarity(a, b, dupThreshold)
	if s < dupThreshold || s >= 1 {
		t.Errorf("got %v, want between %v and 1", s, dupThreshold)
	}
}

func TestSimilarity_Unrelated(t *testing.T) {
	if s := similarity("completely different text here", "nothing alike at all, honestly", dupThreshold); s != 0 {
		t.Errorf("got %v, want 0", s)
	}
}

func TestSimilarity_Symmetric(t *testing.T) {
	a := "the only thing we have to fear is fear itself"
	b := "the only thing we have to fear is fear itselg"
	if similarity(a, b, dupThreshold) != similarity(b, a, dupThreshold) {
		t.Error("similarity should not depend on argument order")
	}
}

func TestSimilarity_UnicodeCountsRunesNotBytes(t *testing.T) {
	a := "tark nagu ööbik laulab ilusasti"
	b := "tark nagu oobik laulab ilusasti"
	if s := similarity(a, b, dupThreshold); s < dupThreshold {
		t.Errorf("got %v, want >= %v (two substitutions in 31 characters)", s, dupThreshold)
	}
}

func TestSimilarity_HigherThresholdRejectsMore(t *testing.T) {
	a := "the only thing we have to fear is fear itself"
	b := "the only thing we have to fear is fear itselg"
	if s := similarity(a, b, 0.99); s != 0 {
		t.Errorf("one typo in 45 characters should not pass 0.99, got %v", s)
	}
}

/* findSimilar */

func TestFindSimilar_ExactAfterNormalizing(t *testing.T) {
	q := testQuotes()
	got := q.findSimilar("BE YOURSELF, everyone else is already taken!", dupThreshold)
	if len(got) != 1 {
		t.Fatalf("expected 1 match, got %d", len(got))
	}
	if got[0].quote.ID != 1 || got[0].score != 1 {
		t.Errorf("got ID %d score %v, want ID 1 score 1", got[0].quote.ID, got[0].score)
	}
}

func TestFindSimilar_NoMatch(t *testing.T) {
	q := testQuotes()
	if got := q.findSimilar("Something entirely different and new", dupThreshold); len(got) != 0 {
		t.Errorf("expected no matches, got %d", len(got))
	}
}

func TestFindSimilar_SortedBestFirst(t *testing.T) {
	q := testQuotes()
	q.QUOTES = append(q.QUOTES, Quote{ID: 4, QUOTE: "Be yourself; everyone else is already keen", AUTHOR: "Typo"})

	got := q.findSimilar("Be yourself, everyone else is already taken", dupThreshold)
	if len(got) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(got))
	}
	if got[0].quote.ID != 1 || got[1].quote.ID != 4 {
		t.Errorf("order = #%d, #%d; want #1, #4", got[0].quote.ID, got[1].quote.ID)
	}
	if got[0].score < got[1].score {
		t.Error("matches should be sorted by score, highest first")
	}
}

func TestFindSimilar_EmptyCollection(t *testing.T) {
	q := &Quotes{}
	if got := q.findSimilar("anything", dupThreshold); len(got) != 0 {
		t.Errorf("expected no matches, got %d", len(got))
	}
}

/* Small helpers */

func TestNextID(t *testing.T) {
	if got := (&Quotes{}).nextID(); got != 1 {
		t.Errorf("empty: got %d, want 1", got)
	}
	q := &Quotes{QUOTES: []Quote{{ID: 1}, {ID: 2}, {ID: 5}}}
	if got := q.nextID(); got != 6 {
		t.Errorf("with gap: got %d, want 6", got)
	}
}

func TestIndexOf(t *testing.T) {
	q := testQuotes()

	i, err := q.indexOf(2)
	if err != nil || i != 1 {
		t.Errorf("indexOf(2) = %d, %v; want 1, nil", i, err)
	}

	if _, err := q.indexOf(99); err == nil {
		t.Error("expected error for unknown ID")
	}
}

func TestParseID(t *testing.T) {
	valid := map[string]int{"5": 5, " 7 ": 7, "123": 123}
	for in, want := range valid {
		got, err := parseID(in)
		if err != nil || got != want {
			t.Errorf("parseID(%q) = %d, %v; want %d", in, got, err, want)
		}
	}
	for _, in := range []string{"", "0", "-1", "abc", "1.5", "1 2"} {
		if _, err := parseID(in); err == nil {
			t.Errorf("parseID(%q) should fail", in)
		}
	}
}

func TestQuoteMatches(t *testing.T) {
	qt := Quote{QUOTE: "Be Yourself", AUTHOR: "Oscar Wilde", LANGUAGE: "English", DATE: "2026-01-02"}

	// matches expects an already lowercased needle
	for _, needle := range []string{"yourself", "wilde", "english", "2026-01", "oscar wilde"} {
		if !qt.matches(needle) {
			t.Errorf("matches(%q) = false, want true", needle)
		}
	}
	for _, needle := range []string{"zzz", "2027", "french"} {
		if qt.matches(needle) {
			t.Errorf("matches(%q) = true, want false", needle)
		}
	}
}

/* Prompt helpers */

func TestPromptRequired_RepromptsUntilNotEmpty(t *testing.T) {
	calls := fakePrompt(t, "", "", "hello")

	var got string
	var err error
	captureStdout(t, func() { got, err = promptRequired("Quote:", "") })

	if err != nil || got != "hello" {
		t.Errorf("got %q, %v; want %q, nil", got, err, "hello")
	}
	if len(*calls) != 3 {
		t.Errorf("prompt called %d times, want 3", len(*calls))
	}
}

func TestPromptRequired_CancelPropagates(t *testing.T) {
	fakePrompt(t) // no answers: behaves like Ctrl+C

	_, err := promptRequired("Quote:", "")
	if !errors.Is(err, errCancelled) {
		t.Errorf("got %v, want errCancelled", err)
	}
}

func TestPromptRequired_KeepsSuggestionOnEnter(t *testing.T) {
	fakePrompt(t, enter)

	got, err := promptRequired("Author:", "Oscar Wilde")
	if err != nil || got != "Oscar Wilde" {
		t.Errorf("got %q, %v; want the suggestion back", got, err)
	}
}

func TestPromptDate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		fakePrompt(t, "2026-03-04")
		got, err := promptDate("Date:", "")
		if err != nil || got != "2026-03-04" {
			t.Errorf("got %q, %v", got, err)
		}
	})

	t.Run("empty is allowed", func(t *testing.T) {
		fakePrompt(t, "")
		got, err := promptDate("Date:", "2026-01-01")
		if err != nil || got != "" {
			t.Errorf("got %q, %v; want empty", got, err)
		}
	})

	t.Run("invalid reprompts with what was typed", func(t *testing.T) {
		calls := fakePrompt(t, "bad", "2026-02-30", "2026-02-03")

		var got string
		var err error
		captureStdout(t, func() { got, err = promptDate("Date:", "2026-01-01") })

		if err != nil || got != "2026-02-03" {
			t.Fatalf("got %q, %v; want 2026-02-03", got, err)
		}
		if len(*calls) != 3 {
			t.Fatalf("prompt called %d times, want 3", len(*calls))
		}
		if (*calls)[1].suggestion != "bad" || (*calls)[2].suggestion != "2026-02-30" {
			t.Errorf("retries should pre-fill the previous input, got %q and %q",
				(*calls)[1].suggestion, (*calls)[2].suggestion)
		}
	})

	t.Run("cancel", func(t *testing.T) {
		fakePrompt(t)
		if _, err := promptDate("Date:", ""); !errors.Is(err, errCancelled) {
			t.Errorf("got %v, want errCancelled", err)
		}
	})
}

func TestPromptEntry_New(t *testing.T) {
	q := testQuotes()
	calls := fakePrompt(t, "A new quote", "Someone", "English", enter)

	got, err := q.promptEntry(Quote{}, true)
	if err != nil {
		t.Fatal(err)
	}

	today := time.Now().Format("2006-01-02")
	want := Quote{ID: 4, QUOTE: "A new quote", AUTHOR: "Someone", LANGUAGE: "English", DATE: today}
	if got != want {
		t.Errorf("got %+v, want %+v", got, want)
	}
	if len(*calls) != 4 || (*calls)[3].suggestion != today {
		t.Errorf("date prompt should be pre-filled with today, calls: %+v", *calls)
	}
}

func TestPromptEntry_UpdateKeepsExistingOnEnter(t *testing.T) {
	q := testQuotes()
	existing := q.QUOTES[0]
	fakePrompt(t, enter, enter, enter, enter)

	got, err := q.promptEntry(existing, false)
	if err != nil {
		t.Fatal(err)
	}
	if got != existing {
		t.Errorf("got %+v, want unchanged %+v", got, existing)
	}
}

func TestPromptEntry_UpdateCanClearDate(t *testing.T) {
	q := testQuotes()
	fakePrompt(t, enter, enter, enter, "")

	got, err := q.promptEntry(q.QUOTES[0], false)
	if err != nil {
		t.Fatal(err)
	}
	if got.DATE != "" {
		t.Errorf("DATE = %q, want empty", got.DATE)
	}
	if got.ID != 1 {
		t.Errorf("ID must not change on update, got %d", got.ID)
	}
}

func TestPromptEntry_UpdateDoesNotAssignNewID(t *testing.T) {
	q := testQuotes()
	fakePrompt(t, enter, enter, enter, enter)

	got, err := q.promptEntry(q.QUOTES[1], false)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != 2 {
		t.Errorf("ID = %d, want 2", got.ID)
	}
}

func TestPromptEntry_CancelMidway(t *testing.T) {
	q := testQuotes()
	fakePrompt(t, "Only a quote") // author prompt gets Ctrl+C

	got, err := q.promptEntry(Quote{}, true)
	if !errors.Is(err, errCancelled) {
		t.Errorf("got %v, want errCancelled", err)
	}
	if got != (Quote{}) {
		t.Errorf("expected zero Quote on error, got %+v", got)
	}
}

/* Add */

func TestAdd_NewQuoteIsSaved(t *testing.T) {
	path := setupTempStorage(t)
	q := testQuotes()
	fakePrompt(t, "A brand new thought nobody wrote before", "Me", "English", "2026-05-06")

	var err error
	captureStdout(t, func() { err = q.Add() })
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	if len(q.QUOTES) != 4 {
		t.Fatalf("expected 4 quotes, got %d", len(q.QUOTES))
	}
	if last := q.QUOTES[3]; last.ID != 4 || last.AUTHOR != "Me" || last.DATE != "2026-05-06" {
		t.Errorf("unexpected new quote: %+v", last)
	}

	saved := loadSaved(t, path)
	if len(saved.QUOTES) != 4 {
		t.Errorf("saved file has %d quotes, want 4", len(saved.QUOTES))
	}
}

func TestAdd_DuplicateDeclined(t *testing.T) {
	path := setupTempStorage(t)
	q := testQuotes()
	calls := fakePrompt(t,
		"Be yourself, everyone else is already taken.", "Someone", "English", enter,
		"n")

	var err error
	out := captureStdout(t, func() { err = q.Add() })
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	if len(q.QUOTES) != 3 {
		t.Errorf("duplicate should not be added, have %d quotes", len(q.QUOTES))
	}
	if !strings.Contains(out, "Similar quote(s) already exist") || !strings.Contains(out, "Not added.") {
		t.Errorf("unexpected output: %q", out)
	}
	if last := (*calls)[len(*calls)-1]; !strings.Contains(last.label, "Add anyway") {
		t.Errorf("last prompt = %q, want the 'Add anyway' question", last.label)
	}
	if _, statErr := os.Stat(path); !errors.Is(statErr, os.ErrNotExist) {
		t.Error("nothing should be saved when the add is declined")
	}
}

func TestAdd_DuplicateDefaultAnswerIsNo(t *testing.T) {
	setupTempStorage(t)
	q := testQuotes()
	fakePrompt(t,
		"Be yourself, everyone else is already taken.", "Someone", "English", enter,
		"") // just Enter

	captureStdout(t, func() { q.Add() })

	if len(q.QUOTES) != 3 {
		t.Errorf("pressing Enter must not add the duplicate, have %d quotes", len(q.QUOTES))
	}
}

func TestAdd_DuplicateConfirmed(t *testing.T) {
	path := setupTempStorage(t)
	q := testQuotes()
	fakePrompt(t,
		"Be yourself, everyone else is already taken.", "Someone Else", "English", enter,
		"y")

	var err error
	captureStdout(t, func() { err = q.Add() })
	if err != nil {
		t.Fatalf("Add: %v", err)
	}

	if len(q.QUOTES) != 4 || q.QUOTES[3].ID != 4 {
		t.Fatalf("duplicate should be added after 'y', quotes: %+v", q.QUOTES)
	}
	if saved := loadSaved(t, path); len(saved.QUOTES) != 4 {
		t.Errorf("saved file has %d quotes, want 4", len(saved.QUOTES))
	}
}

func TestAdd_CancelSavesNothing(t *testing.T) {
	path := setupTempStorage(t)
	q := testQuotes()
	fakePrompt(t) // Ctrl+C immediately

	err := q.Add()
	if !errors.Is(err, errCancelled) {
		t.Errorf("got %v, want errCancelled", err)
	}
	if len(q.QUOTES) != 3 {
		t.Errorf("quote count changed to %d", len(q.QUOTES))
	}
	if _, statErr := os.Stat(path); !errors.Is(statErr, os.ErrNotExist) {
		t.Error("nothing should be saved on cancel")
	}
}

/* Update */

func TestUpdate_ChangesOnlyEditedFieldAndSaves(t *testing.T) {
	path := setupTempStorage(t)
	q := testQuotes()
	withArgs(t, "update", "1")
	fakePrompt(t, enter, "Oscar W.", enter, enter)

	if err := q.Update(); err != nil {
		t.Fatalf("Update: %v", err)
	}

	want := Quote{ID: 1, QUOTE: "Be yourself; everyone else is already taken", AUTHOR: "Oscar W.", LANGUAGE: "English", DATE: "2026-01-02"}
	if q.QUOTES[0] != want {
		t.Errorf("got %+v, want %+v", q.QUOTES[0], want)
	}
	if saved := loadSaved(t, path); saved.QUOTES[0] != want {
		t.Errorf("saved %+v, want %+v", saved.QUOTES[0], want)
	}
}

func TestUpdate_Errors(t *testing.T) {
	setupTempStorage(t)

	t.Run("missing id", func(t *testing.T) {
		withArgs(t, "update")
		if err := testQuotes().Update(); err == nil || !strings.Contains(err.Error(), "usage") {
			t.Errorf("got %v, want usage error", err)
		}
	})
	t.Run("invalid id", func(t *testing.T) {
		withArgs(t, "update", "abc")
		if err := testQuotes().Update(); err == nil {
			t.Error("expected error")
		}
	})
	t.Run("unknown id", func(t *testing.T) {
		withArgs(t, "update", "99")
		if err := testQuotes().Update(); err == nil {
			t.Error("expected error")
		}
	})
}

func TestUpdate_CancelKeepsOriginal(t *testing.T) {
	setupTempStorage(t)
	q := testQuotes()
	orig := q.QUOTES[0]
	withArgs(t, "update", "1")
	fakePrompt(t, "Changed text") // cancelled at the author prompt

	if err := q.Update(); !errors.Is(err, errCancelled) {
		t.Errorf("got %v, want errCancelled", err)
	}
	if q.QUOTES[0] != orig {
		t.Errorf("quote was modified despite cancel: %+v", q.QUOTES[0])
	}
}

/* Delete */

func TestDelete_ConfirmedRemovesAndRenumbers(t *testing.T) {
	path := setupTempStorage(t)
	q := testQuotes()
	withArgs(t, "delete", "2")
	fakePrompt(t, "y")

	var err error
	captureStdout(t, func() { err = q.Delete() })
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}

	if len(q.QUOTES) != 2 {
		t.Fatalf("expected 2 quotes, got %d", len(q.QUOTES))
	}
	if q.QUOTES[0].ID != 1 || q.QUOTES[1].ID != 2 {
		t.Errorf("IDs should be renumbered 1,2 got %d,%d", q.QUOTES[0].ID, q.QUOTES[1].ID)
	}
	if q.QUOTES[1].QUOTE != "Tark nagu öö" {
		t.Errorf("wrong quote kept: %+v", q.QUOTES[1])
	}
	if saved := loadSaved(t, path); len(saved.QUOTES) != 2 {
		t.Errorf("saved file has %d quotes, want 2", len(saved.QUOTES))
	}
}

func TestDelete_AcceptsYes(t *testing.T) {
	setupTempStorage(t)
	q := testQuotes()
	withArgs(t, "delete", "1")
	fakePrompt(t, "YES")

	captureStdout(t, func() { q.Delete() })
	if len(q.QUOTES) != 2 {
		t.Errorf("'YES' should confirm, have %d quotes", len(q.QUOTES))
	}
}

func TestDelete_DeclinedKeepsEverything(t *testing.T) {
	path := setupTempStorage(t)

	for _, answer := range []string{"n", "", "no", "maybe"} {
		q := testQuotes()
		withArgs(t, "delete", "2")
		fakePrompt(t, answer)

		var err error
		out := captureStdout(t, func() { err = q.Delete() })
		if err != nil {
			t.Fatalf("answer %q: %v", answer, err)
		}
		if len(q.QUOTES) != 3 {
			t.Errorf("answer %q: quote was deleted", answer)
		}
		if !strings.Contains(out, "Not deleted.") {
			t.Errorf("answer %q: missing 'Not deleted.' in %q", answer, out)
		}
	}

	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		t.Error("nothing should be saved when the delete is declined")
	}
}

func TestDelete_Errors(t *testing.T) {
	setupTempStorage(t)

	t.Run("missing id", func(t *testing.T) {
		withArgs(t, "delete")
		if err := testQuotes().Delete(); err == nil || !strings.Contains(err.Error(), "usage") {
			t.Errorf("got %v, want usage error", err)
		}
	})
	t.Run("unknown id", func(t *testing.T) {
		withArgs(t, "delete", "99")
		if err := testQuotes().Delete(); err == nil {
			t.Error("expected error")
		}
	})
	t.Run("cancel at confirmation", func(t *testing.T) {
		withArgs(t, "delete", "1")
		fakePrompt(t)
		q := testQuotes()
		var err error
		captureStdout(t, func() { err = q.Delete() })
		if !errors.Is(err, errCancelled) {
			t.Errorf("got %v, want errCancelled", err)
		}
		if len(q.QUOTES) != 3 {
			t.Error("quote deleted despite cancel")
		}
	})
}

/* Search */

func TestSearch_FindsByAuthorCaseInsensitive(t *testing.T) {
	q := testQuotes()
	withArgs(t, "search", "WILDE")

	var err error
	out := captureStdout(t, func() { err = q.Search() })
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "Be yourself") {
		t.Errorf("expected the Wilde quote in output: %q", out)
	}
	if strings.Contains(out, "fear itself") {
		t.Errorf("unrelated quote should not appear: %q", out)
	}
	if !strings.Contains(out, "1 found.") {
		t.Errorf("expected '1 found.' in output: %q", out)
	}
}

func TestSearch_MultiWordQueryWithoutQuotes(t *testing.T) {
	q := testQuotes()
	withArgs(t, "search", "fear", "itself")

	out := captureStdout(t, func() { q.Search() })
	if !strings.Contains(out, "fear is fear itself") || !strings.Contains(out, "1 found.") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestSearch_MatchesDate(t *testing.T) {
	q := testQuotes()
	withArgs(t, "search", "2026-0")

	out := captureStdout(t, func() { q.Search() })
	if !strings.Contains(out, "2 found.") {
		t.Errorf("expected 2 matches by date prefix: %q", out)
	}
}

func TestSearch_NoResults(t *testing.T) {
	q := testQuotes()
	withArgs(t, "search", "zzzzz")

	var err error
	out := captureStdout(t, func() { err = q.Search() })
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "No quotes found.") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestSearch_PromptsWhenNoQuery(t *testing.T) {
	q := testQuotes()
	withArgs(t, "search")
	calls := fakePrompt(t, "roosevelt")

	out := captureStdout(t, func() { q.Search() })
	if len(*calls) != 1 {
		t.Errorf("expected one prompt, got %d", len(*calls))
	}
	if !strings.Contains(out, "fear itself") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestSearch_CancelAtPrompt(t *testing.T) {
	q := testQuotes()
	withArgs(t, "search")
	fakePrompt(t)

	if err := q.Search(); !errors.Is(err, errCancelled) {
		t.Errorf("got %v, want errCancelled", err)
	}
}

/* Read */

func TestRead_EmptyCollection(t *testing.T) {
	q := &Quotes{}
	withArgs(t, "read")

	var err error
	out := captureStdout(t, func() { err = q.Read() })
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "No quotes yet.") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestRead_NextPrevQuit(t *testing.T) {
	q := testQuotes()
	withArgs(t, "read")
	fakePrompt(t, "", "p", "q")

	var err error
	out := captureStdout(t, func() { err = q.Read() })
	if err != nil {
		t.Fatal(err)
	}

	if got := strings.Count(out, "(1/3)"); got != 2 {
		t.Errorf("quote 1 should be shown twice (start and after prev), got %d", got)
	}
	if got := strings.Count(out, "(2/3)"); got != 1 {
		t.Errorf("quote 2 should be shown once, got %d", got)
	}
}

func TestRead_PrevAtStartStaysPut(t *testing.T) {
	q := testQuotes()
	withArgs(t, "read")
	fakePrompt(t, "p", "q")

	out := captureStdout(t, func() { q.Read() })
	if !strings.Contains(out, "This is the first quote.") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestRead_StopsAfterLastQuote(t *testing.T) {
	q := testQuotes()
	withArgs(t, "read", "3")
	fakePrompt(t, "")

	var err error
	out := captureStdout(t, func() { err = q.Read() })
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "(3/3)") || !strings.Contains(out, "That was the last quote.") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestRead_StartsAtGivenID(t *testing.T) {
	q := testQuotes()
	withArgs(t, "read", "2")
	fakePrompt(t, "q")

	out := captureStdout(t, func() { q.Read() })
	if !strings.Contains(out, "(2/3)") || strings.Contains(out, "(1/3)") {
		t.Errorf("should start at quote 2: %q", out)
	}
}

func TestRead_InvalidStartID(t *testing.T) {
	q := testQuotes()

	withArgs(t, "read", "99")
	if err := q.Read(); err == nil {
		t.Error("expected error for unknown ID")
	}

	withArgs(t, "read", "abc")
	if err := q.Read(); err == nil {
		t.Error("expected error for invalid ID")
	}
}

func TestRead_RandomStaysInRange(t *testing.T) {
	q := &Quotes{QUOTES: []Quote{{ID: 1, QUOTE: "Only one", AUTHOR: "A"}}}
	withArgs(t, "read")
	fakePrompt(t, "r", "r", "q")

	var err error
	captureStdout(t, func() { err = q.Read() })
	if err != nil {
		t.Fatal(err)
	}
}

func TestRead_CancelIsNotAnError(t *testing.T) {
	q := testQuotes()
	withArgs(t, "read")
	fakePrompt(t) // Ctrl+C at the first prompt

	var err error
	captureStdout(t, func() { err = q.Read() })
	if err != nil {
		t.Errorf("Ctrl+C while reading should return nil, got %v", err)
	}
}

/* Duplicates */

func dupQuotes() *Quotes {
	return &Quotes{QUOTES: []Quote{
		{ID: 1, QUOTE: "Be yourself; everyone else is already taken", AUTHOR: "Oscar Wilde"},
		{ID: 2, QUOTE: "The only thing we have to fear is fear itself", AUTHOR: "F. D. Roosevelt"},
		{ID: 3, QUOTE: "BE YOURSELF, everyone else is already taken!!", AUTHOR: "Someone"},
		{ID: 4, QUOTE: "The only thing we have to fear is fear itselg", AUTHOR: "Typo"},
		{ID: 5, QUOTE: "A completely unrelated line about gardening", AUTHOR: "Nobody"},
	}}
}

func TestDuplicates_NoneFound(t *testing.T) {
	q := testQuotes()
	withArgs(t, "dub")

	var err error
	out := captureStdout(t, func() { err = q.Duplicates() })
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "No duplicates or similar quotes found.") {
		t.Errorf("unexpected output: %q", out)
	}
}

func TestDuplicates_FindsExactAndNear(t *testing.T) {
	q := dupQuotes()
	withArgs(t, "dub")

	var err error
	out := captureStdout(t, func() { err = q.Duplicates() })
	if err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(out, "2 pair(s) found") {
		t.Errorf("expected 2 pairs: %q", out)
	}
	if !strings.Contains(out, "exact duplicate") {
		t.Errorf("expected an exact duplicate: %q", out)
	}
	if !strings.Contains(out, "98% similar") {
		t.Errorf("expected a 98%% similar pair: %q", out)
	}
	if strings.Contains(out, "gardening") {
		t.Errorf("unrelated quote should not be listed: %q", out)
	}
	// exact matches are listed before near matches
	if strings.Index(out, "exact duplicate") > strings.Index(out, "98% similar") {
		t.Errorf("results should be sorted best first: %q", out)
	}
}

func TestDuplicates_ThresholdArgument(t *testing.T) {
	for _, arg := range []string{"0.99", "99"} { // fraction and percent forms
		q := dupQuotes()
		withArgs(t, "dub", arg)

		var err error
		out := captureStdout(t, func() { err = q.Duplicates() })
		if err != nil {
			t.Fatalf("threshold %q: %v", arg, err)
		}
		if !strings.Contains(out, "1 pair(s) found") {
			t.Errorf("threshold %q should leave only the exact pair: %q", arg, out)
		}
	}
}

func TestDuplicates_InvalidThreshold(t *testing.T) {
	for _, arg := range []string{"abc", "0", "-1", "150"} {
		q := dupQuotes()
		withArgs(t, "dub", arg)
		if err := q.Duplicates(); err == nil {
			t.Errorf("threshold %q should be rejected", arg)
		}
	}
}

func TestDuplicates_EmptyCollection(t *testing.T) {
	q := &Quotes{}
	withArgs(t, "dub")

	var err error
	out := captureStdout(t, func() { err = q.Duplicates() })
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "No duplicates") {
		t.Errorf("unexpected output: %q", out)
	}
}

/* Benchmarks */

// BenchmarkSimilarity measures one comparison of two typical, similar-length
// but different quotes (the common case when scanning for duplicates).
func BenchmarkSimilarity(b *testing.B) {
	x := normalize("The only thing we have to fear is fear itself, and nothing more than that")
	y := normalize("Be yourself; everyone else is already taken, so do not try to be someone")
	for i := 0; i < b.N; i++ {
		similarity(x, y, dupThreshold)
	}
}

// BenchmarkDuplicateScan simulates Duplicates() over 1000 synthetic quotes.
func BenchmarkDuplicateScan(b *testing.B) {
	r := rand.New(rand.NewSource(1))
	words := strings.Fields("life love time world people heart mind truth never always great small " +
		"good bad dream hope fear light dark road home friend enemy peace war day night")

	norm := make([]string, 1000)
	for i := range norm {
		n := 8 + r.Intn(10)
		parts := make([]string, n)
		for j := range parts {
			parts[j] = words[r.Intn(len(words))]
		}
		norm[i] = strings.Join(parts, " ")
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < len(norm); i++ {
			for j := i + 1; j < len(norm); j++ {
				similarity(norm[i], norm[j], dupThreshold)
			}
		}
	}
}

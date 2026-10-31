package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"unicode/utf8"
)

func (q *Quotes) Report() {

	clearScreen()

	var ProgramVersion = "1.24"

	nowPlaying := "Now playing: Flute.mp3"

	stringFormat := `` +
		Cyan + "VK-Quotes" + Reset + " %s" + "\n" + // Program Name
		Purple + "%s" + Reset + "\n" + // Now Playing
		"%s" + // Last Quote
		Yellow + `%s` + Reset + // Commands
		``

	formattedLastQuote := FormatQuote(q.QUOTES[len(q.QUOTES)-1])

	commands := "\n< add, update, delete, read, search, dublicates >\n"

	cli := fmt.Sprintf(stringFormat, ProgramVersion, nowPlaying, formattedLastQuote, commands)

	fmt.Print(cli)

}

const quoteWidth = 72

func FormatQuote(quote Quote) string {
	var b strings.Builder

	prefix := fmt.Sprintf("%d. ", quote.ID)
	indent := strings.Repeat(" ", len(prefix))

	// Room for the prefix and the opening quote mark.
	lines := wrapText(quote.QUOTE, quoteWidth-len(prefix)-1)

	b.WriteString("\n")
	for i, l := range lines {
		if i == 0 {
			b.WriteString(Gray + prefix + Reset + Cyan + "“" + Reset)
		} else {
			b.WriteString(indent + " ")
		}
		b.WriteString(Italic + l + Reset)
		if i == len(lines)-1 {
			b.WriteString(Cyan + "”" + Reset)
		}
		b.WriteString("\n")
	}

	// Attribution line
	b.WriteString(indent + "  " + Green + "— " + quote.AUTHOR + Reset)

	var meta []string
	if quote.DATE != "" {
		meta = append(meta, quote.DATE)
	}
	if quote.LANGUAGE != "" {
		meta = append(meta, quote.LANGUAGE)
	}
	if len(meta) > 0 {
		b.WriteString("  " + Gray + strings.Join(meta, " · ") + Reset)
	}
	b.WriteString("\n")

	return b.String()
}

func wrapText(text string, width int) []string {
	var lines []string

	for _, para := range strings.Split(text, "\n") {
		words := strings.Fields(para)
		if len(words) == 0 {
			lines = append(lines, "")
			continue
		}

		cur := words[0]
		for _, w := range words[1:] {
			if utf8.RuneCountInString(cur)+1+utf8.RuneCountInString(w) > width {
				lines = append(lines, cur)
				cur = w
			} else {
				cur += " " + w
			}
		}
		lines = append(lines, cur)
	}

	return lines
}

func clearScreen() {

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error clearing screen:", err)
	}
}

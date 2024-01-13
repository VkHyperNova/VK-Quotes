package cmd

import (
	"math/rand"
	"strconv"
	"strings"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"

	"github.com/fatih/color"
)

func PrintVKQUOTES(Version string) {
	green := color.New(color.FgGreen)
	boldGreen := green.Add(color.Bold)
	boldGreen.Println("\n<< VK-QUOTES " + Version + " >>")

	util.PrintCyan("\nQuotes: " + strconv.Itoa(len(db.DATABASE)) + "\n")

	if db.SearchIndexByID(db.LastAddID) != -1 {
		PrintQuote(db.SearchIndexByID(db.LastAddID))
	} else if len(db.DATABASE) > 0 {
		PrintRandomQuote()
	}

	util.PrintGray("\n")
	Commands := [6]string{"add", "update", "delete", "showall", "stats", "q"}
	for _, value := range Commands {
		PrintBrackets(value)
	}
	util.PrintCyan("\n=> ")
}

func PrintBrackets(name string) {
	util.PrintCyan("[")
	util.PrintYellow(name)
	util.PrintCyan("] ")
}

func PrintQuote(index int) {
	spaces := strings.Repeat(" ", 5)
	if len(db.DATABASE[index].QUOTE)-len(db.DATABASE[index].AUTHOR) >= 10 {
		spaces = strings.Repeat(" ", len(db.DATABASE[index].QUOTE)-len(db.DATABASE[index].AUTHOR))
	}

	util.PrintBlue("\n(" + strconv.Itoa(db.DATABASE[index].ID) + ") ")
	util.PrintGray("\"" + db.DATABASE[index].QUOTE + "\"")
	util.PrintCyan("\n" + spaces + " By " + db.DATABASE[index].AUTHOR + " (" + db.DATABASE[index].DATE + ")\n")
}

func PrintAllQuotes() {
	util.PrintCyan("\n\n<< All Quotes >>\n")

	for key := range db.DATABASE {
		PrintQuote(key)
	}
	util.PressAnyKey()
	util.ClearScreen()
	CMD()
}

func PrintRandomQuote() {

	randIndex := rand.Intn(len(db.DATABASE))

	for key := range db.DATABASE {
		if key == randIndex {
			PrintQuote(key)
		}

	}
}

func PrintStatistics() {
	util.PrintCyan("\n\n<< Statistics >>\n")

	// Get all authors
	var authors []string
	for _, value := range db.DATABASE {
		if !util.Contains(authors, value.AUTHOR) {
			authors = append(authors, value.AUTHOR)
		}
	}

	// Count quotes by author
	myMap := make(map[string]int)

	for _, author := range authors {
		for _, value := range db.DATABASE {
			if value.AUTHOR == author {
				myMap[author] += 1
			}
		}
	}

	// Print stats by author
	for author, count := range myMap {
		util.PrintGray(author + ": " + strconv.Itoa(count) + "\n")
	}

	// Get all languages
	var languages []string
	for _, value := range db.DATABASE {
		if !util.Contains(languages, value.LANGUAGE) {
			languages = append(languages, value.LANGUAGE)
		}
	}

	myMap2 := make(map[string]int)

	for _, language := range languages {
		for _, value := range db.DATABASE {
			if value.LANGUAGE == language {
				myMap2[language] += 1
			}
		}
	}

	// Print stats by language
	for language, count := range myMap2 {
		util.PrintGray(language + ": " + strconv.Itoa(count) + "\n")
	}



	util.PressAnyKey()
	util.ClearScreen()
	CMD()
}

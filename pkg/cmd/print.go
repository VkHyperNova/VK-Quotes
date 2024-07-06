package cmd

/*
All Print functions
*/

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"

	"github.com/fatih/color"
)

func PrintCLI(quotes *db.Quotes) {

	PrintProgramNameAndVersion()

	/* Message */
	PrintMessage()

	/* Quote */
	PrintRandomQuote(quotes)
	PrintEditedQuote(quotes)
	PrintLastQuote(quotes)

	/* ReadMode */
	PrintReadCounter(quotes)

	/* Commands */
	PrintCommands()
}

func PrintLastQuote(quotes *db.Quotes) {
	if !IsReadMode && MustPrintQuoteID == -1 {
		index := quotes.QuotesCount() - 1
		LastQuoteId, err := quotes.FindIDByIndex(index)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		quotes.PrintQuote(LastQuoteId)
	}
}

func PrintProgramNameAndVersion() {
	white := color.New(color.FgCyan)
	boldGreen := white.Add(color.Bold)
	util.PrintGreen("\n|| ")
	boldGreen.Print("VK-QUOTES " + IsVersion)
	util.PrintGreen(" ||")
}

func PrintCommands() {
	Commands := [7]string{"add", "update", "delete", "read", "showall", "stats", "q"}
	util.PrintCyan("\n\n")
	for _, value := range Commands {
		util.PrintGreen("|")
		util.PrintYellow(value)
		util.PrintGreen("| ")
	}
	util.PrintYellow("=> ")
}

func PrintMessage() {

	if IsMessage != "" {
		length := len(IsMessage) + 5
		dots := ""
		for i := 1; i < length; i++ {
			dots += "-"
		}
		util.PrintGreen("\n" + dots + "\n" + IsMessage + "\n" + dots)
	}
}

func PrintReadCounter(quotes *db.Quotes) {

	if IsReadMode {
		MustPrintQuoteID = -1
		util.PrintGreen("\n[" + strconv.Itoa(IsReadCount) + "] ")
		
		percentage := float64(IsReadCount) / float64(quotes.QuotesCount()) * 100
		util.PrintGray(fmt.Sprintf("%.2f", percentage) + "% ")

		i := 0
		util.PrintGray("|")
		for i < IsReadCount {
			util.PrintGreen("-")
			i++
		}
		util.PrintGray("|")
	}
}

func PrintAllQuotes(quotes *db.Quotes) {
	util.PrintCyan("\n\n<< All Quotes >>\n")
	
	for index := range quotes.QUOTES {
		id, err := quotes.FindIDByIndex(index)
		if err != nil {
			fmt.Println(err)
		}
		quotes.PrintQuote(id)
	}
}

func PrintRandomQuote(quotes *db.Quotes) {

	/* Check if index exists and append to global variable 'IsUsedIndexes' */

	var randomIndex int
	isValid := false

	if quotes.QuotesCount() > 0 {
		for !isValid {
			randomIndex = rand.Intn(quotes.QuotesCount())
			if !util.ArrayContainsInt(IsUsedIndexes, randomIndex) {
				if IsReadMode {
					IsUsedIndexes = append(IsUsedIndexes, randomIndex)
				}
				isValid = true
			}
		}
	}

	/* Find and print random quote */

	if quotes.QuotesCount() > 0 && IsReadMode {
		for index := range quotes.QUOTES {
			if index == randomIndex {
				id, err := quotes.FindIDByIndex(index)
				if err != nil {
					util.PrintRed("Error: " + err.Error())
				}
				quotes.PrintQuote(id)
			}
		}
	}
}

func PrintEditedQuote(quotes *db.Quotes) {

	/* Print last edited quote if exists */

	if MustPrintQuoteID != -1 {
		quotes.PrintQuote(MustPrintQuoteID)
	}
}

func DeleteUsedIndexes(quotes *db.Quotes) {

	/* Empty the Slice if its full */

	if len(IsUsedIndexes) == quotes.QuotesCount() {
		IsMessage = "<< You Have Read Everything! >>"
		IsUsedIndexes = IsUsedIndexes[:0]
		IsReadCount = 0
		IsReadMode = false
	}
}

/* Statistics */

func PrintStatistics(quotes *db.Quotes) {
	util.PrintCyan("\n\n\t<< Statistics >>\n")
	util.PrintCyan("\n--------------------------------------------\n\n")
	PrintAuthors(quotes)
	util.PrintCyan("\n--------------------------------------------\n")
	PrintLanguages(quotes)
	util.PrintCyan("\n\n--------------------------------------------\n")
}

func PrintAuthors(quotes *db.Quotes) {

	/* Get All Author Names */

	var authors []string
	for _, value := range quotes.QUOTES {
		if !util.ArrayContainsString(authors, value.AUTHOR) {
			authors = append(authors, value.AUTHOR)
		}
	}

	/* Count Authors By Name */

	authorsMap := make(map[string]int)
	for _, name := range authors {
		for _, value := range quotes.QUOTES {
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

	/* Sort */

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})

	/* Print */

	for i := 0; i < len(pairs) && i < 10; i++ {
		util.PrintGray("[" + strconv.Itoa(pairs[i].count) + "] ")
		util.PrintGreen(pairs[i].name + "\n")
	}
}

func PrintLanguages(quotes *db.Quotes) {

	languages := []string{"English", "Russian", "Estonian", "German"}

	/* Count */
	languagesMap := make(map[string]int)
	for _, name := range languages {
		for _, value := range quotes.QUOTES {
			if value.LANGUAGE == name {
				languagesMap[name] += 1
			}
		}
	}

	/* Print */
	for name, num := range languagesMap {
		util.PrintGray("\n[" + strconv.Itoa(num) + "] " + name)
	}
}

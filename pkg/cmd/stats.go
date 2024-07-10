package cmd

import (
	"sort"
	"strconv"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

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

	/* Sort by count */

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})

	/* Print top 10 */

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
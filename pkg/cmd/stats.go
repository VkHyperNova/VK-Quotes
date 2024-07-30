package cmd

import (
	"fmt"
	"sort"
	"strconv"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util" 
)

/* Statistics mby move me to methods*/

func printStats(quotes *db.Quotes) {
	util.ClearScreen()
	format := "%s%s%s"
	name := util.Cyan + "Statistics: " + util.Reset
	stats := fmt.Sprintf(format, name, topAuthors(quotes), topLanguages(quotes))
	fmt.Println(stats)
}

func topAuthors(quotes *db.Quotes) string {

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

	// Make a string
	authorsString := ""
	for i := 0; i < len(pairs) && i < 10; i++ {
		authorsString += "\n" + strconv.Itoa(pairs[i].count) + " " + util.Cyan + pairs[i].name + util.Reset
	}

	return authorsString
}

func topLanguages(quotes *db.Quotes) string {

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

	/* Make a string */
	languagesString := ""
	for name, num := range languagesMap {
		languagesString += "\n" + strconv.Itoa(num) + " " + util.Yellow + name + util.Reset
	}
	return languagesString
}

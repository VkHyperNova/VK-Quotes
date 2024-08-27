package db

import (
	"fmt"
	"sort"
	"strconv"
	"vk-quotes/pkg/config"
	"vk-quotes/pkg/util"
)

func (q *Quotes) PrintStatistics() {

	util.ClearScreen()

	format := "%s%s%s"

	name := config.Cyan + "Statistics: " + config.Reset

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
		authorsString += "\n" + strconv.Itoa(pairs[i].count) + " " + config.Cyan + pairs[i].name + config.Reset
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
		languagesString += "\n" + strconv.Itoa(num) + " " + config.Yellow + name + config.Reset
	}
	return languagesString
}

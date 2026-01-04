package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"vk-quotes/pkg/audio"
	"vk-quotes/pkg/config"
	"vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

func CommandLine(quotes *db.Quotes) {

	go func() {
		err := audio.PlayMP3()
		if err != nil {
			log.Printf("MP3 Playback Error: %v", err)
		}
	}()

	for {
		quotes.PrintCLI()

		var input string = ""
		var inputID int = 0

		fmt.Printf("[%d]=> ", len(quotes.QUOTES))

		fmt.Scanln(&input, &inputID)

		input = strings.ToLower(input)

		switch input {

		case "add", "a":
			added := quotes.Add()
			if !added {
				audio.PlayErrorSound()
			}
		case "update", "u":
			err := quotes.Update(inputID)
			if  err != nil {
				fmt.Println(config.Red, err, config.Reset)
				audio.PlayErrorSound()
				util.PressAnyKey()
			}
		case "delete", "d":
			deleted := quotes.Delete(inputID)
			if !deleted {
				audio.PlayErrorSound()
			}

			quotes.ResetIDs(quotes) // Reset all IDs in database
		case "find", "f":
			found := quotes.Find()
			if !found {
				audio.PlayErrorSound()
			}
		case "showall":
			quotes.PrintAllQuotes()
		case "stats":
			quotes.PrintStatistics()
		case "resetids":
			quotes.ResetIDs(quotes)
		case "read":
			quotes.Read()
		case "similarquotes":
			db.FindSimilarQuotes(quotes)
		case "pause":
			audio.PauseMP3()
		case "resume":
			audio.ResumeMP3()
		case "random", "r":
			quotes.PrintRandomQuotes(inputID)
		case "q", "quit":
			quotes.Backup()
			util.ClearScreen()
			os.Exit(0)
		default:
			config.AddMessage("Enter pressed!")
		}
	}
}

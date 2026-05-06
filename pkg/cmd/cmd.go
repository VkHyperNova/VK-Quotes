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

func CommandLine(q *db.Quotes) {

	go func() {
		err := audio.PlayMP3()
		if err != nil {
			log.Printf("MP3 Playback Error: %v", err)
		}
	}()

	for {
		q.PrintCLI()

		var input string = ""
		var inputID int = 0

		fmt.Printf("[%d]=> ", len(q.QUOTES))

		fmt.Scanln(&input, &inputID)

		input = strings.ToLower(input)

		switch input {

		case "add", "a":
			if err := q.Add(); err != nil {
				audio.PlayErrorSound()
				fmt.Println(err)
				util.PressAnyKey()
			}
		case "update", "u":
			if err := q.Update(inputID); err != nil {
				audio.PlayErrorSound()
				fmt.Println(err)
				util.PressAnyKey()
			}
		case "delete", "d":
			if err := q.Delete(inputID); err != nil {
				audio.PlayErrorSound()
				fmt.Println(err)
				util.PressAnyKey()
			}
		case "find", "f":
			found := q.Find()
			if !found {
				audio.PlayErrorSound()
			}
		case "import", "i":
			if err := q.Import(); err != nil {
				fmt.Println(err)
			}
			util.PressAnyKey()
		case "export", "e":
			if err := q.Export(); err != nil {
				fmt.Println(err)
			}
			util.PressAnyKey()
		case "unmount":
			if err := util.UnmountDrive(); err != nil {
				fmt.Println(err)
			}
			util.PressAnyKey()
		case "history", "h":
			q.History()
		case "stats":
			q.PrintStatistics()
		case "read":
			q.Read()
		case "findsimilar":
			db.FindSimilarQuotes(q)
		case "pause":
			audio.PauseMP3()
		case "resume":
			audio.ResumeMP3()
		case "random", "r":
			q.PrintRandomQuotes(inputID)
		case "q", "quit":
			util.ClearScreen()
			os.Exit(0)
		default:
			config.AddMessage("Enter pressed!")
		}
	}
}

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	Bold   = "\033[1m"
	Italic = "\033[3m"

	driveLabel      = "VK DATA"
	driveMountPoint = "/media/veikko/VK DATA"
	file            = "quotes.json"
	baseDB          = "QUOTES"
	baseLocal       = "DATABASES"
	// baseLocal = "DEVELOPMENT"
)

var (
	LocalFile  = filepath.Join(baseLocal, baseDB, file)
	BackupFile = filepath.Join(driveMountPoint, baseLocal, baseDB, file)
)

var notices []string

func notice(msg string) { notices = append(notices, msg) }

const usage = Bold + "Quotes" + Reset + " - manage your quote collection\n" +
	"\n" +
	Yellow + "Usage:" + Reset + "\n" +
	"  quotes <command> [arguments]\n" +
	"\n" +
	Yellow + "Commands:" + Reset + "\n" +
	"  " + Green + "add" + Reset + ", a, log" + Gray + "         Add a new quote (warns about similar ones)" + Reset + "\n" +
	"  " + Green + "update" + Reset + ", u <id>" + Gray + "      Edit the quote with the given ID" + Reset + "\n" +
	"  " + Green + "delete" + Reset + ", d <id>" + Gray + "      Delete the quote with the given ID" + Reset + "\n" +
	"  " + Green + "search" + Reset + ", s <text>" + Gray + "    Search quote, author, language and date" + Reset + "\n" +
	"  " + Green + "read" + Reset + ", r [id]" + Gray + "        Read quotes one by one (optionally from an ID)" + Reset + "\n" +
	"  " + Green + "dublicates" + Reset + ", dub [t]" + Gray + "  Find similar quotes (threshold 0-1, default 0.85)" + Reset + "\n" +
	"  " + Green + "help" + Reset + ", -h, --help" + Gray + "     Show this help" + Reset + "\n" +
	"\n" +
	Yellow + "Examples:" + Reset + "\n" +
	"  quotes add\n" +
	"  quotes search be yourself\n" +
	"  quotes update 42\n" +
	"  quotes read 100\n" +
	"  quotes dub 0.9\n"

func main() {

	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "help", "-h", "--help":
			fmt.Print(usage)
			return
		}
	}

	q := Quotes{}

	if err := q.Load(LocalFile); err != nil {
		log.Fatalf("load from file: %v", err)
	}

	go func() {
		if err := PlayMP3(); err != nil {
			log.Printf("MP3 Playback Error: %v", err)
		}
	}()

	var err error
	showReport := true // redraw the main screen after the command

	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "add", "a", "log":
			err = q.Add()
		case "update", "u":
			err = q.Update()
		case "delete", "d":
			err = q.Delete()
		case "search", "s":
			err = q.Search()
			showReport = false
		case "read", "r":
			err = q.Read()
		case "dublicates", "dub":
			err = q.Duplicates()
			showReport = false
		case "reset", "resetids":
			err = q.ResetIDs()
		default:
			closePrompt()
			fmt.Printf("%serror: unknown command: %s%s\n\n", Red, os.Args[1], Reset)
			fmt.Print(usage)
			os.Exit(1)
		}
	}

	// Restore the terminal before printing anything else.
	closePrompt()

	// Ctrl+C / Ctrl+D in a prompt is a normal way to cancel.
	if errors.Is(err, errCancelled) {
		err = nil
	}

	if showReport {
		q.Report()
	}

	for _, n := range notices {
		fmt.Println(Yellow + n + Reset)
	}

	if err != nil {
		fmt.Printf("%serror: %v%s\n", Red, err, Reset)
		os.Exit(1)
	}
}

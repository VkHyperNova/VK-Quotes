package util

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/peterh/liner"
)

type Settings struct {
	RandomIDs       []int
	ID              int
	UserInputs      []string
	Command         string
	Message         string
	ReadCounter     int
	SaveQuotesPath  string
	SaveSimilarPath string
	SaveFolderPath  string
	Version         string
}

func CommandPrompt(settings *Settings) (string, int) {

	line := liner.NewLiner()
	defer line.Close()

	input, err := line.Prompt("")

	if err != nil {
		settings.Message = "<< Error reading input >>"
	}

	parts := strings.Fields(input)

	if len(parts) > 1 {
		IsInteger, err := strconv.Atoi(parts[1])
		if err != nil {
			return input, -1
		}
		return parts[0], IsInteger
	}

	return input, -1
}

func ClearScreen() {

	if runtime.GOOS == "linux" {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func PressAnyKey() {
	PrintGray("\nPress Any Key To Continue...")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	ClearScreen()
}

func ArrayContainsString(arr []string, name string) bool {
	for _, n := range arr {
		if n == name {
			return true
		}
	}
	return false
}

func ArrayContainsInt(numbers []int, number int) bool {
	for _, n := range numbers {
		if n == number {
			return true
		}
	}
	return false
}

func FillEmptyInput(a, b string) string {
	if a == "" {
		a = b
	}

	return a
}

func MoveBack(a string) bool {
	return a == "b"
}

func SetRandomID(settings *Settings) {

	/* Get random int */
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	randomIndex := r.Intn(len(settings.RandomIDs))

	/* Set id */
	settings.ID = settings.RandomIDs[randomIndex]

	/* Remove this id from list */
	settings.RandomIDs = append(settings.RandomIDs[:randomIndex], settings.RandomIDs[randomIndex+1:]...)
}

// ReadCounter generates a formatted string displaying the read counter and the read percentage.
func ReadCounter(count int, size int) string {
    // Calculate the read percentage.
    percentage := float64(count) / float64(size) * 100
    
    // Return the formatted read counter string.
    return fmt.Sprintf("\n[%d] %.0f%%\n", count, percentage)
}


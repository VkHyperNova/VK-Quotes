package util

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"vk-quotes/pkg/config"

	"github.com/peterh/liner"
)

func CommandPrompt(prompt string) (string, int) {
	/*
		CommandPrompt prompts the user for input and parses it into a command and an optional ID.
		It returns the command as a string and the ID as an integer. If an error occurs,
		it sets an error message in the settings.
	*/

	// Initialize a new line reader for user input

	line := liner.NewLiner()

	// Ensure the line reader is closed to free resources

	defer line.Close()

	// Prompt the user with the given prompt string and read input

	input, err := line.Prompt(prompt)

	// Handle any errors from the prompt

	if err != nil {
		config.Messages = append(config.Messages, err.Error())
	}

	// Initialize the command and commandID variables

	command := input
	commandID := 0

	// Split the input into parts based on whitespace

	parts := strings.Fields(input)

	// Check if the input contains exactly two parts

	if len(parts) == 2 {

		// Assume the first part is the command

		isCommand := parts[0]

		// Try to convert the second part to an integer

		isID, err := strconv.Atoi(parts[1])

		// If the conversion is successful, update the command and commandID

		if err == nil {
			command = isCommand
			commandID = isID
		}
	}

	// Return the parsed command and commandID

	return command, commandID
}

func ClearScreen() {

	cmd := exec.Command("clear")

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}

func PressAnyKey() {
	PrintGray("\nPress Any Key To Continue...\n")
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

func FormattedReadCounter(count int, size int) string {

	percentage := float64(count) / float64(size) * 100

	readCounter := fmt.Sprintf("Reading: [%d] %.0f%%", count, percentage)

	return readCounter
}

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

	line := liner.NewLiner()
	defer line.Close()

	// Prompt the user with the given prompt string and read userInput
	userInput, err := line.Prompt(prompt)
	if err != nil {
		config.AddMessage(err.Error())
		return "", 0
	}

	// Initialize default values
	input := strings.TrimSpace(userInput)
	inputID := 0

	// Split the input into parts based on whitespace
	parts := strings.Fields(userInput)

	// Check if the input contains exactly two parts
	if len(parts) == 2 {

		// Assume the first part is the command
		isCommand := parts[0]

		// Try to convert the second part to an integer
		isID, err := strconv.Atoi(parts[1])

		// If the conversion is successful, update the command and commandID
		if err == nil {
			input = isCommand
			inputID = isID
		}
	}

	// Convert the input to lowercase
	input = strings.ToLower(input)

	return input, inputID
}

func ClearScreen() {

	var cmd *exec.Cmd

	// Determine the command to clear the screen based on the OS
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	// Set the command output to the standard output
	cmd.Stdout = os.Stdout

	// Run the command and handle any errors
	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error clearing screen:", err)
	}
}

func PressAnyKey() {

	fmt.Println("\nPress any key to continue...")

	// Create a new scanner to read from standard input
	scanner := bufio.NewScanner(os.Stdin)

	// Scan the next token (line) from standard input
	if scanner.Scan() {
		// Successfully read input; can handle further actions here if needed
		return
	}

	// Check for scanning errors
	if err := scanner.Err(); err != nil {
		// Print the error message to standard error
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
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

func CreateDirectory() {

	path := config.LocalPath

	if _, err := os.Stat(path); os.IsNotExist(err) {

		_ = os.Mkdir("QUOTES", 0700)

		err = os.WriteFile(path, []byte(`{"quotes": []}`), 0644)
		if err != nil {
			panic(err)
		}

		message := config.Green + "Local Database Created" + config.Reset
		config.AddMessage(message)
	}
}


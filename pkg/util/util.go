package util

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"unicode"
	"vk-quotes/pkg/config"

)

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

		// Create dir 
		_ = os.Mkdir("QUOTES", 0700)

		// Create file
		err = os.WriteFile(path, []byte(`{"quotes": []}`), 0644)
		if err != nil {
			panic(err)
		}

		message := config.Green + "Local Database Created" + config.Reset
		config.AddMessage(message)
	}
}

func Quit(input string) bool {
	if input == "q" {
		message := config.Red + "Previous action aborted by user" + config.Reset
		config.AddMessage(message)
		return false
	}

	return true
}

func AutoDetectLanguage(quote string) string {

	for _, char := range quote {
		if unicode.In(char, unicode.Cyrillic) {
			return "Russian"
		}
	}

	return "English"
}



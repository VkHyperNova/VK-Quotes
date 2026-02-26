package util

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"
	"vk-quotes/pkg/color"
	"vk-quotes/pkg/config"

	"github.com/peterh/liner"
)

func EnsureSentenceEnd(s string) string {
	if strings.HasSuffix(s, ".") || strings.HasSuffix(s, "!") || strings.HasSuffix(s, "?") {
		return s
	}
	return s + "."
}

func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	firstChar := string([]rune(s)[0])              // Get the first character as a string
	capitalizedFirst := strings.ToUpper(firstChar) // Capitalize the first character
	return capitalizedFirst + s[len(firstChar):]   // Concatenate with the rest of the string
}

func ClearScreen() {

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Error clearing screen:", err)
	}
}

func PressAnyKey() {
    fmt.Println("\nPress Enter to continue...")

    reader := bufio.NewReader(os.Stdin)
    _, err := reader.ReadString('\n') // Waits for the user to press Enter
    if err != nil {
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

func Quit(input string) bool {
	if input == "q" {
		message := color.Red + "Previous action aborted by user" + color.Reset
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

func PromptWithSuggestion(name string, edit string) string {

	line := liner.NewLiner()
	defer line.Close()

	input, err := line.PromptWithSuggestion("   "+name+": ", edit, -1)
	if err != nil {
		panic(err)
	}

	return input
}

func Contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func ensureFile(path string, content string) error {

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("error creating directory for %s: %w", path, err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("error creating file %s: %w", path, err)
		}
	}

	return nil
}

func CreateFilesAndFolders() error {

	if err := ensureFile(config.LocalFile, config.DefaultContent); err != nil {
		return err
	}

	if !HardDriveMountCheck() {
		input := Input("Do you want to continue? (y/n) ")
		if strings.ToLower(strings.TrimSpace(input)) != "y" {
			fmt.Println("Exiting program.")
			os.Exit(0)
		}
	} else {
		if err := ensureFile(config.BackupFile, config.DefaultContent); err != nil {
			return err
		}
	}

	return nil
}

func HardDriveMountCheck() bool {
	if runtime.GOOS != "linux" {
		fmt.Println("This program only works on Linux.")
		return false
	}

	mountPoint := "/media/veikko/VK\\040DATA" // match /proc/mounts format

	file, err := os.Open("/proc/mounts")
	if err != nil {
		fmt.Println("Cannot open /proc/mounts:", err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 2 && fields[1] == mountPoint {
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning /proc/mounts:", err)
		return false
	}

	fmt.Println(color.Red + "\nVK DATA is NOT mounted" + color.Reset)
	return false
}

func Input(prompt string) string {

	line := liner.NewLiner()
	defer line.Close()

	userInput, err := line.Prompt(prompt)
	if err != nil {
		panic(err)
	}
	return userInput
}
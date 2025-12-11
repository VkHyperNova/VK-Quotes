package util

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"unicode"
	"vk-quotes/pkg/audio"
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

func PromptWithSuggestion(name string, edit string) string {

	line := liner.NewLiner()
	defer line.Close()

	input, err := line.PromptWithSuggestion("   "+name+": ", edit, -1)
	if err != nil {
		panic(err)
	}

	return input
}

func isMounted(mountPoint string) (bool, error) {
    file, err := os.Open("/proc/mounts")
    if err != nil {
        return false, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        fields := strings.Fields(line)
        if len(fields) >= 2 && fields[1] == mountPoint {
            return true, nil
        }
    }

    return false, scanner.Err()
}

func IsVKDataMounted() {

	if runtime.GOOS != "linux" {
        fmt.Println("This program only works on Linux.")
        return
    }

	mountPoint := "/media/veikko/VK\\040DATA" // change to your actual mount path

    mounted, err := isMounted(mountPoint)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    if mounted {
        fmt.Println(config.Green + "VK DATA is mounted" + config.Reset)
    } else {
        fmt.Println(config.Red + "VK DATA is NOT mounted" + config.Reset)
		audio.PlayErrorSound()
    }
}

func Contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

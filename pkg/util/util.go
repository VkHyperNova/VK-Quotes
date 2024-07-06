package util

import (
	"bufio"
	"os"
	"os/exec"
	"runtime"

	"github.com/peterh/liner"
)

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

func CreateRequiredFiles(DatabasePath string) {
	if _, err := os.Stat(DatabasePath); os.IsNotExist(err) {
		_ = os.Mkdir("database", 0700)
		HandleError(os.WriteFile(DatabasePath, []byte(`{"quotes": []}`), 0644))
		PrintRed("New Database Created!\n")
	}
}

func Abort(a string) bool {
	if a == "q" {
		return true
	}
	return false
}

func ScanOrEditWithLiner(name string, editableString string) string {

	line := liner.NewLiner()
	defer line.Close()

	if editableString != "" {
		input, err := line.PromptWithSuggestion("   "+name+": ", editableString, -1)
		HandleError(err)
		return input
	} else {
		input, err := line.Prompt("   " + name + ": ")
		HandleError(err)
		return input
	}
}

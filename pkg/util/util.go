package util

import (
	"bufio"
	"encoding/json"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
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

func StructToJson(data interface{}) []byte {
	dataBytes, err := json.MarshalIndent(data, "", "  ")
	HandleError(err)

	return dataBytes
}

func GetFormattedDate() string {
	return time.Now().Format("02.01.2006")
}

func PressAnyKey() {
	PrintGray("\nPress Any Key To Continue...")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	ClearScreen()
}

func Contains(arr []string, name string) bool {
	for _, n := range arr {
		if n == name {
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

func PrintBrackets(name string) {
	PrintCyan("[")
	PrintYellow(name)
	PrintCyan("] ")
}


func ScanUserInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	userInput := scanner.Text()
	userInput = strings.TrimSpace(userInput)

	return userInput
}

func Abort(input string) bool {
	return input == "q"
}

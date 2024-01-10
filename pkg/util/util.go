package util

import (
	"bufio"
	"encoding/json"
	"fmt"
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

func GetInput(inputName string) string {
	PrintCyan(inputName)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	NewTaskString := scanner.Text()
	NewTaskString = strings.TrimSpace(NewTaskString)

	return NewTaskString
}

func InterfaceToByte(data interface{}) []byte {
	dataBytes, err := json.MarshalIndent(data, "", "  ")
	HandleError(err)

	return dataBytes
}

func GetFormattedDate() string {
	return time.Now().Format("02.01.2006")
}

func Confirm() bool {

	user_input := Prompt("\n\nThis One?: ")

	if user_input == "n" || user_input == "no" {
		return false
	}
	return true
}

func Prompt(Question string) string {

	PrintCyan(Question)

	var user_input_string string
	fmt.Scanln(&user_input_string)

	return user_input_string
}

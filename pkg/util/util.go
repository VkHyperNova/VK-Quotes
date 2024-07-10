package util

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"time"

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
		err = os.WriteFile(DatabasePath, []byte(`{"quotes": []}`), 0644)
		if err != nil {
			fmt.Println("Error creating database file")
			fmt.Println(err)
		}
		PrintRed("New Database Created!\n")
	}
}

func Abort(a string) bool {
	return a == "q"
}

func MoveBack(a string) bool {
	return a == "b"
}

func PromptWithSuggestion(name string, editableString string) string {

	line := liner.NewLiner()
	defer line.Close()

	input, err := line.PromptWithSuggestion("   "+name+": ", editableString, -1)
	if err != nil {
		fmt.Println("Error reading input: ", err)
		return ""
	}
	return input
}

func Prompt(name string) string {
	line := liner.NewLiner()
	defer line.Close()
	input, err := line.Prompt("   " + name + ": ")
	if err != nil {
		fmt.Println("Error reading input: ", err)
		return ""
	}
	return input
}

func GetRandomNumber(arraySize int) int {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	return r.Intn(arraySize)
}

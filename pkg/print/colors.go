package print

import (
	"fmt"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
)

func PrintRed(a string) {
	fmt.Print(Red + a + Reset)
}

func PrintGreen(a string) {
	fmt.Print(Green + a + Reset)
}

func PrintYellow(a string) {
	fmt.Print(Yellow + a + Reset)
}

func PrintBlue(a string) {
	fmt.Print(Blue + a + Reset)
}

func PrintPurple(a string) {
	fmt.Print(Purple + a + Reset)
}

func PrintCyan(a string) {
	fmt.Print(Cyan + a + Reset)
}

func PrintGray(a string) {
	fmt.Print(Gray + a + Reset)
}

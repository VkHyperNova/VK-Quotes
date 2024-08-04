package config

var ProgramVersion = "1.24"

/* All Paths */

const FolderName = "QUOTES"
const LocalPath = "./QUOTES/quotes.json"
const SimilarPath = "./QUOTES/similiar.json"

var BackupPathLinux = "/media/veikko/VK DATA/DATABASES/QUOTES/quotes"
var BackupPathWindows = "D:\\DATABASES\\QUOTES\\quotes"

var ReadCounter string
var Counter int
var RandomIDs []int

var MainQuoteID int
var Messages []string
var UserInputs []string

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

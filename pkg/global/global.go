package global

var DB []Quotes

var Version = "1.0"

var DatabasePath = "./database/quotes.json"

type Quotes struct {
	ID       int    `json:"id"`
	QUOTE    string `json:"quote"`
	AUTHOR   string `json:"author"`
	LANGUAGE string `json:"language"`
	DATE     string `json:"date"`
}

package db

import (
	"encoding/json"
	"io"
	"os"
	"runtime"
	"strconv"
	"time"
	"vk-quotes/pkg/config"
)

func (q *Quotes) ReadFromFile() {

	path := config.LocalPath

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(byteValue, q)
	if err != nil {
		panic(err)
	}
}

func (q *Quotes) SaveToFile(message string) {

	byteValue, err := json.MarshalIndent(q, "", "  ")
	if err != nil {
		panic(err)
	}

	path := config.LocalPath

	err = os.WriteFile(path, byteValue, 0644)
	if err != nil {
		panic(err)
	}

	config.AddMessage(message)
}

func (q *Quotes) Backup() {

	byteValue, err := json.MarshalIndent(q, "", "  ")
	if err != nil {
		panic(err)
	}

	currentTime := time.Now()
	layout := "(02.01.2006_15-04-05)"

	backupPath := "/media/veikko/VK DATA/DATABASES/QUOTES/" + strconv.Itoa(len(q.QUOTES)) + ". quotes " + currentTime.Format(layout) + ".json"

	if runtime.GOOS == "windows" {
		backupPath = "D:\\DATABASES\\QUOTES\\" + strconv.Itoa(len(q.QUOTES)) + ".json"
	}

	err = os.WriteFile(backupPath, byteValue, 0644)
	if err != nil {
		message := config.Red + "<< No Backup >>" + config.Reset
		config.AddMessage(message)
		config.AddMessage(err.Error())
		return
	}
}

func (q *Quotes) ResetIDs(quotes *Quotes) {

	for key := range q.QUOTES {
		q.QUOTES[key].ID = key + 1
	}

	q.SaveToFile("<< All ID's Reset >>")

	q.SetToDefaultQuote()
}

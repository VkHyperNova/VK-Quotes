package db

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
	"vk-quotes/pkg/color"
	"vk-quotes/pkg/config"
	"vk-quotes/pkg/util"
)

func (q *Quotes) ReadFromFile(path string) error {

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(byteValue, q)
	if err != nil {
		return err
	}

	return nil
}

func (q *Quotes) SaveToFile(message string) error {

	byteValue, err := json.MarshalIndent(q, "", "  ")
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(config.LocalFile, byteValue, 0644)
	if err != nil {
		panic(err)
	}

	var date = time.Now().Format("02.01.2006")

	if util.HardDriveMountCheck() {
		err = os.WriteFile(config.BackupFile, byteValue, 0644)
		if err != nil {
			return err
		}

		backupWithDate := filepath.Join(config.BaseBackup, config.BaseLocal,config.BaseDB,strconv.Itoa(len(q.QUOTES)) + ". quotes " + date + ".json")

		err = os.WriteFile(backupWithDate , byteValue, 0644)
		if err != nil {
			return err
		}
	}

	config.AddMessage(message)
	return nil
}

func (q *Quotes) Backup() {

	byteValue, err := json.MarshalIndent(q, "", "  ")
	if err != nil {
		panic(err)
	}

	currentTime := time.Now()
	layout := "(02.01.2006)"

	backupPath := "/media/veikko/VK DATA/DATABASES/QUOTES/" + strconv.Itoa(len(q.QUOTES)) + ". quotes " + currentTime.Format(layout) + ".json"

	if runtime.GOOS == "windows" {
		backupPath = "D:\\DATABASES\\QUOTES\\" + strconv.Itoa(len(q.QUOTES)) + ".json"
	}

	err = os.WriteFile(backupPath, byteValue, 0644)
	if err != nil {
		config.AddMessage(color.Red + "<< No Backup >>" + color.Reset)
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

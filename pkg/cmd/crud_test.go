package cmd

import (
	"testing"
	db "vk-quotes/pkg/db"
)

var DatabasePathTest = "/home/veikko/Desktop/VK-Quotes/database/quotes.json"

func TestAdd(t *testing.T) {

	DatabaseTest := []db.Quotes{}

	Quote = "Quote"
	Author = "Author"
	Language = "Language"

	for i := 0; i < 5; i++ {
		result := Add(&DatabaseTest, DatabasePathTest)
		CompareBoolean(result, true, t)
	}

	t.Log("10 Quotes Added!")
}

func TestUpdate(t *testing.T) {

	DatabaseTest := db.OpenDB(DatabasePathTest)

	Quote = "Quote Updated"
	Author = "Author Updated"
	Language = "Language Updated"

	for i := 1; i <= 5; i++ {
		result := Update(i, &DatabaseTest, DatabasePathTest)
		CompareBoolean(result, true, t)
	}

	t.Log("5 Quotes Updated!")
}

func TestDelete(t *testing.T) {

	DatabaseTest := db.OpenDB(DatabasePathTest)

	for i := 1; i < 5; i++ {
		result := Delete(i, &DatabaseTest, DatabasePathTest)
		CompareBoolean(result, true, t)
	}

	t.Log("5 Quotes Deleted!")
}

func CompareBoolean(result, expected bool, t *testing.T) {
	if result != expected {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}

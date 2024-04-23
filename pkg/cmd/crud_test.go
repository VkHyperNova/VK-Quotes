package cmd

import (
	"testing"
	db "vk-quotes/pkg/db"
)

var DatabasePathTest = "/home/veikko/Desktop/VK-Quotes/database/quotes.json"

func TestAdd(t *testing.T) {

	DatabaseTest := []db.Quotes{}

	TestQuote1 := []string{"Quote 1", "Top G", "English"}
	TestQuote2 := []string{"Quote 2", "Top A", "RUSSIAN"}
	TestQuote3 := []string{"Quote 3", "Top B", "ESTONIAN"}
	TestQuote4 := []string{"Quote 4", "Top C", "LATVIAN"}
	TestQuote5 := []string{"Quote 5", "Top D", "KAUKAAASIAN"}
	TestQuote6 := []string{"Quote 6", "Top D", "KAUKAAASIAN"}
	TestQuote7 := []string{"Quote 7", "Top D", "KAUKAAASIAN"}
	TestQuote8 := []string{"Quote 8", "Top D", "KAUKAAASIAN"}
	TestQuote9 := []string{"Quote 9", "Top D", "KAUKAAASIAN"}
	TestQuote10 := []string{"Quote 10", "Top D", "KAUKAAASIAN"}

	TestQuoteArray := [][]string{TestQuote1, TestQuote2, TestQuote3, TestQuote4, TestQuote5, TestQuote6, TestQuote7, TestQuote8, TestQuote9, TestQuote10}

	for _, testquote := range TestQuoteArray {
		result := Create(testquote, &DatabaseTest, DatabasePathTest)
		CompareBoolean(result, true, t)
	}

	t.Log("10 Quotes Added!")
}

func TestUpdate(t *testing.T) {

	DatabaseTest := db.OpenDB(DatabasePathTest)

	TestUpdateQuote := []string{"Quote Updated Successfully", "Top G Updated", "Language Updated"}

	for i := 5; i < 11; i++ {
		result := Update(i, TestUpdateQuote, &DatabaseTest, DatabasePathTest)
		CompareBoolean(result, true, t)
	}

	t.Log("5 Quotes Updated!")
}

func TestDelete(t *testing.T) {

	DatabaseTest := db.OpenDB(DatabasePathTest)

	for i := 1; i < 5; i++ {
		Create([]string{"Quote For Deleting", "Top G Del", "English Del"}, &DatabaseTest, DatabasePathTest)
	}

	for i := 11; i < 15; i++ {
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

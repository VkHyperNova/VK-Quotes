package cmd

import (
	"testing"
	db "vk-quotes/pkg/db"
)

var DatabasePathTest = "/home/veikko/Desktop/VK-Quotes/database/quotes.json"

func TestAdd(t *testing.T) {

	t.Log("Testing Adding...")

	DatabaseTest := []db.Quotes{}

	TestQuote1 := []string{"Quote 1", "Top G", "English"}

	result := Create(TestQuote1,&DatabaseTest, DatabasePathTest)
	CompareBoolean(result, true, t)

	TestQuote2 := []string{"Quote 2", "Top A", "RUSSIAN"}

	result = Create(TestQuote2,&DatabaseTest, DatabasePathTest)
	CompareBoolean(result, true, t)

	TestQuote3 := []string{"Quote 2", "Top B", "ESTONIAN"}

	result = Create(TestQuote3,&DatabaseTest, DatabasePathTest)
	CompareBoolean(result, true, t)

	TestQuote4 := []string{"Quote 4", "Top C", "LATVIAN"}

	result = Create(TestQuote4,&DatabaseTest, DatabasePathTest)
	CompareBoolean(result, true, t)

	TestQuote5 := []string{"Quote 5", "Top D", "KAUKAAASIAN"}

	result = Create(TestQuote5,&DatabaseTest, DatabasePathTest)
	CompareBoolean(result, true, t)

}

func TestUpdate(t *testing.T) {

	t.Log("Testing Updating...")

	DatabaseTest := db.OpenDB(DatabasePathTest)

	TestUpdateQuote := []string{"Quote Updated", "Top GG", "English FTW"}

	result := Update(1, TestUpdateQuote, &DatabaseTest, DatabasePathTest)
	CompareBoolean(result, true, t)

	result = Update(2, TestUpdateQuote, &DatabaseTest, DatabasePathTest)
	CompareBoolean(result, true, t)
}

func TestDelete(t *testing.T) {

	t.Log("Testing Deleting...")

	DatabaseTest := db.OpenDB(DatabasePathTest)

	expected := true

	result := Delete(3, &DatabaseTest, DatabasePathTest)
	CompareBoolean(result, expected, t)
}

func CompareBoolean(result, expected bool, t *testing.T) {
	if result != expected {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}

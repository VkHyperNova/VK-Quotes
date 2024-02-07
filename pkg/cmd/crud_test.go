package cmd

import (
	"fmt"
	"testing"
	db "vk-quotes/pkg/db"
)

var DatabasePathTest = "/home/veikko/Desktop/VK-Quotes/database/quotes.json"

func TestAdd(t *testing.T) {

	t.Log("Testing Adding...")

	DatabaseTest := []db.Quotes{}

	expected := true

	// One word
	result := Add(1, "Quote", "Author", "English", &DatabaseTest)
	CompareBoolean(result, expected, t)

	// Multiple words
	result = Add(2, "Multiple Words", "Author Author", "English Russian", &DatabaseTest)
	CompareBoolean(result, expected, t)

	db.SaveDB(&DatabaseTest, DatabasePathTest)
	fmt.Println(DatabaseTest)

}
func TestUpdate(t *testing.T) {

	t.Log("Testing Updating...")
	
	// Load Database
	DatabaseTest := db.LoadDB(DatabasePathTest)

	expected := true

	// One Word
	result := Update(0, "Updated", "Updated", "Updated", DatabasePathTest,&DatabaseTest)
	CompareBoolean(result, expected, t)

	// Multiple Words
	result = Update(1, "Everything works as expected", "Author", "Language", DatabasePathTest,&DatabaseTest)
	CompareBoolean(result, expected, t)

	// Save and Print
	db.SaveDB(&DatabaseTest, DatabasePathTest)
	fmt.Println(DatabaseTest)

}

func TestDelete(t *testing.T) {

	t.Log("Testing Deleting...")

	// Load Database
	DatabaseTest := db.LoadDB(DatabasePathTest)

	expected := true

	result := Delete(0, DatabasePathTest, &DatabaseTest)
	CompareBoolean(result, expected, t)

	// Save Database
	db.SaveDB(&DatabaseTest, DatabasePathTest)
	fmt.Println(DatabaseTest)
}

func CompareBoolean(result, expected bool, t *testing.T) {
	if result != expected {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}

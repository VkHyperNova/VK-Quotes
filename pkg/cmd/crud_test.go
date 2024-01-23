package cmd

import (
	"fmt"
	"testing"
	db "vk-quotes/pkg/db"
	"vk-quotes/pkg/util"
)

var DatabasePathTest = "/home/veikko/Desktop/VK-Quotes/database/quotes.json"



func GetQuotesDetails() db.Quotes {
    return db.Quotes{ID: 1, QUOTE: "Quote", AUTHOR: "Author", LANGUAGE: "English", DATE: util.GetFormattedDate()}
}

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

	// For Updating
	result = Add(3, "This Quote is for updating", "Author Author", "English Russian", &DatabaseTest)
	CompareBoolean(result, expected, t)

	// For Deleting
	result = Add(4, "Delete me", "Author", "English", &DatabaseTest)
	CompareBoolean(result, expected, t)

	fmt.Println(DatabaseTest)
}

// func TestUpdate(t *testing.T) {

// 	t.Log("Testing Updating...")
	
// 	result := Update(2, "Updated Quote", "UPDATED Author", "UPDATED Language", DatabasePathTest)
// 	expected := true
// 	CompareBoolean(result, expected, t)
	
// }

// func TestDelete(t *testing.T) {

// 	t.Log("Testing Deleting...")
	
// 	result := Delete(3, DatabasePathTest)
// 	expected := true
// 	CompareBoolean(result, expected, t)
	
// }

// func TestSaveDB(t *testing.T) {
// 	t.Log("Testing SaveDB...")


// }

func CompareBoolean(result, expected bool, t *testing.T) {
	if result != expected {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}


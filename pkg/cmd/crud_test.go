package cmd

import (
	"fmt"
	"testing"
	"vk-quotes/pkg/db"
)

var DatabasePathTest = "/home/veikko/Desktop/VK-Quotes/database/quotes.json"

func TestAdd(t *testing.T) {

	expected := true

	result := Add("Quote", "Author", "English", DatabasePathTest)
	CompareBoolean(result, expected, t)


	// Multiple words
	result = Add("This Quote has multiple words", "Author Author", "English Russian", DatabasePathTest)
	CompareBoolean(result, expected, t)

	// For Updating
	result = Add("This Quote is for updating", "Author Author", "English Russian", DatabasePathTest)
	CompareBoolean(result, expected, t)


	// For Deleting
	result = Add("This Quote gets deleted", "Author Author", "English Russian", DatabasePathTest)
	CompareBoolean(result, expected, t)

}

func TestUpdate(t *testing.T) {
	
	result := Update(3, "UPDATED Quote!", "UPDATED Author", "UPDATED Language", DatabasePathTest)
	expected := true
	CompareBoolean(result, expected, t)

}

func TestDelete(t *testing.T) {
	
	result := Delete(2, DatabasePathTest)
	expected := true
	CompareBoolean(result, expected, t)
}

func CompareBoolean(result, expected bool, t *testing.T) {
	if result != expected {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}

func PrintInput(LastAddID int) {
	fmt.Println("{")
	fmt.Println("id: ", db.DATABASE[db.SearchIndexByID(db.LastItemID)].ID)
	fmt.Println("quote: " + db.DATABASE[db.SearchIndexByID(db.LastItemID)].QUOTE)
	fmt.Println("author: " + db.DATABASE[db.SearchIndexByID(db.LastItemID)].AUTHOR)
	fmt.Println("language: " + db.DATABASE[db.SearchIndexByID(db.LastItemID)].LANGUAGE)
	fmt.Println("date: " + db.DATABASE[db.SearchIndexByID(db.LastItemID)].DATE)
	fmt.Println("},")
}

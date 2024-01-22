package cmd

import (
	"testing"
)

var DatabasePathTest = "/home/veikko/Desktop/VK-Quotes/database/quotes.json"

func TestAdd(t *testing.T) {

	t.Log("Testing Adding...")

	expected := true

	// One word
	result := Add("Quote", "Author", "English", DatabasePathTest)
	CompareBoolean(result, expected, t)

	// Multiple words
	result = Add("Multiple Words", "Author Author", "English Russian", DatabasePathTest)
	CompareBoolean(result, expected, t)

	// For Updating
	result = Add("This Quote is for updating", "Author Author", "English Russian", DatabasePathTest)
	CompareBoolean(result, expected, t)

	// For Deleting
	result = Add("Quote", "Author", "English", DatabasePathTest)
	CompareBoolean(result, expected, t)
}

func TestUpdate(t *testing.T) {

	t.Log("Testing Updating...")
	
	result := Update(3, "Updated Quote", "UPDATED Author", "UPDATED Language", DatabasePathTest)
	expected := true
	CompareBoolean(result, expected, t)
	
}

func TestDelete(t *testing.T) {

	t.Log("Testing Deleting...")
	
	result := Delete(4, DatabasePathTest)
	expected := true
	CompareBoolean(result, expected, t)
	
}

func CompareBoolean(result, expected bool, t *testing.T) {
	if result != expected {
		t.Errorf("Expected %v but got %v", expected, result)
	}
}


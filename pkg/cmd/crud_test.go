package cmd

import (
	"testing"
)

func TestAdd(t *testing.T) {

	// testQuotes = [3]string{"Test quote", "", ""}
	

	testQuote := "Test quote"
	testAuthor := "Test author"
	testLanguage := "Test language"
	DatabasePath := "/home/veikko/Desktop/VK-Quotes/database/quotes.json"

	result := Add(testQuote, testAuthor, testLanguage, DatabasePath)
	expected := "SUCCESS!"

	if result != expected {
		t.Errorf("Expected %q but got %q", expected, result)
	} 

	testQuote2 := "Test quote"
	testAuthor2 := "Test author"
	testLanguage2 := "Test language"

	result2 := Add(testQuote2, testAuthor2, testLanguage2, DatabasePath)
	expected2 := "SUCCESS!"

	if result2 != expected2 {
		t.Errorf("Expected %q but got %q", expected, result)
	}
}

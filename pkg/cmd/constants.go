package cmd

var Version = "1.1"
var DatabasePath = "./database/quotes.json"
var CurrentQuoteIndex = -1
var SuccessMsg = ""
var ErrorMsg = ""
var ReadCount = 1
var AddCount = 0
var UsedIndexes []int
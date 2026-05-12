package config

import (
	"path/filepath"
)

/* All Paths */
var DefaultContent = `{"quotes": []}`
var	file = "quotes.json"
var BaseDB = "QUOTES"
var BaseLocal = "DATABASES"
var	BaseBackup = "/media/veikko/VK DATA/"

var LocalFile = filepath.Join(BaseLocal, BaseDB, file)
var BackupFile = filepath.Join(BaseBackup, BaseLocal, BaseDB, file)




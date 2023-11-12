package database

import (
	"vk-quotes/pkg/dir"
	"vk-quotes/pkg/global"
	"vk-quotes/pkg/util"
)

func LoadQuotesDatabase() []global.Quotes {
	file := dir.ReadFile(global.DatabasePath)
	data := util.GetQuotesArray(file)

	return data
}

package database

import (
	"vk-quotes/pkg/util"
	"vk-quotes/pkg/dir"
	"vk-quotes/pkg/global"
)

func SaveQuoteDatabase() {
	DatabaseAsByte := util.InterfaceToByte(global.DB)
	dir.WriteDataToFile(global.DatabasePath, DatabaseAsByte)
}

package main

import (
	"Readee-Backend/common/config"
	"Readee-Backend/common/database"
)

func main() {
	config.Init()
	database.Init()
}

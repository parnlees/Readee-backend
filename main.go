package main

import (
	"Readee-Backend/common/config"
	"Readee-Backend/common/database"
	"Readee-Backend/fiber"
)

func main() {
	config.Init()
	database.Init()
	fiber.Init()
}

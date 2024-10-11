package main

import (
	"Readee-Backend/common/config"
	"Readee-Backend/common/database"
	"Readee-Backend/endpoint"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	//"Readee-Backend/fiber"
)

func main() {
	config.Init()
	database.Init()
	//fiber.Init()

	// Initialize Fiber app
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET,POST,PUT,DELETE,PATCH",
	}))

	// Register all routes
	endpoint.RegisterRoutes(app)

	// Start the server
	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}

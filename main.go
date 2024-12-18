package main

import (
	"Readee-Backend/common/config"
	"Readee-Backend/common/database"
	"Readee-Backend/endpoint"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	app.Get("/user/:id", func(c *fiber.Ctx) error {
		userID := c.Params("id")

		// Try to fetch the user from cache or database
		userData, source, err := endpoint.GetUserWithCache(userID)
		if err != nil {
			if source == "database" {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status":  "error",
					"message": "User not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to fetch user",
			})
		}

		return c.JSON(fiber.Map{
			"status": "success",
			"source": source,
			"user":   userData,
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	err := app.Listen("0.0.0.0:" + port)
	if err != nil {
		panic(err)
	}
}

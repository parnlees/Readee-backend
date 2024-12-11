package main

import (
	"Readee-Backend/common/config"
	"Readee-Backend/common/database"
	"Readee-Backend/endpoint"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/patrickmn/go-cache"
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

	app.Get("/user/:id", func(c *fiber.Ctx) error {
		userID := c.Params("id")

		// ลองดึงข้อมูลจากแคชก่อน
		cachedUser, found := config.AppCache.Get(userID)
		if found {
			return c.JSON(fiber.Map{
				"status": "success",
				"source": "cache",
				"user":   cachedUser,
			})
		}

		// หากไม่พบข้อมูลในแคช (Cache Miss)
		userData := fmt.Sprintf("User %s", userID) // จำลองการดึงข้อมูล
		config.AppCache.Set(userID, userData, cache.DefaultExpiration)

		return c.JSON(fiber.Map{
			"status": "success",
			"source": "database",
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

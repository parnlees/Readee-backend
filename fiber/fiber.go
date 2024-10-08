package fiber

import (
	cc "Readee-Backend/common"
	"fmt"

	fiber "github.com/gofiber/fiber/v2"
)

func Init() {
	// Custom config
	app := fiber.New(fiber.Config{
		Prefork: false,
		AppName: "Readee",
	})

	// Set up the basic route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello world 🌈") //string
	})

	//JSON
	app.Get("/info", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"msg":  "welcome to Readee 🌈",
			"bool": true,
			"int":  123,
		})
	})

	// Start the Fiber application
	err := app.Listen(*cc.Config.Address)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

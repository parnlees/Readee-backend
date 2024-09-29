package endpoint

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	// Register user-related routes
	app.Get("/users", GetUsers)
	app.Get("/users/:userId", GetUser)

	// Register book-related routes
	// Add your book routes here

	// Register matching-related routes
	// Add your matching routes here
}

package endpoint

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	// Register user-related routes
	// User
	app.Get("/users", GetUsers)
	app.Get("/users/:userId", GetUserSpecific)

	// Register book-related routes
	// Book
	app.Post("/createBook", CreateBook)
	app.Get("/getBook/:BookId", GetBookSpecific)
	app.Get("/getBook", GetBooks)

	// Register matching-related routes
	// Add your matching routes here
}

package endpoint

import (
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	// User
	app.Get("/users", GetUsers)
	app.Get("/users/:userId", GetUserSpecific)

	// Book
	app.Post("/createBook", CreateBook)
	app.Get("/getBook/:BookId", GetBookSpecific)
	app.Get("/getBooks", GetBooks)
	app.Patch("/editBook/:BookId", EditBook)
	app.Delete("/deleteBook/:BookId", DeleteBook)

	// Log
	app.Post("/books/:bookId/like/:userId", LikeBook)
	app.Post("/books/:bookId/unlike/:userId", UnLikeBook)

	//Match
	app.Post("/match", MatchBook)
}

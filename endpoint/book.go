package endpoint

import (
	//"Readee-backend/endpoint"
	"github.com/gofiber/fiber/v2"
)

var book Book

// Create Book
func CreateBook(c *fiber.Ctx) error {
	if err := c.BodyParser(&book); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	//books = append(books, book)
	return c.JSON(book)
}

// Get Book
func GetBook(c *fiber.Ctx) error {
	//id, _ := strconv.Atoi(c.Params("id"))
	//return c.JSON(books[id])
	return c.JSON(book)
}

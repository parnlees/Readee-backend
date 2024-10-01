package endpoint

import (
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"
	"log"

	"github.com/gofiber/fiber/v2"
)

// Swipe to right = like it -> set value 'liked' to 'true' (liked) in Log table
func LikeBook(c *fiber.Ctx) error {
	var bookTable table.Book
	var logTable table.Log
	BookId := c.Params("BookId")

	// Find the book by ID
	if err := database.DB.First(&bookTable, BookId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Book not found"})
	}

	// Create data in log Table by liked == true
	
	// Set the value of 'liked' to 1
	liked := true
	logTable.Liked = &liked

	// Save the updated log to the database
	if err := database.DB.Save(&logTable).Error; err != nil {
		log.Printf("Error updating log: %v", err) // Log the error
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update log"})
	}

	return c.JSON(logTable)
}

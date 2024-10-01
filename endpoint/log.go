package endpoint

import (
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func BoolPointer(b bool) *bool {
	return &b
}

// Swipe to right = like it -> set value 'liked' to 'true' (liked) in Log table
func LikeBook(c *fiber.Ctx) error {
	var logEntry table.Log

	userId, err := strconv.ParseUint(c.Params("userId"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	bookId, err := strconv.ParseUint(c.Params("bookId"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	logEntry.UserLikeId = &userId
	logEntry.BookLikeId = &bookId
	//set true = 1 in database
	logEntry.Liked = BoolPointer(true)

	if err := database.DB.Create(&logEntry).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to log interaction"})
	}

	return c.Status(201).JSON(logEntry)
}

// Swipe to left = unlike it -> set value 'unliked' to 'false' in Log table
func UnLikeBook(c *fiber.Ctx) error {
	var logEntry table.Log

	userId, err := strconv.ParseUint(c.Params("userId"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	bookId, err := strconv.ParseUint(c.Params("bookId"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	logEntry.UserLikeId = &userId
	logEntry.BookLikeId = &bookId
	//set false = 0 in database
	logEntry.Liked = BoolPointer(false)

	if err := database.DB.Create(&logEntry).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to log interaction"})
	}

	return c.Status(201).JSON(logEntry)
}

package endpoint

import (
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"

	"github.com/gofiber/fiber/v2"
)

// Get all history of a user with book details
func GetHistory(c *fiber.Ctx) error {
	userId := c.Params("userId")
	var histories []table.History

	// Load history with associated OwnerBook, MatchedBook, and Owner details
	if err := database.DB.Preload("OwnerBook").Preload("MatchedBook").Preload("Owner").
		Where("owner_id = ? OR matched_user_id = ?", userId, userId).
		Find(&histories).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve trade history"})
	}

	if len(histories) == 0 {
		return c.Status(200).JSON(fiber.Map{"message": "No history found for this user"})
	}

	// Build the response with the required fields
	response := []fiber.Map{}
	for _, h := range histories {
		response = append(response, fiber.Map{
			"user_book_name":            h.OwnerBook.BookName,
			"user_book_picture":         h.OwnerBook.BookPicture,
			"matched_user_book_name":    h.MatchedBook.BookName,
			"matched_user_book_picture": h.MatchedBook.BookPicture,
			"trade_time":                h.TradeTime,
		})
	}

	return c.Status(200).JSON(fiber.Map{"histories": response})
}

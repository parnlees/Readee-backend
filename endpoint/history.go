package endpoint

import (
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"

	"github.com/gofiber/fiber/v2"
)

// Get all history of a user
func GetHistory(c *fiber.Ctx) error {
	userId := c.Params("userId")
	var histories []table.History

	if err := database.DB.Preload("Book").Preload("Owner").
		Where("owner_id = ? OR owner_match_id = ?", userId, userId). // ตรวจสอบว่าผู้ใช้เป็นเจ้าของหรือผู้แลกเปลี่ยน
		Find(&histories).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve trade history"})
	}

	if len(histories) == 0 {
		return c.Status(200).JSON(fiber.Map{"message": "No history found for this user"})
	}

	response := []fiber.Map{}
	for _, h := range histories {
		var ownerBook table.Book
		var tradedBook table.Book
		if err := database.DB.Where("book_id = ?", h.BookMatchId).First(&ownerBook).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve user's book"})
		}

		if err := database.DB.Where("book_id = ?", h.OwnerMatchId).First(&tradedBook).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve matched user's book"})
		}

		response = append(response, fiber.Map{
			"user_book_name":            ownerBook.BookName,
			"user_book_picture":         ownerBook.BookPicture,
			"matched_user_book_name":    tradedBook.BookName,
			"matched_user_book_picture": tradedBook.BookPicture,
			"match_time":                h.MatchTime,
		})
	}

	return c.Status(200).JSON(fiber.Map{"histories": histories})
}

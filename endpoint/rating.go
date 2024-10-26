package endpoint

import (
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"
	"log"

	"github.com/gofiber/fiber/v2"
)

func GetRatingByUserId(c *fiber.Ctx) error {
	userId := c.Params("userId")

	// Check if the userId parameter is valid
	if userId == "" {
		return c.Status(400).JSON(fiber.Map{"error": "UserId is required"})
	}

	var rating table.Rating

	// Fetch the latest rating given by the specified userId
	if err := database.DB.Where("receiver_id = ?", userId).
		Order("created_at DESC").
		First(&rating).Error; err != nil {
		log.Printf("Error fetching rating: %v", err) // Log the error
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch rating"})
	}

	// Return the latest rating as a JSON response
	return c.Status(200).JSON(rating)
}

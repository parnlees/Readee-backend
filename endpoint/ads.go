package endpoint

import (
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"
	"log"

	// "strconv"

	"github.com/gofiber/fiber/v2"
)

// GetAllAds fetches all active ads from the database
func GetAllAds(c *fiber.Ctx) error {
	var banners []table.Banners

	// Fetch active banners
	result := database.DB.Where("is_active = ?", true).Find(&banners)

	// Check for errors in the query
	if result.Error != nil {
		log.Printf("Error fetching banners: %v", result.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch banners",
		})
	}

	// Return empty array if no active banners are found
	if len(banners) == 0 {
		return c.JSON([]table.Banners{})
	}

	// Return the active banners as JSON
	return c.JSON(banners)
}



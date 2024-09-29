package service

import (
	cc "Readee-Backend/common"
	"Readee-Backend/type/table"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetGenres fetches all genres from the database
func GetGenres(c *fiber.Ctx) error {
	var genres []table.Genre

	// Fetch all genres from the database
	result := cc.DB.Find(&genres)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch genres",
		})
	}

	// Return genres as JSON
	return c.JSON(genres)
}

func GetGenreByID(c *fiber.Ctx) error {
	genreID := c.Params("genre_id")               // Get genre_id from the URL
	id, err := strconv.ParseUint(genreID, 10, 64) // Convert it to uint64
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid genre ID",
		})
	}

	var genre table.Genre

	// Fetch the genre from the database
	result := cc.DB.First(&genre, id)

	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Genre not found",
		})
	}

	// Return the genre as JSON
	return c.JSON(genre)
}

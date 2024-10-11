package endpoint

import (
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetUserGenres(c *fiber.Ctx) error {
	var userGenres []table.UserGenres

	if err := database.DB.Find(&userGenres).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve users"})
	}

	return c.JSON(userGenres)
}

func CreateUserGenres(c *fiber.Ctx) error {
	var request struct {
		User_user_id    *uint64  `json:"User_user_id"`
		Genre_genre_ids []uint64 `json:"Genre_genre_id"` // Expect an array of genre IDs
	}

	// Parse the incoming JSON request
	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Check if user ID and genre IDs are provided
	if request.User_user_id == nil || len(request.Genre_genre_ids) == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "User ID or Genre IDs are missing"})
	}

	// Loop through each genre ID and create a separate record
	for _, genreID := range request.Genre_genre_ids {
		newUserGenre := table.UserGenres{
			User_user_id:   request.User_user_id,
			Genre_genre_id: &genreID, // Use the current genre ID
		}

		// Insert the new record into the database
		if err := database.DB.Create(&newUserGenre).Error; err != nil {
			log.Printf("Error creating userGenres: %v", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create userGenres"})
		}
	}

	return c.Status(201).JSON(fiber.Map{"message": "User genres created successfully"})
}

func GetUserGenreByUserID(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("user_user_id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var logs []table.UserGenres
	result := database.DB.Where("user_user_id = ?", userID).Find(&logs)

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch genre"})
	}

	if len(logs) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "No gemre found for this user"})
	}

	return c.JSON(logs)
}

func EditGenre(c *fiber.Ctx) error {
	var request struct {
		User_user_id    *uint64  `json:"User_user_id"`
		Genre_genre_ids []uint64 `json:"Genre_genre_id"` // Expect an array of genre IDs
	}

	// Parse the incoming JSON request
	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Check if user ID and genre IDs are provided
	if request.User_user_id == nil {
		return c.Status(400).JSON(fiber.Map{"error": "User ID is missing"})
	}

	// Delete all existing genres for the specified user
	if err := database.DB.Where("user_user_id = ?", request.User_user_id).Delete(&table.UserGenres{}).Error; err != nil {
		log.Printf("Error deleting user genres: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete user genres"})
	}

	// Check if new genre IDs are provided
	if len(request.Genre_genre_ids) > 0 {
		// Loop through each genre ID and create a separate record
		for _, genreID := range request.Genre_genre_ids {
			newUserGenre := table.UserGenres{
				User_user_id:   request.User_user_id,
				Genre_genre_id: &genreID, // Use the current genre ID
			}

			// Insert the new record into the database
			if err := database.DB.Create(&newUserGenre).Error; err != nil {
				log.Printf("Error creating userGenres: %v", err)
				return c.Status(500).JSON(fiber.Map{"error": "Failed to create userGenres"})
			}
		}
	}

	return c.Status(200).JSON(fiber.Map{"message": "User genres updated successfully"})
}


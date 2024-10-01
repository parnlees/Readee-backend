package endpoint

import (
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// All
func GetUsers(c *fiber.Ctx) error {
	var users []table.User // Use the correct User model from your table package

	// Query the database to get all users
	if err := database.DB.Find(&users).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve users"})
	}

	return c.JSON(users) // Return the users as JSON
}

// Specific
func GetUserSpecific(c *fiber.Ctx) error {
	userId, err := strconv.ParseUint(c.Params("userId"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var user table.User
	// Find user by userId
	if err := database.DB.First(&user, userId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(user)
}

package endpoint

import (
	"strconv"
	//"Readee-backend/endpoint"
	"github.com/gofiber/fiber/v2"
)

// Add your User struct and mock data here (same as in your current code)

func GetUsers(c *fiber.Ctx) error {
	return c.JSON(users)
}

func GetUser(c *fiber.Ctx) error {
	userId, err := strconv.ParseUint(c.Params("userId"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	for _, user := range users {
		if *user.UserId == userId {
			return c.JSON(user)
		}
	}

	return c.Status(404).JSON(fiber.Map{"error": "User not found"})
}

package endpoint

import (
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"
	"Readee-Backend/util"
	"log"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {
	// Struct to parse JSON input
	var input struct {
		EmailOrUsername string `json:"emailOrUsername"`
		Password        string `json:"password"`
	}

	// Parse request body into input struct
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Check user credentials in database
	var user table.User
	if err := database.DB.Where("email = ? OR username = ?", input.EmailOrUsername, input.EmailOrUsername).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Compare hashed password
	if user.Password == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Password null"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(input.Password)); err != nil {
		log.Println("Hash comparison failed:", err) // Add this line for logging
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate token
	token, err := util.GenerateToken(*user.UserId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	// Return the token and userId in the response
	return c.JSON(fiber.Map{"token": token, "userId": user.UserId, "firstname": user.Firstname, "secKey": user.SecKey})
}

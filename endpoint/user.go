package endpoint

import (
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
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

func CreateUser(c *fiber.Ctx) error {
	var users table.User

	if err := c.BodyParser(&users); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*users.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Could not hash password"})
	}
	hashedPasswordStr := string(hashedPassword)
	users.Password = &hashedPasswordStr

	if err := database.DB.Create(&users).Error; err != nil {
		log.Printf("Error creating users: %v", err) // Log the error
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create users"})
	}
	return c.Status(201).JSON(users)
}

// edit user's information
func EditUser(c *fiber.Ctx) error {
	userId, err := strconv.ParseUint(c.Params("userId"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Retrieve the existing user record
	var user table.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Parse request body into a temporary struct to avoid overwriting unintended fields
	var input struct {
		Email       *string `json:"email"`
		Username    *string `json:"username"`
		PhoneNumber *string `json:"phone_number"`
		ProfileUrl  *string `json:"profile_url"`
		Firstname   *string `json:"firstname"`
		Lastname    *string `json:"lastname"`
		Gender      *string `json:"gender"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Update fields if provided in the request
	if input.Email != nil {
		user.Email = input.Email
	}
	if input.Username != nil {
		user.Username = input.Username
	}
	if input.PhoneNumber != nil {
		user.PhoneNumber = input.PhoneNumber
	}
	if input.ProfileUrl != nil {
		user.ProfileUrl = input.ProfileUrl
	}
	if input.Firstname != nil {
		user.Firstname = input.Firstname
	}
	if input.Lastname != nil {
		user.Lastname = input.Lastname
	}
	if input.Gender != nil {
		user.Gender = input.Gender
	}

	// Save the updated user record
	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update user"})
	}

	return c.JSON(user)
}

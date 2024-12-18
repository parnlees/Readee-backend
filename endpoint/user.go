package endpoint

import (
	"Readee-Backend/common/config"
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
	"golang.org/x/crypto/bcrypt"
)

// Cache
func GetUserWithCache(userID string) (*table.User, string, error) {
	// Try to fetch from cache
	cachedUser, found := config.AppCache.Get(userID)
	if found {
		user, ok := cachedUser.(*table.User)
		if ok {
			return user, "cache", nil
		}
		log.Println("Cache error: cached data invalid")
	}

	// Convert userID to uint64
	userId, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return nil, "", err
	}

	// Fetch from database
	var user table.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		return nil, "database", err
	}

	// Cache the result
	config.AppCache.Set(userID, &user, cache.DefaultExpiration)

	return &user, "database", nil
}

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

	return c.JSON(fiber.Map{
		"status": "success",
		"user":   user,
	})
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
		Email         *string `json:"email"`
		Username      *string `json:"username"`
		PhoneNumber   *string `json:"phone_number"`
		ProfileUrl    *string `json:"profile_url"`
		Firstname     *string `json:"firstname"`
		Lastname      *string `json:"lastname"`
		Gender        *string `json:"gender"`
		SecKey        *string `json:"seckey"`
		RecoverPhrase *string `json:"recoverphrase"`
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
	if input.SecKey != nil {
		user.SecKey = input.SecKey
	}
	if input.RecoverPhrase != nil {
		user.RecoverPhrase = input.RecoverPhrase
	}

	// Save the updated user record
	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update user"})
	}

	return c.JSON(user)
}

func CheckUser(c *fiber.Ctx) error {
	type Request struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	var req Request
	var user table.User

	// Parse the request body to get the username and email
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Check if the username already exists
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{"error": "This username is already taken"})
	}

	// Check if the email already exists
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err == nil {
		return c.Status(409).JSON(fiber.Map{"error": "This email is already exists"})
	}

	// If no conflicts, return success message
	return c.JSON(fiber.Map{"message": "Username and email are available"})
}

func GetUserInfoByEmail(c *fiber.Ctx) error {
	type Request struct {
		Email string `json:"email"` // Expecting email in the request body
	}

	type Response struct {
		UserId        *uint64 `json:"user_id"`
		Email         *string `json:"email"`
		Username      *string `json:"username"`
		SecKey        *string `json:"seckey"`
		RecoverPhrase *string `json:"recover_phrase"`
	}

	var req Request
	var user table.User

	// Parse the request body to get the email
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Check if the email exists in the database
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Build the response
	response := Response{
		UserId:        user.UserId,
		Email:         user.Email,
		Username:      user.Username,
		RecoverPhrase: user.RecoverPhrase,
		SecKey:        user.SecKey,
	}

	// Return the user information
	return c.JSON(response)
}

func ResetPassword(c *fiber.Ctx) error {
	type Request struct {
		NewPassword string `json:"new_password"` // Expecting new password in the request body
	}

	userId, err := strconv.ParseUint(c.Params("userId"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var req Request
	// Parse the request body to get the new password
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Validate the new password
	if len(req.NewPassword) < 8 {
		return c.Status(400).JSON(fiber.Map{"error": "Password must be at least 8 characters long"})
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}

	// Retrieve the existing user record
	var user table.User
	if err := database.DB.First(&user, userId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Update the password field
	hashedPasswordStr := string(hashedPassword)
	user.Password = &hashedPasswordStr

	// Save the updated user record
	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update password"})
	}

	return c.JSON(fiber.Map{"message": "Password updated successfully"})
}

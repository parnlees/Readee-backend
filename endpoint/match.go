package endpoint

import (
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"
	"time"

	//"log"
	//"strconv"

	"github.com/gofiber/fiber/v2"
)

// Match two books when 2 user give liked==1 to each other
// Then add value to Match table
// POST method
func MatchBook(c *fiber.Ctx) error {
	var match table.Match

	// Parse the request body into the match struct
	if err := c.BodyParser(&match); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Unable to parse request body"})
	}

	// Check if both users liked each other's books
	var logOwner, logMatched table.Log
	errOwner := database.DB.Where("user_like_id = ? AND book_like_id = ? AND liked = ?", match.OwnerId, match.MatchedBookId, true).First(&logOwner).Error
	errMatched := database.DB.Where("user_like_id = ? AND book_like_id = ? AND liked = ?", match.MatchedUserId, match.OwnerBookId, true).First(&logMatched).Error

	if errOwner != nil || errMatched != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Match conditions not met"})
	}

	// Set match time and trade time (trade time initially null)
	match.MatchTime = new(time.Time)
	*match.MatchTime = time.Now() // Record match time
	match.TradeTime = nil         // Trade time is null for now

	// Insert the match into the database
	if err := database.DB.Create(&match).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create match"})
	}

	// Return the created match as a response
	return c.Status(201).JSON(match)
}

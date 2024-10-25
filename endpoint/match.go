package endpoint

import (
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetMatchBook retrieves all matched books for a given user by their userId.
// It ensures that ownerBookId corresponds to the userId in the request path.

func GetMatchBook(c *fiber.Ctx) error {
	// Convert the userId from the request parameters
	userId, err := strconv.Atoi(c.Params("userId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Begin the transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		log.Println("tx.Error is " + tx.Error.Error())
		return c.Status(500).JSON(fiber.Map{"error": "Failed to begin transaction"})
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Query to get all matches where the user is either the owner or matched user
	var matches []table.Match
	if err := tx.Where("owner_id = ? OR matched_user_id = ?", userId, userId).Find(&matches).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve matches"})
	}

	// Check if there are no matches
	if len(matches) == 0 {
		return c.Status(200).JSON(fiber.Map{"message": "No matches found for this user"})
	}

	// Commit the transaction
	tx.Commit()

	// Build the response to return only the relevant fields
	response := []fiber.Map{}
	for _, match := range matches {
		// Determine if the user is the owner or the matched user
		var ownerBookId, matchedBookId, matchId *uint64
		var matchTime *time.Time

		if *match.OwnerId == uint64(userId) {
			// User is the owner
			ownerBookId = match.OwnerBookId
			matchedBookId = match.MatchedBookId
			matchId = match.MatchId
			matchTime = match.MatchTime
		} else if *match.MatchedUserId == uint64(userId) {
			// User is the matched user, reverse the roles
			ownerBookId = match.MatchedBookId
			matchedBookId = match.OwnerBookId
			matchId = match.MatchId
			matchTime = match.MatchTime
		}

		// Add the result to the response
		response = append(response, fiber.Map{
			"ownerBookId":   ownerBookId,
			"matchedBookId": matchedBookId,
			"matchId":       matchId,
			"matchTime":     matchTime,
		})
	}

	// Return the match details in the response
	return c.Status(200).JSON(fiber.Map{"matches": response})
}

func GetMatchById(c *fiber.Ctx) error {
	var match table.Match
	MatchId := c.Params("MatchId")

	if err := database.DB.First(&match, MatchId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Match not found"})
	}

	return c.JSON(match)
}

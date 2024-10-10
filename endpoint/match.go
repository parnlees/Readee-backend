package endpoint

import (
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"
	"log"

	"github.com/gofiber/fiber/v2"
)

// GetMatchBook retrieves all matched books for a given user by their userId.
// It ensures that ownerBookId corresponds to the userId in the request path.

func GetMatchBook(c *fiber.Ctx) error {
	// Convert the userId from the request parameters
	userId := c.Params("userId")
	if userId == "" {
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
		var ownerBookId, matchedBookId *uint64

		if match.OwnerId == Uint64Pointer(userId) {
			// User is the owner
			ownerBookId = match.OwnerBookId
			matchedBookId = match.MatchedBookId
		} else if match.MatchedUserId == Uint64Pointer(uint64(userId)) {
			// User is the matched user, reverse the roles
			ownerBookId = match.MatchedBookId
			matchedBookId = match.OwnerBookId
		}

		// Add the result to the response
		response = append(response, fiber.Map{
			"ownerBookId":   ownerBookId,
			"matchedBookId": matchedBookId,
		})
	}

	// Return the match details in the response
	return c.Status(200).JSON(fiber.Map{"matches": response})
}

package endpoint

import (
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"
	"time"

	"strconv"

	"github.com/gofiber/fiber/v2"
)

// 1st : pending
func SendTradeRequest(c *fiber.Ctx) error {
	var match table.Match
	matchId := c.Params("matchId")

	// Fetch the match record
	if err := database.DB.First(&match, matchId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Match not found"})
	}

	// Update the status to "pending"
	match.TradeRequestStatus = "pending"
	initiatorId, _ := strconv.ParseUint(c.Params("initiatorId"), 10, 64) // Assume you're passing it in path
	match.RequestInitiatorId = &initiatorId

	if err := database.DB.Save(&match).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to send trade request"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Trade request sent", "match": match})
}

// 2nd : accept
func AcceptTradeRequest(c *fiber.Ctx) error {
	var match table.Match
	matchId := c.Params("matchId")

	// Fetch the match record
	if err := database.DB.First(&match, matchId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Match not found"})
	}

	// Update the status to "accepted"
	match.TradeRequestStatus = "accepted"
	// Set trade time
	match.TradeTime = TimePointer(time.Now())

	if err := database.DB.Save(&match).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to accept trade request"})
	}

	// Set is_traded in book table == 1 for book struct and match's book struct
	var ownerBook table.Book
	if err := database.DB.First(&ownerBook, match.OwnerBookId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Owner's book not found"})
	}
	ownerBook.IsTraded = BoolPointer(true)
	var matchedBook table.Book
	if err := database.DB.First(&matchedBook, match.MatchedBookId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Matched user's book not found"})
	}
	matchedBook.IsTraded = BoolPointer(true)

	// Save both books
	if err := database.DB.Save(&ownerBook).Error; err != nil || database.DB.Save(&matchedBook).Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update books"})
	}

	// update histories table
	var history table.History
	history.OwnerId = match.OwnerId
	history.MatchedUserId = match.MatchedUserId

	history.OwnerBookId = match.OwnerBookId
	history.MatchedBookId = match.MatchedBookId

	history.TradeTime = TimePointer(time.Now())

	if err := database.DB.Create(&history).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create history"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Trade request accepted", "match": match})
}

func TimePointer(t time.Time) *time.Time {
	return &t
}

// 3rd : reject
// delete from match table and no need to set 'liked' in logs table to 0
// because books in Logs table shouldn't show in home page anymore
func RejectTradeRequest(c *fiber.Ctx) error {
	var match table.Match
	matchId := c.Params("matchId")

	// Fetch the match record
	if err := database.DB.First(&match, matchId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Match not found"})
	}

	// Update the status to "rejected"
	match.TradeRequestStatus = "rejected"
	// Set trade time
	match.TradeTime = TimePointer(time.Now())

	if err := database.DB.Save(&match).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to reject trade request"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Trade request rejected", "match": match})
}

func CancelTradeRequest(c *fiber.Ctx) error {
	var match table.Match
	matchId := c.Params("matchId")

	// Fetch the match record
	if err := database.DB.First(&match, matchId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Match not found"})
	}

	// Reset the trade details
	match.TradeTime = nil
	match.TradeRequestStatus = "none"
	match.RequestInitiatorId = nil

	// Save the updated match record
	if err := database.DB.Save(&match).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to cancel trade request"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Trade request canceled", "match": match})
}

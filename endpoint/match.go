package endpoint

import (
	"Readee-Backend/common/config"
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	//"github.com/patrickmn/go-cache"
)

// GetMatchBook retrieves all matched books for a given user by their userId.
// It ensures that ownerBookId corresponds to the userId in the request path.

func GetMatchBook(c *fiber.Ctx) error {
	userId := c.Params("userId")

	// ลองดึงข้อมูลจาก Cache ก่อน
	cacheKey := "user_" + userId
	cachedData, found := config.AppCache.Get(cacheKey)
	if found {
		return c.JSON(fiber.Map{
			"status":  "success",
			"source":  "cache",
			"matches": cachedData,
		})
	}

	// ดึงข้อมูลจากฐานข้อมูล
	var matches []table.Match
	if err := database.DB.Where("owner_id = ? OR matched_user_id = ?", userId, userId).Find(&matches).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve matches"})
	}

	// หากไม่มีข้อมูล
	if len(matches) == 0 {
		return c.Status(200).JSON(fiber.Map{"message": "No matches found for this user"})
	}

	// เก็บข้อมูลใน Cache โดยการตั้งเวลาหมดอายุ เช่น 5 นาที
	cacheExpiration := 5 * time.Minute
	config.AppCache.Set(cacheKey, matches, cacheExpiration)

	return c.JSON(fiber.Map{
		"status":  "success",
		"source":  "database",
		"matches": matches,
	})
}

func GetMatchById(c *fiber.Ctx) error {
	var match table.Match
	MatchId := c.Params("MatchId")

	if err := database.DB.First(&match, MatchId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Match not found"})
	}

	return c.JSON(match)
}

// Un match - DeleteMatch deletes a match record from the database.
// no need cache
func DeleteMatch(c *fiber.Ctx) error {
	matchId, err := strconv.ParseUint(c.Params("matchId"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid match ID"})
	}

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

	var match table.Match
	if err := tx.First(&match, matchId).Error; err != nil {
		tx.Rollback()
		return c.Status(404).JSON(fiber.Map{"error": "Match not found"})
	}

	if err := tx.Delete(&match).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete match"})
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}
	return c.Status(200).JSON(fiber.Map{"message": "Match deleted successfully"})
}

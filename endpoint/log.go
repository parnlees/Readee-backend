package endpoint

import (
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"
	"log"
	"time"

	"strconv"

	"github.com/gofiber/fiber/v2"
)

func BoolPointer(b bool) *bool {
	return &b
}

func Uint64Pointer(u uint64) *uint64 {
	return &u
}

func LikeBook(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("userId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}
	bookId, err := strconv.Atoi(c.Params("bookId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid book ID"})
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

	var existingLog table.Log
	if err := tx.Where("user_like_id = ? AND book_like_id = ?", userId, bookId).First(&existingLog).Error; err == nil {
		tx.Rollback()
		return c.Status(200).JSON(fiber.Map{"error": "This log already exists"})
	}

	//var Book table.Book
	newLog := table.Log{
		UserLikeId: Uint64Pointer(uint64(userId)),
		BookLikeId: Uint64Pointer(uint64(bookId)),
		Liked:      BoolPointer(true),
		CreatedAt:  TimePointer(time.Now()),
	}

	if err := tx.Create(&newLog).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create log"})
	}
	log.Println("newLog is ", newLog)

	var mutualLog table.Log
	// mutualLog.UserLikeId, Book.BookId
	if err := tx.Where("user_like_id = ? AND book_like_id = ? AND liked = ?", mutualLog.UserLikeId, newLog.BookLikeId, true).First(&mutualLog).Error; err == nil {
		log.Println("mutualLog is ", mutualLog)

		var existingMatch table.Match
		if err := tx.Where("owner_id = ? AND matched_user_id = ?", userId, mutualLog.UserLikeId).First(&existingMatch).Error; err == nil {
			tx.Rollback()
			return c.Status(200).JSON(fiber.Map{"error": "This match already exists"})
		}

		newMatch := table.Match{
			OwnerId:       Uint64Pointer(uint64(userId)),
			MatchedUserId: mutualLog.UserLikeId,
			OwnerBookId:   Uint64Pointer(uint64(bookId)),
			MatchedBookId: mutualLog.BookLikeId,
		}
		if err := tx.Create(&newMatch).Error; err != nil {
			tx.Rollback() // Rollback หากสร้างแมทช์ไม่สำเร็จ
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create match"})
		}
		//return c.Status(201).JSON(fiber.Map{"message": "Match created successfully"})
		log.Printf("Match created: %+v\n", newMatch)
	}
	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Log created successfully"})
}

func UnLikeBook(c *fiber.Ctx) error {
	var logEntry table.Log

	userId, err := strconv.ParseUint(c.Params("userId"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	bookId, err := strconv.ParseUint(c.Params("bookId"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	logEntry.UserLikeId = &userId
	logEntry.BookLikeId = &bookId
	//set false = 0 in database
	logEntry.Liked = BoolPointer(false)

	if err := database.DB.Create(&logEntry).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to log interaction"})
	}

	return c.Status(201).JSON(logEntry)
}

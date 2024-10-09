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
	if err := tx.Where("liker_id = ? AND book_like_id = ?", userId, bookId).First(&existingLog).Error; err == nil {
		tx.Rollback()
		return c.Status(200).JSON(fiber.Map{"error": "This log already exists"})
	}

	newLog := table.Log{
		LikerId:    Uint64Pointer(uint64(userId)),
		BookLikeId: Uint64Pointer(uint64(bookId)),
		Liked:      BoolPointer(true),
		CreatedAt:  TimePointer(time.Now()),
	}

	if err := tx.Create(&newLog).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create log"})
	}
	log.Println("newLog is ", newLog)

	var book table.Book
	if err := tx.Preload("Owner").Where("book_id = ?", *newLog.BookLikeId).First(&book).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Book not found"})
	}
	log.Println("book is ", book)
	// เจ้าของหนังสือ
	bookOwnerId := *book.OwnerId

	// Get all books liked by the book owner
	var likedBooks []table.Log
	if err := tx.Where("liker_id = ? AND liked = ?", bookOwnerId, true).Find(&likedBooks).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve liked books"})
	}
	log.Println("likedBooks are ", likedBooks)

	// check ว่ามีการไลค์จากอีกฝั่งไหม
	for range likedBooks {
		var mutualLogs []table.Log
		// Find all logs where the owner of the book has liked any books of the user
		if err := tx.Joins("JOIN books ON books.book_id = logs.book_like_id").
			Where("logs.liker_id = ? AND books.owner_id = ? AND logs.liked = ?", bookOwnerId, userId, true).
			Find(&mutualLogs).Error; err != nil {
			log.Println("Error fetching mutual likes: ", err)
			continue
		}

		for _, mutualLog := range mutualLogs {
			log.Println("mutualLog is ", mutualLog)

			var existingMatch table.Match
			if err := tx.Where("(owner_book_id = ? AND matched_book_id = ?) OR (owner_book_id = ? AND matched_book_id = ?)",
				bookId, *mutualLog.BookLikeId, *mutualLog.BookLikeId, bookId).
				First(&existingMatch).Error; err == nil {
				log.Println("Match already exists: ", existingMatch)
				continue
			}

			// Create the match
			newMatch := table.Match{
				OwnerId:       Uint64Pointer(uint64(userId)),
				MatchedUserId: mutualLog.LikerId,
				OwnerBookId:   mutualLog.BookLikeId,
				MatchedBookId: Uint64Pointer(uint64(bookId)),
			}
			if err := tx.Create(&newMatch).Error; err != nil {
				tx.Rollback()
				return c.Status(500).JSON(fiber.Map{"error": "Failed to create match"})
			}
			log.Printf("Match created: %+v\n", newMatch)
		}
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

	logEntry.LikerId = &userId
	logEntry.BookLikeId = &bookId
	//set false = 0 in database
	logEntry.Liked = BoolPointer(false)

	if err := database.DB.Create(&logEntry).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to log interaction"})
	}

	return c.Status(201).JSON(logEntry)
}

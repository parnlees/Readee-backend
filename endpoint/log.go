package endpoint

import (
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func BoolPointer(b bool) *bool {
	return &b
}

// Swipe to right = like it -> set value 'liked' to 'true' (liked) in Log table
func LikeBook(c *fiber.Ctx) error {
	var logEntry table.Log

	// ดึง userId และ bookId จากพารามิเตอร์
	userId, err := strconv.ParseUint(c.Params("userId"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	bookId, err := strconv.ParseUint(c.Params("bookId"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid book ID"})
	}

	// บันทึกการกดถูกใจ (like) ของผู้ใช้ปัจจุบันใน Log table
	logEntry.UserLikeId = &userId
	logEntry.BookLikeId = &bookId
	logEntry.Liked = BoolPointer(true)

	if err := database.DB.Create(&logEntry).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to log interaction"})
	}

	// ค้นหาเจ้าของหนังสือเล่มนี้
	var ownerBook table.Book
	if err := database.DB.Where("book_id = ?", bookId).First(&ownerBook).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to find book owner"})
	}

	// Find all books of the owner (A) that B liked
	var booksLikedByB []table.Log
	if err := database.DB.Where("user_like_id = ? AND liked = ?", ownerBook.OwnerId, true).Find(&booksLikedByB).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to find books liked by user B"})
	}

	// Check if A liked any books of B (e.g. bookId == 9)
	var likedBooksByA []table.Log
	if err := database.DB.Where("user_like_id = ? AND liked = ?", userId, true).Find(&likedBooksByA).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to find books liked by user A"})
	}

	// Create matches for all books liked by B and matched with the books liked by A
	matches := []table.Match{}
	for _, bookLikedByB := range booksLikedByB {
		for _, bookLikedByA := range likedBooksByA {
			// Create a match for each combination of books liked by B and A
			for i := 0; i < len(booksLikedByB); i++ {
				for j := 0; j < len(likedBooksByA); j++ {
					match := table.Match{
						OwnerId:            ownerBook.OwnerId,       // Owner of the book
						MatchedUserId:      &userId,                 // Current user who liked the book
						OwnerBookId:        bookLikedByB.BookLikeId, // Book that B liked
						MatchedBookId:      bookLikedByA.BookLikeId, // Book that A liked
						MatchTime:          TimePointer(time.Now()),
						TradeTime:          nil,
						TradeRequestStatus: "none",
					}

					// Save the match to the database
					if err := database.DB.Create(&match).Error; err != nil {
						return c.Status(500).JSON(fiber.Map{"error": "Failed to create match"})
					}
					matches = append(matches, match)
				}
			}
		}
	}

	if len(matches) > 0 {
		return c.Status(201).JSON(fiber.Map{"message": "Matches created successfully", "matches": matches})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Like logged successfully", "log": logEntry})
}

// Swipe to left = unlike it -> set value 'unliked' to 'false' in Log table
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

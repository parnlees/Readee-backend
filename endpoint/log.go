package endpoint

import (
	"Readee-Backend/common/config"
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"
	"fmt"
	"log"
	"time"

	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
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

	// Check if the like already exists in the cache
	cacheKey := fmt.Sprintf("user_%d_likes", userId)
	cachedLikes, found := config.AppCache.Get(cacheKey)
	if found {
		// If the like already exists in the cache, check if it's already present
		if cachedLikes, ok := cachedLikes.([]table.Log); ok {
			for _, log := range cachedLikes {
				if *log.BookLikeId == uint64(bookId) {
					return c.Status(200).JSON(fiber.Map{"error": "This log already exists"})
				}
			}
		} else {
			return c.Status(500).JSON(fiber.Map{"error": "Cache data format mismatch"})
		}
	}

	// Query if log exists in the database
	var existingLog table.Log
	if err := tx.Where("liker_id = ? AND book_like_id = ?", userId, bookId).First(&existingLog).Error; err == nil {
		tx.Rollback()
		return c.Status(200).JSON(fiber.Map{"error": "This log already exists"})
	}

	// Create the new log for the like
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

	// Check if the like was created successfully
	// Add new like to cache
	if found {
		if cachedLikes, ok := cachedLikes.([]table.Log); ok {
			cachedLikes = append(cachedLikes, newLog)
			config.AppCache.Set(cacheKey, cachedLikes, cache.DefaultExpiration)
		} else {
			return c.Status(500).JSON(fiber.Map{"error": "Cache data format mismatch"})
		}
	} else {
		cachedLikes = []table.Log{newLog}
		config.AppCache.Set(cacheKey, cachedLikes, cache.DefaultExpiration)
	}

	// Use cache for checking match
	matchCacheKey := fmt.Sprintf("user_%d_matches", userId)
	cachedMatches, matchFound := config.AppCache.Get(matchCacheKey)
	if matchFound {
		if cachedMatches, ok := cachedMatches.([]table.Match); ok {
			for _, cachedMatch := range cachedMatches {
				if *cachedMatch.OwnerBookId == uint64(bookId) {
					log.Println("Match found in cache: ", cachedMatch)
					// Return cached match (optional, based on your needs)
				}
			}
		} else {
			return c.Status(500).JSON(fiber.Map{"error": "Cache data format mismatch for matches"})
		}
	} else {
		// If no cache, query the database for mutual matches
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

		// Check if there are any mutual likes
		for range likedBooks {
			var mutualLogs []table.Log
			// Find mutual likes
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

				// Cache the new match
				if matchFound {
					if cachedMatches, ok := cachedMatches.([]table.Match); ok {
						cachedMatches = append(cachedMatches, newMatch)
						config.AppCache.Set(matchCacheKey, cachedMatches, cache.DefaultExpiration)
					} else {
						return c.Status(500).JSON(fiber.Map{"error": "Cache data format mismatch"})
					}
				} else {
					cachedMatches = []table.Match{newMatch}
					config.AppCache.Set(matchCacheKey, cachedMatches, cache.DefaultExpiration)
				}
			}
			log.Println("likedBooks are ", likedBooks)
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Log created successfully"})
}

func UnLikeBooks(c *fiber.Ctx) error {
	log.Println("-----------------UnLikeBooks")
	// [{"bookId":59,"status":"unliked"},{"bookId":52,"status":"unliked"}]
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

func GetLogsByUserID(c *fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("liker_id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var logs []table.Log
	result := database.DB.Where("liker_id = ?", userID).Find(&logs)

	if result.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch logs"})
	}

	if len(logs) == 0 {
		return c.Status(404).JSON(fiber.Map{"error": "No logs found for this user"})
	}

	return c.JSON(logs)
}

// unlike no need cache
func UnlikeLogs(c *fiber.Ctx) error {
	// Parse bookLikeId and likerId from the URL parameters
	bookLikeId, err := strconv.ParseUint(c.Params("bookLikeId"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid bookLikeId"})
	}

	likerId, err := strconv.ParseUint(c.Params("likerId"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid likerId"})
	}

	// Retrieve the log entry from the database
	var logEntry table.Log
	result := database.DB.Where("book_like_id = ? AND liker_id = ?", bookLikeId, likerId).First(&logEntry)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return c.Status(404).JSON(fiber.Map{"error": "Log entry not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve log entry"})
	}

	// Set Liked to false
	logEntry.Liked = BoolPointer(false)

	// Save the updated log entry in the database
	if err := database.DB.Save(&logEntry).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update log entry"})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Log updated successfully", "logEntry": logEntry})
}

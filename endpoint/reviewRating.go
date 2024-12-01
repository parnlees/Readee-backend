package endpoint

import (
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Float64Pointer converts a float64 value to a pointer
func Float64Pointer(v float64) *float64 {
	return &v
}

func SubmitRatingAndReview(c *fiber.Ctx) error {
	var req struct {
		GiverId    uint64 `json:"giver_id" validate:"required"`
		ReceiverId uint64 `json:"receiver_id" validate:"required"`
		NewScore   uint64 `json:"new_score" validate:"required,min=1,max=5"`
		TextReview string `json:"text_review" validate:"required"`
	}

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// ตรวจสอบว่าคะแนนต้องอยู่ในช่วง 1-5 เท่านั้น
	if req.NewScore < 1 || req.NewScore > 5 {
		return c.Status(400).JSON(fiber.Map{"error": "Score must be between 1 and 5"})
	}

	// ตรวจสอบว่าเรตติ้งใหม่ถูกให้ไปที่ผู้ใช้ที่ถูกต้อง
	if req.GiverId == req.ReceiverId {
		return c.Status(400).JSON(fiber.Map{"error": "You cannot rate yourself"})
	}

	// สร้าง Review ใหม่
	review := table.Review{
		GiverId:    &req.GiverId,
		ReceiverId: &req.ReceiverId,
		TextReview: req.TextReview,
		CreatedAt:  TimePointer(time.Now()),
	}

	if err := database.DB.Create(&review).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to submit review"})
	}

	// query num_trade in rating table of that user
	var totalNumRate uint64
	if err := database.DB.Model(&table.Rating{}).
		Where("receiver_id = ?", req.ReceiverId).
		Select("COALESCE(COUNT(*), 0) as total_num_rate").
		Scan(&totalNumRate).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get number of ratings"})
	}
	var totalScore uint64
	if err := database.DB.Model(&table.Rating{}).
		Where("receiver_id = ?", req.ReceiverId).
		Select("COALESCE(SUM(score), 0) as total_score").
		Scan(&totalScore).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get total score"})
	}
	totalNumRate++
	log.Printf("totalNumRate: %d", totalNumRate)
	totalScore += req.NewScore
	log.Printf("totalScore: %d", totalScore)
	var newRating = int(math.Round(float64(totalScore) / float64(totalNumRate)))

	rating := table.Rating{
		ReviewId:   review.ReviewId,
		GiverId:    &req.GiverId,
		ReceiverId: &req.ReceiverId,
		Rating:     Uint64Pointer(uint64(newRating)),
		Score:      Uint64Pointer(req.NewScore),
		NumRate:    Uint64Pointer(totalNumRate),
		CreatedAt:  TimePointer(time.Now()),
		UpdatedAt:  TimePointer(time.Now()),
	}

	if err := database.DB.Create(&rating).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to submit rating"})
	}

	return c.Status(201).JSON(fiber.Map{
		"message":        "Rating and review submitted successfully",
		"average_rating": rating.Rating,
	})
}

// Get average rating of a user
func GetAverageRating(c *fiber.Ctx) error {
	userId := c.Params("userId")

	var rating table.Rating
	if err := database.DB.Where("receiver_id = ?", userId).Order("created_at desc").First(&rating).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Rating not found"})
	}

	return c.Status(200).JSON(fiber.Map{
		"average_rating": rating.Rating,
	})
}

// Get review that user got
func GetReceivedReviewsAndRatings(c *fiber.Ctx) error {
	userId := c.Params("userId")

	// Step 1: Retrieve reviews where the user is the receiver
	var reviews []table.Review
	err := database.DB.Preload("Giver").Preload("Receiver").Where("receiver_id = ?", userId).Find(&reviews).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve reviews"})
	}
	if len(reviews) == 0 {
		return c.Status(200).JSON(fiber.Map{"message": "No reviews found for this user"})
	}

	// Step 2: Prepare response with each review and its rating
	response := []fiber.Map{}
	for _, r := range reviews {
		// ตรวจสอบว่า Receiver และ Giver ไม่เป็น nil
		if r.Receiver == nil || r.Giver == nil {
			log.Println("Skipping review with missing Giver or Receiver")
			continue
		}

		// Retrieve rating โดยใช้ ReviewId
		var rating table.Rating
		err := database.DB.Where("review_id = ?", r.ReviewId).Find(&rating).Error
		if err != nil {
			log.Println("No rating found for review_id =", r.ReviewId)
			continue
		}

		// Step 3: Append review and rating to the response
		response = append(response, fiber.Map{
			"user_receiver": r.Receiver.Username,
			"review":        r.TextReview,
			"score":         rating.Score,
			"giver_name":    r.Giver.Username,
			"giver_picture": r.Giver.ProfileUrl,
			"created_at":    r.CreatedAt,
		})
	}

	// Step 4: Return the complete response
	return c.Status(200).JSON(fiber.Map{"reviews": response})
}

// Get rating and review that user gave
func GetGivenReviewsAndRatingsWithTradedBooks(c *fiber.Ctx) error {
	userId := c.Params("userId")
	log.Println("userId:", userId)

	var reviews []table.Review
	// Preload "Receiver" to access receiver information
	if err := database.DB.Preload("Receiver").Where("giver_id = ?", userId).Find(&reviews).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve reviews"})
	}

	// Check if there are reviews
	if len(reviews) == 0 {
		log.Println("No reviews found for user:", userId)
		return c.Status(200).JSON(fiber.Map{"message": "No reviews found for this user"})
	}

	response := []fiber.Map{}
	for _, r := range reviews {
		var rating table.Rating
		var match table.Match

		// Retrieve rating using review_id
		if err := database.DB.Where("review_id = ?", r.ReviewId).First(&rating).Error; err != nil {
			log.Println("No rating found for review_id =", r.ReviewId)
			continue
		}

		// Retrieve match and preload MatchedBook
		err := database.DB.Preload("MatchedBook").Where("matched_user_id = ? AND owner_id = ?", *r.ReceiverId, *r.GiverId).First(&match).Error
		if err != nil {
			// If no match found, use nil for book details
			log.Printf("No match found for matched_user_id = %d and owner_id = %d", *r.ReceiverId, *r.GiverId)
			response = append(response, fiber.Map{
				"receiver_name":         r.Receiver.Username,
				"receiver_picture":      r.Receiver.ProfileUrl,
				"receiver_book_picture": nil,
				"receiver_book_name":    nil,
				"score":                 rating.Score,
				"review":                r.TextReview,
				"created_at":            r.CreatedAt,
			})
			continue
		}

		// Append response for each review and rating
		response = append(response, fiber.Map{
			"receiver_name":         r.Receiver.Username,
			"receiver_picture":      r.Receiver.ProfileUrl,
			"receiver_book_picture": match.MatchedBook.BookPicture,
			"receiver_book_name":    match.MatchedBook.BookName,
			"score":                 rating.Score,
			"review":                r.TextReview,
			"created_at":            r.CreatedAt,
		})
	}

	log.Println("Final response:", response)
	return c.Status(200).JSON(fiber.Map{"reviews": response})
}

// Handler function for fetching the rating and review by giverId and receiverId
func GetReviewRating(c *fiber.Ctx) error {
	// Parse giverId and receiverId from URL params
	giverId, err := strconv.ParseUint(c.Params("giverId"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid giver ID"})
	}

	receiverId, err := strconv.ParseUint(c.Params("receiverId"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid receiver ID"})
	}

	// Fetch the review from the database
	var review table.Review
	if err := database.DB.Where("giver_id = ? AND receiver_id = ?", giverId, receiverId).First(&review).Error; err != nil {
		if err.Error() == "record not found" {
			return c.Status(404).JSON(fiber.Map{"error": "Review not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch review"})
	}

	// Fetch the rating for the review
	var rating table.Rating
	if err := database.DB.Where("giver_id = ? AND receiver_id = ?", giverId, receiverId).First(&rating).Error; err != nil {
		if err.Error() == "record not found" {
			return c.Status(404).JSON(fiber.Map{"error": "Rating not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch rating"})
	}

	// Return the fetched rating and review
	return c.JSON(fiber.Map{
		"rating": rating.Score,
		"review": review.TextReview,
	})
}

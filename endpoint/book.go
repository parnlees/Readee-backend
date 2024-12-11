package endpoint

import (
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Create Book
func CreateBook(c *fiber.Ctx) error {
	var book table.Book

	if err := c.BodyParser(&book); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	if err := database.DB.Create(&book).Error; err != nil {
		log.Printf("Error creating book: %v", err) // Log the error
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create book"})
	}
	return c.Status(201).JSON(book)
}

// Get Book specific
func GetBookSpecific(c *fiber.Ctx) error {
	var book table.Book
	BookId := c.Params("BookId")

	// Find the book by ID
	if err := database.DB.First(&book, BookId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Book not found"})
	}

	return c.JSON(book)
}

// Get all books
func GetBooks(c *fiber.Ctx) error {
	var books []table.Book

	// Query the database to get all books
	if err := database.DB.Find(&books).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve users"})
	}

	return c.JSON(books)
}

// Edit book information
func EditBook(c *fiber.Ctx) error {
	var book table.Book
	BookId := c.Params("BookId")

	// Check if the book exists
	if err := database.DB.First(&book, "book_id = ?", BookId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Book not found"})
	}

	// Parse request body into the book struct
	if err := c.BodyParser(&book); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Update the book information
	if err := database.DB.Model(&book).Where("book_id = ?", BookId).Updates(&book).Error; err != nil {
		log.Printf("Error updating book: %v", err) // Proper logging
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update book"})
	}

	return c.Status(200).JSON(book)
}

// Delete book
func DeleteBook(c *fiber.Ctx) error {
	var book table.Book
	BookId := c.Params("BookId")

	// Check if the book exists
	if err := database.DB.First(&book, "book_id = ?", BookId).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Book not found"})
	}

	// Delete the book by ID
	if err := database.DB.Delete(&table.Book{}, "book_id = ?", BookId).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete book"})
	}

	// Return success response
	return c.Status(200).JSON(fiber.Map{"message": "Book deleted successfully"})
}

func GetBookByOwnerId(c *fiber.Ctx) error {
	var books []table.Book
	OwnerId := c.Params("OwnerId")

	// Query the database to find books by the specified OwnerId
	if err := database.DB.Where("owner_id = ?", OwnerId).Find(&books).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to retrieve books for the specified owner"})
	}

	// If no books are found, return a message indicating this
	if len(books) == 0 {
		return c.Status(404).JSON(fiber.Map{"message": "No books found for this owner"})
	}

	return c.JSON(books)
}

// -----------------------------------------------------------
// Pagination
func getBooksForUser(c *fiber.Ctx) error {
	log.Println("Request received for /books/recommendations/:userId")
	log.Printf("Incoming request: %s %s", c.Method(), c.Path())

	// ตรวจสอบ userID
	userID := c.Params("userID")
	if userID == "" {
		log.Println("Error: userID is required")
		return sendResponse(c, nil, "userID is required", false)
	}
	log.Printf("Received userID: %s", userID)

	if _, err := strconv.Atoi(userID); err != nil {
		log.Printf("Error: Invalid userID: %s", err)
		return sendResponse(c, nil, "Invalid userID", false)
	}

	// ตรวจสอบ query parameters
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page := (offset / limit) + 1
	random := c.Query("random", "false") == "true"

	log.Printf("Pagination - Offset: %d, Limit: %d, Page: %d, Random: %v", offset, limit, page, random)

	var books []table.Book
	// var userGenres []table.UserGenres
	// var likedBooks []table.Log

	// ดึง Books
	query := database.DB.Table("books").
		Where("owner_id != ?", userID).
		Where("is_traded = false").
		Where("is_reported = false").
		Where("book_id NOT IN (?)", database.DB.Table("logs").Select("book_like_id").Where("liker_id = ?", userID)).
		Where("genre_id IN (?)", database.DB.Table("user_genres").Select("genre_genre_id").Where("user_user_id = ?", userID)).
		Where("book_id NOT IN (?)", database.DB.Table("matches").Select("owner_book_id").Where("matched_user_id = ?", userID)).
		Where("book_id NOT IN (?)", database.DB.Table("matches").Select("matched_book_id").Where("owner_id = ?", userID)).
		Offset(offset).
		Limit(limit)

	if random {
		query = query.Order("RAND()")
	}

	// ดึงข้อมูลหนังสือ
	if err := query.Find(&books).Error; err != nil {
		log.Printf("Error fetching books: %s", err)
		return handleError(c, err, "Failed to fetch books", 500)
	}

	// คำนวณ Pagination
	pagination, err := calculatePagination("books", limit, page, database.DB)
	if err != nil {
		log.Printf("Error calculating pagination: %s", err)
		return handleError(c, err, "Failed to calculate pagination", 500)
	}

	// print all of book in 1 stack
	log.Printf("Books: %v", books)
	// ส่ง Response กลับ
	return sendResponse(c, fiber.Map{
		"books":      books,
		"pagination": pagination,
	}, "Books retrieved successfully", true)
}

func getReportBook(c *fiber.Ctx) error {
	// Parse userId from the route parameters
	userId := c.Params("userId")

	// Initialize a slice to store the books
	var books []table.Book

	// Query the database for books where IsReported is true and match the OwnerId
	if err := database.DB.Where("is_reported = ? AND owner_id = ?", true, userId).Find(&books).Error; err != nil {
		log.Printf("Failed to fetch reported books for user %s: %v", userId, err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch reported books"})
	}

	// Check if any books were found
	// if len(books) == 0 {
	//     return c.Status(404).JSON(fiber.Map{"message": "No reported books found for this user"})
	// }

	// Return the reported books
	return c.Status(200).JSON(books)
}

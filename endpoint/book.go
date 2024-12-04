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

// Pagination
func getBooksForUser(c *fiber.Ctx) error {
	userID := c.Query("userID")
	if userID == "" {
		return sendResponse(c, nil, "userID is required", false)
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return sendResponse(c, nil, "Invalid userID", false)
	}

	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page := (offset / limit) + 1
	random := c.Query("random", "false") == "true"

	var userGenres []uint64
	if err := database.DB.Table("user_genres").
		Where("user_user_id = ?", userIDInt).
		Pluck("genre_genre_id", &userGenres).Error; err != nil {
		return handleError(c, err, "Failed to fetch user genres", 500)
	}

	if len(userGenres) == 0 {
		return sendResponse(c, fiber.Map{"books": []table.Book{}}, "No genres found for user", true)
	}

	var likedBookIDs []uint64
	//likedBookIDs = cache.Get("user_liked_books_" + userID)
	if err := database.DB.Table("logs").
		Where("liker_id = ?", userID).
		Pluck("book_like_id", &likedBookIDs).Error; err != nil {
		return handleError(c, err, "Failed to fetch liked books", 500)
	}

	if len(likedBookIDs) == 0 {
		likedBookIDs = []uint64{0}
	}

	var books []table.Book
	query := database.DB.Table("books").
		Where("genre_id IN (?)", userGenres).
		Where("owner_id != ?", userID).
		Where("is_traded = false").
		Where("book_id NOT IN (?)", likedBookIDs).
		Offset(offset).
		Limit(limit)

	if random {
		query = query.Order("RAND()")
	}

	if err := query.Find(&books).Error; err != nil {
		return handleError(c, err, "Failed to fetch books", 500)
	}

	pagination, err := calculatePagination("books", limit, page, database.DB)
	if err != nil {
		return handleError(c, err, "Failed to calculate pagination", 500)
	}

	return sendResponse(c, fiber.Map{
		"books":      books,
		"pagination": pagination,
	}, "Books retrieved successfully", true)
}

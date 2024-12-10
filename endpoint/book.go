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

// Query books for the main page
func getBooksForUser(c *fiber.Ctx) error {
	userID := c.Params("userId")
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	random := c.Query("random", "false") == "true"
	log.Println("Query Parameters:", c.Query("userID"), c.Query("offset"), c.Query("limit"), c.Query("random"))

	// Fetch user's preferred genres
	var userGenres []uint64
	if err := database.DB.Table("user_genres").
		Where("user_user_id = ?", userID).
		Pluck("genre_genre_id", &userGenres).Error; err != nil {
		log.Println("Failed to fetch user genres:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch user genres"})
	}

	if len(userGenres) == 0 {
		log.Println("No genres found for user:", userID)
		return c.Status(200).JSON(fiber.Map{
			"books":   []table.Book{},
			"message": "No genres found for user",
		})
	}

	// Fetch liked book IDs to exclude
	var likedBookIDs []uint64
	if err := database.DB.Table("logs").
		Where("liker_id = ?", userID).
		Pluck("book_like_id", &likedBookIDs).Error; err != nil {
		log.Println("Failed to fetch liked books:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch liked books"})
	}

	// If no liked books are found, use a dummy value to avoid SQL issues
	if len(likedBookIDs) == 0 {
		likedBookIDs = []uint64{0} // Dummy value to prevent SQL errors
	}

	// Fetch books matching the user's preferred genres
	var books []table.Book
	query := database.DB.Table("books").
		Where("genre_id IN (?)", userGenres).      // Filter by user's preferred genres
		Where("owner_id != ?", userID).            // Exclude books owned by the user
		Where("is_traded = false").                // Exclude books that are already traded
		Where("book_id NOT IN (?)", likedBookIDs). // Exclude books the user has liked
		Offset(offset).
		Limit(limit)

	// Apply random sorting if requested
	if random {
		query = query.Order("RAND()")
	}

	if err := query.Find(&books).Error; err != nil {
		log.Println("Failed to fetch books:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch books"})
	}

	// Log the retrieved books for debugging
	log.Printf("Fetched %d books for user %s\n", len(books), userID)

	return c.Status(200).JSON(fiber.Map{"books": books})
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

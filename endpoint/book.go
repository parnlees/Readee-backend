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

// สร้าง struct เพื่อเก็บข้อมูล genre และจำนวนของมัน
type GenreCount struct {
	GenreId    uint64 `json:"genre_id"`
	TotalBooks int    `json:"total_books"`
	LikedBooks int    `json:"liked_books"`
}

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

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		log.Printf("Error: Invalid userID: %s", err)
		return sendResponse(c, nil, "Invalid userID", false)
	}

	// ตรวจสอบ query parameters
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	limit, _ := strconv.Atoi(c.Query("limit", "20")) // 20 เล่มในหนึ่งครั้ง
	page := (offset / limit) + 1
	random := c.Query("random", "false") == "true"

	log.Printf("Pagination - Offset: %d, Limit: %d, Page: %d, Random: %v", offset, limit, page, random)

	// ดึง books ที่ user like แล้ว
	var likedBookIDs []uint64
	if err := database.DB.Table("logs").
		Where("liker_id = ?", userIDInt).
		Pluck("book_like_id", &likedBookIDs).Error; err != nil {
		log.Printf("Error fetching liked books: %s", err)
		return handleError(c, err, "Failed to fetch liked books", 500)
	}
	log.Printf("Liked book IDs: %v", likedBookIDs)

	if len(likedBookIDs) == 0 {
		likedBookIDs = []uint64{0} // ป้องกัน error กรณี likedBookIDs ว่าง
	}

	// ดึง books ที่ user จะได้เห็นใน stack
	var books []table.Book
	query := database.DB.Table("books").
		Where("book_id NOT IN (?)", likedBookIDs). // Exclude already liked books
		Offset(offset).
		Limit(limit)

	if random {
		query = query.Order("RAND()")
	}

	// ดึง books
	if err := query.Find(&books).Error; err != nil {
		log.Printf("Error fetching books: %s", err)
		return handleError(c, err, "Failed to fetch books", 500)
	}
	log.Printf("Books fetched: %v", books)

	// เช็คจำนวนไลค์จาก 20 เล่มใน stack
	likedBooksCount := 0
	for _, book := range books {
		for _, likedBookID := range likedBookIDs {
			if *book.BookId == likedBookID {
				likedBooksCount++
				break
			}
		}
	}

	// เช็คว่าไลค์มากกว่า 50% หรือไม่ (อย่างน้อย 10 เล่ม)
	if likedBooksCount < 10 {
		log.Printf("User liked less than 50%% of books in stack: %d liked out of %d", likedBooksCount, limit)
		// ส่งผลการถามผู้ใช้ว่าจะเปลี่ยน genre หรือไม่
		return sendResponse(c, nil, "Would you like to change genre?", false)
	}

	// คำนวณ Pagination
	pagination, err := calculatePagination("books", limit, page, database.DB)
	if err != nil {
		log.Printf("Error calculating pagination: %s", err)
		return handleError(c, err, "Failed to calculate pagination", 500)
	}

	// ส่ง Response กลับ
	return sendResponse(c, fiber.Map{
		"books":      books,
		"pagination": pagination,
	}, "Books retrieved successfully", true)
}

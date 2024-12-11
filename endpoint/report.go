package endpoint

import (
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateReport(c *fiber.Ctx) error {
	var report table.Report

	// Parse userId and bookId from the route parameters
	userId, err := strconv.ParseUint(c.Params("userId"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid userId"})
	}

	bookId, err := strconv.ParseUint(c.Params("bookId"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid bookId"})
	}

	var existingReport table.Report
	err = database.DB.Where("user_id = ? AND book_id = ?", userId, bookId).First(&existingReport).Error

	if err == nil {
		// Report exists
		return c.Status(409).JSON(fiber.Map{"message": "Report already exists"})
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Some other error occurred
		log.Printf("Error checking existing report: %v", err) // Log the error
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	// Create a new report object
	report.UserId = &userId
	report.BookId = &bookId
	now := time.Now()
	report.ReportAt = &now

	// Save the report to the database
	if err := database.DB.Create(&report).Error; err != nil {
		log.Printf("Error creating report: %v", err) // Log the error
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create report"})
	}

	var reportCount int64
	if err := database.DB.Model(&table.Report{}).Where("book_id = ?", bookId).Count(&reportCount).Error; err != nil {
		log.Printf("Error counting reports: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to count reports"})
	}

	// If the report count is 3 or more, update the book's is_reported column
	if reportCount >= 3 {
		if err := database.DB.Model(&table.Book{}).Where("book_id = ?", bookId).Updates(map[string]interface{}{
			"is_reported":  true,
			"expired_date": time.Now(),
		}).Error; err != nil {
			log.Printf("Error updating book's is_reported column: %v", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update book status"})
		}
	}

	// Return the created report as a response
	return c.Status(201).JSON(report)
}

func GetReportByBookID(c *fiber.Ctx) error {
	// Parse bookId from the route parameters
	bookIdParam := c.Params("bookId")
	bookId, err := strconv.ParseUint(bookIdParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid bookId"})
	}

	// Fetch reports for the given BookId
	var reports []table.Report
	if err := database.DB.Where("book_id = ?", bookId).Find(&reports).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch reports"})
	}

	// Return the fetched reports
	return c.Status(fiber.StatusOK).JSON(reports)
}

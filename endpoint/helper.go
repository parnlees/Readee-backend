package endpoint

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Pagination struct {
	Next          int `json:"next"`            // หน้าถัดไป
	Previous      int `json:"previous"`        // หน้าก่อนหน้า
	RecordPerPage int `json:"record_per_page"` // จำนวนเรคอร์ดต่อหน้า
	CurrentPage   int `json:"current_page"`    // หน้าปัจจุบัน
	TotalPage     int `json:"total_page"`      // จำนวนหน้าทั้งหมด
	TotalRecords  int `json:"total_records"`   // จำนวนเรคอร์ดทั้งหมด
}

func calculatePagination(table string, limit, page int, db *gorm.DB) (Pagination, error) {
	var totalRecords int64

	// คำนวณจำนวนเรคอร์ดทั้งหมดในตาราง
	if err := db.Table(table).Count(&totalRecords).Error; err != nil {
		return Pagination{}, err
	}

	totalPages := (int(totalRecords) + limit - 1) / limit // คำนวณจำนวนหน้าทั้งหมด
	previous := page - 1
	if previous < 1 {
		previous = 0
	}
	next := page + 1
	if next > totalPages {
		next = 0
	}

	return Pagination{
		Next:          next,
		Previous:      previous,
		RecordPerPage: limit,
		CurrentPage:   page,
		TotalPage:     totalPages,
		TotalRecords:  int(totalRecords),
	}, nil
}

func handleError(c *fiber.Ctx, err error, message string, statusCode int) error {
	log.Printf("Error: %s - %v\n", message, err) // Log ข้อผิดพลาด
	return c.Status(statusCode).JSON(fiber.Map{
		"error":   message,
		"details": err.Error(),
	})
}

type APIResponse struct {
	Data    interface{} `json:"data"`    // ข้อมูลที่ต้องการส่งกลับ
	Message string      `json:"message"` // ข้อความ
	Success bool        `json:"success"` // สถานะความสำเร็จ
}

func sendResponse(c *fiber.Ctx, data interface{}, message string, success bool) error {
	return c.JSON(APIResponse{
		Data:    data,
		Message: message,
		Success: success,
	})
}

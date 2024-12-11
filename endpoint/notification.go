package endpoint

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Notification struct represents a notification record in the database
type Notification struct {
	NotiId         uint64    `gorm:"primaryKey;autoIncrement"`
	NotiMessage    string    `gorm:"type:VARCHAR(256);not null"`
	NotiType       string    `gorm:"type:VARCHAR(50);not null"` // e.g., "trade_request", "trade_confirmation", "match_found"
	NotiSenderId   uint64    `gorm:"not null;index"`             // Sender ID (FK to User table)
	NotiReceiverId uint64    `gorm:"not null;index"`             // Receiver ID (FK to User table)
	BookId         uint64    `gorm:"not null;index"`             // Book ID related to the notification
	SendAt         time.Time `gorm:"autoCreateTime"`
}

// CreateNotification creates a new notification record in the database
func CreateNotification(db *gorm.DB, senderId uint64, receiverId uint64, bookId uint64, notiType, message string) (*Notification, error) {
	notification := Notification{
		NotiMessage:    message,
		NotiType:       notiType,
		NotiSenderId:   senderId,
		NotiReceiverId: receiverId,
		BookId:         bookId,
		SendAt:         time.Now(),
	}

	if err := db.Create(&notification).Error; err != nil {
		return nil, err
	}

	return &notification, nil
}

// GetNotificationsForUser retrieves all notifications for a specific user
func GetNotificationsForUser(db *gorm.DB, receiverId uint64, notiType string) ([]map[string]interface{}, error) {
	var notifications []map[string]interface{}

	// Query with joins to get additional details like sender name and book title
	query := db.Table("notifications").
		Select(`
			notifications.noti_id,
			notifications.noti_message,
			notifications.noti_type,
			notifications.send_at,
			users.name AS sender_name,
			books.title AS book_title
		`).
		Joins("JOIN users ON users.user_id = notifications.noti_sender_id").
		Joins("JOIN books ON books.book_id = notifications.book_id").
		Where("notifications.noti_receiver_id = ?", receiverId)

	if notiType != "" {
		query = query.Where("notifications.noti_type = ?", notiType)
	}

	if err := query.Order("notifications.send_at DESC").Scan(&notifications).Error; err != nil {
		return nil, err
	}

	return notifications, nil
}

// CreateNotificationHandler handles the creation of a new notification via API
func CreateNotificationHandler(c *fiber.Ctx) error {
	var req struct {
		SenderId   uint64 `json:"sender_id"`
		ReceiverId uint64 `json:"receiver_id"`
		BookId     uint64 `json:"book_id"`
		NotiType   string `json:"noti_type"`
		Message    string `json:"message"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	db, ok := c.Locals("db").(*gorm.DB)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database connection not available"})
	}

	// Create the notification
	notification, err := CreateNotification(db, req.SenderId, req.ReceiverId, req.BookId, req.NotiType, req.Message)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create notification"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":      "Notification created successfully",
		"notification": notification,
	})
}

// GetNotificationsHandler handles retrieving notifications for a user via API
func GetNotificationsHandler(c *fiber.Ctx) error {
	receiverId := c.Params("receiver_id")
	notiType := c.Query("type") // Optional query parameter to filter by notification type

	// Convert receiver ID to uint64
	receiverIdUint, err := strconv.ParseUint(receiverId, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid receiver ID"})
	}

	db, ok := c.Locals("db").(*gorm.DB)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database connection not available"})
	}

	// Fetch notifications
	notifications, err := GetNotificationsForUser(db, receiverIdUint, notiType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch notifications"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"notifications": notifications,
	})
}

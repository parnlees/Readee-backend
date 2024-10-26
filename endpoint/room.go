package endpoint

import (
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func CreateRoom(c *fiber.Ctx) error {
	// Get parameters from URL
	senderIdParam := c.Params("senderId")
	receiverIdParam := c.Params("receiverId")

	// Convert senderId and receiverId to uint64
	senderId, err := strconv.ParseUint(senderIdParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid sender ID"})
	}
	receiverId, err := strconv.ParseUint(receiverIdParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid receiver ID"})
	}

	// Create new Room instance
	newRoom := Room{
		SenderId:   &senderId,
		ReceiverId: &receiverId,
	}

	// Save newRoom to the database
	if err := database.DB.Create(&newRoom).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create room"})
	}

	// Respond with the created room as JSON
	return c.Status(fiber.StatusOK).JSON(newRoom)
}

func GetRoomId(c *fiber.Ctx) error {
	// Get parameters from URL
	senderIdParam := c.Params("senderId")
	receiverIdParam := c.Params("receiverId")

	// Convert senderId and receiverId to uint64
	senderId, err := strconv.ParseUint(senderIdParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid sender ID"})
	}
	receiverId, err := strconv.ParseUint(receiverIdParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid receiver ID"})
	}

	// Find the RoomId with two-way check
	var room table.Room
	if err := database.DB.Where(
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		senderId, receiverId, receiverId, senderId,
	).First(&room).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Room not found"})
	}

	// Respond with the found RoomId
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"roomId": room.RoomId})
}

func GetAllChatByUserId(c *fiber.Ctx) error {
	// Get userId parameter from URL
	userIdParam := c.Params("userId")

	// Convert userId to uint64
	userId, err := strconv.ParseUint(userIdParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	// Find all rooms where userId is either SenderId or ReceiverId
	var rooms []table.Room
	if err := database.DB.Where("sender_id = ? OR receiver_id = ?", userId, userId).Find(&rooms).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve rooms"})
	}

	// If no rooms found, return an empty list
	if len(rooms) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"rooms": []table.Room{}})
	}

	// Respond with the list of rooms
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"rooms": rooms})
}

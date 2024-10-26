package endpoint

import (
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"
	"log"

	"github.com/gofiber/fiber/v2"
)

func CreateMessage(c *fiber.Ctx) error {
	var message table.Message

	// Parse the request body into the message struct
	if err := c.BodyParser(&message); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Ensure that the required fields are present (you can add more validations as needed)
	if message.RoomId == nil || message.SenderId == nil || message.Message == nil {
		return c.Status(400).JSON(fiber.Map{"error": "RoomId, SenderId, and Message are required"})
	}

	// Create the new message in the database
	if err := database.DB.Create(&message).Error; err != nil {
		log.Printf("Error creating message: %v", err) // Log the error
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create message"})
	}

	return c.Status(201).JSON(message)
}


func GetMessagesByRoomId(c *fiber.Ctx) error {
	roomId := c.Params("roomId")

	// Check if the roomId is valid
	if roomId == "" {
		return c.Status(400).JSON(fiber.Map{"error": "RoomId is required"})
	}

	var messages []table.Message

	// Fetch messages from the database for the specified roomId
	if err := database.DB.Where("room_id = ?", roomId).Find(&messages).Error; err != nil {
		log.Printf("Error fetching messages: %v", err) // Log the error
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch messages"})
	}

	// Return the messages as a JSON response
	return c.Status(200).JSON(messages)
}
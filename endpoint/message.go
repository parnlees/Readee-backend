package endpoint

import (
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"

	//"github.com/gorilla/websocket"
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

// var clients = make(map[*websocket.Conn]bool) // Connected clients
var mu sync.Mutex // Mutex to manage concurrent access
var rooms = make(map[string]map[*websocket.Conn]bool)

func Chat(c *websocket.Conn) {
	roomId := c.Params("roomId")
	if roomId == "" {
		log.Println("Room ID is missing.")
		c.Close()
		return
	}

	log.Printf("WebSocket connection established for room: %s", roomId)

	mu.Lock()
	if rooms[roomId] == nil {
		rooms[roomId] = make(map[*websocket.Conn]bool)
	}
	rooms[roomId][c] = true
	mu.Unlock()

	defer func() {
		log.Printf("WebSocket connection closed for room: %s", roomId)
		mu.Lock()
		delete(rooms[roomId], c)
		if len(rooms[roomId]) == 0 {
			delete(rooms, roomId)
		}
		mu.Unlock()
		c.Close()
	}()

	for {
		var message table.Message
		if err := c.ReadJSON(&message); err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Println("WebSocket closed normally.")
			} else {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		mu.Lock()
		for client := range rooms[roomId] {
			if err := client.WriteJSON(message); err != nil {
				log.Printf("Error sending message: %v", err)
				client.Close()
				delete(rooms[roomId], client)
			}
		}
		mu.Unlock()

		if err := database.DB.Create(&message).Error; err != nil {
			log.Printf("Error saving message to database: %v", err)
		}
	}
}

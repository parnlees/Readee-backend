package endpoint

import (
	"Readee-Backend/common/database"
	"Readee-Backend/type/table"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func UploadImage(file *multipart.FileHeader) (string, error) {
	// Retrieve credentials from environment variables
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT_NAME")
	accountKey := os.Getenv("AZURE_STORAGE_ACCOUNT_KEY")

	if accountName == "" || accountKey == "" {
		log.Println("Azure Storage account credentials not set")
		return "", fmt.Errorf("Azure Storage account credentials not set")
	}

	// Create a credential object for Azure
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Printf("Failed to create Azure credential: %v", err)
		return "", err
	}

	// Define the service URL for Azure Blob Storage
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	// Create a Blob Service client
	client, err := azblob.NewClientWithSharedKeyCredential(serviceURL, cred, nil)
	if err != nil {
		log.Printf("Failed to create Azure blob client: %v", err)
		return "", err
	}

	// Define the container and blob name
	containerName := "chat-images"
	blobName := file.Filename

	// Get the container client
	containerClient := client.ServiceClient().NewContainerClient(containerName)

	// Create a blob client for the file upload
	blobClient := containerClient.NewBlockBlobClient(blobName)

	// Open the file for reading
	fileReader, err := file.Open()
	if err != nil {
		log.Printf("Failed to open file: %v", err)
		return "", err
	}
	defer fileReader.Close()

	// Upload the file to Azure Blob Storage
	_, err = blobClient.Upload(context.Background(), fileReader, nil)
	if err != nil {
		log.Printf("Failed to upload file: %v", err)
		return "", err
	}

	// Generate the URL of the uploaded image
	imageUrl := blobClient.URL()

	// Log the URL of the uploaded image
	log.Printf("Image uploaded successfully url = %s", imageUrl)

	// Return the URL of the uploaded image
	return imageUrl, nil
}

func CreateMessage(c *fiber.Ctx) error {
	// Handle file upload for image messages
	file, err := c.FormFile("file")
	var imageUrl *string

	// If there is a file, upload it to Azure
	if err == nil && file != nil {
		// Log the file details for debugging
		log.Printf("Uploading file: %s", file.Filename)

		uploadedURL, err := UploadImage(file) // Upload the image
		if err != nil {
			log.Printf("Error uploading image: %v", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to upload image"})
		}

		// If the image is uploaded successfully, store the URL
		imageUrl = &uploadedURL // Set the image URL
		log.Printf("Image uploaded successfully, URL: %s", *imageUrl)
	} else {
		log.Println("No file received or error in file handling")
	}

	// Parse the rest of the message data
	var message table.Message
	if err := c.BodyParser(&message); err != nil {
		log.Printf("Error parsing message body: %v", err)
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// If there was an image, add the URL to the message
	if imageUrl != nil {
		message.ImageUrl = imageUrl
		log.Printf("Assigned Image URL to message: %s", *imageUrl)
	}

	// Log the message data before saving it
	log.Printf("Saving message: %+v", message)

	// Save the message (including the image URL) to the database
	if err := database.DB.Create(&message).Error; err != nil {
		log.Printf("Error saving message: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create message"})
	}

	// Log the saved message
	log.Printf("Message saved successfully: %+v", message)

	// Return the message data with the image URL
	return c.Status(201).JSON(fiber.Map{
		"imageUrl": message.ImageUrl, // Include the message data
	})
}

func GetMessagesByRoomId(c *fiber.Ctx) error {
	roomId := c.Params("roomId")
	if roomId == "" {
		return c.Status(400).JSON(fiber.Map{"error": "RoomId is required"})
	}

	var messages []table.Message

	if err := database.DB.Where("room_id = ?", roomId).Find(&messages).Error; err != nil {
		log.Printf("Error fetching messages: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch messages"})
	}

	return c.Status(200).JSON(messages)
}

var clients = make(map[*websocket.Conn]bool)

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
				log.Printf("Error broadcasting message: %v", err)
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

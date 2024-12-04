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

	// Create a credential object
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Printf("Failed to create Azure credential: %v", err)
		return "", err
	}

	// Define the service URL
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	// Create a blob service client
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

	// Create a blob client
	blobClient := containerClient.NewBlockBlobClient(blobName)

	// Open the file to upload
	fileReader, err := file.Open()
	if err != nil {
		log.Printf("Failed to open file: %v", err)
		return "", err
	}
	defer fileReader.Close()

	// Upload the file
	_, err = blobClient.Upload(context.Background(), fileReader, nil)
	if err != nil {
		log.Printf("Failed to upload file: %v", err)
		return "", err
	}

	// Return the blob's URL
	return blobClient.URL(), nil
}

func CreateMessage(c *fiber.Ctx) error {
	// Handle file upload for image messages
	file, err := c.FormFile("file")
	var imageUrl *string

	if err == nil && file != nil {
		uploadedURL, err := UploadImage(file)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to upload image"})
		}
		imageUrl = &uploadedURL
	}

	var message table.Message
	if err := c.BodyParser(&message); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Add the image URL to the message if present
	message.ImageUrl = imageUrl

	// Save the message to the database
	if err := database.DB.Create(&message).Error; err != nil {
		log.Printf("Error creating message: %v", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create message"})
	}

	return c.Status(201).JSON(message)
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
var mu sync.Mutex

func Chat(c *websocket.Conn) {
	mu.Lock()
	clients[c] = true
	mu.Unlock()

	defer func() {
		mu.Lock()
		delete(clients, c)
		mu.Unlock()
		c.Close()
	}()

	for {
		var message table.Message
		if err := c.ReadJSON(&message); err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		mu.Lock()
		for client := range clients {
			if err := client.WriteJSON(message); err != nil {
				log.Printf("Error sending message: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}

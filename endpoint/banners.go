package endpoint

import (
	// "Readee-Backend/common/database"
	// "Readee-Backend/type/table"
	// "log"
	"context"
	"fmt"
	"log"
	"mime/multipart"

	//"net/http"
	"os"

	// "strconv"

	// "Readee-Backend/common/database"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	// "github.com/gin-gonic/gin"
	// "github.com/gofiber/fiber/v2"
	// "github.com/patrickmn/go-cache"
)

func UploadImageBanner(file *multipart.FileHeader) (string, error) {
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
	containerName := "ads-banner"
	// Here we generate the file name as banner-01, banner-02, ..., banner-99
	// You can implement a way to track the last used number, or hardcode the file name for testing purposes.
	// For now, assuming the user wants "banner-01" as an example:
	blobName := "banner-01" // You can update this dynamically as needed

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

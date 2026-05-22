package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

var client *azblob.Client
var containerName string

func InitBlob() {
	connStr := os.Getenv("AZURE_STORAGE_CONNECTION_STRING")
	containerName = os.Getenv("AZURE_STORAGE_CONTAINER")

	if connStr == "" || containerName == "" {
		panic("Storage env vars not set")
	}

	var err error
	client, err = azblob.NewClientFromConnectionString(connStr, nil)
	if err != nil {
		panic("Failed to create blob client: " + err.Error())
	}
}

// Upload file → returns public URL
func UploadFile(file multipart.File, filename string) (string, error) {
	ctx := context.Background()

	blobName := fmt.Sprintf("%d-%s", time.Now().Unix(), filename)

	_, err := client.UploadStream(ctx, containerName, blobName, file, nil)
	if err != nil {
		return "", err
	}

	// Construct URL
	accountName := extractAccountName(client.URL())
	url := fmt.Sprintf("https://%s.blob.core.windows.net/%s/%s",
		accountName,
		containerName,
		blobName,
	)

	return url, nil
}

// Delete blob using URL
func DeleteFile(blobURL string) error {
	ctx := context.Background()

	parts := strings.Split(blobURL, "/")
	blobName := parts[len(parts)-1]

	_, err := client.DeleteBlob(ctx, containerName, blobName, nil)
	return err
}

// helper to extract account name
func extractAccountName(serviceURL string) string {
	// https://account.blob.core.windows.net/
	parts := strings.Split(serviceURL, ".")
	return strings.Replace(parts[0], "https://", "", 1)
}

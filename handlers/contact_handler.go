package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/smituz-for-simform/trainee_backend/config"
	"github.com/smituz-for-simform/trainee_backend/models"

	"github.com/gin-gonic/gin"
)

var phoneRegex = regexp.MustCompile(`^[0-9]{10}$`)

func Ready(c *gin.Context) {
	err := config.DB.Ping(context.Background())
	if err != nil {
		c.JSON(500, gin.H{"status": "db not ready"})
		return
	}
	c.JSON(200, gin.H{"status": "ready"})
}

// ✅ GET all contacts (WITH IMAGE)
func GetContacts(c *gin.Context) {
	rows, err := config.DB.Query(context.Background(),
		"SELECT id, name, phone, COALESCE(image_url, '') FROM contacts")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contacts"})
		return
	}
	defer rows.Close()

	var contacts []models.Contact

	for rows.Next() {
		var contact models.Contact
		if err := rows.Scan(&contact.ID, &contact.Name, &contact.Phone, &contact.ImageURL); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading data"})
			return
		}
		contacts = append(contacts, contact)
	}

	c.JSON(http.StatusOK, contacts)
}

// ✅ CREATE contact (WITH IMAGE)
func CreateContact(c *gin.Context) {
	name := c.PostForm("name")
	phone := c.PostForm("phone")

	if name == "" || phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name and phone are required"})
		return
	}

	if !phoneRegex.MatchString(phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone must be exactly 10 digits"})
		return
	}

	// Duplicate check
	var existingID int
	err := config.DB.QueryRow(context.Background(),
		"SELECT id FROM contacts WHERE name=$1", name).Scan(&existingID)

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Contact with this name already exists"})
		return
	}

	// 🔹 HANDLE IMAGE
	var imageURL string
	file, err := c.FormFile("image")

	if err == nil {
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
		savePath := "./uploads/" + filename

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
			return
		}

		imageURL = "/uploads/" + filename
	}

	// Insert
	_, err = config.DB.Exec(context.Background(),
		"INSERT INTO contacts (name, phone, image_url) VALUES ($1, $2, $3)",
		name, phone, imageURL)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create contact"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact created successfully"})
}

// ✅ UPDATE contact (WITH OPTIONAL IMAGE UPDATE)
func UpdateContact(c *gin.Context) {
	id := c.PostForm("id")
	name := c.PostForm("name")
	phone := c.PostForm("phone")

	if id == "" || name == "" || phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID, name and phone are required"})
		return
	}

	if !phoneRegex.MatchString(phone) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone must be exactly 10 digits"})
		return
	}

	// ✅ FIXED: handle NULL safely
	var existingImage string
	err := config.DB.QueryRow(context.Background(),
		"SELECT COALESCE(image_url, '') FROM contacts WHERE id=$1", id).
		Scan(&existingImage)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	// Duplicate name check
	var duplicateID int
	err = config.DB.QueryRow(context.Background(),
		"SELECT id FROM contacts WHERE name=$1 AND id != $2", name, id).
		Scan(&duplicateID)

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Another contact with this name already exists",
		})
		return
	}

	imageURL := existingImage

	// 🔹 HANDLE NEW IMAGE (optional)
	file, err := c.FormFile("image")
	if err == nil {
		// delete old image if exists
		if existingImage != "" {
			os.Remove("." + existingImage)
		}

		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
		savePath := "./uploads/" + filename

		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save image",
			})
			return
		}

		imageURL = "/uploads/" + filename
	}

	// Update
	result, err := config.DB.Exec(context.Background(),
		"UPDATE contacts SET name=$1, phone=$2, image_url=$3 WHERE id=$4",
		name, phone, imageURL, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update contact",
		})
		return
	}

	// ✅ EXTRA SAFETY (good practice)
	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Contact not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Contact updated successfully",
	})
}

// ✅ DELETE contact (WITH IMAGE DELETE)
func DeleteContact(c *gin.Context) {
	id := c.Param("id")

	var imageURL string

	err := config.DB.QueryRow(context.Background(),
		"SELECT COALESCE(image_url, '') FROM contacts WHERE id=$1", id).
		Scan(&imageURL)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	result, err := config.DB.Exec(context.Background(),
		"DELETE FROM contacts WHERE id=$1", id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete contact"})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}

	// delete image if exists
	if imageURL != "" {
		os.Remove("." + imageURL)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact deleted successfully"})
}

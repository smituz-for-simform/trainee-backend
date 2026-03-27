package handlers

import (
	"context"
	"net/http"
	"regexp"

	"github.com/smituz-for-simform/trainee_backend/config"
	"github.com/smituz-for-simform/trainee_backend/models"

	"github.com/gin-gonic/gin"
)

var phoneRegex = regexp.MustCompile(`^[0-9]{10}$`)

// ✅ GET all contacts
func GetContacts(c *gin.Context) {
	rows, err := config.DB.Query(context.Background(), "SELECT id, name, phone FROM contacts")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch contacts",
		})
		return
	}
	defer rows.Close()

	var contacts []models.Contact

	for rows.Next() {
		var contact models.Contact
		if err := rows.Scan(&contact.ID, &contact.Name, &contact.Phone); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error reading data",
			})
			return
		}
		contacts = append(contacts, contact)
	}

	c.JSON(http.StatusOK, contacts)
}

// ✅ CREATE contact
func CreateContact(c *gin.Context) {
	var contact models.Contact

	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Name and phone are required",
		})
		return
	}

	// Validate phone
	if !phoneRegex.MatchString(contact.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Phone must be exactly 10 digits",
		})
		return
	}

	// Check duplicate name
	var existingID int
	err := config.DB.QueryRow(context.Background(),
		"SELECT id FROM contacts WHERE name=$1", contact.Name).Scan(&existingID)

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Contact with this name already exists",
		})
		return
	}

	// Insert
	_, err = config.DB.Exec(context.Background(),
		"INSERT INTO contacts (name, phone) VALUES ($1, $2)",
		contact.Name, contact.Phone)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create contact",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Contact created successfully",
	})
}

// ✅ UPDATE contact
func UpdateContact(c *gin.Context) {
	var contact models.Contact

	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID, name and phone are required",
		})
		return
	}

	// Validate ID
	if contact.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Valid contact ID is required",
		})
		return
	}

	// Validate phone
	if !phoneRegex.MatchString(contact.Phone) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Phone must be exactly 10 digits",
		})
		return
	}

	// Check if contact exists
	var existingID int
	err := config.DB.QueryRow(context.Background(),
		"SELECT id FROM contacts WHERE id=$1", contact.ID).Scan(&existingID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Contact not found",
		})
		return
	}

	// Optional: prevent duplicate name (excluding self)
	var duplicateID int
	err = config.DB.QueryRow(context.Background(),
		"SELECT id FROM contacts WHERE name=$1 AND id != $2",
		contact.Name, contact.ID).Scan(&duplicateID)

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Another contact with this name already exists",
		})
		return
	}

	// Update
	result, err := config.DB.Exec(context.Background(),
		"UPDATE contacts SET name=$1, phone=$2 WHERE id=$3",
		contact.Name, contact.Phone, contact.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update contact",
		})
		return
	}

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

// ✅ DELETE contact
func DeleteContact(c *gin.Context) {
	id := c.Param("id")

	result, err := config.DB.Exec(context.Background(),
		"DELETE FROM contacts WHERE id=$1", id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete contact",
		})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Contact not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Contact deleted successfully",
	})
}

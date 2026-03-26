package handlers

import (
	"context"
	"net/http"

	"github.com/smituz-for-simform/trainee_backend/config"
	"github.com/smituz-for-simform/trainee_backend/models"

	"github.com/gin-gonic/gin"
)

func GetContacts(c *gin.Context) {
	rows, _ := config.DB.Query(context.Background(), "SELECT id, name, phone FROM contacts")

	var contacts []models.Contact

	for rows.Next() {
		var contact models.Contact
		rows.Scan(&contact.ID, &contact.Name, &contact.Phone)
		contacts = append(contacts, contact)
	}

	c.JSON(http.StatusOK, contacts)
}

func CreateContact(c *gin.Context) {
	var contact models.Contact
	c.BindJSON(&contact)

	_, err := config.DB.Exec(context.Background(),
		"INSERT INTO contacts (name, phone) VALUES ($1, $2)",
		contact.Name, contact.Phone)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "Contact added")
}

func UpdateContact(c *gin.Context) {
	var contact models.Contact
	c.BindJSON(&contact)

	_, err := config.DB.Exec(context.Background(),
		"UPDATE contacts SET name=$1, phone=$2 WHERE id=$3",
		contact.Name, contact.Phone, contact.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "Contact updated")
}

func DeleteContact(c *gin.Context) {
	id := c.Param("id") // get id from URL

	result, err := config.DB.Exec(context.Background(),
		"DELETE FROM contacts WHERE id=$1", id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// Check if anything was actually deleted
	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, "Contact not found")
		return
	}

	c.JSON(http.StatusOK, "Contact deleted")
}

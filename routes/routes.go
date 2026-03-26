package routes

import (
	"github.com/smituz-for-simform/trainee_backend/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/get_contacts", handlers.GetContacts)
	r.POST("/add_contact", handlers.CreateContact)
	r.PUT("/update_contact", handlers.UpdateContact)
	r.DELETE("/del_contact/:id", handlers.DeleteContact)
}

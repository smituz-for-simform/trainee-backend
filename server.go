// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// )

// func main() {
// 	http.HandleFunc("/", handler)
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprint(w, "Hello, World!")
// }

package main

import (
	"github.com/smituz-for-simform/trainee_backend/config"
	"github.com/smituz-for-simform/trainee_backend/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"time"

	"github.com/gin-contrib/cors"
)

func main() {
	godotenv.Load()

	config.ConnectDB()

	r := gin.Default()

	r.Static("/uploads", "./uploads")

	// ✅ Add this
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.SetupRoutes(r)

	r.Run(":8080")
}

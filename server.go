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
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/smituz-for-simform/trainee_backend/config"
	"github.com/smituz-for-simform/trainee_backend/routes"
	"github.com/smituz-for-simform/trainee_backend/utils"
)

func main() {
	//  Load .env ONLY for local development
	if os.Getenv("ENV") != "production" {
		godotenv.Load()
	}

	config.ConnectDB()
	config.InitSchema()
	utils.InitBlob()

	r := gin.Default()

	//  Remove static uploads in cloud
	//  Keep only for local dev
	// if os.Getenv("ENV") != "production" {
	// 	r.Static("/uploads", "./uploads")
	// }
	if os.Getenv("WEBSITE_SITE_NAME") == "" {
		godotenv.Load()
	}

	requiredEnvs := []string{
    "DB_HOST",
    "DB_PORT",
    "DB_NAME",
    "DB_USER",
    "DB_PASSWORD",
   	"AZURE_STORAGE_ACCOUNT_NAME",
	"AZURE_STORAGE_CONTAINER",
}

	for _, env := range requiredEnvs {
		if os.Getenv(env) == "" {
			panic("Missing required env: " + env)
		}
	}

	//  Dynamic CORS
	frontendURL := os.Getenv("FRONTEND_URL")

	var origins []string

	if frontendURL == "" {
		origins = []string{"http://localhost:3000"}
	} else if frontendURL == "*" {
		origins = []string{"*"}
	} else {
		raw := strings.Split(frontendURL, ",")

		for _, o := range raw {
			o = strings.TrimSpace(o)
			if o == "" {
				continue // 🔥 skip empty values
			}
			if !strings.HasPrefix(o, "http://") && !strings.HasPrefix(o, "https://") {
				panic("Invalid FRONTEND_URL entry: " + o)
			}
			origins = append(origins, o)
		}
	}

	fmt.Println("FINAL ORIGINS:", origins)

	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.SetupRoutes(r)

	//  Use PORT env var (required for Azure)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}

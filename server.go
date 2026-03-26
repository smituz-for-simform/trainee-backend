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
)

func main() {
	godotenv.Load()

	config.ConnectDB()

	r := gin.Default()
	routes.SetupRoutes(r)

	r.Run(":8080")
}

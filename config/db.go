// package config

// import (
// 	"context"
// 	"os"

// 	"github.com/jackc/pgx/v5"
// )

// var DB *pgx.Conn

//	func ConnectDB() {
//		conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
//		if err != nil {
//			panic(err)
//		}
//		DB = conn
//	}
package config

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to connect to DB: %v\n", err)
	}

	DB = pool

	log.Println("✅ Connected to database")
}

func InitSchema() {
	sqlBytes, err := os.ReadFile("db-init/init.sql")
	if err != nil {
		log.Fatal("Failed to read init.sql:", err)
	}

	query := string(sqlBytes)

	_, err = DB.Exec(context.Background(), query)
	if err != nil {
		log.Fatal("Failed to execute init.sql:", err)
	}

	log.Println("Database schema initialized successfully")
}

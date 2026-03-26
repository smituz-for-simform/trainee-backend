package config

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func ConnectDB() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		panic(err)
	}
	DB = conn
}

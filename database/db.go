package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func InitDB(connString string) {
	var err error
	DB, err = pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal("Database connection error: ", err)
	}

	err = DB.Ping(context.Background())
	if err != nil {
		log.Fatal("Could not ping database: ", err)
	}

	fmt.Println("Successfully connection")
}

func CloseDB() {
	DB.Close(context.Background())
}

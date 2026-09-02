package main

import (
	"log"
	"os"

	"asset-backend/internal/user/domain"
	sharedDB "asset-backend/internal/shared/db"
)

func main() {
	dsn := os.Getenv("USER_DB_DSN")
	if dsn == "" {
		dsn = "root:vishal123@tcp(127.0.0.1:3306)/user_db?charset=utf8mb4&parseTime=True&loc=Local"
	}

	database, err := sharedDB.Connect(dsn)
	if err != nil {
		log.Fatalf("user-service: failed to connect to database: %v", err)
	}

	if err := database.AutoMigrate(&domain.Employee{}); err != nil {
		log.Fatalf("user-service: failed to run migrations: %v", err)
	}

	log.Println("user-service: connected and migrated successfully")
}

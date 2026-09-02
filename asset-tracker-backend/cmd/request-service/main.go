package main

import (
	"log"
	"os"

	"asset-backend/internal/request/domain"
	sharedDB "asset-backend/internal/shared/db"
)

func main() {
	dsn := os.Getenv("REQUEST_DB_DSN")
	if dsn == "" {
		dsn = "root:vishal123@tcp(127.0.0.1:3306)/request_db?charset=utf8mb4&parseTime=True&loc=Local"
	}

	database, err := sharedDB.Connect(dsn)
	if err != nil {
		log.Fatalf("request-service: failed to connect to database: %v", err)
	}

	if err := database.AutoMigrate(&domain.AssetRequest{}, &domain.Approval{}); err != nil {
		log.Fatalf("request-service: failed to run migrations: %v", err)
	}

	log.Println("request-service: connected and migrated successfully")
}

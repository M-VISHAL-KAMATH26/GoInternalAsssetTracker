package main

import (
	"log"
	"os"

	"asset-backend/internal/inventory/domain"
	sharedDB "asset-backend/internal/shared/db"
)

func main() {
	dsn := os.Getenv("INVENTORY_DB_DSN")
	if dsn == "" {
		dsn = "root:vishal123@tcp(127.0.0.1:3306)/inventory_db?charset=utf8mb4&parseTime=True&loc=Local"
	}

	database, err := sharedDB.Connect(dsn)
	if err != nil {
		log.Fatalf("inventory-service: failed to connect to database: %v", err)
	}

	if err := database.AutoMigrate(&domain.Asset{}, &domain.AssetAssignment{}); err != nil {
		log.Fatalf("inventory-service: failed to run migrations: %v", err)
	}

	log.Println("inventory-service: connected and migrated successfully")

	// Gin/gRPC server setup comes in a later phase.
}

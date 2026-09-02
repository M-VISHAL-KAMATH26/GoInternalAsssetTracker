package main

import (
	"log"
	"os"

	"asset-backend/internal/inventory/domain"
	"asset-backend/internal/inventory/handler"
	"asset-backend/internal/inventory/repository"
	sharedDB "asset-backend/internal/shared/db"

	"github.com/gin-gonic/gin"
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

	// Wire up the repository -> handler chain.
	assetRepo := repository.NewAssetRepository(database)
	assetHandler := handler.NewAssetHandler(assetRepo)

	router := gin.Default()
	handler.RegisterAssetRoutes(router, assetHandler)

	port := os.Getenv("INVENTORY_SERVICE_PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("inventory-service: listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("inventory-service: server failed: %v", err)
	}
}
package main

import (
	"log"
	"os"

	"asset-backend/internal/request/client"
	"asset-backend/internal/request/domain"
	"asset-backend/internal/request/handler"
	"asset-backend/internal/request/repository"
	sharedDB "asset-backend/internal/shared/db"

	"github.com/gin-gonic/gin"
)

func main() {
	dsn := os.Getenv("REQUEST_DB_DSN")
	if dsn == "" {
		dsn = "root:password@tcp(127.0.0.1:3306)/request_db?charset=utf8mb4&parseTime=True&loc=Local"
	}

	database, err := sharedDB.Connect(dsn)
	if err != nil {
		log.Fatalf("request-service: failed to connect to database: %v", err)
	}

	if err := database.AutoMigrate(&domain.AssetRequest{}, &domain.Approval{}); err != nil {
		log.Fatalf("request-service: failed to run migrations: %v", err)
	}
	log.Println("request-service: connected and migrated successfully")

	userAddr := os.Getenv("USER_SERVICE_GRPC_ADDR")
	if userAddr == "" {
		userAddr = "localhost:9091"
	}
	userClient, err := client.NewUserClient(userAddr)
	if err != nil {
		log.Fatalf("request-service: failed to connect to user-service: %v", err)
	}

	inventoryAddr := os.Getenv("INVENTORY_SERVICE_GRPC_ADDR")
	if inventoryAddr == "" {
		inventoryAddr = "localhost:9092"
	}
	inventoryClient, err := client.NewInventoryClient(inventoryAddr)
	if err != nil {
		log.Fatalf("request-service: failed to connect to inventory-service: %v", err)
	}

	requestRepo := repository.NewRequestRepository(database)
	approvalRepo := repository.NewApprovalRepository(database)

	requestHandler := handler.NewRequestHandler(requestRepo, userClient)
	approvalHandler := handler.NewApprovalHandler(requestRepo, approvalRepo, inventoryClient)

	router := gin.Default()
	handler.RegisterRequestRoutes(router, requestHandler, approvalHandler)

	port := os.Getenv("REQUEST_SERVICE_PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("request-service: listening on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("request-service: server failed: %v", err)
	}
}
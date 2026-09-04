package main

import (
	"log"
	"net"
	"os"

	"asset-backend/internal/inventory/domain"
	inventorygrpc "asset-backend/internal/inventory/grpc"
	"asset-backend/internal/inventory/handler"
	"asset-backend/internal/inventory/repository"
	sharedDB "asset-backend/internal/shared/db"
	pb "asset-backend/proto/inventory"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
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

	assetRepo := repository.NewAssetRepository(database)

	go func() {
		grpcPort := os.Getenv("INVENTORY_SERVICE_GRPC_PORT")
		if grpcPort == "" {
			grpcPort = "9092"
		}

		listener, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			log.Fatalf("inventory-service: failed to listen on gRPC port: %v", err)
		}

		grpcServer := grpc.NewServer()
		pb.RegisterInventoryServiceServer(grpcServer, inventorygrpc.NewServer(assetRepo))

		log.Printf("inventory-service: gRPC server listening on :%s", grpcPort)
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("inventory-service: gRPC server failed: %v", err)
		}
	}()

	assetHandler := handler.NewAssetHandler(assetRepo)
	router := gin.Default()
	handler.RegisterAssetRoutes(router, assetHandler)

	httpPort := os.Getenv("INVENTORY_SERVICE_PORT")
	if httpPort == "" {
		httpPort = "8081"
	}

	log.Printf("inventory-service: HTTP server listening on :%s", httpPort)
	if err := router.Run(":" + httpPort); err != nil {
		log.Fatalf("inventory-service: HTTP server failed: %v", err)
	}
}
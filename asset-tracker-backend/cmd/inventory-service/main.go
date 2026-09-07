package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"asset-backend/internal/inventory/domain"
	inventorygrpc "asset-backend/internal/inventory/grpc"
	"asset-backend/internal/inventory/handler"
	"asset-backend/internal/inventory/repository"
	sharedDB "asset-backend/internal/shared/db"
	"asset-backend/internal/shared/config"
	"asset-backend/internal/shared/logger"
	"asset-backend/internal/shared/middleware"
	pb "asset-backend/proto/inventory"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadInventoryServiceConfig()
	log := logger.New("inventory-service")

	database, err := sharedDB.Connect(cfg.DBDSN)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	if err := database.AutoMigrate(&domain.Asset{}, &domain.AssetAssignment{}); err != nil {
		log.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}
	log.Info("connected and migrated successfully")

	assetRepo := repository.NewAssetRepository(database)

	grpcServer := grpc.NewServer()
	pb.RegisterInventoryServiceServer(grpcServer, inventorygrpc.NewServer(assetRepo))

	go func() {
		listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
		if err != nil {
			log.Error("failed to listen on gRPC port", "error", err)
			os.Exit(1)
		}
		log.Info("gRPC server listening", "port", cfg.GRPCPort)
		if err := grpcServer.Serve(listener); err != nil {
			log.Error("gRPC server failed", "error", err)
			os.Exit(1)
		}
	}()

	assetHandler := handler.NewAssetHandler(assetRepo)
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	handler.RegisterAssetRoutes(router, assetHandler)

	httpServer := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	go func() {
		log.Info("HTTP server listening", "port", cfg.HTTPPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("HTTP server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = httpServer.Shutdown(ctx)
	grpcServer.GracefulStop()

	sqlDB, err := database.DB()
	if err == nil {
		_ = sqlDB.Close()
	}
	log.Info("shutdown complete")
}
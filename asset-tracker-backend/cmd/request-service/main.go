package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"asset-backend/internal/request/client"
	"asset-backend/internal/request/domain"
	"asset-backend/internal/request/handler"
	"asset-backend/internal/request/repository"
	sharedDB "asset-backend/internal/shared/db"
	"asset-backend/internal/shared/config"
	"asset-backend/internal/shared/logger"
	"asset-backend/internal/shared/mq"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadRequestServiceConfig()
	log := logger.New("request-service")

	database, err := sharedDB.Connect(cfg.DBDSN)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	if err := database.AutoMigrate(&domain.AssetRequest{}, &domain.Approval{}); err != nil {
		log.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}
	log.Info("connected and migrated successfully")

	userClient, err := client.NewUserClient(cfg.UserGRPCAddr)
	if err != nil {
		log.Error("failed to connect to user-service", "error", err)
		os.Exit(1)
	}

	inventoryClient, err := client.NewInventoryClient(cfg.InventoryGRPCAddr)
	if err != nil {
		log.Error("failed to connect to inventory-service", "error", err)
		os.Exit(1)
	}

	publisher, err := mq.NewPublisher(cfg.RabbitURL)
	if err != nil {
		log.Error("failed to connect to rabbitmq", "error", err)
		os.Exit(1)
	}

	requestRepo := repository.NewRequestRepository(database)
	approvalRepo := repository.NewApprovalRepository(database)

	requestHandler := handler.NewRequestHandler(requestRepo, userClient)
	approvalHandler := handler.NewApprovalHandler(requestRepo, approvalRepo, inventoryClient, publisher)

	router := gin.Default()
	handler.RegisterRequestRoutes(router, requestHandler, approvalHandler)

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
	_ = publisher.Close()

	sqlDB, err := database.DB()
	if err == nil {
		_ = sqlDB.Close()
	}
	log.Info("shutdown complete")
}
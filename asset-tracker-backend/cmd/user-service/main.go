package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"asset-backend/internal/shared/config"
	sharedDB "asset-backend/internal/shared/db"
	"asset-backend/internal/shared/logger"
	"asset-backend/internal/shared/middleware"
	userdomain "asset-backend/internal/user/domain"
	usergrpc "asset-backend/internal/user/grpc"
	"asset-backend/internal/user/handler"
	"asset-backend/internal/user/repository"
	pb "asset-backend/proto/user"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadUserServiceConfig()
	log := logger.New("user-service")

	database, err := sharedDB.Connect(cfg.DBDSN)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	if err := database.AutoMigrate(&userdomain.Employee{}); err != nil {
		log.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}
	log.Info("connected and migrated successfully")

	employeeRepo := repository.NewEmployeeRepository(database)
	userServer := usergrpc.NewServer(employeeRepo)

	listener, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		log.Error("failed to listen", "error", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userServer)

	go func() {
		log.Info("gRPC server listening", "port", cfg.GRPCPort)
		if err := grpcServer.Serve(listener); err != nil {
			log.Error("gRPC server failed", "error", err)
			os.Exit(1)
		}
	}()

	authHandler := handler.NewAuthHandler(employeeRepo, cfg.JWTSecret)
	router := gin.Default()
	router.Use(middleware.CORSMiddleware())
	handler.RegisterAuthRoutes(router, authHandler)

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

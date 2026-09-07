package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"asset-backend/internal/shared/config"
	"asset-backend/internal/shared/logger"
	userdomain "asset-backend/internal/user/domain"
	usergrpc "asset-backend/internal/user/grpc"
	"asset-backend/internal/user/repository"
	sharedDB "asset-backend/internal/shared/db"
	pb "asset-backend/proto/user"

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

	// Graceful shutdown: wait for SIGINT/SIGTERM, then stop cleanly.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down gracefully...")
	grpcServer.GracefulStop()

	sqlDB, err := database.DB()
	if err == nil {
		_ = sqlDB.Close()
	}
	log.Info("shutdown complete")
	_ = context.Background()
}
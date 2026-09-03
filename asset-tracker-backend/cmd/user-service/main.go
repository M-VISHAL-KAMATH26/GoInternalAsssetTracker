package main

import (
	"log"
	"net"
	"os"

	userdomain "asset-backend/internal/user/domain"
	usergrpc "asset-backend/internal/user/grpc"
	"asset-backend/internal/user/repository"
	sharedDB "asset-backend/internal/shared/db"
	pb "asset-backend/proto/user"

	"google.golang.org/grpc"
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

	if err := database.AutoMigrate(&userdomain.Employee{}); err != nil {
		log.Fatalf("user-service: failed to run migrations: %v", err)
	}
	log.Println("user-service: connected and migrated successfully")

	employeeRepo := repository.NewEmployeeRepository(database)
	userServer := usergrpc.NewServer(employeeRepo)

	port := os.Getenv("USER_SERVICE_GRPC_PORT")
	if port == "" {
		port = "9091"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("user-service: failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userServer)

	log.Printf("user-service: gRPC server listening on :%s", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("user-service: server failed: %v", err)
	}
}
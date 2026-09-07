package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Load reads a .env file if present (ignored if missing — real
// deployments use real env vars, not .env files) so local dev doesn't
// require exporting everything by hand each terminal session.
func Load() {
	_ = godotenv.Load()
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

type UserServiceConfig struct {
	DBDSN    string
	GRPCPort string
	JWTSecret string
}

func LoadUserServiceConfig() UserServiceConfig {
	Load()
	return UserServiceConfig{
		DBDSN:     getEnv("USER_DB_DSN", "root:vishal123@tcp(127.0.0.1:3306)/user_db?charset=utf8mb4&parseTime=True&loc=Local"),
		GRPCPort:  getEnv("USER_SERVICE_GRPC_PORT", "9091"),
		JWTSecret: getEnv("JWT_SECRET", ""),
	}
}

type InventoryServiceConfig struct {
	DBDSN    string
	HTTPPort string
	GRPCPort string
	JWTSecret string
}

func LoadInventoryServiceConfig() InventoryServiceConfig {
	Load()
	return InventoryServiceConfig{
		DBDSN:     getEnv("INVENTORY_DB_DSN", "root:vishal123@tcp(127.0.0.1:3306)/inventory_db?charset=utf8mb4&parseTime=True&loc=Local"),
		HTTPPort:  getEnv("INVENTORY_SERVICE_PORT", "8081"),
		GRPCPort:  getEnv("INVENTORY_SERVICE_GRPC_PORT", "9092"),
		JWTSecret: getEnv("JWT_SECRET", ""),
	}
}

type RequestServiceConfig struct {
	DBDSN            string
	HTTPPort         string
	JWTSecret        string
	UserGRPCAddr     string
	InventoryGRPCAddr string
	RabbitURL        string
}

func LoadRequestServiceConfig() RequestServiceConfig {
	Load()
	return RequestServiceConfig{
		DBDSN:             getEnv("REQUEST_DB_DSN", "root:vishal123@tcp(127.0.0.1:3306)/request_db?charset=utf8mb4&parseTime=True&loc=Local"),
		HTTPPort:          getEnv("REQUEST_SERVICE_PORT", "8082"),
		JWTSecret:         getEnv("JWT_SECRET", ""),
		UserGRPCAddr:      getEnv("USER_SERVICE_GRPC_ADDR", "localhost:9091"),
		InventoryGRPCAddr: getEnv("INVENTORY_SERVICE_GRPC_ADDR", "localhost:9092"),
		RabbitURL:         getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
	}
}

type NotificationConsumerConfig struct {
	DBDSN     string
	RabbitURL string
}

func LoadNotificationConsumerConfig() NotificationConsumerConfig {
	Load()
	return NotificationConsumerConfig{
		DBDSN:     getEnv("NOTIFICATION_DB_DSN", "root:vishal123@tcp(127.0.0.1:3306)/notification_db?charset=utf8mb4&parseTime=True&loc=Local"),
		RabbitURL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
	}
}
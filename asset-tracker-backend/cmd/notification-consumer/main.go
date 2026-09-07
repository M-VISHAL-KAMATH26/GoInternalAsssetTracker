package main

import (
	"context"
	"log"
	"os"

	notifdomain "asset-backend/internal/notification/domain"
	"asset-backend/internal/notification/repository"
	"asset-backend/internal/notification/service"
	sharedDB "asset-backend/internal/shared/db"
	"asset-backend/internal/shared/mq"
)

func main() {
	dsn := os.Getenv("NOTIFICATION_DB_DSN")
	if dsn == "" {
		dsn = "root:vishal123@tcp(127.0.0.1:3306)/notification_db?charset=utf8mb4&parseTime=True&loc=Local"
	}

	database, err := sharedDB.Connect(dsn)
	if err != nil {
		log.Fatalf("notification-consumer: failed to connect to database: %v", err)
	}

	if err := database.AutoMigrate(&notifdomain.Notification{}); err != nil {
		log.Fatalf("notification-consumer: failed to run migrations: %v", err)
	}
	log.Println("notification-consumer: connected and migrated successfully")

	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	consumer, err := mq.NewConsumer(rabbitURL)
	if err != nil {
		log.Fatalf("notification-consumer: failed to connect to rabbitmq: %v", err)
	}
	defer consumer.Close()

	notifRepo := repository.NewNotificationRepository(database)
	notifier := service.NewNotifier(notifRepo)

	log.Println("notification-consumer: waiting for messages...")
	err = consumer.Consume(mq.ApprovalDecidedQueue, func(body []byte) error {
		return notifier.Handle(context.Background(), body)
	})
	if err != nil {
		log.Fatalf("notification-consumer: consume failed: %v", err)
	}
}
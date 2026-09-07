package main

import (
	"context"
	"os"

	notifdomain "asset-backend/internal/notification/domain"
	"asset-backend/internal/notification/repository"
	"asset-backend/internal/notification/service"
	sharedDB "asset-backend/internal/shared/db"
	"asset-backend/internal/shared/config"
	"asset-backend/internal/shared/logger"
	"asset-backend/internal/shared/mq"
)

func main() {
	cfg := config.LoadNotificationConsumerConfig()
	log := logger.New("notification-consumer")

	database, err := sharedDB.Connect(cfg.DBDSN)
	if err != nil {
		log.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	if err := database.AutoMigrate(&notifdomain.Notification{}); err != nil {
		log.Error("failed to run migrations", "error", err)
		os.Exit(1)
	}
	log.Info("connected and migrated successfully")

	consumer, err := mq.NewConsumer(cfg.RabbitURL)
	if err != nil {
		log.Error("failed to connect to rabbitmq", "error", err)
		os.Exit(1)
	}
	defer consumer.Close()

	notifRepo := repository.NewNotificationRepository(database)
	notifier := service.NewNotifier(notifRepo)

	log.Info("waiting for messages...")
	err = consumer.Consume(mq.ApprovalDecidedQueue, func(body []byte) error {
		return notifier.Handle(context.Background(), body)
	})
	if err != nil {
		log.Error("consume failed", "error", err)
		os.Exit(1)
	}
}
package repository

import (
	"context"
	"errors"
	"time"

	"asset-backend/internal/notification/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrNotificationNotFound = errors.New("notification not found")

type NotificationRepository interface {
	Create(ctx context.Context, n *domain.Notification) error
	ListByRecipient(ctx context.Context, recipientID uuid.UUID) ([]domain.Notification, error)
	MarkSent(ctx context.Context, id uuid.UUID, sentAt time.Time) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, n *domain.Notification) error {
	return r.db.WithContext(ctx).Create(n).Error
}

func (r *notificationRepository) ListByRecipient(ctx context.Context, recipientID uuid.UUID) ([]domain.Notification, error) {
	var notifications []domain.Notification
	err := r.db.WithContext(ctx).Where("recipient_id = ?", recipientID).Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) MarkSent(ctx context.Context, id uuid.UUID, sentAt time.Time) error {
	return r.db.WithContext(ctx).Model(&domain.Notification{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{"sent": true, "sent_at": sentAt}).Error
}
package domain

import (
	"time"

	"github.com/google/uuid"
)

// NotificationType represents the fixed set of notification kinds this
// service can produce.
type NotificationType string

const (
	NotificationTypeEmployeeApproved   NotificationType = "employee_approved"
	NotificationTypeEmployeeRejected   NotificationType = "employee_rejected"
	NotificationTypeProcurementHandoff NotificationType = "procurement_handoff"
)

// Notification is a persisted record of a notification this service
// has produced — gives an audit trail independent of whatever actual
// delivery mechanism (console log, email, Slack) is used underneath.
type Notification struct {
	ID          uuid.UUID        `gorm:"type:char(36);primaryKey"`
	RecipientID uuid.UUID        `gorm:"type:char(36);not null;index"` // employee ID, or a fixed procurement UUID
	Type        NotificationType `gorm:"type:varchar(50);not null"`
	Message     string           `gorm:"type:text;not null"`
	Sent        bool             `gorm:"not null;default:false"`
	SentAt      *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
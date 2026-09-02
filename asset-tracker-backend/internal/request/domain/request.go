package domain

import (
	"time"

	"github.com/google/uuid"
)

// RequestStatus represents the fixed set of states an asset request can be in.
type RequestStatus string

const (
	RequestStatusPending   RequestStatus = "pending"
	RequestStatusApproved  RequestStatus = "approved"
	RequestStatusRejected  RequestStatus = "rejected"
	RequestStatusFulfilled RequestStatus = "fulfilled"
)

// AssetRequest represents an employee's request for a piece of equipment.
// EmployeeID is a plain reference to the User domain — no GORM
// foreign key/association, since User lives in its own database.
type AssetRequest struct {
	ID            uuid.UUID     `gorm:"type:char(36);primaryKey"`
	EmployeeID    uuid.UUID     `gorm:"type:char(36);not null;index"`
	AssetType     string        `gorm:"type:varchar(100);not null"`
	Category      string        `gorm:"type:varchar(100);not null"`
	Justification string        `gorm:"type:text"`
	Status        RequestStatus `gorm:"type:varchar(20);not null;default:pending"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

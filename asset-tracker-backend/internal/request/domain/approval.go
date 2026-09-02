package domain

import (
	"time"

	"github.com/google/uuid"
)

// ApprovalDecision represents the fixed set of outcomes a manager can record.
type ApprovalDecision string

const (
	DecisionApproved ApprovalDecision = "approved"
	DecisionRejected ApprovalDecision = "rejected"
)

// Approval records a manager's decision on an AssetRequest. RequestID
// is a real GORM foreign key since AssetRequest lives in this same
// domain/database. ManagerID is a plain UUID column — no association —
// since the manager record itself lives in the User domain's database.
type Approval struct {
	ID        uuid.UUID        `gorm:"type:char(36);primaryKey"`
	RequestID uuid.UUID        `gorm:"type:char(36);not null;index"`
	Request   AssetRequest     `gorm:"foreignKey:RequestID;references:ID"`
	ManagerID uuid.UUID        `gorm:"type:char(36);not null;index"`
	Decision  ApprovalDecision `gorm:"type:varchar(20);not null"`
	Comment   string           `gorm:"type:text"`
	DecidedAt time.Time        `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

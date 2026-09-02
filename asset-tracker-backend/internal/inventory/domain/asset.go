package domain

import (
	"time"

	"github.com/google/uuid"
)

// AssetStatus represents the fixed set of states an asset can be in.
type AssetStatus string

const (
	AssetStatusAvailable   AssetStatus = "available"
	AssetStatusAssigned    AssetStatus = "assigned"
	AssetStatusRetired     AssetStatus = "retired"
	AssetStatusMaintenance AssetStatus = "maintenance"
)

// Asset represents a single trackable item of equipment owned by the
// Inventory domain (laptops, monitors, badges, etc).
type Asset struct {
	ID           uuid.UUID   `gorm:"type:char(36);primaryKey"`
	Name         string      `gorm:"type:varchar(255);not null"`
	Type         string      `gorm:"type:varchar(100);not null"`
	Category     string      `gorm:"type:varchar(100);not null"`
	SerialNumber string      `gorm:"type:varchar(255);uniqueIndex;not null"`
	Status       AssetStatus `gorm:"type:varchar(20);not null;default:available"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

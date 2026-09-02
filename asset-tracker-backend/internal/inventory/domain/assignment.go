package domain

import (
	"time"

	"github.com/google/uuid"
)

// AssetAssignment records that an asset was handed to an employee.
// AssetID is a real GORM association since Asset lives in this same
// Inventory database. EmployeeID is a plain UUID reference to the
// User domain — no GORM foreign key/association.
type AssetAssignment struct {
	ID         uuid.UUID  `gorm:"type:char(36);primaryKey"`
	AssetID    uuid.UUID  `gorm:"type:char(36);not null;index"`
	Asset      Asset      `gorm:"foreignKey:AssetID;references:ID"`
	EmployeeID uuid.UUID  `gorm:"type:char(36);not null;index"`
	AssignedAt time.Time  `gorm:"not null"`
	ReturnedAt *time.Time // nil while the employee still holds the asset
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

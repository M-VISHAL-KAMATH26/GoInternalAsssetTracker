package domain

import (
	"time"

	"github.com/google/uuid"
)

// EmployeeRole represents the fixed set of roles an employee can have.
type EmployeeRole string

const (
	RoleEmployee EmployeeRole = "employee"
	RoleManager  EmployeeRole = "manager"
	RoleAdmin    EmployeeRole = "admin"
)

// Employee represents a person in the User domain. ManagerID is
// self-referencing within this same domain/database, so it's a real
// GORM association — unlike cross-domain references (e.g. EmployeeID
// used in Inventory or Request), which are plain UUID columns.
type Employee struct {
	ID         uuid.UUID    `gorm:"type:char(36);primaryKey"`
	Name       string       `gorm:"type:varchar(255);not null"`
	Email      string       `gorm:"type:varchar(255);uniqueIndex;not null"`
	Role       EmployeeRole `gorm:"type:varchar(20);not null;default:employee"`
	Department string       `gorm:"type:varchar(100)"`
	ManagerID  *uuid.UUID   `gorm:"type:char(36);index"` // nil if no manager (e.g. top-level exec)
	Manager    *Employee    `gorm:"foreignKey:ManagerID;references:ID"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

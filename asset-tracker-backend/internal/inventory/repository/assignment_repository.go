package repository

import (
	"context"
	"errors"
	"time"

	"asset-backend/internal/inventory/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrAssignmentNotFound = errors.New("assignment not found")

type AssignmentRepository interface {
	Create(ctx context.Context, assignment *domain.AssetAssignment) error
	GetActiveByAsset(ctx context.Context, assetID uuid.UUID) (*domain.AssetAssignment, error)
	ListActiveByEmployee(ctx context.Context, employeeID uuid.UUID) ([]domain.AssetAssignment, error)
	ReturnAsset(ctx context.Context, id uuid.UUID, returnedAt time.Time) error
}

type assignmentRepository struct {
	db *gorm.DB
}

func NewAssignmentRepository(db *gorm.DB) AssignmentRepository {
	return &assignmentRepository{db: db}
}

func (r *assignmentRepository) Create(ctx context.Context, assignment *domain.AssetAssignment) error {
	return r.db.WithContext(ctx).Create(assignment).Error
}

// GetActiveByAsset returns the current (not yet returned) assignment for an asset, if any.
func (r *assignmentRepository) GetActiveByAsset(ctx context.Context, assetID uuid.UUID) (*domain.AssetAssignment, error) {
	var assignment domain.AssetAssignment
	err := r.db.WithContext(ctx).
		Where("asset_id = ? AND returned_at IS NULL", assetID).
		First(&assignment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAssignmentNotFound
		}
		return nil, err
	}
	return &assignment, nil
}

// ListActiveByEmployee returns every asset an employee currently holds.
func (r *assignmentRepository) ListActiveByEmployee(ctx context.Context, employeeID uuid.UUID) ([]domain.AssetAssignment, error) {
	var assignments []domain.AssetAssignment
	err := r.db.WithContext(ctx).
		Preload("Asset").
		Where("employee_id = ? AND returned_at IS NULL", employeeID).
		Find(&assignments).Error
	return assignments, err
}

func (r *assignmentRepository) ReturnAsset(ctx context.Context, id uuid.UUID, returnedAt time.Time) error {
	return r.db.WithContext(ctx).Model(&domain.AssetAssignment{}).
		Where("id = ?", id).
		Update("returned_at", returnedAt).Error
}

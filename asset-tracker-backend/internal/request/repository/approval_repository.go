package repository

import (
	"context"
	"errors"

	"asset-backend/internal/request/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrApprovalNotFound = errors.New("approval not found")

type ApprovalRepository interface {
	Create(ctx context.Context, approval *domain.Approval) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Approval, error)
	GetByRequestID(ctx context.Context, requestID uuid.UUID) (*domain.Approval, error)
}

type approvalRepository struct {
	db *gorm.DB
}

func NewApprovalRepository(db *gorm.DB) ApprovalRepository {
	return &approvalRepository{db: db}
}

func (r *approvalRepository) Create(ctx context.Context, approval *domain.Approval) error {
	return r.db.WithContext(ctx).Create(approval).Error
}

func (r *approvalRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Approval, error) {
	var approval domain.Approval
	if err := r.db.WithContext(ctx).First(&approval, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrApprovalNotFound
		}
		return nil, err
	}
	return &approval, nil
}

func (r *approvalRepository) GetByRequestID(ctx context.Context, requestID uuid.UUID) (*domain.Approval, error) {
	var approval domain.Approval
	if err := r.db.WithContext(ctx).First(&approval, "request_id = ?", requestID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrApprovalNotFound
		}
		return nil, err
	}
	return &approval, nil
}

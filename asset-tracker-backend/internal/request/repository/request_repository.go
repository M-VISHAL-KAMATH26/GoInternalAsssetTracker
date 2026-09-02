package repository

import (
	"context"
	"errors"

	"asset-backend/internal/request/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrRequestNotFound = errors.New("request not found")

type RequestRepository interface {
	Create(ctx context.Context, request *domain.AssetRequest) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.AssetRequest, error)
	ListByEmployee(ctx context.Context, employeeID uuid.UUID) ([]domain.AssetRequest, error)
	ListByStatus(ctx context.Context, status domain.RequestStatus) ([]domain.AssetRequest, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.RequestStatus) error
}

type requestRepository struct {
	db *gorm.DB
}

func NewRequestRepository(db *gorm.DB) RequestRepository {
	return &requestRepository{db: db}
}

func (r *requestRepository) Create(ctx context.Context, request *domain.AssetRequest) error {
	return r.db.WithContext(ctx).Create(request).Error
}

func (r *requestRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.AssetRequest, error) {
	var request domain.AssetRequest
	if err := r.db.WithContext(ctx).First(&request, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRequestNotFound
		}
		return nil, err
	}
	return &request, nil
}

func (r *requestRepository) ListByEmployee(ctx context.Context, employeeID uuid.UUID) ([]domain.AssetRequest, error) {
	var requests []domain.AssetRequest
	err := r.db.WithContext(ctx).Where("employee_id = ?", employeeID).Find(&requests).Error
	return requests, err
}

func (r *requestRepository) ListByStatus(ctx context.Context, status domain.RequestStatus) ([]domain.AssetRequest, error) {
	var requests []domain.AssetRequest
	err := r.db.WithContext(ctx).Where("status = ?", status).Find(&requests).Error
	return requests, err
}

func (r *requestRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.RequestStatus) error {
	return r.db.WithContext(ctx).Model(&domain.AssetRequest{}).
		Where("id = ?", id).
		Update("status", status).Error
}

package repository

import (
	"context"
	"errors"

	"asset-backend/internal/inventory/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrAssetNotFound = errors.New("asset not found")

type AssetRepository interface {
	Create(ctx context.Context, asset *domain.Asset) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Asset, error)
	GetBySerialNumber(ctx context.Context, serial string) (*domain.Asset, error)
	ListByStatus(ctx context.Context, status domain.AssetStatus) ([]domain.Asset, error)
	List(ctx context.Context) ([]domain.Asset, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.AssetStatus) error
	Update(ctx context.Context, asset *domain.Asset) error
}

type assetRepository struct {
	db *gorm.DB
}

func NewAssetRepository(db *gorm.DB) AssetRepository {
	return &assetRepository{db: db}
}

func (r *assetRepository) Create(ctx context.Context, asset *domain.Asset) error {
	return r.db.WithContext(ctx).Create(asset).Error
}

func (r *assetRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Asset, error) {
	var asset domain.Asset
	if err := r.db.WithContext(ctx).First(&asset, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAssetNotFound
		}
		return nil, err
	}
	return &asset, nil
}

func (r *assetRepository) GetBySerialNumber(ctx context.Context, serial string) (*domain.Asset, error) {
	var asset domain.Asset
	if err := r.db.WithContext(ctx).First(&asset, "serial_number = ?", serial).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAssetNotFound
		}
		return nil, err
	}
	return &asset, nil
}

func (r *assetRepository) ListByStatus(ctx context.Context, status domain.AssetStatus) ([]domain.Asset, error) {
	var assets []domain.Asset
	err := r.db.WithContext(ctx).Where("status = ?", status).Find(&assets).Error
	return assets, err
}

func (r *assetRepository) List(ctx context.Context) ([]domain.Asset, error) {
	var assets []domain.Asset
	err := r.db.WithContext(ctx).Find(&assets).Error
	return assets, err
}

func (r *assetRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.AssetStatus) error {
	return r.db.WithContext(ctx).Model(&domain.Asset{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *assetRepository) Update(ctx context.Context, asset *domain.Asset) error {
	return r.db.WithContext(ctx).Save(asset).Error
}

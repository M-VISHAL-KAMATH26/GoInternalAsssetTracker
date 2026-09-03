package repository

import (
	"context"
	"errors"
	"time"

	"asset-backend/internal/inventory/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrAssetNotFound = errors.New("asset not found")
var ErrNoAssetAvailable = errors.New("no asset available for this type/category")

type AssetRepository interface {
	Create(ctx context.Context, asset *domain.Asset) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Asset, error)
	GetBySerialNumber(ctx context.Context, serial string) (*domain.Asset, error)
	ListByStatus(ctx context.Context, status domain.AssetStatus) ([]domain.Asset, error)
	List(ctx context.Context) ([]domain.Asset, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.AssetStatus) error
	Update(ctx context.Context, asset *domain.Asset) error
	ReserveAvailableAsset(ctx context.Context, assetType, category string, employeeID uuid.UUID) (*domain.Asset, error)
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

// ReserveAvailableAsset finds one available asset matching type/category,
// locks it for the duration of the transaction (SELECT ... FOR UPDATE),
// marks it assigned, and records the assignment — all atomically, so
// two concurrent calls can never reserve the same asset.
func (r *assetRepository) ReserveAvailableAsset(ctx context.Context, assetType, category string, employeeID uuid.UUID) (*domain.Asset, error) {
	var asset domain.Asset

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("type = ? AND category = ? AND status = ?", assetType, category, domain.AssetStatusAvailable).
			Order("created_at").
			First(&asset).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrNoAssetAvailable
			}
			return err
		}

		asset.Status = domain.AssetStatusAssigned
		if err := tx.Save(&asset).Error; err != nil {
			return err
		}

		assignment := domain.AssetAssignment{
			ID:         uuid.New(),
			AssetID:    asset.ID,
			EmployeeID: employeeID,
			AssignedAt: time.Now(),
		}
		return tx.Create(&assignment).Error
	})

	if err != nil {
		return nil, err
	}
	return &asset, nil
}
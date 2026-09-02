package handler

import (
	"time"

	"asset-backend/internal/inventory/domain"

	"github.com/google/uuid"
)

// CreateAssetRequest is the expected JSON body for POST /assets.
// Validation tags are intentionally stricter than the DB model —
// e.g. required fields here aren't necessarily NOT NULL at the DB level.
type CreateAssetRequest struct {
	Name         string `json:"name" binding:"required,min=2,max=255"`
	Type         string `json:"type" binding:"required,max=100"`
	Category     string `json:"category" binding:"required,max=100"`
	SerialNumber string `json:"serial_number" binding:"required,max=255"`
}

// UpdateAssetRequest is the expected JSON body for PUT /assets/:id.
// All fields are required here — this is a full replace, matching
// the repository's Update (Save) semantics, not a partial patch.
type UpdateAssetRequest struct {
	Name         string             `json:"name" binding:"required,min=2,max=255"`
	Type         string             `json:"type" binding:"required,max=100"`
	Category     string             `json:"category" binding:"required,max=100"`
	SerialNumber string             `json:"serial_number" binding:"required,max=255"`
	Status       domain.AssetStatus `json:"status" binding:"required,oneof=available assigned retired maintenance"`
}

// AssetResponse is what we send back over HTTP. Deliberately separate
// from domain.Asset so the API shape can diverge from the DB shape
// (e.g. dropping fields, renaming, adding computed values) without
// touching the GORM model.
type AssetResponse struct {
	ID           uuid.UUID          `json:"id"`
	Name         string             `json:"name"`
	Type         string             `json:"type"`
	Category     string             `json:"category"`
	SerialNumber string             `json:"serial_number"`
	Status       domain.AssetStatus `json:"status"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

// toAssetResponse maps a domain.Asset (DB model) to an AssetResponse (API model).
func toAssetResponse(a *domain.Asset) AssetResponse {
	return AssetResponse{
		ID:           a.ID,
		Name:         a.Name,
		Type:         a.Type,
		Category:     a.Category,
		SerialNumber: a.SerialNumber,
		Status:       a.Status,
		CreatedAt:    a.CreatedAt,
		UpdatedAt:    a.UpdatedAt,
	}
}

// toAssetResponseList maps a slice of domain.Asset to a slice of AssetResponse.
func toAssetResponseList(assets []domain.Asset) []AssetResponse {
	responses := make([]AssetResponse, 0, len(assets))
	for _, a := range assets {
		responses = append(responses, toAssetResponse(&a))
	}
	return responses
}
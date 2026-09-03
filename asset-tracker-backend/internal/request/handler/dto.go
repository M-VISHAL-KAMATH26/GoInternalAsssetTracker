package handler

import (
	"time"

	"asset-backend/internal/request/domain"

	"github.com/google/uuid"
)

// CreateRequestRequest is the expected JSON body for POST /requests.
// Note: no employee_id field here — the employee is identified from
// the JWT, never trusted from the request body.
type CreateRequestRequest struct {
	AssetType     string `json:"asset_type" binding:"required,max=100"`
	Category      string `json:"category" binding:"required,max=100"`
	Justification string `json:"justification" binding:"max=1000"`
}

// RequestResponse is what we send back over HTTP.
type RequestResponse struct {
	ID            uuid.UUID            `json:"id"`
	EmployeeID    uuid.UUID            `json:"employee_id"`
	AssetType     string               `json:"asset_type"`
	Category      string               `json:"category"`
	Justification string               `json:"justification"`
	Status        domain.RequestStatus `json:"status"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

func toRequestResponse(r *domain.AssetRequest) RequestResponse {
	return RequestResponse{
		ID:            r.ID,
		EmployeeID:    r.EmployeeID,
		AssetType:     r.AssetType,
		Category:      r.Category,
		Justification: r.Justification,
		Status:        r.Status,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
	}
}

func toRequestResponseList(requests []domain.AssetRequest) []RequestResponse {
	responses := make([]RequestResponse, 0, len(requests))
	for _, r := range requests {
		responses = append(responses, toRequestResponse(&r))
	}
	return responses
}
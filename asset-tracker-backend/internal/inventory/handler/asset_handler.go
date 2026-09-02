package handler

import (
	"errors"
	"net/http"

	"asset-backend/internal/inventory/domain"
	"asset-backend/internal/inventory/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AssetHandler holds dependencies needed by the asset HTTP handlers.
// It depends on the AssetRepository INTERFACE, not the concrete GORM
// struct — same reasoning as before: lets tests inject a fake repo.
type AssetHandler struct {
	repo repository.AssetRepository
}

func NewAssetHandler(repo repository.AssetRepository) *AssetHandler {
	return &AssetHandler{repo: repo}
}

// CreateAsset handles POST /assets.
func (h *AssetHandler) CreateAsset(c *gin.Context) {
	var req CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	asset := &domain.Asset{
		ID:           uuid.New(),
		Name:         req.Name,
		Type:         req.Type,
		Category:     req.Category,
		SerialNumber: req.SerialNumber,
		Status:       domain.AssetStatusAvailable,
	}

	if err := h.repo.Create(c.Request.Context(), asset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create asset"})
		return
	}

	c.JSON(http.StatusCreated, toAssetResponse(asset))
}

// ListAssets handles GET /assets, with an optional ?status= query filter.
func (h *AssetHandler) ListAssets(c *gin.Context) {
	statusFilter := c.Query("status")

	var (
		assets []domain.Asset
		err    error
	)

	if statusFilter != "" {
		assets, err = h.repo.ListByStatus(c.Request.Context(), domain.AssetStatus(statusFilter))
	} else {
		assets, err = h.repo.List(c.Request.Context())
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list assets"})
		return
	}

	c.JSON(http.StatusOK, toAssetResponseList(assets))
}

// GetAsset handles GET /assets/:id.
func (h *AssetHandler) GetAsset(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset id"})
		return
	}

	asset, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrAssetNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "asset not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get asset"})
		return
	}

	c.JSON(http.StatusOK, toAssetResponse(asset))
}

// UpdateAsset handles PUT /assets/:id — a full replace of the asset's fields.
func (h *AssetHandler) UpdateAsset(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset id"})
		return
	}

	var req UpdateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrAssetNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "asset not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get asset"})
		return
	}

	existing.Name = req.Name
	existing.Type = req.Type
	existing.Category = req.Category
	existing.SerialNumber = req.SerialNumber
	existing.Status = req.Status

	if err := h.repo.Update(c.Request.Context(), existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update asset"})
		return
	}

	c.JSON(http.StatusOK, toAssetResponse(existing))
}

// RetireAsset handles PATCH /assets/:id/retire — a narrow status-only update.
func (h *AssetHandler) RetireAsset(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid asset id"})
		return
	}

	if _, err := h.repo.GetByID(c.Request.Context(), id); err != nil {
		if errors.Is(err, repository.ErrAssetNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "asset not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get asset"})
		return
	}

	if err := h.repo.UpdateStatus(c.Request.Context(), id, domain.AssetStatusRetired); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retire asset"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "asset retired", "id": id})
}
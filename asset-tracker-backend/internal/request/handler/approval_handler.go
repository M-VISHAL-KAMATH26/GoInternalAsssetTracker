package handler

import (
	"errors"
	"log"
	"net/http"

	"asset-backend/internal/request/client"
	"asset-backend/internal/request/domain"
	"asset-backend/internal/request/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ApprovalHandler struct {
	requestRepo     repository.RequestRepository
	approvalRepo    repository.ApprovalRepository
	inventoryClient client.InventoryClient
}

func NewApprovalHandler(requestRepo repository.RequestRepository, approvalRepo repository.ApprovalRepository, inventoryClient client.InventoryClient) *ApprovalHandler {
	return &ApprovalHandler{requestRepo: requestRepo, approvalRepo: approvalRepo, inventoryClient: inventoryClient}
}

// ApproveRequest handles PATCH /requests/:id/approve — manager/admin only.
func (h *ApprovalHandler) ApproveRequest(c *gin.Context) {
	requestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	managerID, err := currentEmployeeID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid manager identity"})
		return
	}

	var body ApprovalActionRequest
	_ = c.ShouldBindJSON(&body) // comment is optional; ignore bind error on empty body

	assetRequest, err := h.requestRepo.GetByID(c.Request.Context(), requestID)
	if err != nil {
		if errors.Is(err, repository.ErrRequestNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get request"})
		return
	}

	if assetRequest.Status != domain.RequestStatusPending {
		c.JSON(http.StatusConflict, gin.H{"error": "request is not pending"})
		return
	}

	// gRPC call to Inventory Service — reserve an available asset.
	assetID, _, err := h.inventoryClient.ReserveAsset(c.Request.Context(), assetRequest.AssetType, assetRequest.Category, assetRequest.EmployeeID)
	if err != nil {
		if errors.Is(err, client.ErrAssetUnavailable) {
			c.JSON(http.StatusConflict, gin.H{"error": "no asset available to fulfill this request"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reserve asset"})
		return
	}

	approval := &domain.Approval{
		ID:        uuid.New(),
		RequestID: requestID,
		ManagerID: managerID,
		Decision:  domain.DecisionApproved,
		Comment:   body.Comment,
	}
	if err := h.approvalRepo.Create(c.Request.Context(), approval); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record approval"})
		return
	}

	if err := h.requestRepo.UpdateStatus(c.Request.Context(), requestID, domain.RequestStatusApproved); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update request status"})
		return
	}

	// Phase 7 will replace this log line with a real queue publish.
	log.Printf("TODO: publish ApprovalDecided event for request %s (approved, asset %s)", requestID, assetID)

	c.JSON(http.StatusOK, ApprovalResponse{
		RequestID: requestID.String(),
		Decision:  string(domain.DecisionApproved),
		AssetID:   assetID.String(),
		Message:   "request approved and asset reserved",
	})
}

// RejectRequest handles PATCH /requests/:id/reject — manager/admin only.
func (h *ApprovalHandler) RejectRequest(c *gin.Context) {
	requestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	managerID, err := currentEmployeeID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid manager identity"})
		return
	}

	var body ApprovalActionRequest
	_ = c.ShouldBindJSON(&body)

	assetRequest, err := h.requestRepo.GetByID(c.Request.Context(), requestID)
	if err != nil {
		if errors.Is(err, repository.ErrRequestNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get request"})
		return
	}

	if assetRequest.Status != domain.RequestStatusPending {
		c.JSON(http.StatusConflict, gin.H{"error": "request is not pending"})
		return
	}

	approval := &domain.Approval{
		ID:        uuid.New(),
		RequestID: requestID,
		ManagerID: managerID,
		Decision:  domain.DecisionRejected,
		Comment:   body.Comment,
	}
	if err := h.approvalRepo.Create(c.Request.Context(), approval); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record approval"})
		return
	}

	if err := h.requestRepo.UpdateStatus(c.Request.Context(), requestID, domain.RequestStatusRejected); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update request status"})
		return
	}

	log.Printf("TODO: publish ApprovalDecided event for request %s (rejected)", requestID)

	c.JSON(http.StatusOK, ApprovalResponse{
		RequestID: requestID.String(),
		Decision:  string(domain.DecisionRejected),
		Message:   "request rejected",
	})
}
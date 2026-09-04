package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	"asset-backend/internal/request/client"
	"asset-backend/internal/request/domain"
	"asset-backend/internal/request/repository"
	"asset-backend/internal/shared/mq"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ApprovalHandler struct {
	requestRepo     repository.RequestRepository
	approvalRepo    repository.ApprovalRepository
	inventoryClient client.InventoryClient
	publisher       *mq.Publisher
}

func NewApprovalHandler(requestRepo repository.RequestRepository, approvalRepo repository.ApprovalRepository, inventoryClient client.InventoryClient, publisher *mq.Publisher) *ApprovalHandler {
	return &ApprovalHandler{
		requestRepo:     requestRepo,
		approvalRepo:    approvalRepo,
		inventoryClient: inventoryClient,
		publisher:       publisher,
	}
}

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
		DecidedAt: time.Now(),
	}
	if err := h.approvalRepo.Create(c.Request.Context(), approval); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record approval"})
		return
	}

	if err := h.requestRepo.UpdateStatus(c.Request.Context(), requestID, domain.RequestStatusApproved); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update request status"})
		return
	}

	// Publish async event — the manager's response does NOT wait on
	// this succeeding. A publish failure is logged, not surfaced to
	// the caller, since the approval itself already succeeded.
	event := mq.ApprovalDecidedEvent{
		RequestID:  requestID,
		EmployeeID: assetRequest.EmployeeID,
		Decision:   string(domain.DecisionApproved),
		AssetID:    assetID.String(),
		Comment:    body.Comment,
		DecidedAt:  time.Now(),
	}
	if err := h.publisher.PublishApprovalDecided(context.Background(), event); err != nil {
		// TODO: proper logging/metrics — for now this is the only
		// visibility into a failed publish.
		println("failed to publish approval decided event:", err.Error())
	}

	c.JSON(http.StatusOK, ApprovalResponse{
		RequestID: requestID.String(),
		Decision:  string(domain.DecisionApproved),
		AssetID:   assetID.String(),
		Message:   "request approved and asset reserved",
	})
}

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
		DecidedAt: time.Now(),
	}
	if err := h.approvalRepo.Create(c.Request.Context(), approval); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to record approval"})
		return
	}

	if err := h.requestRepo.UpdateStatus(c.Request.Context(), requestID, domain.RequestStatusRejected); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update request status"})
		return
	}

	event := mq.ApprovalDecidedEvent{
		RequestID:  requestID,
		EmployeeID: assetRequest.EmployeeID,
		Decision:   string(domain.DecisionRejected),
		Comment:    body.Comment,
		DecidedAt:  time.Now(),
	}
	if err := h.publisher.PublishApprovalDecided(context.Background(), event); err != nil {
		println("failed to publish approval decided event:", err.Error())
	}

	c.JSON(http.StatusOK, ApprovalResponse{
		RequestID: requestID.String(),
		Decision:  string(domain.DecisionRejected),
		Message:   "request rejected",
	})
}
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
	userClient      client.UserClient
	publisher       *mq.Publisher
}

func NewApprovalHandler(requestRepo repository.RequestRepository, approvalRepo repository.ApprovalRepository, inventoryClient client.InventoryClient, userClient client.UserClient, publisher *mq.Publisher) *ApprovalHandler {
	return &ApprovalHandler{
		requestRepo:     requestRepo,
		approvalRepo:    approvalRepo,
		inventoryClient: inventoryClient,
		userClient:      userClient,
		publisher:       publisher,
	}
}

func (h *ApprovalHandler) canReview(ctx context.Context, requesterID, reviewerID uuid.UUID, role string) (bool, error) {
	if role == "admin" {
		return true, nil
	}

	manager, err := h.userClient.GetManagerOf(ctx, requesterID)
	if err != nil {
		return false, err
	}
	return manager.ID == reviewerID, nil
}

func (h *ApprovalHandler) ListPendingApprovals(c *gin.Context) {
	reviewerID, err := currentEmployeeID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid manager identity"})
		return
	}

	requests, err := h.requestRepo.ListByStatus(c.Request.Context(), domain.RequestStatusPending)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list pending requests"})
		return
	}

	if currentRole(c) == "admin" {
		c.JSON(http.StatusOK, toRequestResponseList(requests))
		return
	}

	assigned := make([]domain.AssetRequest, 0)
	for i := range requests {
		canReview, err := h.canReview(c.Request.Context(), requests[i].EmployeeID, reviewerID, currentRole(c))
		if errors.Is(err, client.ErrEmployeeNotFound) {
			continue
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify request manager"})
			return
		}
		if canReview {
			assigned = append(assigned, requests[i])
		}
	}

	c.JSON(http.StatusOK, toRequestResponseList(assigned))
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

	canReview, err := h.canReview(c.Request.Context(), assetRequest.EmployeeID, managerID, currentRole(c))
	if err != nil && !errors.Is(err, client.ErrEmployeeNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify request manager"})
		return
	}
	if !canReview {
		c.JSON(http.StatusForbidden, gin.H{"error": "request is not assigned to this manager"})
		return
	}

	// Claim the decision before reserving inventory: whoever wins this
	// update owns the request, and every other reviewer gets a conflict.
	claimed, err := h.requestRepo.ClaimPending(c.Request.Context(), requestID, domain.RequestStatusApproved)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update request status"})
		return
	}
	if !claimed {
		c.JSON(http.StatusConflict, gin.H{"error": "this request has already been decided"})
		return
	}

	assetID, _, err := h.inventoryClient.ReserveAsset(c.Request.Context(), assetRequest.AssetType, assetRequest.Category, assetRequest.EmployeeID)
	if err != nil {
		// Reservation failed, so release the claim and let it be retried.
		_ = h.requestRepo.UpdateStatus(c.Request.Context(), requestID, domain.RequestStatusPending)

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

	canReview, err := h.canReview(c.Request.Context(), assetRequest.EmployeeID, managerID, currentRole(c))
	if err != nil && !errors.Is(err, client.ErrEmployeeNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify request manager"})
		return
	}
	if !canReview {
		c.JSON(http.StatusForbidden, gin.H{"error": "request is not assigned to this manager"})
		return
	}

	claimed, err := h.requestRepo.ClaimPending(c.Request.Context(), requestID, domain.RequestStatusRejected)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update request status"})
		return
	}
	if !claimed {
		c.JSON(http.StatusConflict, gin.H{"error": "this request has already been decided"})
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

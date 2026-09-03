package handler

import (
	"errors"
	"net/http"

	"asset-backend/internal/request/domain"
	"asset-backend/internal/request/repository"
	"asset-backend/internal/shared/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RequestHandler struct {
	repo repository.RequestRepository
}

func NewRequestHandler(repo repository.RequestRepository) *RequestHandler {
	return &RequestHandler{repo: repo}
}

// currentEmployeeID pulls the authenticated employee's ID out of the
// Gin context, where AuthMiddleware placed it after validating the JWT.
func currentEmployeeID(c *gin.Context) (uuid.UUID, error) {
	raw, exists := c.Get(middleware.ContextKeyEmployeeID)
	if !exists {
		return uuid.UUID{}, errors.New("employee id not found in context")
	}
	return uuid.Parse(raw.(string))
}

func currentRole(c *gin.Context) string {
	raw, exists := c.Get(middleware.ContextKeyRole)
	if !exists {
		return ""
	}
	return raw.(string)
}

// CreateRequest handles POST /requests.
func (h *RequestHandler) CreateRequest(c *gin.Context) {
	employeeID, err := currentEmployeeID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid employee identity"})
		return
	}

	var req CreateRequestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assetRequest := &domain.AssetRequest{
		ID:            uuid.New(),
		EmployeeID:    employeeID,
		AssetType:     req.AssetType,
		Category:      req.Category,
		Justification: req.Justification,
		Status:        domain.RequestStatusPending,
	}

	if err := h.repo.Create(c.Request.Context(), assetRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create request"})
		return
	}

	c.JSON(http.StatusCreated, toRequestResponse(assetRequest))
}

// ListMyRequests handles GET /requests — always scoped to the caller's own requests.
func (h *RequestHandler) ListMyRequests(c *gin.Context) {
	employeeID, err := currentEmployeeID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid employee identity"})
		return
	}

	requests, err := h.repo.ListByEmployee(c.Request.Context(), employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list requests"})
		return
	}

	c.JSON(http.StatusOK, toRequestResponseList(requests))
}

// GetRequest handles GET /requests/:id. Employees can only view their
// own request; managers/admins can view any request (foreshadowing
// the approval workflow, where a manager needs to see requests that
// aren't their own).
func (h *RequestHandler) GetRequest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}

	employeeID, err := currentEmployeeID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid employee identity"})
		return
	}

	assetRequest, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrRequestNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get request"})
		return
	}

	role := currentRole(c)
	if role == "employee" && assetRequest.EmployeeID != employeeID {
		// Deliberately 404, not 403 — don't reveal that a request with
		// this ID exists at all to someone who doesn't own it.
		c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
		return
	}

	c.JSON(http.StatusOK, toRequestResponse(assetRequest))
}
package handler

import (
	"errors"
	"net/http"

	"asset-backend/internal/request/client"
	"asset-backend/internal/request/domain"
	"asset-backend/internal/request/repository"
	"asset-backend/internal/shared/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RequestHandler struct {
	repo       repository.RequestRepository
	userClient client.UserClient
}

func NewRequestHandler(repo repository.RequestRepository, userClient client.UserClient) *RequestHandler {
	return &RequestHandler{repo: repo, userClient: userClient}
}

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

	_, err = h.userClient.GetEmployee(c.Request.Context(), employeeID)
	if err != nil {
		if errors.Is(err, client.ErrEmployeeNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "employee not recognized"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify employee"})
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
		c.JSON(http.StatusNotFound, gin.H{"error": "request not found"})
		return
	}

	c.JSON(http.StatusOK, toRequestResponse(assetRequest))
}
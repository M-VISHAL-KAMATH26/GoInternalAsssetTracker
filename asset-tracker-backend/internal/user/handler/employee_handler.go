package handler

import (
	"errors"
	"net/http"
	"strings"

	"asset-backend/internal/shared/auth"
	"asset-backend/internal/shared/middleware"
	"asset-backend/internal/user/domain"
	"asset-backend/internal/user/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EmployeeHandler struct {
	repo repository.EmployeeRepository
}

func NewEmployeeHandler(repo repository.EmployeeRepository) *EmployeeHandler {
	return &EmployeeHandler{repo: repo}
}

func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var req CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))
	managerEmail := strings.ToLower(strings.TrimSpace(req.ManagerEmail))
	if managerEmail != "" && managerEmail == email {
		c.JSON(http.StatusBadRequest, gin.H{"error": "an employee cannot be their own manager"})
		return
	}
	if _, err := h.repo.GetByEmail(c.Request.Context(), email); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "an employee with this email already exists"})
		return
	} else if !errors.Is(err, repository.ErrEmployeeNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check employee email"})
		return
	}

	var managerID *uuid.UUID
	if req.ManagerEmail != "" {
		manager, err := h.repo.GetByEmail(c.Request.Context(), managerEmail)
		if errors.Is(err, repository.ErrEmployeeNotFound) {
			if req.ManagerName == "" || req.ManagerPassword == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "manager_name and manager_password are required when the manager does not exist",
				})
				return
			}

			managerHash, hashErr := auth.HashPassword(req.ManagerPassword)
			if hashErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to secure manager password"})
				return
			}

			adminID, parseErr := uuid.Parse(c.GetString(middleware.ContextKeyEmployeeID))
			if parseErr != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid admin identity"})
				return
			}

			manager = &domain.Employee{
				ID:           uuid.New(),
				Name:         strings.TrimSpace(req.ManagerName),
				Email:        managerEmail,
				PasswordHash: managerHash,
				Role:         domain.RoleManager,
				ManagerID:    &adminID,
			}
			if err := h.repo.Create(c.Request.Context(), manager); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create manager"})
				return
			}
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find manager"})
			return
		}

		if manager.Role != domain.RoleManager && manager.Role != domain.RoleAdmin {
			c.JSON(http.StatusBadRequest, gin.H{"error": "approving manager must have manager or admin role"})
			return
		}
		managerID = &manager.ID
	}

	if req.Role != string(domain.RoleAdmin) && managerID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "manager_email is required for employees and managers"})
		return
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to secure password"})
		return
	}

	employee := &domain.Employee{
		ID:           uuid.New(),
		Name:         strings.TrimSpace(req.Name),
		Email:        email,
		PasswordHash: passwordHash,
		Role:         domain.EmployeeRole(req.Role),
		ManagerID:    managerID,
	}
	if err := h.repo.Create(c.Request.Context(), employee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create employee"})
		return
	}

	c.JSON(http.StatusCreated, toEmployeeResponse(employee))
}

func (h *EmployeeHandler) ListEmployees(c *gin.Context) {
	employees, err := h.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list employees"})
		return
	}

	responses := make([]EmployeeResponse, 0, len(employees))
	for i := range employees {
		responses = append(responses, toEmployeeResponse(&employees[i]))
	}
	c.JSON(http.StatusOK, responses)
}

func toEmployeeResponse(employee *domain.Employee) EmployeeResponse {
	return EmployeeResponse{
		ID:        employee.ID,
		Name:      employee.Name,
		Email:     employee.Email,
		Role:      string(employee.Role),
		ManagerID: employee.ManagerID,
	}
}

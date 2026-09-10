package handler

import (
	"net/http"

	"asset-backend/internal/shared/auth"
	"asset-backend/internal/user/repository"

	"github.com/gin-gonic/gin"
)

const invalidCredentials = "invalid email or password"

// AuthHandler holds dependencies for authentication HTTP handlers.
type AuthHandler struct {
	repo      repository.EmployeeRepository
	jwtSecret string
}

func NewAuthHandler(repo repository.EmployeeRepository, jwtSecret string) *AuthHandler {
	return &AuthHandler{repo: repo, jwtSecret: jwtSecret}
}

// Login handles POST /login.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employee, err := h.repo.GetByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": invalidCredentials})
		return
	}

	if err := auth.ComparePassword(employee.PasswordHash, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": invalidCredentials})
		return
	}

	token, err := auth.GenerateToken(employee.ID, string(employee.Role), h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token:      token,
		EmployeeID: employee.ID,
		Role:       string(employee.Role),
		Name:       employee.Name,
	})
}

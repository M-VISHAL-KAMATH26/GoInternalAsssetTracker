package handler

import (
	"asset-backend/internal/shared/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes wires public authentication endpoints onto the given
// router. Login is the auth entry point, so it is not wrapped in JWT middleware.
func RegisterAuthRoutes(router *gin.Engine, h *AuthHandler) {
	router.POST("/login", h.Login)
}

func RegisterEmployeeRoutes(router *gin.Engine, h *EmployeeHandler) {
	employees := router.Group("/employees")
	employees.Use(middleware.AuthMiddleware(), middleware.RequireRole("admin"))
	{
		employees.POST("", h.CreateEmployee)
		employees.GET("", h.ListEmployees)
	}
}

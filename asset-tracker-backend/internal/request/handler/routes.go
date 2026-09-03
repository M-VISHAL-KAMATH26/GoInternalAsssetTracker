package handler

import (
	"asset-backend/internal/shared/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRequestRoutes wires up /requests endpoints, open to any
// authenticated employee (not admin-only, unlike the asset routes).
func RegisterRequestRoutes(router *gin.Engine, h *RequestHandler) {
	requests := router.Group("/requests")
	requests.Use(middleware.AuthMiddleware(), middleware.RequireRole("employee", "manager", "admin"))
	{
		requests.POST("", h.CreateRequest)
		requests.GET("", h.ListMyRequests)
		requests.GET("/:id", h.GetRequest)
	}
}
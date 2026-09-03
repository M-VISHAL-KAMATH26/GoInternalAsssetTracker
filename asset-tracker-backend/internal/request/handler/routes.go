package handler

import (
	"asset-backend/internal/shared/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRequestRoutes(router *gin.Engine, h *RequestHandler, approvalHandler *ApprovalHandler) {
	requests := router.Group("/requests")
	requests.Use(middleware.AuthMiddleware(), middleware.RequireRole("employee", "manager", "admin"))
	{
		requests.POST("", h.CreateRequest)
		requests.GET("", h.ListMyRequests)
		requests.GET("/:id", h.GetRequest)
	}

	approvals := router.Group("/requests")
	approvals.Use(middleware.AuthMiddleware(), middleware.RequireRole("manager", "admin"))
	{
		approvals.PATCH("/:id/approve", approvalHandler.ApproveRequest)
		approvals.PATCH("/:id/reject", approvalHandler.RejectRequest)
	}
}
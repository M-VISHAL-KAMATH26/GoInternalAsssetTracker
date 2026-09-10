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
	}

	adminRequests := router.Group("/requests")
	adminRequests.Use(middleware.AuthMiddleware(), middleware.RequireRole("admin"))
	{
		adminRequests.GET("/all", h.ListAllRequests)
	}

	requests.GET("/:id", h.GetRequest)

	approvals := router.Group("/requests")
	approvals.Use(middleware.AuthMiddleware(), middleware.RequireRole("manager", "admin"))
	{
		approvals.PATCH("/:id/approve", approvalHandler.ApproveRequest)
		approvals.PATCH("/:id/reject", approvalHandler.RejectRequest)
	}

	pendingApprovals := router.Group("/approvals")
	pendingApprovals.Use(middleware.AuthMiddleware(), middleware.RequireRole("manager", "admin"))
	{
		pendingApprovals.GET("/pending", approvalHandler.ListPendingApprovals)
	}
}

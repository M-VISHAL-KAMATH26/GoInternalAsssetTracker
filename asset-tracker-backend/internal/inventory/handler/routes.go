package handler

import (
	"asset-backend/internal/shared/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterAssetRoutes wires up all /assets endpoints onto the given
// router, protected by JWT auth + admin-only RBAC.
func RegisterAssetRoutes(router *gin.Engine, h *AssetHandler) {
	assets := router.Group("/assets")
	assets.Use(middleware.AuthMiddleware(), middleware.RequireRole("admin"))
	{
		assets.POST("", h.CreateAsset)
		assets.GET("", h.ListAssets)
		assets.GET("/:id", h.GetAsset)
		assets.PUT("/:id", h.UpdateAsset)
		assets.PATCH("/:id/retire", h.RetireAsset)
	}
}
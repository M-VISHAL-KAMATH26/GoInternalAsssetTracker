package handler

import (
	"asset-backend/internal/shared/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterAssetRoutes wires up all /assets endpoints onto the given
// router. Everything is admin-only except the catalog, which any
// authenticated employee needs in order to raise a request.
func RegisterAssetRoutes(router *gin.Engine, h *AssetHandler) {
	catalog := router.Group("/assets")
	catalog.Use(middleware.AuthMiddleware(), middleware.RequireRole("employee", "manager", "admin"))
	{
		catalog.GET("/catalog", h.ListAssetCatalog)
	}

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

package handler

import "github.com/gin-gonic/gin"

// RegisterAuthRoutes wires public authentication endpoints onto the given
// router. Login is the auth entry point, so it is not wrapped in JWT middleware.
func RegisterAuthRoutes(router *gin.Engine, h *AuthHandler) {
	router.POST("/login", h.Login)
}

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireRole returns a middleware that only allows the request to
// proceed if the role attached to the context (by AuthMiddleware)
// matches one of the given allowed roles. Must run AFTER AuthMiddleware
// in the chain — it relies on ContextKeyRole already being set.
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(allowedRoles))
	for _, r := range allowedRoles {
		allowed[r] = true
	}

	return func(c *gin.Context) {
		role, exists := c.Get(ContextKeyRole)
		if !exists {
			// AuthMiddleware didn't run before this — a wiring mistake,
			// not a client error, but we still fail closed (deny access).
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "role not found in context"})
			return
		}

		roleStr, ok := role.(string)
		if !ok || !allowed[roleStr] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}

		c.Next()
	}
}
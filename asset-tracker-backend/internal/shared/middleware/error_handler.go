package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler is a helper for consistent error responses across every
// handler — logs the underlying error with context (path, method) and
// returns a uniform JSON shape, so callers can always expect
// {"error": "..."} regardless of which handler failed.
func RespondError(c *gin.Context, log *slog.Logger, status int, publicMessage string, internalErr error) {
	if internalErr != nil {
		log.Error(publicMessage,
			"error", internalErr.Error(),
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
			"status", status,
		)
	}
	c.JSON(status, gin.H{"error": publicMessage})
}

// RecoveryLogger returns a Gin middleware that logs panics with
// structured context before Gin's default recovery converts them to a
// 500 — use alongside gin.Recovery(), not instead of it.
func RecoveryLogger(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Error("panic recovered",
					"panic", r,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
		}()
		c.Next()
	}
}
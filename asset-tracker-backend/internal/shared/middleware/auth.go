package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Context keys used to pass identity from AuthMiddleware to downstream
// middleware/handlers via Gin's request context.
const (
	ContextKeyEmployeeID = "employeeID"
	ContextKeyRole       = "role"
)

// Claims represents the custom fields embedded in our JWT, alongside
// the standard registered claims (expiry, issued-at, etc).
type Claims struct {
	EmployeeID string `json:"employee_id"`
	Role       string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates the JWT on the Authorization header and, if
// valid, attaches the employee's ID and role to the request context
// for downstream middleware/handlers to use. If the token is missing,
// malformed, expired, or has a bad signature, the request is aborted
// here with 401 — no handler below this ever executes.
func AuthMiddleware() gin.HandlerFunc {
	secret := []byte(os.Getenv("JWT_SECRET"))

	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header must be a Bearer token"})
			return
		}
		tokenString := parts[1]

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		if _, err := uuid.Parse(claims.EmployeeID); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token contains invalid employee id"})
			return
		}

		c.Set(ContextKeyEmployeeID, claims.EmployeeID)
		c.Set(ContextKeyRole, claims.Role)

		c.Next()
	}
}
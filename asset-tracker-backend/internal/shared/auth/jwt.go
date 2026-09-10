package auth

import (
	"time"

	"asset-backend/internal/shared/middleware"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// GenerateToken issues a signed HS256 JWT for the given employee, using
// the same Claims shape AuthMiddleware expects when validating requests.
func GenerateToken(employeeID uuid.UUID, role string, secret string) (string, error) {
	claims := middleware.Claims{
		EmployeeID: employeeID.String(),
		Role:       role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

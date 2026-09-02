package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"asset-backend/internal/shared/middleware"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func main() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET env var must be set")
	}

	role := "admin"
	if len(os.Args) > 1 {
		role = os.Args[1] // e.g. `go run ./cmd/gen-test-token employee`
	}

	claims := middleware.Claims{
		EmployeeID: uuid.New().String(),
		Role:       role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatalf("failed to sign token: %v", err)
	}

	fmt.Println(signed)
}
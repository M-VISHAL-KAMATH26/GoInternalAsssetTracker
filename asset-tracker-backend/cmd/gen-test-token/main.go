package main

import (
	"fmt"
	"log"
	"os"

	"asset-backend/internal/shared/auth"

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

	signed, err := auth.GenerateToken(uuid.New(), role, secret)
	if err != nil {
		log.Fatalf("failed to sign token: %v", err)
	}

	fmt.Println(signed)
}

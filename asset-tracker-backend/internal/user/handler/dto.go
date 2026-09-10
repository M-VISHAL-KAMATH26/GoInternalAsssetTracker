package handler

import "github.com/google/uuid"

// LoginRequest is the expected JSON body for POST /login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse is returned after a successful login.
type LoginResponse struct {
	Token      string    `json:"token"`
	EmployeeID uuid.UUID `json:"employee_id"`
	Role       string    `json:"role"`
	Name       string    `json:"name"`
}

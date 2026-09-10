package handler

import "github.com/google/uuid"

// LoginRequest is the expected JSON body for POST /login.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse is returned after a successful login.
type LoginResponse struct {
	Token          string    `json:"token"`
	EmployeeID     uuid.UUID `json:"employee_id"`
	Role           string    `json:"role"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	OfficeLocation string    `json:"office_location"`
	AvatarURL      string    `json:"avatar_url"`
}

// CreateEmployeeRequest provisions a user. If ManagerEmail does not already
// exist, ManagerName and ManagerPassword are used to create that manager.
type CreateEmployeeRequest struct {
	Name                  string `json:"name" binding:"required,min=2,max=255"`
	Email                 string `json:"email" binding:"required,email"`
	Password              string `json:"password" binding:"required,min=8"`
	Role                  string `json:"role" binding:"required,oneof=employee manager admin"`
	OfficeLocation        string `json:"office_location" binding:"required,max=100"`
	AvatarURL             string `json:"avatar_url" binding:"omitempty,url,max=512"`
	ManagerEmail          string `json:"manager_email" binding:"omitempty,email"`
	ManagerName           string `json:"manager_name" binding:"omitempty,min=2,max=255"`
	ManagerPassword       string `json:"manager_password" binding:"omitempty,min=8"`
	ManagerOfficeLocation string `json:"manager_office_location" binding:"omitempty,max=100"`
	ManagerAvatarURL      string `json:"manager_avatar_url" binding:"omitempty,url,max=512"`
}

type EmployeeResponse struct {
	ID             uuid.UUID  `json:"id"`
	Name           string     `json:"name"`
	Email          string     `json:"email"`
	Role           string     `json:"role"`
	OfficeLocation string     `json:"office_location"`
	AvatarURL      string     `json:"avatar_url"`
	ManagerID      *uuid.UUID `json:"manager_id,omitempty"`
}

package model

import "time"

var ValidUserRoles = map[string]bool{
	"employee": true, "procurement": true, "ceo": true, "admin": true,
}

type CreateUserRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	Role      string `json:"role" binding:"required"`
}

type UpdateUserRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email" binding:"omitempty,email"`
	Role      *string `json:"role"`
}

type UserWriteResponse struct {
	ID        string     `json:"id"`
	FirstName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type Role string

const (
	RoleEmployee    Role = "employee"
	RoleProcurement Role = "procurement"
	RoleCeo         Role = "ceo"
	RoleAdmin       Role = "admin"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewErrorResponse(code int, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:    code,
		Message: message,
	}
}

type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewSuccessResponse(code int, message string, data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func FromStringToRole(roleStr string) Role {
	switch roleStr {
	case "employee":
		return RoleEmployee
	case "procurement":
		return RoleProcurement
	case "ceo":
		return RoleCeo
	case "admin":
		return RoleAdmin
	default:
		return ""
	}
}

package auth

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleStudent UserRole = "student"
	RoleTeacher UserRole = "teacher"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	Role      UserRole  `json:"user_role"`
}

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Role     UserRole  `json:"role"`
}

type Session struct {
	Token     string    `json:"token"`
	UserId    uuid.UUID `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

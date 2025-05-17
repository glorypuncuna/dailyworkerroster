package model

import "time"

const (
	// User roles
	ROLE_ADMIN  = "ADMIN"
	ROLE_WORKER = "WORKER"
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"` // ADMIN, WORKER
	JWTToken  string    `json:"jwt_token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

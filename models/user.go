package models

import "time"

type LoginRequest struct {
	Email    string `json:"email" `
	Password string `json:"password" `
}
type User struct {
	User_id   int       `json:"user_id"`    // Primary key
	Name      string    `json:"name"`       // User's full name
	Email     string    `json:"email"`      // Email address
	Password  string    `json:"password"`   // Hashed password
	Role      string    `json:"role"`       // User role (admin, user, etc.)
	CreatedAt time.Time `json:"created_at"` // Date when the account was created
	UpdatedAt time.Time `json:"updated_at"` // Date when the account was last updated
}

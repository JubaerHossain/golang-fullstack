package entity

import (
	"time"

	"github.com/JubaerHossain/rootx/pkg/core/entity"
)

// User represents the user entity
type User struct {
	ID        uint      `json:"id"` // Primary key
	Name      string    `json:"name" validate:"required,min=3,max=100"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Password  string    `json:"password" validate:"required,min=6,max=100"`
	Email     string    `json:"email" validate:"required,email"`
	Role      uint      `json:"role"`
	Status    bool      `json:"status"`
}

// UpdateUser represents the user update request
type UpdateUser struct {
	Name      string    `json:"name" validate:"omitempty,min=3,max=100"`
	Status    bool      `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ResponseUser represents the user response
type ResponseUser struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    bool      `json:"status"`
}

type UserResponsePagination struct {
	Data       []*ResponseUser   `json:"data"`
	Pagination entity.Pagination `json:"pagination"`
}

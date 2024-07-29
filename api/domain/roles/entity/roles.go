package entity

import (
	"time"
    "github.com/JubaerHossain/rootx/pkg/core/entity"
)

// Role represents the role entity
type Role struct {
	ID        uint          `json:"id"` // Primary key
	Name      string        `json:"name" validate:"required,min=3,max=100"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Status    bool          `json:"status"`
}

// UpdateRole represents the role update request
type UpdateRole struct {
	Name   string        `json:"name" validate:"omitempty,min=3,max=100"`
	Status bool          `json:"status"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ResponseRole represents the role response
type ResponseRole struct {
	ID        uint          `json:"id"`
	Name      string        `json:"name"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Status    bool          `json:"status"`
}

type RoleResponsePagination struct {
	Data       []*ResponseRole   `json:"data"`
	Pagination entity.Pagination `json:"pagination"`
}

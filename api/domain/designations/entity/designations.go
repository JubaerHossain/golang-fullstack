package entity

import (
	"time"
    "github.com/JubaerHossain/rootx/pkg/core/entity"
)

// Designation represents the designation entity
type Designation struct {
	ID        uint          `json:"id"` // Primary key
	Name      string        `json:"name" validate:"required,min=3,max=100"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Status    bool          `json:"status"`
}

// UpdateDesignation represents the designation update request
type UpdateDesignation struct {
	Name   string        `json:"name" validate:"omitempty,min=3,max=100"`
	Status bool          `json:"status"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ResponseDesignation represents the designation response
type ResponseDesignation struct {
	ID        uint          `json:"id"`
	Name      string        `json:"name"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Status    bool          `json:"status"`
}

type DesignationResponsePagination struct {
	Data       []*ResponseDesignation   `json:"data"`
	Pagination entity.Pagination `json:"pagination"`
}

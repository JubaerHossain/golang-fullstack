package entity

import (
	"time"

	"github.com/JubaerHossain/rootx/pkg/core/entity"
)

// Category represents the category entity
type Category struct {
	ID         uint64     `json:"id" validate:"required"`
	Title      *string    `json:"title" validate:"max=191"`
	Slug       string     `json:"slug"`
	Order      int        `json:"order" validate:"gte=0"`
	Label      int        `json:"label" validate:"gte=0"`
	IsFeatured bool       `json:"is_featured" validate:"oneof=0 1"`
	ParentID   *uint64    `json:"parent_id"`
	ViewCount  uint32     `json:"view_count" validate:"gte=0"`
	StatusID   uint64     `json:"status_id" validate:"required,gte=1"`
	CreatedBy  *uint64    `json:"created_by"`
	UpdatedBy  *uint64    `json:"updated_by"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

// UpdateCategory represents the category update request
type UpdateCategory struct {
	Title    *string `json:"title" validate:"max=191"`
	StatusID uint64  `json:"status_id" validate:"required,gte=1"`
}

// ResponseCategory represents the category response
type ResponseCategoryDetails struct {
	ID         uint64  `json:"id"`
	Title      string  `json:"title"`
	Slug       string  `json:"slug"`
	Order      int     `json:"order"`
	Label      int     `json:"label"`
	IsFeatured bool    `json:"is_featured"`
	ParentID   *uint64 `json:"parent_id"`
	ViewCount  uint32  `json:"view_count"`
	StatusID   uint64  `json:"status_id"`
}

// ResponseCategory represents the category response
type ResponseCategory struct {
	ID            uint64              `json:"id"`
	Title         string              `json:"title"`
	Slug          string              `json:"slug"`
	Order         int                 `json:"order"`
	StatusID      uint64              `json:"status_id"`
	ParentID      *uint64             `json:"parent_id"`
	ChildCategory []*ResponseCategory `json:"child_category"`
}

type CategoryResponsePagination struct {
	Data       []*ResponseCategory `json:"data"`
	Pagination entity.Pagination   `json:"pagination"`
}
type Response struct {
	Success bool                          `json:"success"`
	Message string                        `json:"message,omitempty"`
	Data    []*CategoryResponsePagination `json:"data"`
}

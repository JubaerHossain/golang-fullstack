package entity

import (
	"time"

	"github.com/JubaerHossain/rootx/pkg/core/entity"
)

// Department represents the department entity
type Department struct {
	ID        uint      `json:"id"` // Primary key
	Title     string    `json:"title" validate:"required,min=3,max=100"`
	Slug      string    `json:"slug" validate:"required,min=3,max=100"`
	CreatedBy uint      `json:"created_by"`
	UpdatedBy uint      `json:"updated_by"`
	StatusID  uint      `json:"status_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeleteAt  time.Time `json:"delete_at"`
}

// UpdateDepartment represents the department update request
type UpdateDepartment struct {
	Title     string    `json:"title" validate:"required,min=3,max=100"`
	Slug      string    `json:"slug" validate:"required,min=3,max=100"`
	UpdatedBy uint      `json:"updated_by"`
	StatusID  uint      `json:"status_id"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ResponseDepartment represents the department response
type ResponseDepartment struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	CreatedBy uint   `json:"created_by"`
	UpdatedBy uint   `json:"updated_by"`
	StatusID  uint   `json:"status_id"`
}

type DepartmentResponsePagination struct {
	Data       []*ResponseDepartment `json:"data"`
	Pagination entity.Pagination     `json:"pagination"`
}

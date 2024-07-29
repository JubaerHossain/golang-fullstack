package entity

import (
	"time"

	"github.com/JubaerHossain/rootx/pkg/core/entity"
)

// News represents the news entity
type News struct {
	ID        uint      `json:"id"` // Primary key
	Name      string    `json:"name" validate:"required,min=3,max=100"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    bool      `json:"status"`
}

type Tags struct {
	ID   uint   `json:"id"` // Primary key
	Name string `json:"name" validate:"required,min=3,max=100"`
}

type ScrollNews struct {
	ID           uint     `json:"id"`
	Title        string   `json:"title"`
	Slug         string   `json:"slug"`
	Type         string   `json:"type"`
	SubTitle     string   `json:"sub_title"`
	Tags         []string `json:"tags"`
	Content      string   `json:"content"`
	MetaTitle    string   `json:"meta_title"`
	MetaDesc     string   `json:"meta_description"`
	MetaKeywords []string `json:"meta_keywords"`
	UpdatedAt    string   `json:"updated_at"`
	Author       string   `json:"author"`
	URL          string   `json:"url"`
	PathSmall    string   `json:"path_small"`
	PathMedium   string   `json:"path_medium"`
	PathLarge    string   `json:"path_large"`
	Loading      string   `json:"loading"`
	Status       string   `json:"status"`
	Category     string   `json:"category"`
}

// UpdateNews represents the news update request
type UpdateNews struct {
	Name      string    `json:"name" validate:"omitempty,min=3,max=100"`
	Status    bool      `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ResponseNews represents the news response
type ResponseNews struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    bool      `json:"status"`
}

type NewsResponsePagination struct {
	Data       []*ResponseNews   `json:"data"`
	Pagination entity.Pagination `json:"pagination"`
}

type ScrollNewsResponse struct {
	Data       []*ScrollNews     `json:"data"`
	Pagination entity.Pagination `json:"pagination"`
}

type ThumbnailNewsResponse struct {
	Data       []*ScrollNews     `json:"data"`
	Pagination entity.Pagination `json:"pagination"`
}

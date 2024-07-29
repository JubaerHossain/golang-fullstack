package repository

import (
	"net/http"

	"github.com/JubaerHossain/cn-api/domain/categories/entity"
)


// CategoryRepository defines methods for category data access
type CategoryRepository interface {
	GetCategories(r *http.Request) (*entity.CategoryResponsePagination, error)
	GetCategoryByID(categoryID uint) (*entity.Category, error)
	GetCategory(categoryID uint) (*entity.ResponseCategory, error)
	CreateCategory(category *entity.Category, r *http.Request)  error
	UpdateCategory(oldCategory *entity.Category, category *entity.UpdateCategory, r *http.Request) error
	DeleteCategory(category *entity.Category, r *http.Request) error
}
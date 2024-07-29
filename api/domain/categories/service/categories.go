package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/JubaerHossain/cn-api/domain/categories/entity"
	"github.com/JubaerHossain/cn-api/domain/categories/infrastructure/persistence"
	"github.com/JubaerHossain/cn-api/domain/categories/repository"
	"github.com/JubaerHossain/rootx/pkg/core/app"
	"go.uber.org/zap"
)

type Service struct {
	app  *app.App
	repo repository.CategoryRepository
}

func NewService(app *app.App) *Service {
	repo := persistence.NewCategoryRepository(app)
	return &Service{
		app:  app,
		repo: repo,
	}
}

func (s *Service) GetCategories(r *http.Request) (*entity.CategoryResponsePagination, error) {
	// Call repository to get all categories
	categories, categoryErr := s.repo.GetCategories(r)
	if categoryErr != nil {
		s.app.Logger.Error("Error getting category", zap.Error(categoryErr))
		return nil, categoryErr
	}
	return categories, nil
}



// CreateCategory creates a new category
func (s *Service) CreateCategory(category *entity.Category, r *http.Request)  error {
	// Add any validation or business logic here before creating the category
    if err := s.repo.CreateCategory(category, r); err != nil {
		s.app.Logger.Error("Error creating category", zap.Error(err))
        return err
    }
	return nil
}

func (s *Service) GetCategoryByID(r *http.Request) (*entity.Category, error) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID")
	}
	category, categoryErr := s.repo.GetCategoryByID(uint(id))
	if categoryErr != nil {
		s.app.Logger.Error("Error getting category by ID", zap.Error(categoryErr))
		return nil, categoryErr
	}
	return category, nil
}

// GetCategoryDetails retrieves a category by ID
func (s *Service) GetCategoryDetails(r *http.Request) (*entity.ResponseCategory, error) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid category ID")
	}
	category, categoryErr := s.repo.GetCategory(uint(id))
	if categoryErr != nil {
		s.app.Logger.Error("Error getting category details", zap.Error(categoryErr))
		return nil, categoryErr
	}
	return category, nil
}

// UpdateCategory updates an existing category
func (s *Service) UpdateCategory(r *http.Request, category *entity.UpdateCategory)  error {
	// Call repository to update category
	oldCategory, err := s.GetCategoryByID(r)
	if err != nil {
		return err
	}

	err2 := s.repo.UpdateCategory(oldCategory, category, r)
	if err2 != nil {
		s.app.Logger.Error("Error updating category", zap.Error(err2))
		return err2
	}
	return  nil
}

// DeleteCategory deletes a category by ID
func (s *Service) DeleteCategory(r *http.Request) error {
	// Call repository to delete category
	category, err := s.GetCategoryByID(r)
	if err != nil {
		return err
	}

	err2 := s.repo.DeleteCategory(category, r)
	if err2 != nil {
		s.app.Logger.Error("Error deleting category", zap.Error(err2))
		return err2
	}

	return nil
}

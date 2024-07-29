package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/JubaerHossain/cn-api/domain/designations/entity"
	"github.com/JubaerHossain/cn-api/domain/designations/infrastructure/persistence"
	"github.com/JubaerHossain/cn-api/domain/designations/repository"
	"github.com/JubaerHossain/rootx/pkg/core/app"
)

type Service struct {
	app  *app.App
	repo repository.DesignationRepository
}

func NewService(app *app.App) *Service {
	repo := persistence.NewDesignationRepository(app)
	return &Service{
		app:  app,
		repo: repo,
	}
}

func (s *Service) GetDesignations(r *http.Request) (*entity.DesignationResponsePagination, error) {
	// Call repository to get all designations
	designations, designationErr := s.repo.GetDesignations(r)
	if designationErr != nil {
		return nil, designationErr
	}
	return designations, nil
}



// CreateDesignation creates a new designation
func (s *Service) CreateDesignation(designation *entity.Designation, r *http.Request)  error {
	// Add any validation or business logic here before creating the designation
    if err := s.repo.CreateDesignation(designation, r); err != nil {
        return err
    }
	return nil
}

func (s *Service) GetDesignationByID(r *http.Request) (*entity.Designation, error) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid designation ID")
	}
	designation, designationErr := s.repo.GetDesignationByID(uint(id))
	if designationErr != nil {
		return nil, designationErr
	}
	return designation, nil
}

// GetDesignationDetails retrieves a designation by ID
func (s *Service) GetDesignationDetails(r *http.Request) (*entity.ResponseDesignation, error) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid designation ID")
	}
	designation, designationErr := s.repo.GetDesignation(uint(id))
	if designationErr != nil {
		return nil, designationErr
	}
	return designation, nil
}

// UpdateDesignation updates an existing designation
func (s *Service) UpdateDesignation(r *http.Request, designation *entity.UpdateDesignation)  error {
	// Call repository to update designation
	oldDesignation, err := s.GetDesignationByID(r)
	if err != nil {
		return err
	}

	err2 := s.repo.UpdateDesignation(oldDesignation, designation, r)
	if err2 != nil {
		return err2
	}
	return  nil
}

// DeleteDesignation deletes a designation by ID
func (s *Service) DeleteDesignation(r *http.Request) error {
	// Call repository to delete designation
	designation, err := s.GetDesignationByID(r)
	if err != nil {
		return err
	}

	err2 := s.repo.DeleteDesignation(designation, r)
	if err2 != nil {
		return err2
	}

	return nil
}

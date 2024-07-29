package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/JubaerHossain/cn-api/domain/departments/entity"
	"github.com/JubaerHossain/cn-api/domain/departments/infrastructure/persistence"
	"github.com/JubaerHossain/cn-api/domain/departments/repository"
	"github.com/JubaerHossain/rootx/pkg/core/app"
)

type Service struct {
	app  *app.App
	repo repository.DepartmentRepository
}

func NewService(app *app.App) *Service {
	repo := persistence.NewDepartmentRepository(app)
	return &Service{
		app:  app,
		repo: repo,
	}
}

func (s *Service) GetDepartments(r *http.Request) (*entity.DepartmentResponsePagination, error) {
	// Call repository to get all departments
	departments, departmentErr := s.repo.GetDepartments(r)
	if departmentErr != nil {
		return nil, departmentErr
	}
	return departments, nil
}



// CreateDepartment creates a new department
func (s *Service) CreateDepartment(department *entity.Department, r *http.Request)  error {
	// Add any validation or business logic here before creating the department
    if err := s.repo.CreateDepartment(department, r); err != nil {
        return err
    }
	return nil
}

func (s *Service) GetDepartmentByID(r *http.Request) (*entity.Department, error) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid department ID")
	}
	department, departmentErr := s.repo.GetDepartmentByID(uint(id))
	if departmentErr != nil {
		return nil, departmentErr
	}
	return department, nil
}

// GetDepartmentDetails retrieves a department by ID
func (s *Service) GetDepartmentDetails(r *http.Request) (*entity.ResponseDepartment, error) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid department ID")
	}
	department, departmentErr := s.repo.GetDepartment(uint(id))
	if departmentErr != nil {
		return nil, departmentErr
	}
	return department, nil
}

// UpdateDepartment updates an existing department
func (s *Service) UpdateDepartment(r *http.Request, department *entity.UpdateDepartment)  error {
	// Call repository to update department
	oldDepartment, err := s.GetDepartmentByID(r)
	if err != nil {
		return err
	}

	err2 := s.repo.UpdateDepartment(oldDepartment, department, r)
	if err2 != nil {
		return err2
	}
	return  nil
}

// DeleteDepartment deletes a department by ID
func (s *Service) DeleteDepartment(r *http.Request) error {
	// Call repository to delete department
	department, err := s.GetDepartmentByID(r)
	if err != nil {
		return err
	}

	err2 := s.repo.DeleteDepartment(department, r)
	if err2 != nil {
		return err2
	}

	return nil
}

package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/JubaerHossain/cn-api/domain/roles/entity"
	"github.com/JubaerHossain/cn-api/domain/roles/infrastructure/persistence"
	"github.com/JubaerHossain/cn-api/domain/roles/repository"
	"github.com/JubaerHossain/rootx/pkg/core/app"
)

type Service struct {
	app  *app.App
	repo repository.RoleRepository
}

func NewService(app *app.App) *Service {
	repo := persistence.NewRoleRepository(app)
	return &Service{
		app:  app,
		repo: repo,
	}
}

func (s *Service) GetRoles(r *http.Request) (*entity.RoleResponsePagination, error) {
	// Call repository to get all roles
	roles, roleErr := s.repo.GetRoles(r)
	if roleErr != nil {
		return nil, roleErr
	}
	return roles, nil
}



// CreateRole creates a new role
func (s *Service) CreateRole(role *entity.Role, r *http.Request)  error {
	// Add any validation or business logic here before creating the role
    if err := s.repo.CreateRole(role, r); err != nil {
        return err
    }
	return nil
}

func (s *Service) GetRoleByID(r *http.Request) (*entity.Role, error) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid role ID")
	}
	role, roleErr := s.repo.GetRoleByID(uint(id))
	if roleErr != nil {
		return nil, roleErr
	}
	return role, nil
}

// GetRoleDetails retrieves a role by ID
func (s *Service) GetRoleDetails(r *http.Request) (*entity.ResponseRole, error) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid role ID")
	}
	role, roleErr := s.repo.GetRole(uint(id))
	if roleErr != nil {
		return nil, roleErr
	}
	return role, nil
}

// UpdateRole updates an existing role
func (s *Service) UpdateRole(r *http.Request, role *entity.UpdateRole)  error {
	// Call repository to update role
	oldRole, err := s.GetRoleByID(r)
	if err != nil {
		return err
	}

	err2 := s.repo.UpdateRole(oldRole, role, r)
	if err2 != nil {
		return err2
	}
	return  nil
}

// DeleteRole deletes a role by ID
func (s *Service) DeleteRole(r *http.Request) error {
	// Call repository to delete role
	role, err := s.GetRoleByID(r)
	if err != nil {
		return err
	}

	err2 := s.repo.DeleteRole(role, r)
	if err2 != nil {
		return err2
	}

	return nil
}

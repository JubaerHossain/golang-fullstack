package service

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/JubaerHossain/cn-api/domain/users/entity"
	"github.com/JubaerHossain/cn-api/domain/users/infrastructure/persistence"
	"github.com/JubaerHossain/cn-api/domain/users/repository"
	"github.com/JubaerHossain/rootx/pkg/core/app"
)

type Service struct {
	app  *app.App
	repo repository.UserRepository
}

func NewService(app *app.App) *Service {
	repo := persistence.NewUserRepository(app)
	return &Service{
		app:  app,
		repo: repo,
	}
}

func (s *Service) GetUsers(r *http.Request) (*entity.UserResponsePagination, error) {
	// Call repository to get all users
	users, userErr := s.repo.GetUsers(r)
	if userErr != nil {
		return nil, userErr
	}
	return users, nil
}



// CreateUser creates a new user
func (s *Service) CreateUser(user *entity.User, r *http.Request)  error {
	// Add any validation or business logic here before creating the user
    if err := s.repo.CreateUser(user, r); err != nil {
        return err
    }
	return nil
}

func (s *Service) GetUserByID(r *http.Request) (*entity.User, error) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}
	user, userErr := s.repo.GetUserByID(uint(id))
	if userErr != nil {
		return nil, userErr
	}
	return user, nil
}

// GetUserDetails retrieves a user by ID
func (s *Service) GetUserDetails(r *http.Request) (*entity.ResponseUser, error) {
	id, err := strconv.ParseUint(r.PathValue("id"), 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID")
	}
	user, userErr := s.repo.GetUser(uint(id))
	if userErr != nil {
		return nil, userErr
	}
	return user, nil
}

// UpdateUser updates an existing user
func (s *Service) UpdateUser(r *http.Request, user *entity.UpdateUser)  error {
	// Call repository to update user
	oldUser, err := s.GetUserByID(r)
	if err != nil {
		return err
	}

	err2 := s.repo.UpdateUser(oldUser, user, r)
	if err2 != nil {
		return err2
	}
	return  nil
}

// DeleteUser deletes a user by ID
func (s *Service) DeleteUser(r *http.Request) error {
	// Call repository to delete user
	user, err := s.GetUserByID(r)
	if err != nil {
		return err
	}

	err2 := s.repo.DeleteUser(user, r)
	if err2 != nil {
		return err2
	}

	return nil
}

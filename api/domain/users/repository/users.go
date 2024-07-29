package repository

import (
	"net/http"

	"github.com/JubaerHossain/cn-api/domain/users/entity"
)


// UserRepository defines methods for user data access
type UserRepository interface {
	GetUsers(r *http.Request) (*entity.UserResponsePagination, error)
	GetUserByID(userID uint) (*entity.User, error)
	GetUser(userID uint) (*entity.ResponseUser, error)
	CreateUser(user *entity.User, r *http.Request)  error
	UpdateUser(oldUser *entity.User, user *entity.UpdateUser, r *http.Request) error
	DeleteUser(user *entity.User, r *http.Request) error
}
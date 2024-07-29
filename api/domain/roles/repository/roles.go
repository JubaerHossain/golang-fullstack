package repository

import (
	"net/http"

	"github.com/JubaerHossain/cn-api/domain/roles/entity"
)


// RoleRepository defines methods for role data access
type RoleRepository interface {
	GetRoles(r *http.Request) (*entity.RoleResponsePagination, error)
	GetRoleByID(roleID uint) (*entity.Role, error)
	GetRole(roleID uint) (*entity.ResponseRole, error)
	CreateRole(role *entity.Role, r *http.Request)  error
	UpdateRole(oldRole *entity.Role, role *entity.UpdateRole, r *http.Request) error
	DeleteRole(role *entity.Role, r *http.Request) error
}
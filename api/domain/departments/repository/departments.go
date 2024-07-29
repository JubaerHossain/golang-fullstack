package repository

import (
	"net/http"

	"github.com/JubaerHossain/cn-api/domain/departments/entity"
)


// DepartmentRepository defines methods for department data access
type DepartmentRepository interface {
	GetDepartments(r *http.Request) (*entity.DepartmentResponsePagination, error)
	GetDepartmentByID(departmentID uint) (*entity.Department, error)
	GetDepartment(departmentID uint) (*entity.ResponseDepartment, error)
	CreateDepartment(department *entity.Department, r *http.Request)  error
	UpdateDepartment(oldDepartment *entity.Department, department *entity.UpdateDepartment, r *http.Request) error
	DeleteDepartment(department *entity.Department, r *http.Request) error
}
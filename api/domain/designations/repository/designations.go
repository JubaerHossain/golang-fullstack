package repository

import (
	"net/http"

	"github.com/JubaerHossain/cn-api/domain/designations/entity"
)


// DesignationRepository defines methods for designation data access
type DesignationRepository interface {
	GetDesignations(r *http.Request) (*entity.DesignationResponsePagination, error)
	GetDesignationByID(designationID uint) (*entity.Designation, error)
	GetDesignation(designationID uint) (*entity.ResponseDesignation, error)
	CreateDesignation(designation *entity.Designation, r *http.Request)  error
	UpdateDesignation(oldDesignation *entity.Designation, designation *entity.UpdateDesignation, r *http.Request) error
	DeleteDesignation(designation *entity.Designation, r *http.Request) error
}
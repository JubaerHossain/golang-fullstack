package departmentHttp

import (
	"net/http"

	"github.com/JubaerHossain/cn-api/domain/departments/entity"
	"github.com/JubaerHossain/cn-api/domain/departments/service"
	"github.com/JubaerHossain/rootx/pkg/core/app"
	utilQuery "github.com/JubaerHossain/rootx/pkg/query"
	"github.com/JubaerHossain/rootx/pkg/utils"
)

// Handler handles API requests
type Handler struct {
	App *service.Service
}

// NewHandler creates a new instance of Handler
func NewHandler(app *app.App) *Handler {
	return &Handler{
		App: service.NewService(app),
	}
}

// @Summary Get all departments
// @Description Get details of all departments
// @Tags departments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param search query string false "Search query"
// @Param sort query string false "Sort by"
// @Param status query bool false "Filter by status"
// @Success 200 {object} entity.DepartmentResponsePagination
// @Router /departments [get]
func (h *Handler) GetDepartments(w http.ResponseWriter, r *http.Request) {
	// Implement GetDepartments handler
	departments, err := h.App.GetDepartments(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to fetch departments")
		return
	}
	// Write response
	utils.JsonResponse(w, http.StatusOK, map[string]interface{}{
		"results": departments,
	})
}

// @Summary Create a new Department
// @Description Create a new Department
// @Tags departments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 201 {object} map[string]interface{}
// @Param department body entity.Department true "The Department to be created"
// @Router /departments [post]
func (h *Handler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	// Implement CreateDepartment handler
	var newDepartment entity.Department

	pareErr := utilQuery.BodyParse(&newDepartment, w, r, true) // Parse request body and validate it
	if pareErr != nil {
		return
	}

	// Call the CreateDepartment function to create the role
	err := h.App.CreateDepartment(&newDepartment, r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Write response
	utils.WriteJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"message": "Department created successfully",
	})
}


func (h *Handler) GetDepartmentByID(w http.ResponseWriter, r *http.Request) {
	department, err := h.App.GetDepartmentByID(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Write response
	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Department fetched successfully",
		"results": department,
	})

}

// @Summary Get detailed information about a Department by ID
// @Description Get detailed information about a Department by ID
// @Tags departments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} entity.ResponseDepartment
// @Param id path string true "The ID of the Department"
// @Router /departments/{id}/details [get]
func (h *Handler) GetDepartmentDetails(w http.ResponseWriter, r *http.Request) {
	department, err := h.App.GetDepartmentDetails(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Write response
	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Department fetched successfully",
		"results": department,
	})

}

// @Summary Update an existing Department
// @Description Update an existing Department
// @Tags departments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 201 {object} map[string]interface{}
// @Param id path string true "The ID of the Department"
// @Param department body entity.UpdateDepartment true "Updated Department object"
// @Router /departments/{id} [put]
func (h *Handler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	// Implement UpdateDepartment handler
	var updateDepartment entity.UpdateDepartment
	pareErr := utilQuery.BodyParse(&updateDepartment, w, r, true) // Parse request body and validate it
	if pareErr != nil {
		return
	}

	// Call the CreateDepartment function to create the department
	err := h.App.UpdateDepartment(r, &updateDepartment)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Write response
	utils.WriteJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"message": "Department updated successfully",
	})
}

// @Summary Delete a Department
// @Description Delete a Department
// @Tags departments
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Param id path string true "The ID of the Department"
// @Router /departments/{id} [delete]
func (h *Handler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	// Implement DeleteDepartment handler
	err := h.App.DeleteDepartment(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Write response
	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Department deleted successfully",
	})
}

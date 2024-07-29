package roleHttp

import (
	"net/http"

	"github.com/JubaerHossain/cn-api/domain/roles/entity"
	"github.com/JubaerHossain/cn-api/domain/roles/service"
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

// @Summary Get all roles
// @Description Get details of all roles
// @Tags roles
// @Accept json
// @Produce json
// @Success 200 {object} entity.RoleResponsePagination
// @Router /roles [get]
func (h *Handler) GetRoles(w http.ResponseWriter, r *http.Request) {
	// Implement GetRoles handler
	roles, err := h.App.GetRoles(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to fetch roles")
		return
	}
	// Write response
	utils.JsonResponse(w, http.StatusOK, map[string]interface{}{
		"results": roles,
	})
}

// @Summary Create a new Role
// @Description Create a new Role
// @Tags roles
// @Accept json
// @Produce json
// @Param role body entity.Role true "The Role to be created"
// @Router /roles [post]
func (h *Handler) CreateRole(w http.ResponseWriter, r *http.Request) {
	// Implement CreateRole handler
	var newRole entity.Role

	pareErr := utilQuery.BodyParse(&newRole, w, r, true) // Parse request body and validate it
	if pareErr != nil {
		return
	}

	// Call the CreateRole function to create the role
	err := h.App.CreateRole(&newRole, r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Write response
	utils.WriteJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"message": "Role created successfully",
	})
}


func (h *Handler) GetRoleByID(w http.ResponseWriter, r *http.Request) {
	role, err := h.App.GetRoleByID(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Write response
	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Role fetched successfully",
		"results": role,
	})

}

// @Summary Get detailed information about a Role by ID
// @Description Get detailed information about a Role by ID
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "The ID of the Role"
// @Router /roles/{id}/details [get]
func (h *Handler) GetRoleDetails(w http.ResponseWriter, r *http.Request) {
	role, err := h.App.GetRoleDetails(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Write response
	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Role fetched successfully",
		"results": role,
	})

}

// @Summary Update an existing Role
// @Description Update an existing Role
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "The ID of the Role"
// @Param role body entity.UpdateRole true "Updated Role object"
// @Router /roles/{id} [put]
func (h *Handler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	// Implement UpdateRole handler
	var updateRole entity.UpdateRole
	pareErr := utilQuery.BodyParse(&updateRole, w, r, true) // Parse request body and validate it
	if pareErr != nil {
		return
	}

	// Call the CreateRole function to create the role
	err := h.App.UpdateRole(r, &updateRole)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Write response
	utils.WriteJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"message": "Role updated successfully",
	})
}

// @Summary Delete a Role
// @Description Delete a Role
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "The ID of the Role"
// @Router /roles/{id} [delete]
func (h *Handler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	// Implement DeleteRole handler
	err := h.App.DeleteRole(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Write response
	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Role deleted successfully",
	})
}

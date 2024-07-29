package designationHttp

import (
	"net/http"

	"github.com/JubaerHossain/cn-api/domain/designations/entity"
	"github.com/JubaerHossain/cn-api/domain/designations/service"
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

// @Summary Get all designations
// @Description Get details of all designations
// @Tags designations
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param search query string false "Search query"
// @Param sort query string false "Sort by"
// @Param status query bool false "Filter by status"
// @Success 200 {object} entity.DesignationResponsePagination
// @Router /designations [get]
func (h *Handler) GetDesignations(w http.ResponseWriter, r *http.Request) {
	// Implement GetDesignations handler
	designations, err := h.App.GetDesignations(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to fetch designations")
		return
	}
	// Write response
	utils.JsonResponse(w, http.StatusOK, map[string]interface{}{
		"results": designations,
	})
}

// @Summary Create a new Designation
// @Description Create a new Designation
// @Tags designations
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 201 {object} map[string]interface{}
// @Param designation body entity.Designation true "The Designation to be created"
// @Router /designations [post]
func (h *Handler) CreateDesignation(w http.ResponseWriter, r *http.Request) {
	// Implement CreateDesignation handler
	var newDesignation entity.Designation

	pareErr := utilQuery.BodyParse(&newDesignation, w, r, true) // Parse request body and validate it
	if pareErr != nil {
		return
	}

	// Call the CreateDesignation function to create the role
	err := h.App.CreateDesignation(&newDesignation, r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Write response
	utils.WriteJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"message": "Designation created successfully",
	})
}


func (h *Handler) GetDesignationByID(w http.ResponseWriter, r *http.Request) {
	designation, err := h.App.GetDesignationByID(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Write response
	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Designation fetched successfully",
		"results": designation,
	})

}

// @Summary Get detailed information about a Designation by ID
// @Description Get detailed information about a Designation by ID
// @Tags designations
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} entity.ResponseDesignation
// @Param id path string true "The ID of the Designation"
// @Router /designations/{id}/details [get]
func (h *Handler) GetDesignationDetails(w http.ResponseWriter, r *http.Request) {
	designation, err := h.App.GetDesignationDetails(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Write response
	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Designation fetched successfully",
		"results": designation,
	})

}

// @Summary Update an existing Designation
// @Description Update an existing Designation
// @Tags designations
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 201 {object} map[string]interface{}
// @Param id path string true "The ID of the Designation"
// @Param designation body entity.UpdateDesignation true "Updated Designation object"
// @Router /designations/{id} [put]
func (h *Handler) UpdateDesignation(w http.ResponseWriter, r *http.Request) {
	// Implement UpdateDesignation handler
	var updateDesignation entity.UpdateDesignation
	pareErr := utilQuery.BodyParse(&updateDesignation, w, r, true) // Parse request body and validate it
	if pareErr != nil {
		return
	}

	// Call the CreateDesignation function to create the designation
	err := h.App.UpdateDesignation(r, &updateDesignation)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Write response
	utils.WriteJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"message": "Designation updated successfully",
	})
}

// @Summary Delete a Designation
// @Description Delete a Designation
// @Tags designations
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Param id path string true "The ID of the Designation"
// @Router /designations/{id} [delete]
func (h *Handler) DeleteDesignation(w http.ResponseWriter, r *http.Request) {
	// Implement DeleteDesignation handler
	err := h.App.DeleteDesignation(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Write response
	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Designation deleted successfully",
	})
}

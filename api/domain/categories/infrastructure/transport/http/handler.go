package categoryHttp

import (
	"net/http"

	"github.com/JubaerHossain/cn-api/domain/categories/entity"
	"github.com/JubaerHossain/cn-api/domain/categories/service"
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

// @Summary Get all categories
// @Description Get details of all categories
// @Tags categories
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} entity.CategoryResponsePagination
// @Router /categories [get]
func (h *Handler) GetCategories(w http.ResponseWriter, r *http.Request) {
	// Implement GetCategories handler
	categories, err := h.App.GetCategories(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to fetch categories")
		return
	}
	// Write response
	utils.JsonResponse(w, http.StatusOK, categories.Data)
}

// @Summary Create a new Category
// @Description Create a new Category
// @Tags categories
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "Category created successfully"
// @Param category body entity.Category true "The Category to be created"
// @Router /categories [post]
func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	// Implement CreateCategory handler
	var newCategory entity.Category

	pareErr := utilQuery.BodyParse(&newCategory, w, r, true) // Parse request body and validate it
	if pareErr != nil {
		return
	}

	// Call the CreateCategory function to create the role
	err := h.App.CreateCategory(&newCategory, r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Write response
	utils.WriteJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"message": "Category created successfully",
	})
}

func (h *Handler) GetCategoryByID(w http.ResponseWriter, r *http.Request) {
	category, err := h.App.GetCategoryByID(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Write response
	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Category fetched successfully",
		"results": category,
	})

}

// @Summary Get detailed information about a Category by ID
// @Description Get detailed information about a Category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} entity.ResponseCategory
// @Param id path string true "The ID of the Category"
// @Router /categories/{id}/details [get]
func (h *Handler) GetCategoryDetails(w http.ResponseWriter, r *http.Request) {
	category, err := h.App.GetCategoryDetails(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Write response
	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Category fetched successfully",
		"results": category,
	})

}

// @Summary Update an existing Category
// @Description Update an existing Category
// @Tags categories
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "Category updated successfully"
// @Param id path string true "The ID of the Category"
// @Param category body entity.UpdateCategory true "Updated Category object"
// @Router /categories/{id} [put]
func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	// Implement UpdateCategory handler
	var updateCategory entity.UpdateCategory
	pareErr := utilQuery.BodyParse(&updateCategory, w, r, true) // Parse request body and validate it
	if pareErr != nil {
		return
	}

	// Call the CreateCategory function to create the category
	err := h.App.UpdateCategory(r, &updateCategory)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Write response
	utils.WriteJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"message": "Category updated successfully",
	})
}

// @Summary Delete a Category
// @Description Delete a Category
// @Tags categories
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "Category deleted successfully"
// @Param id path string true "The ID of the Category"
// @Router /categories/{id} [delete]
func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	// Implement DeleteCategory handler
	err := h.App.DeleteCategory(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Write response
	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Category deleted successfully",
	})
}

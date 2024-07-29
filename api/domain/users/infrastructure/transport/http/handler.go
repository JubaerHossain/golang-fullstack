package userHttp

import (
	"net/http"

	"github.com/JubaerHossain/cn-api/domain/users/entity"
	"github.com/JubaerHossain/cn-api/domain/users/service"
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

// @Summary Get all users
// @Description Get details of all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} entity.UserResponsePagination
// @Router /users [get]
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Implement GetUsers handler
	users, err := h.App.GetUsers(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}
	// Write response
	utils.JsonResponse(w, http.StatusOK, map[string]interface{}{
		"results": users,
	})
}

// @Summary Create a new User
// @Description Create a new User
// @Tags users
// @Accept json
// @Produce json
// @Param user body entity.User true "The User to be created"
// @Router /users [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Implement CreateUser handler
	var newUser entity.User

	pareErr := utilQuery.BodyParse(&newUser, w, r, true) // Parse request body and validate it
	if pareErr != nil {
		return
	}

	// Call the CreateUser function to create the role
	err := h.App.CreateUser(&newUser, r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Write response
	utils.WriteJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"message": "User created successfully",
	})
}


func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	user, err := h.App.GetUserByID(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Write response
	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "User fetched successfully",
		"results": user,
	})

}

// @Summary Get detailed information about a User by ID
// @Description Get detailed information about a User by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "The ID of the User"
// @Router /users/{id}/details [get]
func (h *Handler) GetUserDetails(w http.ResponseWriter, r *http.Request) {
	user, err := h.App.GetUserDetails(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Write response
	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "User fetched successfully",
		"results": user,
	})

}

// @Summary Update an existing User
// @Description Update an existing User
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "The ID of the User"
// @Param user body entity.UpdateUser true "Updated User object"
// @Router /users/{id} [put]
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Implement UpdateUser handler
	var updateUser entity.UpdateUser
	pareErr := utilQuery.BodyParse(&updateUser, w, r, true) // Parse request body and validate it
	if pareErr != nil {
		return
	}

	// Call the CreateUser function to create the user
	err := h.App.UpdateUser(r, &updateUser)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Write response
	utils.WriteJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"message": "User updated successfully",
	})
}

// @Summary Delete a User
// @Description Delete a User
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "The ID of the User"
// @Router /users/{id} [delete]
func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Implement DeleteUser handler
	err := h.App.DeleteUser(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Write response
	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "User deleted successfully",
	})
}

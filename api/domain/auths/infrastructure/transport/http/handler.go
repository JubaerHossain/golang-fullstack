package authHttp

import (
	"net/http"

	"github.com/JubaerHossain/cn-api/domain/auths/entity"
	"github.com/JubaerHossain/cn-api/domain/auths/service"
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

// @Summary  Get sign-in
// @Description  Get sign-in
// @Tags sign-in
// @Accept json
// @Produce json
// @Param auth body entity.LoginUser true "The Auth to be created"
// @Success 200 {object} map[string]interface{}
// @Router /auth/sign-in [post]
func (h *Handler) GetSignIn(w http.ResponseWriter, r *http.Request) {

	var newUser entity.LoginUser

	pareErr := utilQuery.BodyParse(&newUser, w, r, true) // Parse request body and validate it
	if pareErr != nil {
		return
	}
	auth, err := h.App.GetSignIn(&newUser, r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Write response
	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Sign-in successful",
		"results": auth,
	})

}

// @Summary  Get refresh-token
// @Description  Get refresh-token
// @Tags refresh-token
// @Accept json
// @Produce json
// @Param auth body entity.RefreshToken true "The Auth to be created"
// @Success 200 {object} map[string]interface{}
// @Router /auth/refresh-token [post]
func (h *Handler) GetRefreshToken(w http.ResponseWriter, r *http.Request) {
	var refreshToken entity.RefreshToken
	if pareErr := utilQuery.BodyParse(&refreshToken, w, r, true); pareErr != nil {
		return
	}
	auth, err := h.App.GetRefreshToken(&refreshToken, r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Write response
	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "Refresh-token successful",
		"results": auth,
	})

}

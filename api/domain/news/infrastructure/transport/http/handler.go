package newsHttp

import (
	"net/http"

	"github.com/JubaerHossain/cn-api/domain/news/entity"
	"github.com/JubaerHossain/cn-api/domain/news/service"
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

// @Summary Get all news
// @Description Get details of all news
// @Tags news
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} entity.NewsResponsePagination
// @Router /news [get]
func (h *Handler) GetNewses(w http.ResponseWriter, r *http.Request) {
	// Implement GetNewses handler
	news, err := h.App.GetNewses(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to fetch news")
		return
	}
	// Write response
	utils.JsonResponse(w, http.StatusOK, map[string]interface{}{
		"results": news,
	})
}

// @Summary Create a new News
// @Description Create a new News
// @Tags news
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "News created successfully"
// @Param news body entity.News true "The News to be created"
// @Router /news [post]
func (h *Handler) CreateNews(w http.ResponseWriter, r *http.Request) {
	// Implement CreateNews handler
	var newNews entity.News

	pareErr := utilQuery.BodyParse(&newNews, w, r, true) // Parse request body and validate it
	if pareErr != nil {
		return
	}

	// Call the CreateNews function to create the role
	err := h.App.CreateNews(&newNews, r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Write response
	utils.WriteJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"message": "News created successfully",
	})
}


func (h *Handler) GetNewsByID(w http.ResponseWriter, r *http.Request) {
	news, err := h.App.GetNewsByID(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Write response
	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "News fetched successfully",
		"results": news,
	})

}

// @Summary Get detailed information about a News by ID
// @Description Get detailed information about a News by ID
// @Tags news
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} entity.ResponseNews
// @Param id path string true "The ID of the News"
// @Router /news/{id}/details [get]
func (h *Handler) GetNewsDetails(w http.ResponseWriter, r *http.Request) {
	news, err := h.App.GetNewsDetails(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Write response
	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "News fetched successfully",
		"results": news,
	})

}

// @Summary Update an existing News
// @Description Update an existing News
// @Tags news
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "News updated successfully"
// @Param id path string true "The ID of the News"
// @Param news body entity.UpdateNews true "Updated News object"
// @Router /news/{id} [put]
func (h *Handler) UpdateNews(w http.ResponseWriter, r *http.Request) {
	// Implement UpdateNews handler
	var updateNews entity.UpdateNews
	pareErr := utilQuery.BodyParse(&updateNews, w, r, true) // Parse request body and validate it
	if pareErr != nil {
		return
	}

	// Call the CreateNews function to create the news
	err := h.App.UpdateNews(r, &updateNews)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Write response
	utils.WriteJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"message": "News updated successfully",
	})
}

// @Summary Delete a News
// @Description Delete a News
// @Tags news
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{} "News deleted successfully"
// @Param id path string true "The ID of the News"
// @Router /news/{id} [delete]
func (h *Handler) DeleteNews(w http.ResponseWriter, r *http.Request) {
	// Implement DeleteNews handler
	err := h.App.DeleteNews(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Write response
	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"message": "News deleted successfully",
	})
}

// @Summary Get all scroll news
// @Description Get details of all scroll news
// @Tags news
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} entity.ScrollNewsResponse
// @Router /public/v1/breaking-scrolling-news [get]
func (h *Handler) GetBreakingScrollingNews(w http.ResponseWriter, r *http.Request) {
	// Implement GetBreakingScrollingNews handler
	news, err := h.App.GetBreakingScrollingNews(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to fetch news")
		return
	}
	// Write response
	utils.JsonResponse(w, http.StatusOK,news)
}

// @Summary Get all thumbnail news
// @Description Get details of all thumbnail news
// @Tags news
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} entity.ThumbnailNewsResponse
// @Router /public/v1/breaking-thumbnail-news [get]
func (h *Handler) GetBreakingThumbnailNews(w http.ResponseWriter, r *http.Request) {
	// Implement GetBreakingThumbnailNews handler
	news, err := h.App.GetBreakingThumbnailNews(r)
	if err != nil {
		utils.WriteJSONError(w, http.StatusInternalServerError, "Failed to fetch news")
		return
	}
	// Write response
	utils.JsonResponse(w, http.StatusOK,news)
}

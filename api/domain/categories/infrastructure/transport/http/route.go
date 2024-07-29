package categoryHttp

import (
	"net/http"

	"github.com/JubaerHossain/rootx/pkg/core/app"
	"github.com/JubaerHossain/rootx/pkg/core/middleware"
)

// CategoryRouter registers routes for API endpoints
func CategoryRouter(router *http.ServeMux, application *app.App) http.Handler {

	handler := NewHandler(application)
	// Register category routes

	router.Handle("GET /categories", middleware.LimiterMiddleware(http.HandlerFunc(handler.GetCategories)))
	// router.Handle("POST /categories", middleware.LimiterMiddleware(http.HandlerFunc(handler.CreateCategory)))
	// router.Handle("GET /categories/{id}", middleware.LimiterMiddleware(http.HandlerFunc(handler.GetCategoryDetails)))
	// router.Handle("PUT /categories/{id}", middleware.LimiterMiddleware(http.HandlerFunc(handler.UpdateCategory)))
	// router.Handle("DELETE /categories/{id}", middleware.LimiterMiddleware(http.HandlerFunc(handler.DeleteCategory)))

	return router
}

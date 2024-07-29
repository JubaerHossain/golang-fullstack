package newsHttp

import (
	"net/http"
	"github.com/JubaerHossain/rootx/pkg/core/app"
    "github.com/JubaerHossain/rootx/pkg/core/middleware"
)

// NewsRouter registers routes for API endpoints
func NewsRouter(router *http.ServeMux, application *app.App) http.Handler {

	
	handler := NewHandler(application)
	// Register news routes

	router.Handle("GET /news", middleware.LimiterMiddleware(http.HandlerFunc(handler.GetNewses)))


	router.Handle("GET /breaking-scrolling-news", middleware.LimiterMiddleware(http.HandlerFunc(handler.GetBreakingScrollingNews)))
	router.Handle("GET /breaking-thumbnail-news", middleware.LimiterMiddleware(http.HandlerFunc(handler.GetBreakingThumbnailNews)))
	
   

	return router
}

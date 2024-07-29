package roleHttp

import (
	"net/http"
	"github.com/JubaerHossain/rootx/pkg/core/app"
    "github.com/JubaerHossain/rootx/pkg/core/middleware"
)

// RoleRouter registers routes for API endpoints
func RoleRouter(application *app.App) http.Handler {
	router := http.NewServeMux()

	
	handler := NewHandler(application)
	// Register role routes

	router.Handle("GET /roles", middleware.LimiterMiddleware(http.HandlerFunc(handler.GetRoles)))
	router.Handle("POST /roles", middleware.LimiterMiddleware(http.HandlerFunc(handler.CreateRole)))
	router.Handle("GET /roles/{id}", middleware.LimiterMiddleware(http.HandlerFunc(handler.GetRoleDetails)))
	router.Handle("PUT /roles/{id}", middleware.LimiterMiddleware(http.HandlerFunc(handler.UpdateRole)))
	router.Handle("DELETE /roles/{id}", middleware.LimiterMiddleware(http.HandlerFunc(handler.DeleteRole)))
   

	return router
}

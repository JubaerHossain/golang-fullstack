package userHttp

import (
	"net/http"
	"github.com/JubaerHossain/rootx/pkg/core/app"
    "github.com/JubaerHossain/rootx/pkg/core/middleware"
)

// UserRouter registers routes for API endpoints
func UserRouter(application *app.App) http.Handler {
	router := http.NewServeMux()

	
	handler := NewHandler(application)
	// Register user routes

	router.Handle("GET /users", middleware.LimiterMiddleware(http.HandlerFunc(handler.GetUsers)))
	router.Handle("POST /users", middleware.LimiterMiddleware(http.HandlerFunc(handler.CreateUser)))
	router.Handle("GET /users/{id}", middleware.LimiterMiddleware(http.HandlerFunc(handler.GetUserDetails)))
	router.Handle("PUT /users/{id}", middleware.LimiterMiddleware(http.HandlerFunc(handler.UpdateUser)))
	router.Handle("DELETE /users/{id}", middleware.LimiterMiddleware(http.HandlerFunc(handler.DeleteUser)))
   

	return router
}

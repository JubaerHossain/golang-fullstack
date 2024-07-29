package designationHttp

import (
	"net/http"
	"github.com/JubaerHossain/rootx/pkg/core/app"
    "github.com/JubaerHossain/rootx/pkg/core/middleware"
)

// DesignationRouter registers routes for API endpoints
func DesignationRouter(application *app.App) http.Handler {
	router := http.NewServeMux()

	
	handler := NewHandler(application)
	// Register designation routes

	router.Handle("GET /designations", middleware.LimiterMiddleware(http.HandlerFunc(handler.GetDesignations)))
	router.Handle("POST /designations", middleware.LimiterMiddleware(http.HandlerFunc(handler.CreateDesignation)))
	router.Handle("GET /designations/{id}", middleware.LimiterMiddleware(http.HandlerFunc(handler.GetDesignationDetails)))
	router.Handle("PUT /designations/{id}", middleware.LimiterMiddleware(http.HandlerFunc(handler.UpdateDesignation)))
	router.Handle("DELETE /designations/{id}", middleware.LimiterMiddleware(http.HandlerFunc(handler.DeleteDesignation)))
   

	return router
}

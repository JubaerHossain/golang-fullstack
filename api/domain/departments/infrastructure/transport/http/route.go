package departmentHttp

import (
	"net/http"

	"github.com/JubaerHossain/rootx/pkg/core/app"
	"github.com/JubaerHossain/rootx/pkg/core/middleware"
)

// DepartmentRouter registers routes for API endpoints
func DepartmentRouter(router *http.ServeMux, application *app.App) http.Handler {

	handler := NewHandler(application)
	// Register department routes

	router.Handle("GET /departments", middleware.LimiterMiddleware(http.HandlerFunc(handler.GetDepartments)))
	router.Handle("POST /departments", middleware.LimiterMiddleware(http.HandlerFunc(handler.CreateDepartment)))
	router.Handle("GET /departments/{id}", middleware.LimiterMiddleware(http.HandlerFunc(handler.GetDepartmentDetails)))
	router.Handle("PUT /departments/{id}", middleware.LimiterMiddleware(http.HandlerFunc(handler.UpdateDepartment)))
	router.Handle("DELETE /departments/{id}", middleware.LimiterMiddleware(http.HandlerFunc(handler.DeleteDepartment)))

	return router
}

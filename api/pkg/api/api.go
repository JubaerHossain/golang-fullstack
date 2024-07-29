package api

import (
	"net/http"

	departmentHttp "github.com/JubaerHossain/cn-api/domain/departments/infrastructure/transport/http"
	"github.com/JubaerHossain/rootx/pkg/core/app"
)

func APIRouter(application *app.App) http.Handler {
	router := http.NewServeMux()
	//Register department routes
	departmentHttp.DepartmentRouter(router, application)

	return router
}

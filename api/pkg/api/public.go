package api

import (
	"net/http"

	categoryHttp "github.com/JubaerHossain/cn-api/domain/categories/infrastructure/transport/http"
	newsHttp "github.com/JubaerHossain/cn-api/domain/news/infrastructure/transport/http"
	"github.com/JubaerHossain/rootx/pkg/core/app"
)

func PublicAPIRouter(application *app.App) http.Handler {
	router := http.NewServeMux()

	//public routes
	categoryHttp.CategoryRouter(router, application)
	newsHttp.NewsRouter(router, application)

	return router
}

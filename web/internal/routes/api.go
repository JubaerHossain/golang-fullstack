package routes

import (
	"net/http"

	"github.com/JubaerHossain/golang-htmx-starter/pkg/core"
	"github.com/labstack/echo/v4"
)

func BindApiRoute(a *core.App) {
	apiGroup := a.Echo.Group("/api")
	apiGroup.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

}

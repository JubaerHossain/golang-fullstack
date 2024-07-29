package routes

import (
	"github.com/JubaerHossain/golang-htmx-starter/internal/handler"
	"github.com/JubaerHossain/golang-htmx-starter/pkg/core"
	"github.com/labstack/echo/v4"
)

func BindWebRoute(a *core.App) {
	a.Echo.GET("/", func(c echo.Context) error {
		return handler.Home(c, a)
	})

}

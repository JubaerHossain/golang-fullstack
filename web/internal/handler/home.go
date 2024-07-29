package handler

import (
	"net/http"

	"github.com/JubaerHossain/golang-htmx-starter/pkg/core"
	"github.com/labstack/echo/v4"
)

// RootPath handles the root route to display index.gohtml
func Home(c echo.Context, a *core.App) error {

	// Create a data map to pass to the renderer
	data := map[string]interface{}{
		"Title": "Welcome to Htmx",
	}

	// Use Echo's built-in rendering method
	return c.Render(http.StatusOK, "index.html", data)
}

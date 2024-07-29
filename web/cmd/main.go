package main

import (
	"log"
	"strconv"

	"github.com/JubaerHossain/golang-htmx-starter/internal/routes"
	"github.com/JubaerHossain/golang-htmx-starter/pkg/core"
)

func main() {
	// Initialize logger

	// Create new app
	app, err := core.StartApp()
	if err != nil {
		log.Fatalf("failed to create new app: %v", err)
	}

	// Start the app
	router(app)

	// Or you can perform any other actions with the app instance
}

func router(app *core.App) {
	// Bind web routes
	routes.BindWebRoute(app)

	// Bind API routes
	routes.BindApiRoute(app)

	// Start the server
	app.Echo.Logger.Fatal(app.Echo.Start(":" + strconv.Itoa(app.HttpPort)))
}

// Bindings are created in the serve package and listed here

package html

import (
	"embed"

	"github.com/labstack/echo/v4"
)

// Embedding the public directory
//
//go:embed public/*
var public embed.FS
var PublicFS = echo.MustSubFS(public, "public/templates")
var StaticFS = echo.MustSubFS(public, "public/assets")

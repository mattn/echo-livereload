package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mattn/echo-livereload"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	})
	e.Use(middleware.Static("assets"))
	e.Use(livereload.LiveReload())
	e.Logger.Fatal(e.Start(":8989"))
}

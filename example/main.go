package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mattn/echo-livereload"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Static("assets"))
	e.Use(livereload.LiveReload())
	e.Logger.Fatal(e.Start(":8080"))
}

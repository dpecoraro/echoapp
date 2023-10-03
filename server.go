package main

import (
	"echo-app/configs"
	"echo-app/routes"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, &echo.Map{"data": "Hello from Echo & mongoDB"})
	})

	//Routes
	routes.UserRoute(e)

	configs.ConnectDB()
	e.Logger.Fatal(e.Start(":6000"))
}

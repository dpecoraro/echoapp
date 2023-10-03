package routes

import (
	"echo-app/controllers"

	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Echo) {
	//All routes related to user context
	e.POST("/users", controllers.CreateUser)
	e.GET("/users/:userId", controllers.GetUser)
	e.GET("/users", controllers.GetAllUsers)
	e.PUT("/users/:userId", controllers.EditUser)
	e.DELETE("/users/:userId", controllers.DeleteUser)
}

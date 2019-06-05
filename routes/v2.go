package routes

import (
	"github.com/labstack/echo"

	"github.com/kilfu0701/echo_app/controllers/api"
)

func v2(e *echo.Echo) {
	v2 := e.Group("/v2")

	// auth
	v2.GET("/hello", api.Hello)
	v2.GET("/auth", api.Auth)
	v2.GET("/auth/:token", api.AuthWithToken)
}

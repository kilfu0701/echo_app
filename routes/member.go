package routes

import (
	"github.com/labstack/echo"

	"github.com/kilfu0701/echo_app/controllers"
)

func member(e *echo.Echo) {
	member := e.Group("/member")

	member.GET("/top", controllers.MemberTop)
	member.GET("/register", controllers.MemberRegister)
}

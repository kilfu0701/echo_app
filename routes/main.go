package routes

import (
	"github.com/labstack/echo"
)

func Load(e *echo.Echo) {
	member(e)
	v2(e)
}

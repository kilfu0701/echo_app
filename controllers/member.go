package controllers

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/kilfu0701/echo_app/core"
)

// http://localhost:1323/member/top
func MemberTop(c echo.Context) error {
	ac := c.(*core.AppContext)
	res := map[string]interface{}{}
	return ac.JSON(http.StatusOK, res)
}

// http://localhost:1323/member/register
func MemberRegister(c echo.Context) error {
	ac := c.(*core.AppContext)
	res := map[string]interface{}{}
	return ac.JSON(http.StatusOK, res)
}

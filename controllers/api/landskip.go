package api

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/kilfu0701/echo_app/core"
)

func GetLandskipList(c echo.Context) error {
	ac := c.(*core.AppContext)
	res := map[string]interface{}{
		"code":    200,
		"message": "ok",
	}
	return ac.JSON(http.StatusOK, res)
}

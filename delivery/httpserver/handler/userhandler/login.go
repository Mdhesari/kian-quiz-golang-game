package userhandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) Login(c echo.Context) error {
	return c.JSON(http.StatusOK, `{"status": "ok"}`)
}
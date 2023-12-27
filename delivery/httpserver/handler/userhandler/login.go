package userhandler

import (
	"mdhesari/kian-quiz-golang-game/param"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) Login(c echo.Context) error {
	req := param.LoginRequest{}

	c.Bind(&req)

	return c.JSON(http.StatusOK, param.LoginResponse{})
}
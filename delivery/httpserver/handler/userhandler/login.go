package userhandler

import (
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) Login(c echo.Context) error {
	var req param.LoginRequest

	c.Bind(&req)

	if fields, err := h.userValidator.ValidateLoginRequest(req); err != nil {
		msg, code := richerror.Error(err)

		return c.JSON(code, echo.Map{
			"message": msg,
			"errors":  fields,
		})
	}

	res, err := h.userSrv.Login(req)
	if err != nil {
		msg, code := richerror.Error(err)

		return c.JSON(code, echo.Map{
			"message": msg,
		})
	}

	return c.JSON(http.StatusOK, res)
}

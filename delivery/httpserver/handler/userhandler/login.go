package userhandler

import (
	"fmt"
	"mdhesari/kian-quiz-golang-game/param"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) Login(c echo.Context) error {
	req := param.LoginRequest{}

	c.Bind(&req)

	res := h.userSrv.Login(req)

	fmt.Println(res)

	return c.JSON(http.StatusOK, res)
}

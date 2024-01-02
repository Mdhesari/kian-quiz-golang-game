package userhandler

import (
	"fmt"
	"log"
	"mdhesari/kian-quiz-golang-game/param"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) Login(c echo.Context) error {
	var req param.LoginRequest

	c.Bind(&req)

	log.Println(req)

	res := h.userSrv.Login(req)

	fmt.Println(res)

	return c.JSON(http.StatusOK, res)
}

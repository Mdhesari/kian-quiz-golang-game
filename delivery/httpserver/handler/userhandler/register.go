package userhandler

import (
	"log"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/service/userservice"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) Register(c echo.Context) error {
	var req param.RegisterRequest

	if err := c.Bind(&req); err != nil {
		log.Println("Register eror: ", err)

		return c.JSON(http.StatusBadRequest, param.RegisterResponse{})
	}

	res := h.userSrv.Register(userservice.UserForm{
		Name:     req.Name,
		Email:    req.Email,
		Mobile:   req.Mobile,
		Password: req.Password,
	})

	return c.JSON(http.StatusOK, res)
}

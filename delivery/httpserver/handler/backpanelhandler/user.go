package backpanelhandler

import (
	"mdhesari/kian-quiz-golang-game/entity"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) ListUsers(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"Items": []entity.User{
			{
				ID:     [12]byte{0},
				Name:   "User1",
				Email:  "user1@gmail.com",
				Mobile: "09121234567",
			},
		},
	})
}

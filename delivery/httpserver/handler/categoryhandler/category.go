package categoryhandler

import (
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) GetAll(c echo.Context) error {
	res, err := h.categorySrv.GetAll(c.Request().Context(), param.CategoryParam{})
	if err != nil {
		msg, code := richerror.Error(err)

		return c.JSON(code, echo.Map{
			"message": msg,
		})
	}

	return c.JSON(http.StatusOK, res)
}

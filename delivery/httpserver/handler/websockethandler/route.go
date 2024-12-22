package websockethandler

import (
	// "mdhesari/kian-quiz-golang-game/delivery/httpserver/middleware"

	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(r *echo.Echo) {
	group := r.Group("/websocket")

	group.GET("/:channel", h.Channel)
}
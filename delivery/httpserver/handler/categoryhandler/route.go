package categoryhandler

import (
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/middleware"

	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(r *echo.Echo) {
	group := r.Group("/categories")

	group.GET("", h.GetAll, middleware.Auth(h.authSrv, h.authConfig), middleware.Presence(h.presenceSrv))
}
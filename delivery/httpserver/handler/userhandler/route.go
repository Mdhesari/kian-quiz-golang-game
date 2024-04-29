package userhandler

import (
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/middleware"

	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(r *echo.Echo) {
	group := r.Group("/users")

	group.GET("/profile", h.Profile, middleware.Auth(h.authSrv, h.authConfig))
	group.POST("/login", h.Login)
	group.POST("/register", h.Register)
}

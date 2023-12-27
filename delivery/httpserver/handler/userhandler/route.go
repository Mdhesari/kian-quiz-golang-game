package userhandler

import (
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(r *echo.Echo) {
	group := r.Group("/users")

	group.POST("/login", h.Login)
}
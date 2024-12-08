package prometheushandler

import (
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
)

func (h Handler) Metrics(c echo.Context) error {
	handler := echoprometheus.NewHandler()

	return handler(c)
}

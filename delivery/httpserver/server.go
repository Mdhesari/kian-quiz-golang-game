package httpserver

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handler interface {
	SetRoutes(r *echo.Echo)
}

type HTTPServer struct {
	Port int `koanf:"port"`
}

type Config struct {
	HTTPServer HTTPServer `koanf:"http_server"`
}

type Server struct {
	config   Config
	handlers []Handler
}

func New(c Config, h []Handler) Server {
	return Server{
		handlers: h,
		config:   c,
	}
}

func (s Server) Serve() {
	echo := echo.New()

	echo.Use(middleware.Logger())
	echo.Use(middleware.Recover())

	// config handlers
	for _, h := range s.handlers {
		h.SetRoutes(echo)
	}

	// Start server
	address := fmt.Sprintf(":%d", s.config.HTTPServer.Port)

	fmt.Printf("start echo server on %s\n", address)

	if err := echo.Start(address); err != nil {
		fmt.Println("router start error", err)
	}
}

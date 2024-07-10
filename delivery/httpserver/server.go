package httpserver

import (
	"fmt"
	"log"

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
	Router   *echo.Echo
	handlers []Handler
}

func New(c Config, h []Handler) Server {
	return Server{
		Router:   echo.New(),
		handlers: h,
		config:   c,
	}
}

func (s Server) Serve() {
	s.Router.Use(middleware.Logger())
	s.Router.Use(middleware.Recover())

	// config handlers
	for _, h := range s.handlers {
		h.SetRoutes(s.Router)
	}

	// Start server
	address := fmt.Sprintf(":%d", s.config.HTTPServer.Port)

	fmt.Printf("start echo server on %s\n", address)

	err := s.Router.Start(address)
	if err != nil {
		log.Println(err)
	}
}

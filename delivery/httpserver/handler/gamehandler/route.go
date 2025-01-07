package gamehandler

import (
	"mdhesari/kian-quiz-golang-game/delivery/httpserver/middleware"

	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(r *echo.Echo) {
	g := r.Group("/games")

	g.GET("", h.GetAllGames, middleware.Auth(h.authSrv, h.authCfg), middleware.Presence(h.presenceSrv))
	g.GET("/:game_id", h.GetGame, middleware.Auth(h.authSrv, h.authCfg), middleware.Presence(h.presenceSrv))
	g.POST("/:game_id/answer-question", h.AnswerQuestion, middleware.Auth(h.authSrv, h.authCfg), middleware.Presence(h.presenceSrv))
	g.GET("/:game_id/next-question", h.GetNextQuestion, middleware.Auth(h.authSrv, h.authCfg), middleware.Presence(h.presenceSrv))
}

package websockethandler

import (
	"mdhesari/kian-quiz-golang-game/adapter/redisadapter"
	"mdhesari/kian-quiz-golang-game/service/authservice"
)

type Handler struct {
	redisAdap *redisadapter.Adapter
	authSrv   *authservice.Service
	authCfg   *authservice.Config
}

func New(redisAdap *redisadapter.Adapter, authSrv *authservice.Service, authCfg *authservice.Config) Handler {
	return Handler{
		redisAdap: redisAdap,
		authSrv:   authSrv,
		authCfg:   authCfg,
	}
}

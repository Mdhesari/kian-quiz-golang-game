package gamehandler

import (
	"mdhesari/kian-quiz-golang-game/service/authservice"
	"mdhesari/kian-quiz-golang-game/service/gameservice"
	"mdhesari/kian-quiz-golang-game/service/presenceservice"
)

type Handler struct {
	gameSrv     *gameservice.Service
	presenceSrv *presenceservice.Service
	authSrv     *authservice.Service
	authCfg     authservice.Config
}

func New(gameSrv *gameservice.Service, presenceSrv *presenceservice.Service, authSrv *authservice.Service, authCfg authservice.Config) Handler {
	return Handler{
		gameSrv:     gameSrv,
		authSrv:     authSrv,
		authCfg:     authCfg,
		presenceSrv: presenceSrv,
	}
}

package gamehandler

import (
	"mdhesari/kian-quiz-golang-game/delivery/validator/gamevalidator"
	"mdhesari/kian-quiz-golang-game/service/authservice"
	"mdhesari/kian-quiz-golang-game/service/gameservice"
	"mdhesari/kian-quiz-golang-game/service/presenceservice"
)

type Handler struct {
	gameValidator *gamevalidator.Validator
	gameSrv       *gameservice.Service
	presenceSrv   *presenceservice.Service
	authSrv       *authservice.Service
	authCfg       authservice.Config
}

func New(gameValidator *gamevalidator.Validator, gameSrv *gameservice.Service, presenceSrv *presenceservice.Service, authSrv *authservice.Service, authCfg authservice.Config) Handler {
	return Handler{
		gameValidator: gameValidator,
		gameSrv:       gameSrv,
		presenceSrv:   presenceSrv,
		authSrv:       authSrv,
		authCfg:       authCfg,
	}
}

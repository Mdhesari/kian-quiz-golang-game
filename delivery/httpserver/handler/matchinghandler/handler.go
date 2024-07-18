package matchinghandler

import (
	"mdhesari/kian-quiz-golang-game/delivery/validator/matchingvalidator"
	"mdhesari/kian-quiz-golang-game/service/authservice"
	"mdhesari/kian-quiz-golang-game/service/matchingservice"
	"mdhesari/kian-quiz-golang-game/service/presenceservice"
)

type Handler struct {
	authConfig        authservice.Config
	authSrv           *authservice.Service
	matchingSrv       matchingservice.Service
	matchingValidator matchingvalidator.Validator
	presenceSrv       presenceservice.Service
}

func New(authConfig authservice.Config, authSrv *authservice.Service, matchingSrv matchingservice.Service, matchingValidator matchingvalidator.Validator, presenceSrv presenceservice.Service) Handler {
	return Handler{
		authConfig:        authConfig,
		authSrv:           authSrv,
		matchingSrv:       matchingSrv,
		matchingValidator: matchingValidator,
		presenceSrv:       presenceSrv,
	}
}

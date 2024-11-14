package categoryhandler

import (
	"mdhesari/kian-quiz-golang-game/service/authservice"
	"mdhesari/kian-quiz-golang-game/service/categoryservice"
	"mdhesari/kian-quiz-golang-game/service/presenceservice"
)

type Handler struct {
	categorySrv *categoryservice.Service
	presenceSrv *presenceservice.Service
	authSrv     *authservice.Service
	authConfig  authservice.Config
}

func New(categorySrv *categoryservice.Service, presenceSrv *presenceservice.Service, authSrv *authservice.Service, authConfig authservice.Config) Handler {
	return Handler{
		categorySrv: categorySrv,
		presenceSrv: presenceSrv,
		authSrv:     authSrv,
		authConfig:  authConfig,
	}
}

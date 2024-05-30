package backpanelhandler

import (
	"mdhesari/kian-quiz-golang-game/service/authservice"
	"mdhesari/kian-quiz-golang-game/service/rbacservice"
	"mdhesari/kian-quiz-golang-game/service/userservice"
)

type Handler struct {
	userSrv *userservice.Service
	rbacSrv *rbacservice.Service
	authSrv *authservice.Service
	authConfig authservice.Config
}

func New(userSrv *userservice.Service, rbacSrv *rbacservice.Service, authSrv *authservice.Service, authConf authservice.Config) Handler {
	return Handler{
		userSrv: userSrv,
		rbacSrv: rbacSrv,
		authSrv: authSrv,
		authConfig: authConf,
	}
}

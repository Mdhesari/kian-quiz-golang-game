package userhandler

import (
	"mdhesari/kian-quiz-golang-game/service/authservice"
	"mdhesari/kian-quiz-golang-game/service/userservice"
)

type Handler struct {
	userSrv userservice.Service
	authSrv authservice.Service
}

func New(userSrv userservice.Service, authSrv authservice.Service) Handler {
	return Handler{
		userSrv: userSrv,
		authSrv: authSrv,
	}
}

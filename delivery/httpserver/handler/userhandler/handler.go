package userhandler

import "mdhesari/kian-quiz-golang-game/service/userservice"

type Handler struct{
	userSrv userservice.Service
}

func New(userSrv userservice.Service) Handler {
	return Handler{
		userSrv: userSrv,
	}
}

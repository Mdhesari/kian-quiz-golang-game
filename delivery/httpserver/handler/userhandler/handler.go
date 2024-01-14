package userhandler

import (
	"mdhesari/kian-quiz-golang-game/delivery/validator/uservalidator"
	"mdhesari/kian-quiz-golang-game/service/authservice"
	"mdhesari/kian-quiz-golang-game/service/userservice"
)

type Handler struct {
	userSrv       userservice.Service
	authSrv       authservice.Service
	userValidator uservalidator.Validator
}

func New(userSrv userservice.Service, authSrv authservice.Service, userValidator uservalidator.Validator) Handler {
	return Handler{
		userSrv:       userSrv,
		authSrv:       authSrv,
		userValidator: userValidator,
	}
}

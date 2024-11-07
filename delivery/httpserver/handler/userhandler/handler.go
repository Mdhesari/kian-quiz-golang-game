package userhandler

import (
	"mdhesari/kian-quiz-golang-game/delivery/validator/uservalidator"
	"mdhesari/kian-quiz-golang-game/service/authservice"
	"mdhesari/kian-quiz-golang-game/service/presenceservice"
	"mdhesari/kian-quiz-golang-game/service/rbacservice"
	"mdhesari/kian-quiz-golang-game/service/userservice"
)

type Handler struct {
	userSrv       *userservice.Service
	authSrv       *authservice.Service
	rbacSrv       *rbacservice.Service
	presenceSrv   *presenceservice.Service
	authConfig    authservice.Config
	userValidator uservalidator.Validator
}

func New(userSrv *userservice.Service, authSrv *authservice.Service, rbacSrv *rbacservice.Service, presenceSrv *presenceservice.Service, authConfig authservice.Config, userValidator uservalidator.Validator) Handler {
	return Handler{
		userSrv:       userSrv,
		authSrv:       authSrv,
		rbacSrv:       rbacSrv,
		presenceSrv:   presenceSrv,
		authConfig:    authConfig,
		userValidator: userValidator,
	}
}

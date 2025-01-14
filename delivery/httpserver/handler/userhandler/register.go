package userhandler

import (
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/constant"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
	"mdhesari/kian-quiz-golang-game/service/userservice"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func (h Handler) Register(c echo.Context) error {
	var req param.RegisterRequest

	if err := c.Bind(&req); err != nil {
		logger.L().Error("Register eror: ", zap.Error(err))

		return c.JSON(http.StatusBadRequest, param.RegisterResponse{})
	}

	if fields, err := h.userValidator.ValidateRegisterRequest(req); err != nil {
		msg, code := richerror.Error(err)

		return c.JSON(code, echo.Map{
			"message": msg,
			"errors":  fields,
		})
	}

	role, err := h.rbacSrv.GetRole(constant.RoleUser)
	if err != nil {
		msg, code := richerror.Error(err)

		logger.L().Error("Could not get role.", zap.Error(err))

		return echo.NewHTTPError(code, msg)
	}

	var roleID *primitive.ObjectID
	if role != nil {
		roleID = &role.ID
	}

	res, err := h.userSrv.Register(userservice.UserForm{
		Name:     req.Name,
		Email:    req.Email,
		Mobile:   req.Mobile,
		Password: req.Password,
		RoleID:   roleID,
	})
	if err != nil {
		logger.L().Error("User handler: could not register user.", zap.Error(err))
		msg, code := richerror.Error(err)

		return c.JSON(code, echo.Map{
			"message": msg,
		})
	}

	return c.JSON(http.StatusCreated, res)
}

package middleware

import (
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/constant"
	"mdhesari/kian-quiz-golang-game/service/authservice"
	"mdhesari/kian-quiz-golang-game/service/rbacservice"
	"mdhesari/kian-quiz-golang-game/service/userservice"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RBAC(userSrv *userservice.Service, rbacSrv *rbacservice.Service) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims := authservice.GetClaims(c)
			res, err := userSrv.GetByID(claims.UserID)
			if err != nil {

				return err
			}

			if res.User.RoleID == nil {
				return c.JSON(http.StatusUnauthorized, param.GetDefaultUnAuthorizedResponse())
			}
			// TODO - res with rbace srv
			perms := []string{constant.PermissionListUsers}
			hasPerm, err := rbacSrv.HasPermissions(*res.User.RoleID, perms...)
			if err != nil {

				panic(err)
			}
			if !hasPerm {
				return c.JSON(http.StatusUnauthorized, param.GetDefaultUnAuthorizedResponse())
			}

			return next(c)
		}
	}
}

package middleware

import (
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

			// TODO - res with rbace srv
			hasPerm, err := rbacSrv.HasPermissions(res.User.RoleID, "list-users")
			if err != nil {

				panic(err)
			}
			if !hasPerm {

				return c.JSON(http.StatusUnauthorized, echo.Map{})
			}

			return next(c)
		}
	}
}

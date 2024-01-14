package richerror

import (
	"mdhesari/kian-quiz-golang-game/pkg/errmsg"
	"net/http"
)

func Error(err error) (string, int) {
	switch err.(type) {
	case RichError:
		re := err.(RichError)
		msg := re.Message()
		code := GetHttpCodeFromKind(re.Kind())

		if code >= 500 {
			msg = errmsg.ErrInternalServer
		}

		return msg, code
	default:
		return err.Error(), http.StatusBadRequest
	}
}

func GetHttpCodeFromKind(kind Kind) int {
	switch kind {
	case KindForbidden:
		return http.StatusForbidden
	case KindUnAthorized:
		return http.StatusUnauthorized
	case KindNotFound:
		return http.StatusNotFound
	case KindUnexpected:
		return http.StatusInternalServerError
	case KindInvalid:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusBadRequest
	}
}

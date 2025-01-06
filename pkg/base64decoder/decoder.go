package base64decoder

import (
	"encoding/base64"
	"mdhesari/kian-quiz-golang-game/logger"

	"go.uber.org/zap"
)

func Decode(s string) []byte {
	var res []byte = make([]byte, 0)

	res, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		logger.L().Error("Could not decode base64.", zap.Error(err), zap.String("str", s))

		return res
	}

	return res
}

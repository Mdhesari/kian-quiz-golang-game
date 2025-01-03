package jsonencoder

import (
	"encoding/json"
	"mdhesari/kian-quiz-golang-game/logger"

	"go.uber.org/zap"
)

func EncodeMessage(msg interface{}) string {
	res, err := json.Marshal(msg)
	if err != nil {
		logger.L().Error("Could not encode msg to json.", zap.Error(err), zap.Any("msg", msg))

		return ""
	}

	return string(res)
}

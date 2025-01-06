package mongoutils

import (
	"mdhesari/kian-quiz-golang-game/logger"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func HexToObjectID(h string) primitive.ObjectID {
	var id primitive.ObjectID

	id, err := primitive.ObjectIDFromHex(h)
	if err != nil {
		logger.L().Error("Could not convert hex to object id.", zap.Error(err))
	}

	return id
}

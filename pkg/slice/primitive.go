package slice

import (
	"mdhesari/kian-quiz-golang-game/logger"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

func MapFromPrimitiveObjectIDToHexString(ids []primitive.ObjectID) []string {
	var mappedIds []string

	for _, id := range ids {
		mappedIds = append(mappedIds, id.Hex())
	}

	return mappedIds
}

func MapFromHexIDStringToPrimitiveObject(ids []string) []primitive.ObjectID {
	var mappedIds []primitive.ObjectID

	for _, id := range ids {
		oId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			// TODO - update metrics
			logger.L().Error("Could convert hex to obj id {%v}\n", zap.Error(err))
		}

		mappedIds = append(mappedIds, oId)
	}

	return mappedIds
}

package slice

import (
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

			log.Fatalf("Could convert hex to obj id {%v}\n", err)
		}
		mappedIds = append(mappedIds, oId)
	}

	return mappedIds
}

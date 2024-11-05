package slice

import "go.mongodb.org/mongo-driver/bson/primitive"

func MapFromPrimitiveObjectIDToHexString(ids []primitive.ObjectID) []string {
	var mappedIds []string

	for _, id := range ids {
		mappedIds = append(mappedIds, id.Hex())
	}

	return mappedIds
}

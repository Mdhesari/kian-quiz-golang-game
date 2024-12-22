package protobufdecoder

import (
	"encoding/base64"
	"log"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/pkg/slice"
	"mdhesari/kian-quiz-golang-game/protobuf/golang/matching"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/encoding/protojson"
)

func DecodePlayersMatchedEvent(s string) entity.PlayersMatched {
	var pbE matching.PlayersMatched

	res, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		// TOOD - update metrics

		panic(err)
	}
	if err := protojson.Unmarshal(res, &pbE); err != nil {
		// TODO - update metrics

		panic(err)
	}

	categoryId, err := primitive.ObjectIDFromHex(pbE.Category.ID)
	if err != nil {
		// TODO - update metrics

		log.Fatalf("could not convert category hex to obj id {%v}.\n", err)
	}

	return entity.PlayersMatched{
		PlayerIDs: slice.MapFromHexIDStringToPrimitiveObject(pbE.PlayerIds),
		Category: entity.Category{
			ID:          categoryId,
			Title:       pbE.Category.Ttile,
			Description: pbE.Category.Description,
		},
	}
}

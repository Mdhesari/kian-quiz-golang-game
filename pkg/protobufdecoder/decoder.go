package protobufdecoder

import (
	"encoding/base64"
	"mdhesari/kian-quiz-golang-game/entity"
)

func DecodeUsersMatchedEvent(s string) entity.PlayersMatched {
	_, err := base64.StdEncoding.DecodeString(s)
	if err != nil {

		panic(err)
	}


	// TODO - complete this

	return entity.PlayersMatched{}
}
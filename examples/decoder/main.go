package main

import (
	"fmt"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/pkg/mongoutils"
	"mdhesari/kian-quiz-golang-game/pkg/protobufdecoder"
	"mdhesari/kian-quiz-golang-game/pkg/protobufencoder"
)

func main() {
	e := entity.PlayerAnswered{
		UserID:     mongoutils.HexToObjectID("6778de7c7dc0f926b37702ad"),
		GameID:     mongoutils.HexToObjectID("677cf8cbd393163176d9b29b"),
		QuestionID: mongoutils.HexToObjectID("6778de90f0766018d7eb7b94"),
		Answer: entity.Answer{
			Title: "Uno",
		},
	}
	s := protobufencoder.EncodePlayerAnswered(e)

	w := entity.WebsocketMsg{
		Type:    string(entity.GamePlayerAnsweredEvent),
		Payload: s,
	}

	fmt.Println(protobufdecoder.DecodeGamePlayerAnsweredEvent(s))

	fmt.Println(protobufencoder.EncodeWebSocketMsg(w))
}

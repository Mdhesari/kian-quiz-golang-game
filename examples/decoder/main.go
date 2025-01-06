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
		UserID:     mongoutils.HexToObjectID("6778de7c7dc0f926b37702ae"),
		GameID:     mongoutils.HexToObjectID("677bc3ff5a3adad24fcff543"),
		QuestionID: mongoutils.HexToObjectID("6778de90f0766018d7eb7bac"),
		Answer: entity.Answer{
			Title: "Cancer",
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

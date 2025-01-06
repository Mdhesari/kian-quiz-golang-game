package main

import (
	"fmt"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/pkg/mongoutils"
	"mdhesari/kian-quiz-golang-game/pkg/protobufencoder"
)

func main() {
	e := entity.PlayerAnswered{
		GameID:     mongoutils.HexToObjectID("677b8887bd248b137d75fea1"),
		QuestionID: mongoutils.HexToObjectID("677b8887bd248b137d75fea1"),
		Answer: entity.Answer{
			Title: "test answer",
		},
	}

	s := protobufencoder.EncodePlayerAnswered(e)

	w := entity.WebsocketMsg{
		Type:    string(entity.GamePlayerAnsweredEvent),
		Payload: s,
	}

	fmt.Println(protobufencoder.EncodeWebSocketMsg(w))
}

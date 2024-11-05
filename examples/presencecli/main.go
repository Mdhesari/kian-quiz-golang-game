package main

import (
	"context"
	"fmt"
	"log"
	"mdhesari/kian-quiz-golang-game/protobuf/golang/presence"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":8089", grpc.WithInsecure())
	if err != nil {

		log.Fatal("Could not dial grpc.")
	}
	defer conn.Close()

	cli := presence.NewPresenceServiceClient(conn)
	res, err := cli.GetPresence(context.Background(), &presence.GetPresenceRequest{
		UserId: []string{"000000000000000000000000"}	,
	})
	if err != nil {
		log.Fatalf("Clould not get presence: %v\n", err)
	}

	fmt.Println(res, "hi")
}

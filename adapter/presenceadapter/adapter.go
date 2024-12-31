package presenceadapter

import (
	"context"
	"fmt"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/protobufmapper"
	"mdhesari/kian-quiz-golang-game/pkg/slice"
	"mdhesari/kian-quiz-golang-game/protobuf/golang/presence"

	"google.golang.org/grpc"
)

type Client struct {
	presenceSrv presence.PresenceServiceClient
}

func New(addr string) Client {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {

		panic(fmt.Sprintf("Could not connect to grpc server: %v", err))
	}

	return Client{
		presenceSrv: presence.NewPresenceServiceClient(conn),
	}
}

func (c Client) GetPresence(ctx context.Context, req param.PresenceRequest) (param.PresenceResponse, error) {
	res, err := c.presenceSrv.GetPresence(ctx, &presence.GetPresenceRequest{
		UserId: slice.MapFromPrimitiveObjectIDToHexString(req.UserIds),
	})
	if err != nil {

		return param.PresenceResponse{}, err
	}

	return protobufmapper.MapFromProtobufPresenceResponseToParam(res), nil
}

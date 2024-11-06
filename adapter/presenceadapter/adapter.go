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

func New(conn *grpc.ClientConn) Client {
	return Client{
		presenceSrv: presence.NewPresenceServiceClient(conn),
	}
}

func (c Client) GetPresence(ctx context.Context, req param.PresenceRequest) (param.PresenceResponse, error) {
	fmt.Println("preeeeeeeeeeeee")
	res, err := c.presenceSrv.GetPresence(ctx, &presence.GetPresenceRequest{
		UserId: slice.MapFromPrimitiveObjectIDToHexString(req.UserIds),
	})
	if err != nil {

		return param.PresenceResponse{}, err
	}

	return protobufmapper.MapFromProtobufPresenceResponseToParam(res), nil
}

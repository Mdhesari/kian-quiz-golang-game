package redisadapter

import (
	"context"
	"fmt"
)

func (a Adapter) Publish(ctx context.Context, topic string, payload string) {
	res, err := a.cli.Publish(ctx, topic, payload).Result()
	if err != nil {

		panic(res)
	}

	fmt.Println("result", res)
}

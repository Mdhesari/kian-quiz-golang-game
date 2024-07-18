package redispresence

import (
	"context"
	"time"
)

func (db DB) Upsert(ctx context.Context, key string, timestamp int64, exp time.Duration) error {
	_, err := db.adapter.Cli().Set(ctx, key, timestamp, exp).Result()
	if err != nil {

		return err
	}

	return nil
}
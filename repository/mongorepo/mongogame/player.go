package mongogame

import (
	"context"
	"errors"
	"fmt"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/pkg/errmsg"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (d *DB) UpdatePlayerStatus(ctx context.Context, gameId, userId primitive.ObjectID, status entity.PlayerStatus) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			fmt.Sprintf("players.%s.status", userId.Hex()): status,
		},
	}
	res, err := d.collection.UpdateOne(
		ctx,
		bson.M{"_id": gameId},
		update,
	)
	if err != nil {

		return false, err
	}
	if res.MatchedCount == 0 {

		return false, errors.New(errmsg.ErrNotFound)
	}

	return res.ModifiedCount > 0, nil
}

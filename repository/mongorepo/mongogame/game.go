package mongogame

import (
	"context"
	"errors"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/pkg/errmsg"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (d *DB) Create(ctx context.Context, game entity.Game) (entity.Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.cli.QueryTimeout)
	defer cancel()
	result, err := d.cli.Conn().Collection("games").InsertOne(ctx, game)
	if err != nil {

		return game, err
	}
	if result.InsertedID == nil {

		return game, errors.New(errmsg.ErrGameNotCreated)
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {

		return game, errors.New(errmsg.ErrGameIDNotConverted)
	}
	game.ID = id

	return game, nil
}

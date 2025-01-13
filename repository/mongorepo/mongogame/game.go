package mongogame

import (
	"context"
	"errors"
	"fmt"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/pkg/errmsg"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *DB) Create(ctx context.Context, game entity.Game) (entity.Game, error) {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()
	result, err := d.collection.InsertOne(ctx, game)
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

func (d *DB) GetGameById(ctx context.Context, id primitive.ObjectID) (entity.Game, error) {
	var game entity.Game

	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	res := d.collection.FindOne(ctx, bson.M{
		"_id": id,
	})
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {

			return game, nil
		}

		return game, res.Err()
	}

	res.Decode(&game)

	return game, nil
}

func (d *DB) UpdateGameEndtime(ctx context.Context, gameId primitive.ObjectID, endTime time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"end_time": endTime,
		},
	}
	res, err := d.collection.UpdateOne(
		ctx,
		bson.M{"_id": gameId},
		update,
	)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New(errmsg.ErrGameNotFound)
	}
	if res.ModifiedCount == 0 {
		return errors.New(errmsg.ErrGameNotModified)
	}

	return nil
}

func (d *DB) GetAllGames(ctx context.Context, userID primitive.ObjectID) ([]entity.Game, error) {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	filter := bson.M{
		fmt.Sprintf("players.%s", userID.Hex()): bson.M{"$exists": true},
	}
	opts := options.Find().SetSort(bson.D{{Key: "start_time", Value: -1}})
	cursor, err := d.collection.Find(ctx, filter, opts)
	if err != nil {
		return []entity.Game{}, fmt.Errorf("failed to execute aggregation: %w", err)
	}
	defer cursor.Close(ctx)

	var games []entity.Game
	if err := cursor.All(ctx, &games); err != nil {
		return nil, fmt.Errorf("failed to decode games: %w", err)
	}

	return games, nil
}

func (d *DB) CreateQuestionAnswer(ctx context.Context, userId primitive.ObjectID, gameId primitive.ObjectID, playerAnswer entity.PlayerAnswer) (entity.PlayerAnswer, error) {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	res, err := d.collection.UpdateByID(ctx, gameId, bson.M{
		"$push": bson.M{
			fmt.Sprintf("players.%s.answers", userId.Hex()): playerAnswer,
		},
	})
	if err != nil {

		return playerAnswer, err
	}
	if res.MatchedCount == 0 {

		return playerAnswer, errors.New(errmsg.ErrGameNotFound)
	}

	return playerAnswer, nil
}

func (d *DB) UpdatePlayer(ctx context.Context, gameId primitive.ObjectID, userId primitive.ObjectID, player entity.Player) error {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			fmt.Sprintf("players.%s", userId.Hex()): player,
		},
	}
	result, err := d.collection.UpdateOne(
		ctx,
		bson.M{"_id": gameId},
		update,
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New(errmsg.ErrGameNotFound)
	}

	return nil
}

func (d *DB) UpdateGameStatus(ctx context.Context, gameId primitive.ObjectID, status entity.GameStatus) error {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}
	result, err := d.collection.UpdateOne(
		ctx,
		bson.M{"_id": gameId},
		update,
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New(errmsg.ErrGameNotFound)
	}

	return nil
}

func (d *DB) UpdateGameWinner(ctx context.Context, gameId primitive.ObjectID, winner entity.Player) error {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"winner_player": winner,
		},
	}
	result, err := d.collection.UpdateOne(
		ctx,
		bson.M{"_id": gameId},
		update,
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {

		return errors.New(errmsg.ErrGameNotFound)
	}

	return nil
}

func (d *DB) IncPlayerScore(ctx context.Context, gameId primitive.ObjectID, userId primitive.ObjectID, score entity.Score) error {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	update := bson.M{
		"$inc": bson.M{
			fmt.Sprintf("players.%s.score", userId.Hex()): score,
		},
	}
	result, err := d.collection.UpdateOne(
		ctx,
		bson.M{"_id": gameId},
		update,
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {

		return errors.New(errmsg.ErrGameNotFound)
	}

	return nil
}

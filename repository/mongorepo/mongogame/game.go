package mongogame

import (
	"context"
	"errors"
	"fmt"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/pkg/errmsg"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func (d *DB) Update(ctx context.Context, game entity.Game) error {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"category_id":  game.CategoryID,
			"question_ids": game.QuestionIDs,
			"player_ids":   game.PlayerIDs,
			"start_time":   game.StartTime,
			"updated_at":   game.UpdatedAt,
		},
	}

	result, err := d.collection.UpdateOne(
		ctx,
		bson.M{"_id": game.ID},
		update,
	)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New(errmsg.ErrGameNotFound)
	}

	if result.ModifiedCount == 0 {
		return errors.New(errmsg.ErrGameNotModified)
	}

	return nil
}

func (d *DB) GetAllGames(ctx context.Context, categoryID primitive.ObjectID, userID primitive.ObjectID) ([]entity.Game, error) {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	pipeline := mongo.Pipeline{}

	if !categoryID.IsZero() {
		pipeline = append(pipeline, bson.D{
			{"$match", bson.D{{"category_id", categoryID}}},
		})
	}

	pipeline = append(pipeline, bson.D{
		{"$lookup", bson.D{
			{"from", "players"},
			{"localField", "player_ids"},
			{"foreignField", "_id"},
			{"as", "players"},
		}},
	})

	pipeline = append(pipeline, bson.D{{"$unwind", "$players"}})

	if !userID.IsZero() {
		pipeline = append(pipeline, bson.D{
			{"$match", bson.D{{"players.user_id", userID}}},
		})
	}

	pipeline = append(pipeline, bson.D{
		{"$group", bson.D{
			{"_id", "$_id"},
			{"category_id", bson.D{{"$first", "$category_id"}}},
			{"question_ids", bson.D{{"$first", "$question_ids"}}},
			{"player_ids", bson.D{{"$first", "$player_ids"}}},
			{"start_time", bson.D{{"$first", "$start_time"}}},
			{"created_at", bson.D{{"$first", "$created_at"}}},
			{"updated_at", bson.D{{"$first", "$updated_at"}}},
		}},
	})

	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("failed to execute aggregation: %w", err)
	}
	defer cursor.Close(ctx)

	var games []entity.Game
	if err := cursor.All(ctx, &games); err != nil {
		return nil, fmt.Errorf("failed to decode games: %w", err)
	}

	return games, nil
}

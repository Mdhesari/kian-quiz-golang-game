package mongoplayer

import (
	"context"
	"errors"
	"fmt"
	"mdhesari/kian-quiz-golang-game/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func New(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		collection: db.Collection("players"),
	}

}

func (r *MongoRepository) Create(ctx context.Context, player entity.Player) (entity.Player, error) {
	result, err := r.collection.InsertOne(ctx, player)
	if err != nil {

		return entity.Player{}, fmt.Errorf("failed to create player: %w", err)
	}

	player.ID = result.InsertedID.(primitive.ObjectID)

	return player, nil
}

func (r *MongoRepository) GetByID(ctx context.Context, id primitive.ObjectID) (entity.Player, error) {
	var player entity.Player
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&player)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {

			return entity.Player{}, fmt.Errorf("player not found: %w", err)
		}

		return entity.Player{}, fmt.Errorf("failed to get player: %w", err)
	}

	return player, nil
}

func (r *MongoRepository) Update(ctx context.Context, player entity.Player) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": player.ID},
		bson.M{"$set": player},
	)
	if err != nil {

		return fmt.Errorf("failed to update player: %w", err)
	}

	return nil
}

func (r *MongoRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {

		return fmt.Errorf("failed to delete player: %w", err)
	}

	return nil
}

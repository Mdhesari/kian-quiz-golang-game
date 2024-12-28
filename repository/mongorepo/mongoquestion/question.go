package mongoquestion

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d *DB) GetRandomByCategory(ctx context.Context, categoryId primitive.ObjectID, count int) ([]entity.Question, error) {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	var questions []entity.Question = make([]entity.Question, count)

	// TODO - The $sample stage is efficient but may not scale well for large datasets, as MongoDB loads documents into memory to perform the sampling.
	pipeline := mongo.Pipeline{
		{{Key: "$sample", Value: bson.D{{Key: "size", Value: count}}}},
	}
	cursor, err := d.collection.Aggregate(ctx, pipeline)
	if err != nil {

		return questions, err
	}

	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &questions); err != nil {

		return questions, err
	}

	return questions, nil
}

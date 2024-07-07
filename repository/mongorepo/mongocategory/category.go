package mongocategory

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d *DB) FindById(ctx context.Context, id primitive.ObjectID) (*entity.Category, error) {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	res := d.cli.Conn().Collection("categories").FindOne(ctx, bson.M{
		"_id": id,
	})
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {

			return nil, nil
		}

		return nil, res.Err()
	}

	var category entity.Category
	res.Decode(&category)

	return &category, nil
}

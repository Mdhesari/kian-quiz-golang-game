package mongocategory

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *DB) FindById(ctx context.Context, id primitive.ObjectID) (entity.Category, error) {
	var category entity.Category

	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	res := d.cli.Conn().Collection("categories").FindOne(ctx, bson.M{
		"_id": id,
	})
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {

			return category, nil
		}

		return category, res.Err()
	}

	res.Decode(&category)

	return category, nil
}

func (d *DB) GetAll(ctx context.Context) ([]entity.Category, error) {
	var categories []entity.Category

	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	cur, err := d.cli.Conn().Collection("categories").Find(ctx, bson.D{}, options.Find())
	if err != nil {

		return categories, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var category entity.Category
		if err := cur.Decode(&category); err != nil {

			return categories, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}

func (d *DB) Exists(ctx context.Context, id primitive.ObjectID) (bool, error) {
	cat, err := d.FindById(ctx, id)
	if err != nil {

		return false, err
	}

	return !cat.ID.IsZero(), nil
}

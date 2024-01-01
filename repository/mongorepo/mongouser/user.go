package mongouser

import (
	"context"
	"errors"
	"mdhesari/kian-quiz-golang-game/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *DB) GetAll(ctx context.Context) ([]entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	var users []entity.User
	cur, err := d.cli.Conn().Collection("users").Find(ctx, bson.D{}, options.Find())
	if err != nil {
		return nil, err
	}

	for cur.Next(context.Background()) {
		var u entity.User

		if err := cur.Decode(&u); err != nil {
			panic(err)
		}

		users = append(users, u)
	}

	return users, nil
}

func (d DB) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	var user entity.User
	filter := bson.D{{"email", email}}
	res := d.cli.Conn().Collection("users").FindOne(ctx, filter)
	if res.Err() != nil {
		return nil, res.Err()
	}

	res.Decode(&user)

	return &user, nil
}

func (d DB) Register(ctx context.Context, u *entity.User) error {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	hash, err := d.cli.Hash.Generate(u.Password)
	if err != nil {
		return err
	}
	u.Password = hash

	result, err := d.cli.Conn().Collection("users").InsertOne(ctx, u)
	if err != nil {
		return err
	}

	if result.InsertedID == nil {
		return errors.New("Could not create a new user")
	}

	return nil
}
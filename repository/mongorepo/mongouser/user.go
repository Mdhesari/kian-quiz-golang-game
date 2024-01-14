package mongouser

import (
	"context"
	"errors"
	"mdhesari/kian-quiz-golang-game/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	filter := bson.M{"email": email}
	res := d.cli.Conn().Collection("users").FindOne(ctx, filter)
	if res.Err() != nil {
		return nil, res.Err()
	}

	res.Decode(&user)

	return &user, nil
}

func (d DB) IsMobileUnique(mobile string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.cli.QueryTimeout)
	defer cancel()

	filter := bson.M{"mobile": mobile}
	res := d.cli.Conn().Collection("users").FindOne(ctx, filter)
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {

			return true, nil
		}

		return false, res.Err()
	}

	return false, nil
}

func (d DB) IsEmailUnique(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.cli.QueryTimeout)
	defer cancel()

	filter := bson.M{"email": email}
	res := d.cli.Conn().Collection("users").FindOne(ctx, filter)
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {

			return true, nil
		}

		return false, res.Err()
	}

	return false, nil
}

func (d DB) FindByID(id primitive.ObjectID) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.cli.QueryTimeout)
	defer cancel()

	var user entity.User
	filter := bson.M{"_id": id}
	res := d.cli.Conn().Collection("users").FindOne(ctx, filter)
	if res.Err() != nil {
		return nil, res.Err()
	}

	res.Decode(&user)

	return &user, nil
}

func (d DB) Register(ctx context.Context, u entity.User) (entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	hash, err := d.cli.Hash.Generate(u.Password)
	if err != nil {
		return entity.User{}, err
	}
	u.Password = hash

	result, err := d.cli.Conn().Collection("users").InsertOne(ctx, u)
	if err != nil {
		return entity.User{}, err
	}

	if result.InsertedID == nil {
		return entity.User{}, errors.New("Could not create a new user")
	}

	uid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		u.ID = uid
	} else {
		return entity.User{}, errors.New("Could not assert object id.")
	}

	return u, nil
}

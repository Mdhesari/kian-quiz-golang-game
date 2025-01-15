package mongouser

import (
	"context"
	"errors"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/pkg/errmsg"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *DB) GetAll(ctx context.Context) ([]entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	var users []entity.User
	cur, err := d.collection.Find(ctx, bson.D{}, options.Find())
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

func (d *DB) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	var user entity.User
	filter := bson.M{"email": email}
	res := d.collection.FindOne(ctx, filter)
	if res.Err() != nil {
		return nil, res.Err()
	}

	res.Decode(&user)

	return &user, nil
}

func (d *DB) IsMobileUnique(mobile string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.cli.QueryTimeout)
	defer cancel()

	filter := bson.M{"mobile": mobile}
	res := d.collection.FindOne(ctx, filter)
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {

			return true, nil
		}

		return false, res.Err()
	}

	return false, nil
}

func (d *DB) IsEmailUnique(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.cli.QueryTimeout)
	defer cancel()

	filter := bson.M{"email": email}
	res := d.collection.FindOne(ctx, filter)
	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {

			return true, nil
		}

		return false, res.Err()
	}

	return false, nil
}

func (d *DB) FindByID(id primitive.ObjectID) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), d.cli.QueryTimeout)
	defer cancel()

	var user entity.User
	filter := bson.M{"_id": id}
	res := d.collection.FindOne(ctx, filter)
	if res.Err() != nil {

		return nil, res.Err()
	}

	res.Decode(&user)

	return &user, nil
}

func (d *DB) Register(ctx context.Context, u entity.User) (entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	result, err := d.collection.InsertOne(ctx, u)
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

func (d *DB) IncrementScore(ctx context.Context, id primitive.ObjectID, score entity.Score) error {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	filter := bson.M{"_id": id}
	update := bson.M{"$inc": bson.M{"score": score}}
	res, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {

		return err
	}

	if res.MatchedCount == 0 {

		return errors.New(errmsg.ErrUserNotFound)
	}

	return nil
}

func (d *DB) FindManyById(ctx context.Context, ids []primitive.ObjectID) ([]entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, d.cli.QueryTimeout)
	defer cancel()

	var users []entity.User
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cur, err := d.collection.Find(ctx, filter)
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

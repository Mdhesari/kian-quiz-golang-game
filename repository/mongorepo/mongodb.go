package mongorepo

import (
	"context"
	"errors"
	"mdhesari/kian-quiz-golang-game/entity"
	"time"

	"github.com/hellofresh/janus/pkg/plugin/basic/encrypt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoQueryTimeout = 10 * time.Second
)

type Repository struct {
	collection *mongo.Collection
	hash       encrypt.Hash
}

func New(c *mongo.Collection) (Repository, error) {
	return Repository{
		collection: c,
		hash:       encrypt.Hash{},
	}, nil
}

func (r Repository) GetAll() ([]entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()

	var users []entity.User
	cur, err := r.collection.Find(ctx, bson.D{}, options.Find())
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

func (r Repository) FindByEmail(email string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()

	var user entity.User
	filter := bson.D{{"email", email}}
	res := r.collection.FindOne(ctx, filter)
	if res.Err() != nil {
		return nil, res.Err()
	}

	res.Decode(&user)

	return &user, nil
}

func (r Repository) Register(u *entity.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()

	hash, err := r.hash.Generate(u.Password)
	if err != nil {
		return err
	}
	u.Password = hash

	result, err := r.collection.InsertOne(ctx, u)
	if err != nil {
		return err
	}

	if result.InsertedID == nil {
		return errors.New("Could not create a new user")
	}

	return nil
}

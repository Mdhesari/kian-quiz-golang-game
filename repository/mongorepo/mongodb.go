package mongorepo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Username        string `koanf:"username"`
	Password        string `koanf:"password"`
	Port            int    `koanf:"port"`
	Host            string `koanf:"host"`
	DBName          string `koanf:"db_name"`
	Migrations      string `koanf:"migrations"`
	DurationSeconds int    `koanf:"duration_seconds"`
}

type MongoDB struct {
	config       Config
	conn         *mongo.Database
	QueryTimeout time.Duration
}

func New(c Config) (*MongoDB, error) {
	url := fmt.Sprintf("mongodb://%s:%s@%s:%d/", c.Username, c.Password, c.Host, c.Port)
	clientOptions := options.Client().ApplyURI(url)

	cli, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {

		return nil, err
	}

	err = cli.Ping(context.Background(), nil)
	if err != nil {

		return nil, err
	}

	return &MongoDB{
		conn:         cli.Database(c.DBName),
		QueryTimeout: time.Duration(c.DurationSeconds * int(time.Second)),
		config:       c,
	}, nil
}

func (m *MongoDB) Conn() *mongo.Database {
	return m.conn
}

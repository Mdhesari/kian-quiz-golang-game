package redisadapter

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string `koanf:"host"`
	Username string `koanf:"username"`
	Port     int    `koanf:"port"`
	Password string `koanf:"password"`
	DB       int    `koanf:"db"`
}

type Adapter struct {
	cli *redis.Client
}

func New(c Config) Adapter {
	cli := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Username: c.Username,
		Password: c.Password,
		DB:       c.DB,
	})

	return Adapter{
		cli: cli,
	}
}

func (a Adapter) Cli() *redis.Client {
	return a.cli
}

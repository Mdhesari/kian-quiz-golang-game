package redispresence

import "mdhesari/kian-quiz-golang-game/adapter/redisadapter"

type DB struct {
	adapter *redisadapter.Adapter
}

func New(adapter *redisadapter.Adapter) DB {
	return DB{
		adapter: adapter,
	}
}

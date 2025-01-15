package redismatching

import "mdhesari/kian-quiz-golang-game/adapter/redisadapter"

const (
	defaultWaitingListPrefix string = "waitinglist"
)

type DB struct {
	adapter           *redisadapter.Adapter
	waitingListPrefix string
}

type option func(*DB)

func WithCustomPrefix(prefix string) option {
	return func(db *DB) {
		db.waitingListPrefix = prefix
	}
}

func New(adapter *redisadapter.Adapter, opts ...option) DB {
	db := DB{
		adapter:           adapter,
		waitingListPrefix: defaultWaitingListPrefix,
	}

	for _, opt := range opts {
		opt(&db)
	}

	return db
}

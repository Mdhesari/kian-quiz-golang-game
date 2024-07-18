package presenceservice

import (
	"context"
	"fmt"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
	"time"
)

type Service struct {
	repo   Repo
	config Config
}

type Config struct {
	Prefix     string        `koanf:"prefix"`
	Expiration time.Duration `kaonf:"expiration"`
}

type Repo interface {
	Upsert(ctx context.Context, key string, timestamp int64, exp time.Duration) error
}

func New(cfg Config, repo Repo) Service {
	return Service{
		repo:   repo,
		config: cfg,
	}
}

func (s Service) Upsert(ctx context.Context, req param.PresenceUpsertRequest) (*param.PrescenceUpsertResponse, error) {
	op := "Prescence service:upsert"

	key := fmt.Sprintf("%s:%s", s.config.Prefix, req.UserId.Hex())
	fmt.Println(key)
	err := s.repo.Upsert(ctx, key, req.Timestamp, s.config.Expiration)
	if err != nil {

		return nil, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return &param.PrescenceUpsertResponse{}, nil
}

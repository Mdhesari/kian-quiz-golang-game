package presenceservice

import (
	"context"
	"fmt"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	GetPresence(ctx context.Context, prefix string, userIds []primitive.ObjectID) (map[primitive.ObjectID]int64, error)
}

func New(cfg Config, repo Repo) Service {
	return Service{
		repo:   repo,
		config: cfg,
	}
}

func (s *Service) GetPresence(ctx context.Context, req param.PresenceRequest) (*param.PresenceResponse, error) {
	res := param.PresenceResponse{}
	op := "Presence service:get"

	presenceList, err := s.repo.GetPresence(ctx, s.config.Prefix, req.UserIds)
	if err != nil {

		return nil, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	for userId, timestamp := range presenceList {
		res.Items = append(res.Items, param.PresenceItem{
			UserId:    userId,
			Timestamp: timestamp,
		})
	}

	return &res, nil
}

func (s *Service) Upsert(ctx context.Context, req param.PresenceUpsertRequest) (*param.PrescenceUpsertResponse, error) {
	op := "Presence service:upsert"

	key := fmt.Sprintf("%s:%s", s.config.Prefix, req.UserId.Hex())
	err := s.repo.Upsert(ctx, key, req.Timestamp, s.config.Expiration)
	if err != nil {

		return nil, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return &param.PrescenceUpsertResponse{}, nil
}

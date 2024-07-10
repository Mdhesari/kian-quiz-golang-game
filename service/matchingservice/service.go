package matchingservice

import (
	"context"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo Repo
}

type Repo interface {
	AddToWaitingList(ctx context.Context, userId primitive.ObjectID, categoryId primitive.ObjectID) error
}

func New(repo Repo) Service {
	return Service{
		repo: repo,
	}
}

func (s Service) AddToWaitingList(req param.MatchingAddToWaitingListRequest) (*param.MatchingAddToWaitingListResponse, error) {
	op := "Matching Service: Add to waiting list."

	err := s.repo.AddToWaitingList(context.Background(), req.UserID, req.CategoryID)
	if err != nil {

		return nil, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return &param.MatchingAddToWaitingListResponse{
		Timeout: 1000 * time.Nanosecond,
	}, nil
}

func (s Service) MatchWaitedUsers(req param.MatchingWaitedUsersRequest) (*param.MatchingWaitedUsersResponse, error) {
	return nil, nil
}

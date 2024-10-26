package matchingservice

import (
	"context"
	"log"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
	"mdhesari/kian-quiz-golang-game/pkg/timestamp"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	repo           Repo
	categoryRepo   CategoryRepo
	presenceClient PresenceClient
}

type PresenceClient interface {
	GetByUserIds(ctx context.Context, req param.PresenceRequest) (param.PresenceResponse, error)
}

type Repo interface {
	AddToWaitingList(ctx context.Context, userId primitive.ObjectID, categoryId primitive.ObjectID) error
	GetWaitingListByCategory(ctx context.Context, category entity.Category) ([]entity.WaitingMember, error)
}

type CategoryRepo interface {
	GetAll(ctx context.Context) ([]entity.Category, error)
}

func New(repo Repo, categoryRepo CategoryRepo) Service {
	return Service{
		repo:         repo,
		categoryRepo: categoryRepo,
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

func (s Service) MatchWaitedUsers(ctx context.Context, req param.MatchingWaitedUsersRequest) (*param.MatchingWaitedUsersResponse, error) {
	op := "Match waited users."

	categories, err := s.categoryRepo.GetAll(ctx)
	if err != nil {

		return nil, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	var wg sync.WaitGroup
	for _, category := range categories {
		wg.Add(1)
		go s.Match(ctx, category, &wg)
	}
	wg.Wait()

	return &param.MatchingWaitedUsersResponse{}, nil
}

func (s Service) Match(ctx context.Context, category entity.Category, wg *sync.WaitGroup) {
	defer wg.Done()

	// get list of scores:category
	waitingList, err := s.repo.GetWaitingListByCategory(ctx, category)
	if err != nil {
		log.Printf("Get waiting list by category error: %v\n", err)

		return
	}

	var userIds []primitive.ObjectID
	for _, m := range waitingList {
		userIds = append(userIds, m.UserId)
	}

	// exclude users that have been offline for a long period of time
	presenceReq := param.PresenceRequest{
		UserIds: userIds,
	}
	presenceList, err := s.presenceClient.GetByUserIds(ctx, presenceReq)

	var finalList []entity.WaitingMember
	for _, m := range waitingList {
		lastOnlineTimestamp, ok := getPresenceItem(presenceList, m.UserId)
		if ok && lastOnlineTimestamp > timestamp.Add(-60*time.Second) && m.Timestamp > timestamp.Add(-300*time.Second) {
			finalList = append(finalList, m)
		}
	}

	// match the list by oldest request and publish matched message to the broker
	// for _, item := range finalList {
		
	// }

	// remove the users from waiting list
}

func getPresenceItem(presenceList param.PresenceResponse, userId primitive.ObjectID) (int64, bool) {
	for _, p := range presenceList.Items {
		if userId == p.UserId {

			return p.Timestamp, true
		}
	}

	return 0, false
}

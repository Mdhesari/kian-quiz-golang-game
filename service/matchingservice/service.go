package matchingservice

import (
	"context"
	"log"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/protobufencoder"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
	"mdhesari/kian-quiz-golang-game/pkg/slice"
	"mdhesari/kian-quiz-golang-game/pkg/timestamp"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Config struct {
	MatchingTimeoutSeconds uint `koanf:"matching_timeout_seconds"`
}

type Service struct {
	cfg            Config
	repo           Repo
	categoryRepo   CategoryRepo
	presenceClient PresenceClient
	pub            Publisher
}

type Publisher interface {
	Publish(ctx context.Context, topic string, payload string)
}

type PresenceClient interface {
	GetPresence(ctx context.Context, req param.PresenceRequest) (param.PresenceResponse, error)
}

type Repo interface {
	AddToWaitingList(ctx context.Context, userId primitive.ObjectID, category entity.Category) error
	GetWaitingListByCategory(ctx context.Context, category entity.Category) ([]entity.WaitingMember, error)
	RemoveUsersFromWaitingList(ctx context.Context, categroy entity.Category, userIds []string) error
}

type CategoryRepo interface {
	FindById(ctx context.Context, categoryId primitive.ObjectID) (entity.Category, error)
	GetAll(ctx context.Context) ([]entity.Category, error)
}

func New(cfg Config, repo Repo, categoryRepo CategoryRepo, presenceCli PresenceClient, pub Publisher) Service {
	return Service{
		cfg:            cfg,
		repo:           repo,
		categoryRepo:   categoryRepo,
		presenceClient: presenceCli,
		pub:            pub,
	}
}

func (s Service) AddToWaitingList(req param.MatchingAddToWaitingListRequest) (*param.MatchingAddToWaitingListResponse, error) {
	op := "Matching Service: Add to waiting list."

	ctx := context.Background()
	category, err := s.categoryRepo.FindById(ctx, req.CategoryID)
	if err != nil {

		return nil, err
	}

	err = s.repo.AddToWaitingList(ctx, req.UserID, category)
	if err != nil {

		return nil, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return &param.MatchingAddToWaitingListResponse{
		Timeout: s.cfg.MatchingTimeoutSeconds,
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

	waitingList, err := s.repo.GetWaitingListByCategory(ctx, category)
	if err != nil {
		log.Printf("Get waiting list by category error: %v\n", err)

		return
	}
	if len(waitingList) < 1 {

		return
	}

	var userIds []primitive.ObjectID
	for _, m := range waitingList {
		userIds = append(userIds, m.UserId)
	}

	presenceReq := param.PresenceRequest{
		UserIds: userIds,
	}

	presenceList, err := s.presenceClient.GetPresence(ctx, presenceReq)
	if err != nil {
		log.Fatalf("Get presence failed: %v\n", err)

		return
	}

	var usersToBeRemoved []string = make([]string, 0)
	// exclude users that have been offline for a long period of time
	var finalList []entity.WaitingMember
	for _, m := range waitingList {
		lastOnlineTimestamp, ok := getPresenceItem(presenceList, m.UserId)
		// TODO - add time to config
		if ok && lastOnlineTimestamp > timestamp.Add(-1000*time.Hour) && m.Timestamp > timestamp.Add(-3000*time.Hour) {
			finalList = append(finalList, m)
		} else {
			usersToBeRemoved = append(usersToBeRemoved, m.UserId.Hex())
		}
	}

	// match the list by oldest request and publish matched message to the broker
	for i := 1; i < len(finalList); i += 2 {
		e := entity.PlayersMatched{
			Players:  []primitive.ObjectID{finalList[i].UserId, finalList[i-1].UserId},
			Category: category,
		}

		s.pub.Publish(ctx, string(entity.UsersMatchedEvent), protobufencoder.EncodeUsersMatchedEvent(e))

		usersToBeRemoved = append(usersToBeRemoved, slice.MapFromPrimitiveObjectIDToHexString(e.Players)...)
	}

	// remove the users from waiting list
	go s.removeUsersFromWaitingList(category, usersToBeRemoved)
}

func getPresenceItem(presenceList param.PresenceResponse, userId primitive.ObjectID) (int64, bool) {
	for _, p := range presenceList.Items {
		if userId == p.UserId {

			return p.Timestamp, true
		}
	}

	return 0, false
}

func (s Service) removeUsersFromWaitingList(category entity.Category, userIds []string) {
	removeUsersCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.repo.RemoveUsersFromWaitingList(removeUsersCtx, category, userIds); err != nil {
		// TODO - update metrics
	}
}

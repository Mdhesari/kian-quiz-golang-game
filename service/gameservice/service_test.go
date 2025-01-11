package gameservice

import (
	"context"
	"errors"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/param"
	"mdhesari/kian-quiz-golang-game/pkg/errmsg"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

type MockPub struct {
	mock.Mock
}

func (mpub *MockPub) Publish(ctx context.Context, topic string, payload string) {
	mpub.Called(ctx, topic, payload)
}

func (m *MockRepository) Create(ctx context.Context, game entity.Game) (entity.Game, error) {
	args := m.Called(ctx, game)
	return args.Get(0).(entity.Game), args.Error(1)
}

func (m *MockRepository) UpdatePlayer(ctx context.Context, gameId primitive.ObjectID, userId primitive.ObjectID, player entity.Player) error {
	args := m.Called(ctx, ctx, gameId, userId, player)
	return args.Error(0)
}

func (m *MockRepository) GetAllGames(ctx context.Context, userID primitive.ObjectID) ([]entity.Game, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]entity.Game), args.Error(1)
}

func (m *MockRepository) GetGameById(ctx context.Context, id primitive.ObjectID) (entity.Game, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(entity.Game), args.Error(1)
}

func (m *MockRepository) CreateQuestionAnswer(ctx context.Context, userId primitive.ObjectID, gameId primitive.ObjectID, playerAnswer entity.PlayerAnswer) (entity.PlayerAnswer, error) {
	args := m.Called(ctx, userId, gameId, playerAnswer)
	return args.Get(0).(entity.PlayerAnswer), args.Error(1)
}

func (m *MockRepository) UpdateGameStatus(ctx context.Context, gameId primitive.ObjectID, status entity.GameStatus) error {
	args := m.Called(ctx, gameId, status)
	return args.Error(0)
}

func (m *MockRepository) UpdateGameEndtime(ctx context.Context, gameId primitive.ObjectID, endTime time.Time) error {
	args := m.Called(ctx, gameId, endTime)
	return args.Error(0)
}

func (m *MockRepository) IncPlayerScore(ctx context.Context, gameId primitive.ObjectID, userId primitive.ObjectID, score entity.Score) error {
	args := m.Called(ctx, gameId, userId, score)
	return args.Error(0)
}

func (m *MockRepository) UpdateGameWinner(ctx context.Context, gameId primitive.ObjectID, player entity.Player) error {
	args := m.Called(ctx, gameId, player)
	return args.Error(0)
}

// Implement other methods of the Repository interface as needed

func TestService_AnswerQuestion_Successful(t *testing.T) {
	mockRepo := new(MockRepository)
	mockPub := new(MockPub)
	service := New(mockRepo, mockPub)

	ctx := context.Background()
	gameID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	questionID := primitive.NewObjectID()

	game := entity.Game{
		ID:     gameID,
		Status: entity.GameStatusInProgress,
		Players: map[string]entity.Player{
			userID.Hex(): {
				LastQuestionStartTime: time.Now().Add(-time.Minute),
			},
		},
		Questions: []entity.Question{
			{
				ID: questionID,
				Answers: []entity.Answer{
					{Title: "Correct", IsCorrect: true},
				},
			},
		},
	}

	mockRepo.On("GetGameById", ctx, gameID).Return(game, nil)
	mockRepo.On("CreateQuestionAnswer", ctx, userID, gameID, mock.AnythingOfType("entity.PlayerAnswer")).Return(entity.PlayerAnswer{}, nil)

	req := param.GameAnswerQuestionRequest{
		UserId:     userID,
		GameId:     gameID,
		QuestionId: questionID,
		Answer:     "Correct",
	}

	_, err := service.AnswerQuestion(ctx, req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestService_AnswerQuestion_GameNotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	mockPub := new(MockPub)
	service := New(mockRepo, mockPub)

	ctx := context.Background()
	gameID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	questionID := primitive.NewObjectID()

	expectedErr := errors.New(errmsg.ErrNotFound)

	mockRepo.On("GetGameById", ctx, gameID).Return(entity.Game{}, expectedErr)

	req := param.GameAnswerQuestionRequest{
		UserId:     userID,
		GameId:     gameID,
		QuestionId: questionID,
		Answer:     "Answer",
	}

	_, err := service.AnswerQuestion(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), errmsg.ErrNotFound)
	mockRepo.AssertExpectations(t)
}

func TestService_AnswerQuestion_GameNotInProgress(t *testing.T) {
	mockRepo := new(MockRepository)
	mockPub := new(MockPub)
	service := New(mockRepo, mockPub)

	ctx := context.Background()
	gameID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	questionID := primitive.NewObjectID()

	game := entity.Game{
		ID:     gameID,
		Status: entity.GameStatusFinished,
	}

	mockRepo.On("GetGameById", ctx, gameID).Return(game, nil)

	req := param.GameAnswerQuestionRequest{
		UserId:     userID,
		GameId:     gameID,
		QuestionId: questionID,
		Answer:     "Answer",
	}

	_, err := service.AnswerQuestion(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), errmsg.ErrGameNotInProgress)
	mockRepo.AssertExpectations(t)
}

func TestService_AnswerQuestion_PlayerNotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	mockPub := new(MockPub)
	service := New(mockRepo, mockPub)

	ctx := context.Background()
	gameID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	questionID := primitive.NewObjectID()

	game := entity.Game{
		ID:      gameID,
		Status:  entity.GameStatusInProgress,
		Players: map[string]entity.Player{},
	}

	mockRepo.On("GetGameById", ctx, gameID).Return(game, nil)

	req := param.GameAnswerQuestionRequest{
		UserId:     userID,
		GameId:     gameID,
		QuestionId: questionID,
		Answer:     "Answer",
	}

	_, err := service.AnswerQuestion(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), errmsg.ErrGamePlayerNotFound)
	mockRepo.AssertExpectations(t)
}

func TestService_AnswerQuestion_AlreadyAnswered(t *testing.T) {
	mockRepo := new(MockRepository)
	mockPub := new(MockPub)
	service := New(mockRepo, mockPub)

	ctx := context.Background()
	gameID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	questionID := primitive.NewObjectID()

	game := entity.Game{
		ID:     gameID,
		Status: entity.GameStatusInProgress,
		Players: map[string]entity.Player{
			userID.Hex(): {
				Answers: []entity.PlayerAnswer{
					{QuestionID: questionID},
				},
			},
		},
	}

	mockRepo.On("GetGameById", ctx, gameID).Return(game, nil)

	req := param.GameAnswerQuestionRequest{
		UserId:     userID,
		GameId:     gameID,
		QuestionId: questionID,
		Answer:     "Answer",
	}

	_, err := service.AnswerQuestion(ctx, req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), errmsg.ErrAlreadyAnswered)
	mockRepo.AssertExpectations(t)
}

func TestService_GetNextQuestion_Successful(t *testing.T) {
	mockRepo := new(MockRepository)
	mockPub := new(MockPub)
	service := New(mockRepo, mockPub)

	ctx := context.Background()
	gameID := primitive.NewObjectID()
	userID := primitive.NewObjectID()
	question1ID := primitive.NewObjectID()
	question2ID := primitive.NewObjectID()

	game := entity.Game{
		ID:     gameID,
		Status: entity.GameStatusInProgress,
		Players: map[string]entity.Player{
			userID.Hex(): {
				Answers: []entity.PlayerAnswer{
					{QuestionID: question1ID},
				},
			},
		},
		Questions: []entity.Question{
			{
				ID: question1ID,
			},
			{
				ID: question2ID,
			},
		},
	}

	mockRepo.On("GetGameById", ctx, gameID).Return(game, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("entity.Game")).Return(nil)

	req := param.GameGetNextQuestionRequest{
		GameId: gameID,
		UserId: userID,
	}

	response, err := service.GetNextQuestion(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, question2ID, response.Question.ID)

	// Verify that the player's information was updated
	mockRepo.AssertCalled(t, "Update", ctx, mock.MatchedBy(func(g entity.Game) bool {
		player := g.Players[userID.Hex()]
		return player.LastQuestionID == question2ID && !player.LastQuestionStartTime.IsZero()
	}))

	mockRepo.AssertExpectations(t)
}

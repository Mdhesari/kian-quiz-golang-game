package playerservice

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/param"
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

func (m *MockRepository) Create(ctx context.Context, player entity.Player) (entity.Player, error) {
	args := m.Called(ctx, player)
	return args.Get(0).(entity.Player), args.Error(1)
}

func (m *MockRepository) GetByID(ctx context.Context, id primitive.ObjectID) (entity.Player, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.Player), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, player entity.Player) error {
	args := m.Called(ctx, player)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreatePlayer(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(Config{}, mockRepo)

	ctx := context.Background()
	now := time.Now()
	req := param.PlayerCreateRequest{
		UserID:    primitive.NewObjectID(),
		GameID:    primitive.NewObjectID(),
		CreatedAt: now,
	}

	expectedPlayer := entity.Player{
		ID:        primitive.NewObjectID(),
		UserID:    req.UserID,
		GameID:    req.GameID,
		Answers:   []entity.PlayerAnswer{},
		Score:     0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("entity.Player")).Return(expectedPlayer, nil)

	resp, err := service.CreatePlayer(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, expectedPlayer, resp.Player)
	mockRepo.AssertExpectations(t)
}

func TestGetPlayerByID(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(Config{}, mockRepo)

	ctx := context.Background()
	playerID := primitive.NewObjectID()
	req := param.PlayerGetRequest{
		ID: playerID,
	}

	expectedPlayer := entity.Player{
		ID:        playerID,
		UserID:    primitive.NewObjectID(),
		GameID:    primitive.NewObjectID(),
		Answers:   []entity.PlayerAnswer{},
		Score:     0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("GetByID", ctx, playerID).Return(expectedPlayer, nil)

	resp, err := service.GetPlayerByID(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, expectedPlayer, resp.Player)
	mockRepo.AssertExpectations(t)
}

func TestUpdatePlayer(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(Config{}, mockRepo)

	ctx := context.Background()
	playerID := primitive.NewObjectID()
	now := time.Now()
	req := param.PlayerUpdateRequest{
		ID: playerID,
		Answers: []entity.PlayerAnswer{
			{
				QuestionID: primitive.NewObjectID(),
				Answer: entity.Answer{
					Title:     "Sample answer",
					IsCorrect: true,
				},
				Score:     5,
				StartTime: now.Add(-5 * time.Minute),
				EndTime:   now,
			},
		},
		Score:     5,
		UpdatedAt: now,
	}

	existingPlayer := entity.Player{
		ID:        playerID,
		UserID:    primitive.NewObjectID(),
		GameID:    primitive.NewObjectID(),
		Answers:   []entity.PlayerAnswer{},
		Score:     0,
		CreatedAt: now.Add(-1 * time.Hour),
		UpdatedAt: now.Add(-1 * time.Hour),
	}

	updatedPlayer := existingPlayer
	updatedPlayer.Answers = req.Answers
	updatedPlayer.Score = req.Score
	updatedPlayer.UpdatedAt = req.UpdatedAt

	mockRepo.On("GetByID", ctx, playerID).Return(existingPlayer, nil)
	mockRepo.On("Update", ctx, updatedPlayer).Return(nil)

	err := service.UpdatePlayer(ctx, req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeletePlayer(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(Config{}, mockRepo)

	ctx := context.Background()
	playerID := primitive.NewObjectID()
	req := param.PlayerDeleteRequest{
		ID: playerID,
	}

	mockRepo.On("Delete", ctx, playerID).Return(nil)

	err := service.DeletePlayer(ctx, req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

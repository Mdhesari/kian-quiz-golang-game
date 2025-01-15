package questionservice

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Config struct {
	QuestionsCount int `koanf:"questions_count"`
}

type Service struct {
	cfg  Config
	repo Repository
}

type Repository interface {
	GetRandomByCategory(ctx context.Context, categoryId primitive.ObjectID, count int) ([]entity.Question, error)
}

func New(cfg Config, repo Repository) Service {
	return Service{
		cfg:  cfg,
		repo: repo,
	}
}

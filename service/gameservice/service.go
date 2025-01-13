package gameservice

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Config struct {
	GameTimeout         time.Duration `koanf:"game_timeout"`
	MaxQuestionTimeout  time.Duration `koanf:"max_question_timeout"`
	MaxScorePerQuestion entity.Score  `koanf:"max_score_per_question"`
}

type Repository interface {
	Create(ctx context.Context, game entity.Game) (entity.Game, error)
	GetGameById(ctx context.Context, id primitive.ObjectID) (entity.Game, error)
	UpdatePlayer(ctx context.Context, gameId primitive.ObjectID, userId primitive.ObjectID, player entity.Player) error
	GetAllGames(ctx context.Context, userID primitive.ObjectID) ([]entity.Game, error)
	CreateQuestionAnswer(ctx context.Context, userId primitive.ObjectID, gameId primitive.ObjectID, playerAnswer entity.PlayerAnswer) (entity.PlayerAnswer, error)
	UpdateGameStatus(ctx context.Context, gameId primitive.ObjectID, status entity.GameStatus) error
	UpdateGameEndtime(ctx context.Context, gameId primitive.ObjectID, endTime time.Time) error
	UpdateGameWinner(ctx context.Context, gameId primitive.ObjectID, player entity.Player) error
	IncPlayerScore(ctx context.Context, gameId primitive.ObjectID, userId primitive.ObjectID, score entity.Score) error
	UpdatePlayerStatus(ctx context.Context, gameId, userId primitive.ObjectID, status entity.PlayerStatus) (bool, error)
}

type Publisher interface {
	Publish(ctx context.Context, topic string, payload string)
}

type Service struct {
	cfg  *Config
	repo Repository
	pub  Publisher
}

func New(cfg *Config, repo Repository, pub Publisher) Service {
	return Service{
		cfg:  cfg,
		repo: repo,
		pub:  pub,
	}
}

package authservice

import (
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/logger"
	"mdhesari/kian-quiz-golang-game/pkg/errmsg"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type Config struct {
	Secret         []byte        `koanf:"secret"`
	ExpireDuration time.Duration `koanf:"expire_duration"`
}

type Service struct {
	config Config
}

func New(c Config) Service {
	return Service{
		config: c,
	}
}

type Claims struct {
	UserID primitive.ObjectID `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *Service) GenerateToken(user *entity.User, iss string) (string, error) {
	op := "Generate token"

	mySigningKey := []byte(s.config.Secret)

	logger.L().Info("creating a new jwt token.", zap.Any("duration", s.config.ExpireDuration))

	expireDu := time.Now().Add(s.config.ExpireDuration)
	// Create the Claims
	claims := Claims{
		user.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireDu),
			Issuer:    iss,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	if err != nil {

		return "", richerror.New(op, errmsg.ErrSignKey).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return ss, nil
}

func (s *Service) VerifyToken(tokenString string) (*Claims, error) {
	op := "User verify token"

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {

		return s.config.Secret, nil
	})
	if err != nil {

		return nil, richerror.New(op, errmsg.ErrAuthorization).WithKind(richerror.KindUnAthorized)
	}

	if claims, ok := token.Claims.(*Claims); ok {

		return claims, nil
	}

	return nil, richerror.New(op, errmsg.ErrClaimAssertion).WithKind(richerror.KindUnexpected)
}

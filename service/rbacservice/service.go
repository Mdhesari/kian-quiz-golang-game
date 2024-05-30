package rbacservice

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccessRepo interface {
	GetRole(ctx context.Context, name string) (*entity.Role, error)
	HasPermissions(ctx context.Context, roleID primitive.ObjectID, permissionIDs ...primitive.ObjectID) (bool, error)
	GetPermissionIds(ctx context.Context, perms ...string) ([]primitive.ObjectID, error) 
}

type Service struct {
	accessRepo AccessRepo
}

func New(accessRepo AccessRepo) Service {
	return Service{
		accessRepo: accessRepo,
	}
}
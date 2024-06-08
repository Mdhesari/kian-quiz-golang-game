package rbacservice

import (
	"context"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s Service) GetPermissionIds(ctx context.Context, perms ...string) ([]primitive.ObjectID, error) {
	op := "RBAC Service: Get permissions."

	permission, err := s.accessRepo.GetPermissionIds(ctx, perms...)
	if err != nil {

		return nil, richerror.New(op, err.Error()).WithErr(err)
	}

	return permission, nil
}

func (s Service) HasPermissions(roleID primitive.ObjectID, permissions ...string) (bool, error) {
	op := "RBAC Service: Has permissoins."

	perms, err := s.accessRepo.GetPermissionIds(context.TODO(), permissions...)
	if err != nil {

		return false, err
	}
	res, err := s.accessRepo.HasPermissions(context.TODO(), roleID, perms...)
	if err != nil {

		return false, richerror.New(op, err.Error()).WithKind(richerror.KindUnexpected)
	}

	return res, nil
}

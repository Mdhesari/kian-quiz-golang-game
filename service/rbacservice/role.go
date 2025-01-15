package rbacservice

import (
	"context"
	"mdhesari/kian-quiz-golang-game/entity"
	"mdhesari/kian-quiz-golang-game/pkg/richerror"
)

func (s *Service) GetRole(name string) (*entity.Role, error) {
	op := "RBAC Service: Get Role"

	role, err := s.accessRepo.GetRole(context.TODO(), name)
	if err != nil {

		return nil, richerror.New(op, err.Error()).WithErr(err).WithKind(richerror.KindUnexpected)
	}

	return role, nil
}

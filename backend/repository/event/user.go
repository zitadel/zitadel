package event

import (
	"context"

	"github.com/zitadel/zitadel/backend/repository"
)

func (s *store) CreateUser(ctx context.Context, user *repository.User) (*repository.User, error) {
	err := s.es.Push(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

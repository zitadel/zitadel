package event

import (
	"context"
	"log"

	"github.com/zitadel/zitadel/backend/repository"
)

func (s *store) CreateUser(ctx context.Context, user *repository.User) (*repository.User, error) {
	log.Println("event.user.create")
	err := s.es.Push(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

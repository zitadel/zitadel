package repository

import (
	"context"
	"log"
)

func (s *eventStore) CreateUser(ctx context.Context, user *User) (*User, error) {
	log.Println("event.user.create")
	err := s.es.Push(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

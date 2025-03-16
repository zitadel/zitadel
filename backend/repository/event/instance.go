package event

import (
	"context"
	"log"

	"github.com/zitadel/zitadel/backend/repository"
)

func (s *store) CreateInstance(ctx context.Context, instance *repository.Instance) (*repository.Instance, error) {
	log.Println("event.instance.create")
	err := s.es.Push(ctx, instance)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

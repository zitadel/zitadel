package event

import (
	"context"

	"github.com/zitadel/zitadel/backend/repository"
)

func (s *store) CreateInstance(ctx context.Context, instance *repository.Instance) (*repository.Instance, error) {
	err := s.es.Push(ctx, instance)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

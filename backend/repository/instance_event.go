package repository

import (
	"context"
	"log"
)

func (s *eventStore) CreateInstance(ctx context.Context, instance *Instance) (*Instance, error) {
	log.Println("event.instance.create")
	err := s.es.Push(ctx, instance)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

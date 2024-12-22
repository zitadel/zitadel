package instance

import (
	"context"

	"github.com/zitadel/zitadel/internal/v3/storage"
)

type InstanceStorage interface {
	WriteInstanceAdded(ctx context.Context, tx storage.Transaction, instance *AddInstanceRequest) error
}

type BusinessLogic struct {
	client  storage.Client
	storage InstanceStorage
}

func NewBusinessLogic(client storage.Client, stores ...InstanceStorage) *BusinessLogic {
	return &BusinessLogic{
		client:  client,
		storage: chainedStorage[InstanceStorage](stores),
	}
}

type chainedStorage[S InstanceStorage] []S

func (cs chainedStorage[S]) WriteInstanceAdded(ctx context.Context, tx storage.Transaction, instance *AddInstanceRequest) error {
	for _, store := range cs {
		if err := store.WriteInstanceAdded(ctx, tx, instance); err != nil {
			return err
		}
	}
	return nil
}

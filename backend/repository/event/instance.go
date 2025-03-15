package event

import (
	"context"

	"github.com/zitadel/zitadel/backend/repository"
	"github.com/zitadel/zitadel/backend/repository/orchestrate/handler"
	"github.com/zitadel/zitadel/backend/storage/database"
	"github.com/zitadel/zitadel/backend/storage/eventstore"
)

func SetUpInstance(
	client database.Executor,
	next handler.Handle[*repository.Instance, *repository.Instance],
) handler.Handle[*repository.Instance, *repository.Instance] {
	es := eventstore.New(client)
	return func(ctx context.Context, instance *repository.Instance) (*repository.Instance, error) {
		instance, err := next(ctx, instance)
		if err != nil {
			return nil, err
		}

		err = es.Push(ctx, instance)
		if err != nil {
			return nil, err
		}
		return instance, nil
	}
}

func SetUpInstanceWithout(client database.Executor) handler.Handle[*repository.Instance, *repository.Instance] {
	es := eventstore.New(client)
	return func(ctx context.Context, instance *repository.Instance) (*repository.Instance, error) {
		err := es.Push(ctx, instance)
		if err != nil {
			return nil, err
		}
		return instance, nil
	}
}

func SetUpInstanceDecorated(
	client database.Executor,
	next handler.Handle[*repository.Instance, *repository.Instance],
	decorate handler.Decorate[*repository.Instance, *repository.Instance],
) handler.Handle[*repository.Instance, *repository.Instance] {
	es := eventstore.New(client)
	return func(ctx context.Context, instance *repository.Instance) (*repository.Instance, error) {
		instance, err := next(ctx, instance)
		if err != nil {
			return nil, err
		}

		return decorate(ctx, instance, func(ctx context.Context, instance *repository.Instance) (*repository.Instance, error) {
			err = es.Push(ctx, instance)
			if err != nil {
				return nil, err
			}
			return instance, nil
		})
	}
}

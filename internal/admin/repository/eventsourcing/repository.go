package eventsourcing

import (
	"context"

	"github.com/zitadel/zitadel/internal/admin/repository/eventsourcing/eventstore"
	admin_view "github.com/zitadel/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/database"
	eventstore2 "github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/static"
)

type Config struct {
	SearchLimit uint64
}

type EsRepository struct {
	eventstore.AdministratorRepo
}

func Start(ctx context.Context, conf Config, static static.Storage, dbClient *database.DB, esV2 *eventstore2.Eventstore) (*EsRepository, error) {
	view, err := admin_view.StartView(dbClient)
	if err != nil {
		return nil, err
	}

	return &EsRepository{
		AdministratorRepo: eventstore.AdministratorRepo{
			View: view,
		},
	}, nil
}

func (repo *EsRepository) Health(ctx context.Context) error {
	return nil
}

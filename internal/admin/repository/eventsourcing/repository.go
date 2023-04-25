package eventsourcing

import (
	"context"

	"github.com/zitadel/zitadel/internal/admin/repository/eventsourcing/eventstore"
	auth_handler "github.com/zitadel/zitadel/internal/admin/repository/eventsourcing/handler"
	admin_view "github.com/zitadel/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/static"
)

type Config struct {
	Spooler auth_handler.Config
}

type EsRepository struct {
	eventstore.AdministratorRepo
}

func Start(ctx context.Context, conf Config, static static.Storage, dbClient *database.DB) (*EsRepository, error) {
	view, err := admin_view.StartView(dbClient)
	if err != nil {
		return nil, err
	}

	auth_handler.Register(ctx, conf.Spooler, view, static)

	return &EsRepository{
		AdministratorRepo: eventstore.AdministratorRepo{
			View: view,
		},
	}, nil
}

func (repo *EsRepository) Health(ctx context.Context) error {
	return nil
}

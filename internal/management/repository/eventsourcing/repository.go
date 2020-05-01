package eventsourcing

import (
	"context"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	es_spol "github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/eventstore"
	"github.com/caos/zitadel/internal/management/repository/eventsourcing/spooler"
	mgmt_view "github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/view"

	es_int "github.com/caos/zitadel/internal/eventstore"
	es_proj "github.com/caos/zitadel/internal/project/repository/eventsourcing"
)

type Config struct {
	Eventstore es_int.Config
	View       view.ViewConfig
	Spooler    spooler.SpoolerConfig
}

type EsRepository struct {
	spooler *es_spol.Spooler
	eventstore.ProjectRepo
}

func Start(conf Config, systemDefaults sd.SystemDefaults) (*EsRepository, error) {
	es, err := es_int.Start(conf.Eventstore)
	if err != nil {
		return nil, err
	}

	view, sql, err := mgmt_view.StartView(conf.View)
	if err != nil {
		return nil, err
	}

	conf.Spooler.View = view
	conf.Spooler.SQL = sql
	spool := spooler.StartSpooler(conf.Spooler, es)

	project, err := es_proj.StartProject(es_proj.ProjectConfig{
		Eventstore: es,
		Cache:      conf.Eventstore.Cache,
	}, systemDefaults)
	if err != nil {
		return nil, err
	}

	return &EsRepository{
		spool,
		eventstore.ProjectRepo{project, view},
	}, nil
}

func (repo *EsRepository) Health() error {
	return repo.ProjectEvents.Health(context.Background())
}

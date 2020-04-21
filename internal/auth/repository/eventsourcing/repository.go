package eventsourcing

import (
	"context"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_user_agent "github.com/caos/zitadel/internal/user_agent/repository/eventsourcing"
)

type Config struct {
	Eventstore es_int.Config
	//View       view.ViewConfig
	//Spooler    spooler.SpoolerConfig
}

type EsRepository struct {
	//spooler *es_spooler.Spooler
	UserAgentRepo
}

func Start(conf Config) (*EsRepository, error) {
	es, err := es_int.Start(conf.Eventstore)
	if err != nil {
		return nil, err
	}

	//view, sql, err := mgmt_view.StartView(conf.View)
	//if err != nil {
	//	return nil, err
	//}

	//conf.Spooler.View = view
	//conf.Spooler.EsClient = es.Client
	//conf.Spooler.SQL = sql
	//spool := spooler.StartSpooler(conf.Spooler)

	useragent, err := es_user_agent.StartUserAgent(es_user_agent.UserAgentConfig{Eventstore: es})
	if err != nil {
		return nil, err
	}

	return &EsRepository{
		UserAgentRepo{useragent},
	}, nil
}

func (repo *EsRepository) Health(ctx context.Context) error {
	return repo.UserAgentEvents.Health(ctx)
}

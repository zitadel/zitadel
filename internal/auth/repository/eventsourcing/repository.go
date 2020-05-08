package eventsourcing

import (
	"context"
	"github.com/caos/zitadel/internal/auth_request/repository/cache"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_user "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type Config struct {
	Eventstore  es_int.Config
	AuthRequest cache.Config
	//View       view.ViewConfig
	//Spooler    spooler.SpoolerConfig
}

type EsRepository struct {
	//spooler *es_spooler.Spooler
	UserRepo        UserRepo
	AuthRequestRepo AuthRequestRepo
}

func Start(conf Config, systemDefaults sd.SystemDefaults) (*EsRepository, error) {
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

	user, err := es_user.StartUser(
		es_user.UserConfig{
			Eventstore: es,
			Cache:      conf.Eventstore.Cache,
		},
		systemDefaults,
	)
	if err != nil {
		return nil, err
	}
	authReq, err := cache.Start(conf.AuthRequest)
	if err != nil {
		return nil, err
	}

	return &EsRepository{
		UserRepo{user},
		AuthRequestRepo{
			UserEvents:               user,
			AuthRequests:             authReq,
			PasswordCheckLifeTime:    systemDefaults.VerificationLifetimes.PasswordCheck.Duration,
			MfaInitSkippedLifeTime:   systemDefaults.VerificationLifetimes.MfaInitSkip.Duration,
			MfaSoftwareCheckLifeTime: systemDefaults.VerificationLifetimes.MfaSoftwareCheck.Duration,
			MfaHardwareCheckLifeTime: systemDefaults.VerificationLifetimes.MfaHardwareCheck.Duration,
		},
	}, nil
}

func (repo *EsRepository) Health(ctx context.Context) error {
	if err := repo.UserRepo.Health(ctx); err != nil {
		return err
	}
	return repo.AuthRequestRepo.Health(ctx)
}

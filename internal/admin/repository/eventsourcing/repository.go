package eventsourcing

import (
	"context"

	"github.com/caos/zitadel/internal/admin/repository/eventsourcing/eventstore"
	"github.com/caos/zitadel/internal/admin/repository/eventsourcing/handler"
	"github.com/caos/zitadel/internal/admin/repository/eventsourcing/spooler"
	admin_view "github.com/caos/zitadel/internal/admin/repository/eventsourcing/view"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_spol "github.com/caos/zitadel/internal/eventstore/spooler"
	es_iam "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	es_org "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	es_usr "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	iam_business "github.com/caos/zitadel/internal/v2/business/iam"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/member"
)

type Config struct {
	SearchLimit uint64
	Eventstore  es_int.Config
	View        types.SQL
	Spooler     spooler.SpoolerConfig
	Domain      string
}

type EsRepository struct {
	spooler *es_spol.Spooler
	eventstore.OrgRepo
	eventstore.IAMRepository
	eventstore.AdministratorRepo
	eventstore.UserRepo
}

func Start(ctx context.Context, conf Config, systemDefaults sd.SystemDefaults, roles []string) (*EsRepository, error) {
	es, err := es_int.Start(conf.Eventstore)
	if err != nil {
		return nil, err
	}
	esV2 := es.V2()
	esV2.RegisterFilterEventMapper(iam.MemberAddedEventType, member.AddedEventMapper).
		RegisterFilterEventMapper(iam.MemberChangedEventType, member.ChangedEventMapper).
		RegisterFilterEventMapper(iam.MemberRemovedEventType, member.RemovedEventMapper).
		RegisterFilterEventMapper(iam.IDPConfigAddedEventType, iam.IDPConfigAddedEventMapper).
		RegisterFilterEventMapper(iam.IDPConfigChangedEventType, iam.IDPConfigChangedEventMapper).
		RegisterFilterEventMapper(iam.IDPConfigDeactivatedEventType, iam.IDPConfigDeactivatedEventMapper).
		RegisterFilterEventMapper(iam.IDPConfigReactivatedEventType, iam.IDPConfigReactivatedEventMapper).
		RegisterFilterEventMapper(iam.IDPConfigRemovedEventType, iam.IDPConfigRemovedEventMapper).
		RegisterFilterEventMapper(iam.IDPOIDCConfigAddedEventType, iam.IDPOIDCConfigAddedEventMapper).
		RegisterFilterEventMapper(iam.IDPOIDCConfigChangedEventType, iam.IDPOIDCConfigChangedEventMapper)

	iam, err := es_iam.StartIAM(es_iam.IAMConfig{
		Eventstore: es,
		Cache:      conf.Eventstore.Cache,
	}, systemDefaults)
	if err != nil {
		return nil, err
	}

	org := es_org.StartOrg(es_org.OrgConfig{Eventstore: es, IAMDomain: conf.Domain}, systemDefaults)

	user, err := es_usr.StartUser(es_usr.UserConfig{
		Eventstore: es,
		Cache:      conf.Eventstore.Cache,
	}, systemDefaults)
	if err != nil {
		return nil, err
	}
	sqlClient, err := conf.View.Start()
	if err != nil {
		return nil, err
	}
	view, err := admin_view.StartView(sqlClient)
	if err != nil {
		return nil, err
	}
	iamV2, err := iam_business.StartRepository(&iam_business.Config{Eventstore: esV2, SystemDefaults: systemDefaults})
	if err != nil {
		return nil, err
	}

	spool := spooler.StartSpooler(conf.Spooler, es, view, sqlClient, handler.EventstoreRepos{UserEvents: user, OrgEvents: org, IamEvents: iam}, systemDefaults)

	return &EsRepository{
		spooler: spool,
		OrgRepo: eventstore.OrgRepo{
			Eventstore:     es,
			OrgEventstore:  org,
			UserEventstore: user,
			View:           view,
			SearchLimit:    conf.SearchLimit,
			SystemDefaults: systemDefaults,
		},
		IAMRepository: eventstore.IAMRepository{
			IAMEventstore:  iam,
			OrgEvents:      org,
			UserEvents:     user,
			View:           view,
			SystemDefaults: systemDefaults,
			SearchLimit:    conf.SearchLimit,
			Roles:          roles,
			IAMV2:          iamV2,
		},
		AdministratorRepo: eventstore.AdministratorRepo{
			View: view,
		},
		UserRepo: eventstore.UserRepo{
			UserEvents:     user,
			OrgEvents:      org,
			View:           view,
			SystemDefaults: systemDefaults,
		},
	}, nil
}

func (repo *EsRepository) Health(ctx context.Context) error {
	err := repo.Eventstore.Health(ctx)
	if err != nil {
		return err
	}
	err = repo.UserEventstore.Health(ctx)
	if err != nil {
		return err
	}
	return repo.OrgEventstore.Health(ctx)
}

package eventsourcing

import (
	"context"

	es_user "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	"github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/member"
	"github.com/caos/zitadel/internal/v2/repository/policy"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/auth_request/repository/cache"
	"github.com/caos/zitadel/internal/authz/repository/eventsourcing/eventstore"
	"github.com/caos/zitadel/internal/authz/repository/eventsourcing/handler"
	"github.com/caos/zitadel/internal/authz/repository/eventsourcing/spooler"
	authz_view "github.com/caos/zitadel/internal/authz/repository/eventsourcing/view"
	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	es_int "github.com/caos/zitadel/internal/eventstore"
	es_spol "github.com/caos/zitadel/internal/eventstore/spooler"
	es_iam "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	"github.com/caos/zitadel/internal/id"
	es_key "github.com/caos/zitadel/internal/key/repository/eventsourcing"
	es_proj "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	iam_business "github.com/caos/zitadel/internal/v2/business/iam"
)

type Config struct {
	Eventstore  es_int.Config
	AuthRequest cache.Config
	View        types.SQL
	Spooler     spooler.SpoolerConfig
	KeyConfig   es_key.KeyConfig
}

type EsRepository struct {
	spooler *es_spol.Spooler
	eventstore.UserGrantRepo
	eventstore.IamRepo
	eventstore.TokenVerifierRepo
}

func Start(conf Config, authZ authz.Config, systemDefaults sd.SystemDefaults) (*EsRepository, error) {
	es, err := es_int.Start(conf.Eventstore)
	if err != nil {
		return nil, err
	}
	esV2 := es.V2()
	esV2.RegisterFilterEventMapper(iam.SetupStartedEventType, iam.SetupStepMapper).
		RegisterFilterEventMapper(iam.SetupDoneEventType, iam.SetupStepMapper).
		RegisterFilterEventMapper(iam.GlobalOrgSetEventType, iam.GlobalOrgSetMapper).
		RegisterFilterEventMapper(iam.ProjectSetEventType, iam.ProjectSetMapper).
		RegisterFilterEventMapper(iam.LabelPolicyAddedEventType, policy.LabelPolicyAddedEventMapper).
		RegisterFilterEventMapper(iam.LabelPolicyChangedEventType, policy.LabelPolicyChangedEventMapper).
		RegisterFilterEventMapper(iam.LoginPolicyAddedEventType, policy.LoginPolicyAddedEventMapper).
		RegisterFilterEventMapper(iam.LoginPolicyChangedEventType, policy.LoginPolicyChangedEventMapper).
		RegisterFilterEventMapper(iam.OrgIAMPolicyAddedEventType, policy.OrgIAMPolicyAddedEventMapper).
		RegisterFilterEventMapper(iam.PasswordAgePolicyAddedEventType, policy.PasswordAgePolicyAddedEventMapper).
		RegisterFilterEventMapper(iam.PasswordAgePolicyChangedEventType, policy.PasswordAgePolicyChangedEventMapper).
		RegisterFilterEventMapper(iam.PasswordComplexityPolicyAddedEventType, policy.PasswordComplexityPolicyAddedEventMapper).
		RegisterFilterEventMapper(iam.PasswordComplexityPolicyChangedEventType, policy.PasswordComplexityPolicyChangedEventMapper).
		RegisterFilterEventMapper(iam.PasswordLockoutPolicyAddedEventType, policy.PasswordLockoutPolicyAddedEventMapper).
		RegisterFilterEventMapper(iam.PasswordLockoutPolicyChangedEventType, policy.PasswordLockoutPolicyChangedEventMapper).
		RegisterFilterEventMapper(iam.MemberAddedEventType, member.AddedEventMapper).
		RegisterFilterEventMapper(iam.MemberChangedEventType, member.ChangedEventMapper).
		RegisterFilterEventMapper(iam.MemberRemovedEventType, member.RemovedEventMapper)

	sqlClient, err := conf.View.Start()
	if err != nil {
		return nil, err
	}

	idGenerator := id.SonyFlakeGenerator
	view, err := authz_view.StartView(sqlClient, idGenerator)
	if err != nil {
		return nil, err
	}

	iam, err := es_iam.StartIAM(es_iam.IAMConfig{
		Eventstore: es,
		Cache:      conf.Eventstore.Cache,
	}, systemDefaults)
	if err != nil {
		return nil, err
	}
	project, err := es_proj.StartProject(es_proj.ProjectConfig{
		Eventstore: es,
		Cache:      conf.Eventstore.Cache,
	}, systemDefaults)
	if err != nil {
		return nil, err
	}
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
	iamV2, err := iam_business.StartRepository(&iam_business.Config{Eventstore: esV2, SystemDefaults: systemDefaults})
	if err != nil {
		return nil, err
	}

	repos := handler.EventstoreRepos{IamEvents: iam}
	spool := spooler.StartSpooler(conf.Spooler, es, view, sqlClient, repos, systemDefaults)

	return &EsRepository{
		spool,
		eventstore.UserGrantRepo{
			View:      view,
			IamID:     systemDefaults.IamID,
			Auth:      authZ,
			IamEvents: iam,
		},
		eventstore.IamRepo{
			IAMID:     systemDefaults.IamID,
			IAMEvents: iam,
			IAMV2:     iamV2,
		},
		eventstore.TokenVerifierRepo{
			//TODO: Add Token Verification Key
			IAMID:         systemDefaults.IamID,
			IAMEvents:     iam,
			ProjectEvents: project,
			UserEvents:    user,
			View:          view,
		},
	}, nil
}

func (repo *EsRepository) Health(ctx context.Context) error {
	if err := repo.UserGrantRepo.Health(); err != nil {
		return err
	}
	return nil
}

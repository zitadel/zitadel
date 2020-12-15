package command

import (
	"context"
	"github.com/caos/zitadel/internal/v2/business/query"

	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

type CommandSide struct {
	eventstore   *eventstore.Eventstore
	idGenerator  id.Generator
	secretCrypto crypto.Crypto
}

type Config struct {
	Eventstore     *eventstore.Eventstore
	SystemDefaults sd.SystemDefaults
}

func StartCommandSide(config *Config) (repo *CommandSide, err error) {
	repo = &CommandSide{
		eventstore:  config.Eventstore,
		idGenerator: id.SonyFlakeGenerator,
	}
	iam_repo.RegisterEventMappers(repo.eventstore)

	repo.secretCrypto, err = crypto.NewAESCrypto(config.SystemDefaults.IDPConfigVerificationKey)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *CommandSide) iamByID(ctx context.Context, id string) (_ *query.ReadModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	readModel := query.NewReadModel(id)
	err = r.eventstore.FilterToQueryReducer(ctx, readModel)
	if err != nil {
		return nil, err
	}

	return readModel, nil
}

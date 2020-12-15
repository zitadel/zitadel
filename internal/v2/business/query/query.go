package query

import (
	"context"

	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore/v2"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

type QuerySide struct {
	eventstore   *eventstore.Eventstore
	idGenerator  id.Generator
	secretCrypto crypto.Crypto
}

type Config struct {
	Eventstore     *eventstore.Eventstore
	SystemDefaults sd.SystemDefaults
}

func StartQuerySide(config *Config) (repo *QuerySide, err error) {
	repo = &QuerySide{
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

func (r *QuerySide) IAMByID(ctx context.Context, id string) (_ *iam_model.IAM, err error) {
	readModel, err := r.iamByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return readModelToIAM(readModel), nil
}

func (r *QuerySide) iamByID(ctx context.Context, id string) (_ *iam_repo.ReadModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	readModel := iam_repo.NewReadModel(id)
	err = r.eventstore.FilterToQueryReducer(ctx, readModel)
	if err != nil {
		return nil, err
	}

	return readModel, nil
}

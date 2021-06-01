package query

import (
	"context"

	sd "github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/config/types"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/eventstore"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/id"
	iam_repo "github.com/caos/zitadel/internal/repository/iam"
	usr_repo "github.com/caos/zitadel/internal/repository/user"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

type Queries struct {
	iamID        string
	eventstore   *eventstore.Eventstore
	idGenerator  id.Generator
	secretCrypto crypto.Crypto
}

type Config struct {
	Eventstore types.SQLUser
}

func StartQueries(eventstore *eventstore.Eventstore, defaults sd.SystemDefaults) (repo *Queries, err error) {
	repo = &Queries{
		iamID:       defaults.IamID,
		eventstore:  eventstore,
		idGenerator: id.SonyFlakeGenerator,
	}
	iam_repo.RegisterEventMappers(repo.eventstore)
	usr_repo.RegisterEventMappers(repo.eventstore)

	repo.secretCrypto, err = crypto.NewAESCrypto(defaults.IDPConfigVerificationKey)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *Queries) IAMByID(ctx context.Context, id string) (_ *iam_model.IAM, err error) {
	readModel, err := r.iamByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return readModelToIAM(readModel), nil
}

func (r *Queries) iamByID(ctx context.Context, id string) (_ *ReadModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	readModel := NewReadModel(id)
	err = r.eventstore.FilterToQueryReducer(ctx, readModel)
	if err != nil {
		return nil, err
	}

	return readModel, nil
}

package iam

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/tracing"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

type Repository struct {
	eventstore *eventstore.Eventstore
}

type Config struct {
	Eventstore *eventstore.Eventstore
}

func StartRepository(config *Config) *Repository {
	return &Repository{
		eventstore: config.Eventstore,
	}
}

func (r *Repository) IAMByID(ctx context.Context, id string) (_ *iam_model.IAM, err error) {
	readModel, err := r.iamByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return readModelToIAM(readModel), nil
}

func (r *Repository) iamByID(ctx context.Context, id string) (_ *iam_repo.ReadModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	readModel := iam_repo.NewReadModel(id)
	err = r.eventstore.FilterToQueryReducer(ctx, readModel)
	if err != nil {
		return nil, err
	}

	return readModel, nil
}

func (r *Repository) memberWriteModelByID(ctx context.Context, iamID, userID string) (member *iam_repo.MemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel := iam_repo.NewMemberReadModel(iamID, userID)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}

	if writeModel.IsRemoved {
		return nil, errors.ThrowNotFound(nil, "IAM-D8JxR", "Errors.NotFound")
	}

	return writeModel, nil
}

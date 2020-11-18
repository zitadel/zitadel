package iam

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/tracing"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
	"github.com/caos/zitadel/internal/v2/repository/member"
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

	query := eventstore.NewSearchQueryFactory(eventstore.ColumnsEvent, iam_repo.AggregateType).AggregateIDs(id)

	readModel := new(iam_repo.ReadModel)
	err = r.eventstore.FilterToReducer(ctx, query, readModel)
	if err != nil {
		return nil, err
	}

	return readModel, nil
}

func (r *Repository) memberWriteModelByID(ctx context.Context, iamID, userID string) (member *iam_repo.MemberWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query := eventstore.NewSearchQueryFactory(eventstore.ColumnsEvent, iam_repo.AggregateType).AggregateIDs(iamID)

	writeModel := new(memberWriteModel)
	err = r.eventstore.FilterToReducer(ctx, query, writeModel)
	if err != nil {
		return nil, err
	}

	if writeModel.isDeleted {
		return nil, errors.ThrowNotFound(nil, "IAM-D8JxR", "Errors.NotFound")
	}

	return &writeModel.MemberWriteModel, nil
}

type memberWriteModel struct {
	iam_repo.MemberWriteModel

	userID    string
	isDeleted bool
}

func (wm *memberWriteModel) AppendEvents(events ...eventstore.EventReader) error {
	for _, event := range events {
		switch e := event.(type) {
		case *member.AddedEvent:
			if e.UserID == wm.userID {
				wm.isDeleted = false
				wm.MemberWriteModel.AppendEvents(event)
			}
		case *member.ChangedEvent:
			if e.UserID == wm.userID {
				wm.MemberWriteModel.AppendEvents(event)
			}
		case *member.RemovedEvent:
			if e.UserID == wm.userID {
				wm.isDeleted = true
				wm.MemberWriteModel = iam_repo.MemberWriteModel{}
			}
		}
	}

	return nil
}

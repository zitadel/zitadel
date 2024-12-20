package query

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (q *Queries) HasHumanEmailCode(ctx context.Context, userID string) (_ bool, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if userID == "" {
		return false, zerrors.ThrowInvalidArgument(nil, "QUERY-4Mfsf", "Errors.User.UserIDMissing")
	}
	model, err := q.emailReadModel(ctx, userID, "")
	if err != nil {
		return false, err
	}
	if model.UserState == domain.UserStateUnspecified || model.UserState == domain.UserStateDeleted {
		return false, zerrors.ThrowNotFound(nil, "QUERY-ieJ2e", "Errors.User.NotFound")
	}
	if model.UserState == domain.UserStateInitial {
		return false, zerrors.ThrowPreconditionFailed(nil, "QUERY-uz0Uu", "Errors.User.NotInitialised")
	}
	ctxData := authz.GetCtxData(ctx)
	if ctxData.UserID != userID {
		if err := q.checkPermission(ctx, domain.PermissionUserRead, model.ResourceOwner, model.AggregateID); err != nil {
			return false, err
		}
	}
	return model.Code != nil, nil
}

func (q *Queries) emailReadModel(ctx context.Context, userID, resourceOwner string) (readModel *HumanEmailReadModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	readModel = NewHumanEmailReadModel(userID, resourceOwner)
	err = q.eventstore.FilterToQueryReducer(ctx, readModel)
	if err != nil {
		return nil, err
	}
	return readModel, nil
}

type HumanEmailReadModel struct {
	*eventstore.ReadModel

	Email           domain.EmailAddress
	IsEmailVerified bool

	Code             *crypto.CryptoValue
	CodeCreationDate time.Time
	CodeExpiry       time.Duration
	AuthRequestID    string

	UserState domain.UserState
}

func NewHumanEmailReadModel(userID, resourceOwner string) *HumanEmailReadModel {
	return &HumanEmailReadModel{
		ReadModel: &eventstore.ReadModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *HumanEmailReadModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			wm.Email = e.EmailAddress
			wm.UserState = domain.UserStateActive
		case *user.HumanRegisteredEvent:
			wm.Email = e.EmailAddress
			wm.UserState = domain.UserStateActive
		case *user.HumanInitialCodeAddedEvent:
			wm.UserState = domain.UserStateInitial
		case *user.HumanInitializedCheckSucceededEvent:
			wm.UserState = domain.UserStateActive
		case *user.HumanEmailChangedEvent:
			wm.Email = e.EmailAddress
			wm.IsEmailVerified = false
			wm.Code = nil
		case *user.HumanEmailCodeAddedEvent:
			wm.Code = e.Code
			wm.CodeCreationDate = e.CreationDate()
			wm.CodeExpiry = e.Expiry
			wm.AuthRequestID = e.AuthRequestID
		case *user.HumanEmailVerifiedEvent:
			wm.IsEmailVerified = true
			wm.Code = nil
		case *user.UserRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		}
	}
	return wm.ReadModel.Reduce()
}

func (wm *HumanEmailReadModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(user.UserV1AddedType,
			user.HumanAddedType,
			user.UserV1RegisteredType,
			user.HumanRegisteredType,
			user.UserV1InitialCodeAddedType,
			user.HumanInitialCodeAddedType,
			user.UserV1InitializedCheckSucceededType,
			user.HumanInitializedCheckSucceededType,
			user.UserV1EmailChangedType,
			user.HumanEmailChangedType,
			user.UserV1EmailCodeAddedType,
			user.HumanEmailCodeAddedType,
			user.UserV1EmailVerifiedType,
			user.HumanEmailVerifiedType,
			user.UserRemovedType).
		Builder()

	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}

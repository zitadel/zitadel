package query

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type HumanPasswordReadModel struct {
	*eventstore.ReadModel

	EncodedHash          string
	SecretChangeRequired bool

	Code                     *crypto.CryptoValue
	CodeCreationDate         time.Time
	CodeExpiry               time.Duration
	PasswordCheckFailedCount uint64

	UserState domain.UserState
}

func (q *Queries) GetHumanPassword(ctx context.Context, orgID, userID string) (encodedHash string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if userID == "" {
		return "", zerrors.ThrowInvalidArgument(nil, "QUERY-4Mfsf", "Errors.User.UserIDMissing")
	}
	existingPassword, err := q.passwordReadModel(ctx, userID, orgID)
	if err != nil {
		return "", zerrors.ThrowInternal(nil, "QUERY-p1k1n2i", "Errors.User.NotFound")
	}
	if existingPassword.UserState == domain.UserStateUnspecified || existingPassword.UserState == domain.UserStateDeleted {
		return "", zerrors.ThrowPreconditionFailed(nil, "QUERY-3n77z", "Errors.User.NotFound")
	}
	return existingPassword.EncodedHash, nil
}

func (q *Queries) passwordReadModel(ctx context.Context, userID, resourceOwner string) (readModel *HumanPasswordReadModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	readModel = NewHumanPasswordReadModel(userID, resourceOwner)
	err = q.eventstore.FilterToQueryReducer(ctx, readModel)
	if err != nil {
		return nil, err
	}
	return readModel, nil
}

func NewHumanPasswordReadModel(userID, resourceOwner string) *HumanPasswordReadModel {
	return &HumanPasswordReadModel{
		ReadModel: &eventstore.ReadModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (rm *HumanPasswordReadModel) AppendEvents(events ...eventstore.Event) {
	rm.ReadModel.AppendEvents(events...)
}

func (wm *HumanPasswordReadModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			wm.EncodedHash = crypto.SecretOrEncodedHash(e.Secret, e.EncodedHash)
			wm.SecretChangeRequired = e.ChangeRequired
			wm.UserState = domain.UserStateActive
		case *user.HumanRegisteredEvent:
			wm.EncodedHash = crypto.SecretOrEncodedHash(e.Secret, e.EncodedHash)
			wm.SecretChangeRequired = e.ChangeRequired
			wm.UserState = domain.UserStateActive
		case *user.HumanInitialCodeAddedEvent:
			wm.UserState = domain.UserStateInitial
		case *user.HumanInitializedCheckSucceededEvent:
			wm.UserState = domain.UserStateActive
		case *user.HumanPasswordChangedEvent:
			wm.EncodedHash = crypto.SecretOrEncodedHash(e.Secret, e.EncodedHash)
			wm.SecretChangeRequired = e.ChangeRequired
			wm.Code = nil
			wm.PasswordCheckFailedCount = 0
		case *user.HumanPasswordCodeAddedEvent:
			wm.Code = e.Code
			wm.CodeCreationDate = e.CreationDate()
			wm.CodeExpiry = e.Expiry
		case *user.HumanEmailVerifiedEvent:
			if wm.UserState == domain.UserStateInitial {
				wm.UserState = domain.UserStateActive
			}
		case *user.HumanPasswordCheckFailedEvent:
			wm.PasswordCheckFailedCount += 1
		case *user.HumanPasswordCheckSucceededEvent:
			wm.PasswordCheckFailedCount = 0
		case *user.UserLockedEvent:
			wm.UserState = domain.UserStateLocked
		case *user.UserUnlockedEvent:
			wm.PasswordCheckFailedCount = 0
			if wm.UserState != domain.UserStateDeleted {
				wm.UserState = domain.UserStateActive
			}
		case *user.UserRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		case *user.HumanPasswordHashUpdatedEvent:
			wm.EncodedHash = e.EncodedHash
		}
	}
	return wm.ReadModel.Reduce()
}

func (wm *HumanPasswordReadModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AwaitOpenTransactions().
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(user.HumanAddedType,
			user.HumanRegisteredType,
			user.HumanInitialCodeAddedType,
			user.HumanInitializedCheckSucceededType,
			user.HumanPasswordChangedType,
			user.HumanPasswordCodeAddedType,
			user.HumanEmailVerifiedType,
			user.HumanPasswordCheckFailedType,
			user.HumanPasswordCheckSucceededType,
			user.HumanPasswordHashUpdatedType,
			user.UserRemovedType,
			user.UserLockedType,
			user.UserUnlockedType,
			user.UserV1AddedType,
			user.UserV1RegisteredType,
			user.UserV1InitialCodeAddedType,
			user.UserV1InitializedCheckSucceededType,
			user.UserV1PasswordChangedType,
			user.UserV1PasswordCodeAddedType,
			user.UserV1EmailVerifiedType,
			user.UserV1PasswordCheckFailedType,
			user.UserV1PasswordCheckSucceededType,
		).
		Builder()

	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}

package query

import (
	"context"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"time"
)

type HumanPasswordWriteModel struct {
	eventstore.WriteModel

	Secret               *crypto.CryptoValue
	SecretChangeRequired bool

	Code                     *crypto.CryptoValue
	CodeCreationDate         time.Time
	CodeExpiry               time.Duration
	PasswordCheckFailedCount uint64

	UserState domain.UserState
}

func (q *Queries) GetHumanPassword(ctx context.Context, orgID, userID string) (passwordHash []byte, algorithm string, err error) {
	if userID == "" {
		return nil, "", caos_errs.ThrowInvalidArgument(nil, "QUERY-4Mfsf", "Errors.User.UserIDMissing")
	}
	existingPassword, err := q.passwordWriteModel(ctx, userID, orgID)
	if existingPassword.UserState == domain.UserStateUnspecified || existingPassword.UserState == domain.UserStateDeleted {
		return nil, "", caos_errs.ThrowPreconditionFailed(nil, "QUERY-3n77z", "Errors.User.NotFound")
	}

	if existingPassword.Secret != nil && existingPassword.Secret.Crypted != nil {
		return existingPassword.Secret.Crypted, existingPassword.Secret.Algorithm, nil
	}

	return nil, "", caos_errs.ThrowNotFound(nil, "QUERY-lp1m2i", "Errors.User.NotFound")
}

func (q *Queries) passwordWriteModel(ctx context.Context, userID, resourceOwner string) (writeModel *HumanPasswordWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanPasswordWriteModel(userID, resourceOwner)
	err = q.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

func NewHumanPasswordWriteModel(userID, resourceOwner string) *HumanPasswordWriteModel {
	return &HumanPasswordWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *HumanPasswordWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			wm.Secret = e.Secret
			wm.SecretChangeRequired = e.ChangeRequired
			wm.UserState = domain.UserStateActive
		case *user.HumanRegisteredEvent:
			wm.Secret = e.Secret
			wm.SecretChangeRequired = e.ChangeRequired
			wm.UserState = domain.UserStateActive
		case *user.HumanInitialCodeAddedEvent:
			wm.UserState = domain.UserStateInitial
		case *user.HumanInitializedCheckSucceededEvent:
			wm.UserState = domain.UserStateActive
		case *user.HumanPasswordChangedEvent:
			wm.Secret = e.Secret
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
		case *user.UserUnlockedEvent:
			wm.PasswordCheckFailedCount = 0
		case *user.UserRemovedEvent:
			wm.UserState = domain.UserStateDeleted
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanPasswordWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
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
			user.UserRemovedType,
			user.UserUnlockedType,
			user.UserV1AddedType,
			user.UserV1RegisteredType,
			user.UserV1InitialCodeAddedType,
			user.UserV1InitializedCheckSucceededType,
			user.UserV1PasswordChangedType,
			user.UserV1PasswordCodeAddedType,
			user.UserV1EmailVerifiedType,
			user.UserV1PasswordCheckFailedType,
			user.UserV1PasswordCheckSucceededType).
		Builder()

	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}

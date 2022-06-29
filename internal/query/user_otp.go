package query

import (
	"context"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (q *Queries) GetHumanOTPSecret(ctx context.Context, userID, resourceowner string) (string, error) {
	if userID == "" {
		return "", caos_errs.ThrowPreconditionFailed(nil, "QUERY-8N9ds", "Errors.User.UserIDMissing")
	}
	existingOTP, err := q.otpWriteModelByID(ctx, userID, resourceowner)
	if err != nil {
		return "", err
	}
	if existingOTP.State != domain.MFAStateReady {
		return "", caos_errs.ThrowNotFound(nil, "QUERY-01982h", "Errors.User.NotFound")
	}

	return crypto.DecryptString(existingOTP.Secret, q.multifactors.OTP.CryptoMFA)
}

func (q *Queries) otpWriteModelByID(ctx context.Context, userID, resourceOwner string) (writeModel *HumanOTPWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanOTPWriteModel(userID, resourceOwner)
	err = q.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

type HumanOTPWriteModel struct {
	eventstore.WriteModel

	State  domain.MFAState
	Secret *crypto.CryptoValue
}

func NewHumanOTPWriteModel(userID, resourceOwner string) *HumanOTPWriteModel {
	return &HumanOTPWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *HumanOTPWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanOTPAddedEvent:
			wm.Secret = e.Secret
			wm.State = domain.MFAStateNotReady
		case *user.HumanOTPVerifiedEvent:
			wm.State = domain.MFAStateReady
		case *user.HumanOTPRemovedEvent:
			wm.State = domain.MFAStateRemoved
		case *user.UserRemovedEvent:
			wm.State = domain.MFAStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *HumanOTPWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(user.HumanMFAOTPAddedType,
			user.HumanMFAOTPVerifiedType,
			user.HumanMFAOTPRemovedType,
			user.UserRemovedType,
			user.UserV1MFAOTPAddedType,
			user.UserV1MFAOTPVerifiedType,
			user.UserV1MFAOTPRemovedType).
		Builder()

	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}

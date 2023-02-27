package query

import (
	"context"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
)

func (q *Queries) GetHumanOTPSecret(ctx context.Context, userID, resourceowner string) (_ string, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if userID == "" {
		return "", errors.ThrowPreconditionFailed(nil, "QUERY-8N9ds", "Errors.User.UserIDMissing")
	}
	existingOTP, err := q.otpReadModelByID(ctx, userID, resourceowner)
	if err != nil {
		return "", err
	}
	if existingOTP.State != domain.MFAStateReady {
		return "", errors.ThrowNotFound(nil, "QUERY-01982h", "Errors.User.NotFound")
	}

	return crypto.DecryptString(existingOTP.Secret, q.multifactors.OTP.CryptoMFA)
}

func (q *Queries) otpReadModelByID(ctx context.Context, userID, resourceOwner string) (readModel *HumanOTPReadModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	readModel = NewHumanOTPReadModel(userID, resourceOwner)
	err = q.eventstore.FilterToQueryReducer(ctx, readModel)
	if err != nil {
		return nil, err
	}
	return readModel, nil
}

type HumanOTPReadModel struct {
	*eventstore.ReadModel

	State  domain.MFAState
	Secret *crypto.CryptoValue
}

func (rm *HumanOTPReadModel) AppendEvents(events ...eventstore.Event) {
	rm.ReadModel.AppendEvents(events...)
}

func NewHumanOTPReadModel(userID, resourceOwner string) *HumanOTPReadModel {
	return &HumanOTPReadModel{
		ReadModel: &eventstore.ReadModel{
			AggregateID:   userID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *HumanOTPReadModel) Reduce() error {
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
	return wm.ReadModel.Reduce()
}

func (wm *HumanOTPReadModel) Query() *eventstore.SearchQueryBuilder {
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

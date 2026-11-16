package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type HumanPasswordWriteModel struct {
	eventstore.WriteModel

	EncodedHash          string
	SecretChangeRequired bool

	Code                     *crypto.CryptoValue
	CodeCreationDate         time.Time
	CodeExpiry               time.Duration
	PasswordCheckFailedCount uint64
	GeneratorID              string
	VerificationID           string

	UserState domain.UserState

	// checksAfterSequence limits the password check events to the ones after this (aggregate) sequence.
	// It must be the sequence of the last event resetting the PasswordCheckFailedCount
	// (see [lastPasswordCheckResetSequence]); 0 loads all password checks.
	// The sequence is used instead of the position, as historical positions can disagree with the sequence order.
	checksAfterSequence uint64
}

// lastPasswordCheckResetSequence returns the aggregate sequence of the last event of the user,
// which resets the [HumanPasswordWriteModel.PasswordCheckFailedCount], or 0 if there is none.
func lastPasswordCheckResetSequence(ctx context.Context, es *eventstore.Eventstore, userID string) (uint64, error) {
	events, err := es.Filter(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		OrderDesc().
		Limit(1).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(userID).
		EventTypes(humanPasswordCheckResetEventTypes...).
		Builder(),
	)
	if err != nil || len(events) == 0 {
		return 0, err
	}
	return events[0].Sequence(), nil
}

// loadHumanPasswordWriteModel loads the password state of the user,
// but only the password checks since the last reset of the failed count,
// so the cost doesn't grow with every login of the user.
func loadHumanPasswordWriteModel(ctx context.Context, es *eventstore.Eventstore, userID, resourceOwner string) (*HumanPasswordWriteModel, error) {
	checksAfterSequence, err := lastPasswordCheckResetSequence(ctx, es, userID)
	if err != nil {
		return nil, err
	}
	wm := NewHumanPasswordWriteModel(userID, resourceOwner)
	wm.checksAfterSequence = checksAfterSequence
	if err = es.FilterToQueryReducer(ctx, wm); err != nil {
		return nil, err
	}
	return wm, nil
}

func (wm *HumanPasswordWriteModel) GetUserState() domain.UserState {
	return wm.UserState
}

func (wm *HumanPasswordWriteModel) GetPasswordCheckFailedCount() uint64 {
	return wm.PasswordCheckFailedCount
}

func (wm *HumanPasswordWriteModel) GetEncodedHash() string {
	return wm.EncodedHash
}

func (wm *HumanPasswordWriteModel) GetResourceOwner() string {
	return wm.ResourceOwner
}

func (wm *HumanPasswordWriteModel) GetWriteModel() *eventstore.WriteModel {
	return &wm.WriteModel
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
			wm.Code = nil
			wm.CodeCreationDate = time.Time{}
			wm.CodeExpiry = 0
			wm.PasswordCheckFailedCount = 0
			wm.GeneratorID = ""
			wm.VerificationID = ""
			wm.EncodedHash = crypto.SecretOrEncodedHash(e.Secret, e.EncodedHash)
			wm.SecretChangeRequired = e.ChangeRequired
			wm.UserState = domain.UserStateActive
		case *user.HumanRegisteredEvent:
			wm.Code = nil
			wm.CodeCreationDate = time.Time{}
			wm.CodeExpiry = 0
			wm.PasswordCheckFailedCount = 0
			wm.GeneratorID = ""
			wm.VerificationID = ""
			wm.EncodedHash = crypto.SecretOrEncodedHash(e.Secret, e.EncodedHash)
			wm.SecretChangeRequired = e.ChangeRequired
			wm.UserState = domain.UserStateActive
		case *user.MachineAddedEvent:
			wm.Code = nil
			wm.CodeCreationDate = time.Time{}
			wm.CodeExpiry = 0
			wm.PasswordCheckFailedCount = 0
			wm.GeneratorID = ""
			wm.VerificationID = ""
			wm.EncodedHash = ""
			wm.SecretChangeRequired = false
			wm.UserState = domain.UserStateUnspecified
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
			wm.GeneratorID = e.GeneratorID
		case *user.HumanPasswordCodeSentEvent:
			wm.GeneratorID = e.GeneratorInfo.GetID()
			wm.VerificationID = e.GeneratorInfo.GetVerificationID()
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
			wm.EncodedHash = ""
			wm.SecretChangeRequired = false
			wm.Code = nil
			wm.CodeCreationDate = time.Time{}
			wm.CodeExpiry = 0
			wm.PasswordCheckFailedCount = 0
			wm.GeneratorID = ""
			wm.VerificationID = ""
			wm.UserState = domain.UserStateDeleted
		case *user.HumanPasswordHashUpdatedEvent:
			wm.EncodedHash = e.EncodedHash
		}
	}
	return wm.WriteModel.Reduce()
}

// humanPasswordStateEventTypes are the events which define the password state of a user,
// except for the password checks (see [humanPasswordCheckEventTypes]).
var humanPasswordStateEventTypes = []eventstore.EventType{
	user.HumanAddedType,
	user.HumanRegisteredType,
	user.MachineAddedEventType,
	user.HumanInitialCodeAddedType,
	user.HumanInitializedCheckSucceededType,
	user.HumanPasswordChangedType,
	user.HumanPasswordCodeAddedType,
	user.HumanPasswordCodeSentType,
	user.HumanEmailVerifiedType,
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
	user.UserV1PasswordCodeSentType,
	user.UserV1EmailVerifiedType,
}

// humanPasswordCheckEventTypes are the password check events.
// They are created on every password check and therefore grow with every login of a user.
var humanPasswordCheckEventTypes = []eventstore.EventType{
	user.HumanPasswordCheckFailedType,
	user.HumanPasswordCheckSucceededType,
	user.UserV1PasswordCheckFailedType,
	user.UserV1PasswordCheckSucceededType,
}

// humanPasswordCheckResetEventTypes are the events which reset the
// [HumanPasswordWriteModel.PasswordCheckFailedCount] to 0.
var humanPasswordCheckResetEventTypes = []eventstore.EventType{
	user.HumanAddedType,
	user.HumanRegisteredType,
	user.MachineAddedEventType,
	user.HumanPasswordChangedType,
	user.HumanPasswordCheckSucceededType,
	user.UserRemovedType,
	user.UserUnlockedType,
	user.UserV1AddedType,
	user.UserV1RegisteredType,
	user.UserV1PasswordChangedType,
	user.UserV1PasswordCheckSucceededType,
}

func (wm *HumanPasswordWriteModel) Query() *eventstore.SearchQueryBuilder {
	// Only the password checks after the last reset of the failed count are relevant.
	// This prevents loading the whole (constantly growing) login history of the user.
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(humanPasswordStateEventTypes...).
		Builder().
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(humanPasswordCheckEventTypes...).
		SequenceGreater(wm.checksAfterSequence).
		Builder()

	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	if wm.WriteModel.ProcessedSequence != 0 {
		query.SequenceGreater(wm.WriteModel.ProcessedSequence)
	}
	return query
}

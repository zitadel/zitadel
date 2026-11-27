package command

import (
	"context"
	"slices"
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
	// passwordChangedSequence limits the password hash events to the ones from this (aggregate) sequence on.
	// It must be the sequence of the last password change of the user (see [humanPasswordLoadBounds]);
	// 0 loads all password hash events.
	passwordChangedSequence uint64
}

// humanPasswordBoundsLookupLimit is the number of most recent reset events
// [lastHumanPasswordLoadBounds] looks at to find the last password change.
// Every successful password check is a reset as well, so the last change is only found if it is among the most recent resets.
// The limit keeps the lookup short for users who log in often but never change their password:
// without it, the lookup would scan their whole history to find out there is no change.
const humanPasswordBoundsLookupLimit = 10

// humanPasswordLoadBounds are the (aggregate) sequences of the last events of a user,
// which make earlier events irrelevant for the [HumanPasswordWriteModel].
type humanPasswordLoadBounds struct {
	events []eventstore.Event

	// lastResetSequence is the sequence of the last event resetting the [HumanPasswordWriteModel.PasswordCheckFailedCount].
	lastResetSequence uint64
	// lastPasswordChangedSequence is the sequence of the last password change, if it is one of the looked up resets.
	lastPasswordChangedSequence uint64
}

func (b *humanPasswordLoadBounds) AppendEvents(events ...eventstore.Event) {
	b.events = append(b.events, events...)
}

func (b *humanPasswordLoadBounds) Reduce() error {
	for _, event := range b.events {
		b.lastResetSequence = max(b.lastResetSequence, event.Sequence())
		if slices.Contains(humanPasswordChangedEventTypes, event.Type()) {
			b.lastPasswordChangedSequence = max(b.lastPasswordChangedSequence, event.Sequence())
		}
	}
	b.events = b.events[:0]
	return nil
}

// lastHumanPasswordLoadBounds returns the sequence of the last reset of the [HumanPasswordWriteModel.PasswordCheckFailedCount]
// (like [lastPasswordCheckResetSequence]) and of the last password change of the user in a single lookup.
// A password change is a reset as well. If it is not among the last [humanPasswordBoundsLookupLimit] resets,
// its sequence is 0 and all password hash events are loaded.
func lastHumanPasswordLoadBounds(ctx context.Context, es *eventstore.Eventstore, userID string) (*humanPasswordLoadBounds, error) {
	bounds := new(humanPasswordLoadBounds)
	err := es.FilterToReducer(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		OrderDesc().
		Limit(humanPasswordBoundsLookupLimit).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(userID).
		EventTypes(humanPasswordCheckResetEventTypes...).
		Builder(),
		bounds,
	)
	return bounds, err
}

// lastPasswordCheckResetSequence returns the aggregate sequence of the last event of the user,
// which resets the [HumanPasswordWriteModel.PasswordCheckFailedCount], or 0 if there is none.
func lastPasswordCheckResetSequence(ctx context.Context, es *eventstore.Eventstore, userID string) (uint64, error) {
	// the [eventstore.WriteModel] keeps the sequence of the reduced event as ProcessedSequence
	var reset eventstore.WriteModel
	err := es.FilterToReducer(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		OrderDesc().
		Limit(1).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(userID).
		EventTypes(humanPasswordCheckResetEventTypes...).
		Builder(),
		&reset,
	)
	return reset.ProcessedSequence, err
}

// loadHumanPasswordWriteModel loads the password state of the user,
// but only the password checks since the last reset of the failed count
// and the password hash events since the last password change,
// so the cost doesn't grow with every login or password change of the user.
func loadHumanPasswordWriteModel(ctx context.Context, es *eventstore.Eventstore, userID, resourceOwner string) (*HumanPasswordWriteModel, error) {
	bounds, err := lastHumanPasswordLoadBounds(ctx, es, userID)
	if err != nil {
		return nil, err
	}
	wm := NewHumanPasswordWriteModel(userID, resourceOwner)
	wm.checksAfterSequence = bounds.lastResetSequence
	wm.passwordChangedSequence = bounds.lastPasswordChangedSequence
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
// except for the password hash changes (see [humanPasswordHashEventTypes])
// and the password checks (see [humanPasswordCheckEventTypes]).
var humanPasswordStateEventTypes = []eventstore.EventType{
	user.HumanAddedType,
	user.HumanRegisteredType,
	user.MachineAddedEventType,
	user.HumanInitialCodeAddedType,
	user.HumanInitializedCheckSucceededType,
	user.HumanPasswordCodeAddedType,
	user.HumanPasswordCodeSentType,
	user.HumanEmailVerifiedType,
	user.UserRemovedType,
	user.UserLockedType,
	user.UserUnlockedType,
	user.UserV1AddedType,
	user.UserV1RegisteredType,
	user.UserV1InitialCodeAddedType,
	user.UserV1InitializedCheckSucceededType,
	user.UserV1PasswordCodeAddedType,
	user.UserV1PasswordCodeSentType,
	user.UserV1EmailVerifiedType,
}

// humanPasswordChangedEventTypes are the password changes.
// A change overwrites everything the earlier hash events set in [HumanPasswordWriteModel.Reduce]
// (the hash, whether a change is required, the code and the failed check count).
var humanPasswordChangedEventTypes = []eventstore.EventType{
	user.HumanPasswordChangedType,
	user.UserV1PasswordChangedType,
}

// humanPasswordHashEventTypes are the events which only set the password hash of a user (and the related state).
// They can be created on every login (e.g. an application setting the password on every login),
// so only the ones from the last password change on are loaded.
var humanPasswordHashEventTypes = append(slices.Clone(humanPasswordChangedEventTypes),
	user.HumanPasswordHashUpdatedType,
)

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
	// Only the password checks after the last reset of the failed count
	// and the password hash events from the last password change on are relevant.
	// This prevents loading the whole (constantly growing) login and password history of the user.
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(humanPasswordStateEventTypes...).
		Builder().
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(humanPasswordHashEventTypes...).
		// the last password change itself must be loaded
		SequenceGreater(max(wm.passwordChangedSequence, 1) - 1).
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
	// On a reload (e.g. the recheck after a password verification) only new events are loaded.
	// The events up to the last reset are already reduced or not relevant,
	// so the reload doesn't need to scan the checks between the last reduced event and the reset.
	if wm.ProcessedSequence != 0 {
		query.SequenceGreater(max(wm.ProcessedSequence, wm.checksAfterSequence))
	}
	return query
}

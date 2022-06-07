package projection

import (
	"context"
	"database/sql"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type userProjection struct {
	crdb.StatementHandler
}

const (
	UserTable        = "zitadel.projections.users"
	UserHumanTable   = UserTable + "_" + UserHumanSuffix
	UserMachineTable = UserTable + "_" + UserMachineSuffix
)

func newUserProjection(ctx context.Context, config crdb.StatementHandlerConfig) *userProjection {
	p := &userProjection{}
	config.ProjectionName = UserTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

const (
	UserIDCol            = "id"
	UserCreationDateCol  = "creation_date"
	UserChangeDateCol    = "change_date"
	UserResourceOwnerCol = "resource_owner"
	UserStateCol         = "state"
	UserSequenceCol      = "sequence"
	UserUsernameCol      = "username"
	UserTypeCol          = "type"
)

const (
	UserHumanSuffix = "humans"
	HumanUserIDCol  = "user_id"

	// profile
	HumanFirstNameCol         = "first_name"
	HumanLastNameCol          = "last_name"
	HumanNickNameCol          = "nick_name"
	HumanDisplayNameCol       = "display_name"
	HumanPreferredLanguageCol = "preferred_language"
	HumanGenderCol            = "gender"
	HumanAvatarURLCol         = "avatar_key"

	// email
	HumanEmailCol           = "email"
	HumanIsEmailVerifiedCol = "is_email_verified"

	// phone
	HumanPhoneCol           = "phone"
	HumanIsPhoneVerifiedCol = "is_phone_verified"
)

const (
	UserMachineSuffix = "machines"
	MachineUserIDCol  = "user_id"

	MachineNameCol        = "name"
	MachineDescriptionCol = "description"
)

func (p *userProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: user.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  user.UserV1AddedType,
					Reduce: p.reduceHumanAdded,
				},
				{
					Event:  user.HumanAddedType,
					Reduce: p.reduceHumanAdded,
				},
				{
					Event:  user.UserV1RegisteredType,
					Reduce: p.reduceHumanRegistered,
				},
				{
					Event:  user.HumanRegisteredType,
					Reduce: p.reduceHumanRegistered,
				},
				{
					Event:  user.HumanInitialCodeAddedType,
					Reduce: p.reduceHumanInitCodeAdded,
				},
				{
					Event:  user.UserV1InitialCodeAddedType,
					Reduce: p.reduceHumanInitCodeAdded,
				},
				{
					Event:  user.HumanInitializedCheckSucceededType,
					Reduce: p.reduceHumanInitCodeSucceeded,
				},
				{
					Event:  user.UserV1InitializedCheckSucceededType,
					Reduce: p.reduceHumanInitCodeSucceeded,
				},
				{
					Event:  user.UserLockedType,
					Reduce: p.reduceUserLocked,
				},
				{
					Event:  user.UserUnlockedType,
					Reduce: p.reduceUserUnlocked,
				},
				{
					Event:  user.UserDeactivatedType,
					Reduce: p.reduceUserDeactivated,
				},
				{
					Event:  user.UserReactivatedType,
					Reduce: p.reduceUserReactivated,
				},
				{
					Event:  user.UserRemovedType,
					Reduce: p.reduceUserRemoved,
				},
				{
					Event:  user.UserUserNameChangedType,
					Reduce: p.reduceUserNameChanged,
				},
				{
					Event:  user.HumanProfileChangedType,
					Reduce: p.reduceHumanProfileChanged,
				},
				{
					Event:  user.UserV1ProfileChangedType,
					Reduce: p.reduceHumanProfileChanged,
				},
				{
					Event:  user.HumanPhoneChangedType,
					Reduce: p.reduceHumanPhoneChanged,
				},
				{
					Event:  user.UserV1PhoneChangedType,
					Reduce: p.reduceHumanPhoneChanged,
				},
				{
					Event:  user.HumanPhoneRemovedType,
					Reduce: p.reduceHumanPhoneRemoved,
				},
				{
					Event:  user.UserV1PhoneRemovedType,
					Reduce: p.reduceHumanPhoneRemoved,
				},
				{
					Event:  user.HumanPhoneVerifiedType,
					Reduce: p.reduceHumanPhoneVerified,
				},
				{
					Event:  user.UserV1PhoneVerifiedType,
					Reduce: p.reduceHumanPhoneVerified,
				},
				{
					Event:  user.HumanEmailChangedType,
					Reduce: p.reduceHumanEmailChanged,
				},
				{
					Event:  user.UserV1EmailChangedType,
					Reduce: p.reduceHumanEmailChanged,
				},
				{
					Event:  user.HumanEmailVerifiedType,
					Reduce: p.reduceHumanEmailVerified,
				},
				{
					Event:  user.UserV1EmailVerifiedType,
					Reduce: p.reduceHumanEmailVerified,
				},
				{
					Event:  user.HumanAvatarAddedType,
					Reduce: p.reduceHumanAvatarAdded,
				},
				{
					Event:  user.HumanAvatarRemovedType,
					Reduce: p.reduceHumanAvatarRemoved,
				},
				{
					Event:  user.MachineAddedEventType,
					Reduce: p.reduceMachineAdded,
				},
				{
					Event:  user.MachineChangedEventType,
					Reduce: p.reduceMachineChanged,
				},
			},
		},
	}
}

func (p *userProjection) reduceHumanAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-Cw9BX", "seq", event.Sequence(), "expectedType", user.HumanAddedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Ebynp", "reduce.wrong.event.type")
	}
	return crdb.NewMultiStatement(
		e,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(UserIDCol, e.Aggregate().ID),
				handler.NewCol(UserCreationDateCol, e.CreationDate()),
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserResourceOwnerCol, e.Aggregate().ResourceOwner),
				handler.NewCol(UserStateCol, domain.UserStateActive),
				handler.NewCol(UserSequenceCol, e.Sequence()),
				handler.NewCol(UserUsernameCol, e.UserName),
				handler.NewCol(UserTypeCol, domain.UserTypeHuman),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(HumanUserIDCol, e.Aggregate().ID),
				handler.NewCol(HumanFirstNameCol, e.FirstName),
				handler.NewCol(HumanLastNameCol, e.LastName),
				handler.NewCol(HumanNickNameCol, &sql.NullString{String: e.NickName, Valid: e.NickName != ""}),
				handler.NewCol(HumanDisplayNameCol, &sql.NullString{String: e.DisplayName, Valid: e.DisplayName != ""}),
				handler.NewCol(HumanPreferredLanguageCol, &sql.NullString{String: e.PreferredLanguage.String(), Valid: !e.PreferredLanguage.IsRoot()}),
				handler.NewCol(HumanGenderCol, &sql.NullInt16{Int16: int16(e.Gender), Valid: e.Gender.Specified()}),
				handler.NewCol(HumanEmailCol, e.EmailAddress),
				handler.NewCol(HumanPhoneCol, &sql.NullString{String: e.PhoneNumber, Valid: e.PhoneNumber != ""}),
			},
			crdb.WithTableSuffix(UserHumanSuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanRegistered(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanRegisteredEvent)
	if !ok {
		logging.LogWithFields("HANDL-qoZyC", "seq", event.Sequence(), "expectedType", user.HumanRegisteredType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-xE53M", "reduce.wrong.event.type")
	}
	return crdb.NewMultiStatement(
		e,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(UserIDCol, e.Aggregate().ID),
				handler.NewCol(UserCreationDateCol, e.CreationDate()),
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserResourceOwnerCol, e.Aggregate().ResourceOwner),
				handler.NewCol(UserStateCol, domain.UserStateActive),
				handler.NewCol(UserSequenceCol, e.Sequence()),
				handler.NewCol(UserUsernameCol, e.UserName),
				handler.NewCol(UserTypeCol, domain.UserTypeHuman),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(HumanUserIDCol, e.Aggregate().ID),
				handler.NewCol(HumanFirstNameCol, e.FirstName),
				handler.NewCol(HumanLastNameCol, e.LastName),
				handler.NewCol(HumanNickNameCol, &sql.NullString{String: e.NickName, Valid: e.NickName != ""}),
				handler.NewCol(HumanDisplayNameCol, &sql.NullString{String: e.DisplayName, Valid: e.DisplayName != ""}),
				handler.NewCol(HumanPreferredLanguageCol, &sql.NullString{String: e.PreferredLanguage.String(), Valid: !e.PreferredLanguage.IsRoot()}),
				handler.NewCol(HumanGenderCol, &sql.NullInt16{Int16: int16(e.Gender), Valid: e.Gender.Specified()}),
				handler.NewCol(HumanEmailCol, e.EmailAddress),
				handler.NewCol(HumanPhoneCol, &sql.NullString{String: e.PhoneNumber, Valid: e.PhoneNumber != ""}),
			},
			crdb.WithTableSuffix(UserHumanSuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanInitCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanInitialCodeAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-DSfe2", "seq", event.Sequence(), "expectedType", user.HumanInitialCodeAddedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Dvgws", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(UserStateCol, domain.UserStateInitial),
		},
		[]handler.Condition{
			handler.NewCond(UserIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *userProjection) reduceHumanInitCodeSucceeded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanInitializedCheckSucceededEvent)
	if !ok {
		logging.LogWithFields("HANDL-Dgff2", "seq", event.Sequence(), "expectedType", user.HumanInitializedCheckSucceededType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Dfvwq", "reduce.wrong.event.type")
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(UserStateCol, domain.UserStateActive),
		},
		[]handler.Condition{
			handler.NewCond(UserIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *userProjection) reduceUserLocked(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserLockedEvent)
	if !ok {
		logging.LogWithFields("HANDL-c6irw", "seq", event.Sequence(), "expectedType", user.UserLockedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-exyBF", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(UserChangeDateCol, e.CreationDate()),
			handler.NewCol(UserStateCol, domain.UserStateLocked),
			handler.NewCol(UserSequenceCol, e.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(UserIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *userProjection) reduceUserUnlocked(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserUnlockedEvent)
	if !ok {
		logging.LogWithFields("HANDL-eyHv5", "seq", event.Sequence(), "expectedType", user.UserUnlockedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-JIyRl", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(UserChangeDateCol, e.CreationDate()),
			handler.NewCol(UserStateCol, domain.UserStateActive),
			handler.NewCol(UserSequenceCol, e.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(UserIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *userProjection) reduceUserDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserDeactivatedEvent)
	if !ok {
		logging.LogWithFields("HANDL-EqbaJ", "seq", event.Sequence(), "expectedType", user.UserDeactivatedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-6BNjj", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(UserChangeDateCol, e.CreationDate()),
			handler.NewCol(UserStateCol, domain.UserStateInactive),
			handler.NewCol(UserSequenceCol, e.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(UserIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *userProjection) reduceUserReactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserReactivatedEvent)
	if !ok {
		logging.LogWithFields("HANDL-kAaBr", "seq", event.Sequence(), "expectedType", user.UserReactivatedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-IoF6j", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(UserChangeDateCol, e.CreationDate()),
			handler.NewCol(UserStateCol, domain.UserStateActive),
			handler.NewCol(UserSequenceCol, e.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(UserIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *userProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-n2JMe", "seq", event.Sequence(), "expectedType", user.UserRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-BQB2t", "reduce.wrong.event.type")
	}

	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(UserIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *userProjection) reduceUserNameChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UsernameChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-7J5xL", "seq", event.Sequence(), "expectedType", user.UserUserNameChangedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-QNKyV", "reduce.wrong.event.type")
	}

	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(UserChangeDateCol, e.CreationDate()),
			handler.NewCol(UserUsernameCol, e.UserName),
			handler.NewCol(UserSequenceCol, e.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(UserIDCol, e.Aggregate().ID),
		},
	), nil
}

func (p *userProjection) reduceHumanProfileChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanProfileChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-Dwfyn", "seq", event.Sequence(), "expectedType", user.HumanProfileChangedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-769v4", "reduce.wrong.event.type")
	}
	cols := make([]handler.Column, 0, 6)
	if e.FirstName != "" {
		cols = append(cols, handler.NewCol(HumanFirstNameCol, e.FirstName))
	}

	if e.LastName != "" {
		cols = append(cols, handler.NewCol(HumanLastNameCol, e.LastName))
	}

	if e.NickName != nil {
		cols = append(cols, handler.NewCol(HumanNickNameCol, *e.NickName))
	}

	if e.DisplayName != nil {
		cols = append(cols, handler.NewCol(HumanDisplayNameCol, *e.DisplayName))
	}

	if e.PreferredLanguage != nil {
		cols = append(cols, handler.NewCol(HumanPreferredLanguageCol, e.PreferredLanguage.String()))
	}

	if e.Gender != nil {
		cols = append(cols, handler.NewCol(HumanGenderCol, *e.Gender))
	}

	return crdb.NewMultiStatement(
		e,
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserSequenceCol, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(UserIDCol, e.Aggregate().ID),
			},
		),
		crdb.AddUpdateStatement(
			cols,
			[]handler.Condition{
				handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
			},
			crdb.WithTableSuffix(UserHumanSuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanPhoneChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPhoneChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-pnRqf", "seq", event.Sequence(), "expectedType", user.HumanPhoneChangedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-xOGIA", "reduce.wrong.event.type")
	}

	return crdb.NewMultiStatement(
		e,
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserSequenceCol, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(UserIDCol, e.Aggregate().ID),
			},
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(HumanPhoneCol, e.PhoneNumber),
				handler.NewCol(HumanIsPhoneVerifiedCol, false),
			},
			[]handler.Condition{
				handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
			},
			crdb.WithTableSuffix(UserHumanSuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanPhoneRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPhoneRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-eMpOG", "seq", event.Sequence(), "expectedType", user.HumanPhoneRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-JI4S1", "reduce.wrong.event.type")
	}

	return crdb.NewMultiStatement(
		e,
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserSequenceCol, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(UserIDCol, e.Aggregate().ID),
			},
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(HumanPhoneCol, nil),
				handler.NewCol(HumanIsPhoneVerifiedCol, nil),
			},
			[]handler.Condition{
				handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
			},
			crdb.WithTableSuffix(UserHumanSuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanPhoneVerified(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPhoneVerifiedEvent)
	if !ok {
		logging.LogWithFields("HANDL-GhFOY", "seq", event.Sequence(), "expectedType", user.HumanPhoneVerifiedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-LBnqG", "reduce.wrong.event.type")
	}

	return crdb.NewMultiStatement(
		e,
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserSequenceCol, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(UserIDCol, e.Aggregate().ID),
			},
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(HumanIsPhoneVerifiedCol, true),
			},
			[]handler.Condition{
				handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
			},
			crdb.WithTableSuffix(UserHumanSuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanEmailChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanEmailChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-MDfHX", "seq", event.Sequence(), "expectedType", user.HumanEmailChangedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-KwiHa", "reduce.wrong.event.type")
	}

	return crdb.NewMultiStatement(
		e,
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserSequenceCol, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(UserIDCol, e.Aggregate().ID),
			},
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(HumanEmailCol, e.EmailAddress),
				handler.NewCol(HumanIsEmailVerifiedCol, false),
			},
			[]handler.Condition{
				handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
			},
			crdb.WithTableSuffix(UserHumanSuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanEmailVerified(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanEmailVerifiedEvent)
	if !ok {
		logging.LogWithFields("HANDL-FdN0b", "seq", event.Sequence(), "expectedType", user.HumanEmailVerifiedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-JzcDq", "reduce.wrong.event.type")
	}

	return crdb.NewMultiStatement(
		e,
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserSequenceCol, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(UserIDCol, e.Aggregate().ID),
			},
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(HumanIsEmailVerifiedCol, true),
			},
			[]handler.Condition{
				handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
			},
			crdb.WithTableSuffix(UserHumanSuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanAvatarAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanAvatarAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-naQue", "seq", event.Sequence(), "expectedType", user.HumanAvatarAddedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-eDEdt", "reduce.wrong.event.type")
	}

	return crdb.NewMultiStatement(
		e,
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserSequenceCol, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(UserIDCol, e.Aggregate().ID),
			},
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(HumanAvatarURLCol, e.StoreKey),
			},
			[]handler.Condition{
				handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
			},
			crdb.WithTableSuffix(UserHumanSuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanAvatarRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanAvatarRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-c6zoV", "seq", event.Sequence(), "expectedType", user.HumanAvatarRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-KhETX", "reduce.wrong.event.type")
	}

	return crdb.NewMultiStatement(
		e,
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserSequenceCol, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(UserIDCol, e.Aggregate().ID),
			},
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(HumanAvatarURLCol, nil),
			},
			[]handler.Condition{
				handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
			},
			crdb.WithTableSuffix(UserHumanSuffix),
		),
	), nil
}

func (p *userProjection) reduceMachineAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MachineAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-8xr78", "seq", event.Sequence(), "expectedType", user.MachineAddedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-q7ier", "reduce.wrong.event.type")
	}

	return crdb.NewMultiStatement(
		e,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(UserIDCol, e.Aggregate().ID),
				handler.NewCol(UserCreationDateCol, e.CreationDate()),
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserResourceOwnerCol, e.Aggregate().ResourceOwner),
				handler.NewCol(UserStateCol, domain.UserStateActive),
				handler.NewCol(UserSequenceCol, e.Sequence()),
				handler.NewCol(UserUsernameCol, e.UserName),
				handler.NewCol(UserTypeCol, domain.UserTypeMachine),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(MachineUserIDCol, e.Aggregate().ID),
				handler.NewCol(MachineNameCol, e.Name),
				handler.NewCol(MachineDescriptionCol, &sql.NullString{String: e.Description, Valid: e.Description != ""}),
			},
			crdb.WithTableSuffix(UserMachineSuffix),
		),
	), nil
}

func (p *userProjection) reduceMachineChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MachineChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-uUFCy", "seq", event.Sequence(), "expectedType", user.MachineChangedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-qYHvj", "reduce.wrong.event.type")
	}

	cols := make([]handler.Column, 0, 2)
	if e.Name != nil {
		cols = append(cols, handler.NewCol(MachineNameCol, *e.Name))
	}
	if e.Description != nil {
		cols = append(cols, handler.NewCol(MachineDescriptionCol, *e.Description))
	}
	if len(cols) == 0 {
		return crdb.NewNoOpStatement(e), nil
	}

	return crdb.NewMultiStatement(
		e,
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserSequenceCol, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(UserIDCol, e.Aggregate().ID),
			},
		),
		crdb.AddUpdateStatement(
			cols,
			[]handler.Condition{
				handler.NewCond(MachineUserIDCol, e.Aggregate().ID),
			},
			crdb.WithTableSuffix(UserMachineSuffix),
		),
	), nil
}

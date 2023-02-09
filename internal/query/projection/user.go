package projection

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type userProjection struct {
	crdb.StatementHandler
}

const (
	UserTable        = "projections.users8"
	UserHumanTable   = UserTable + "_" + UserHumanSuffix
	UserMachineTable = UserTable + "_" + UserMachineSuffix
	UserNotifyTable  = UserTable + "_" + UserNotifySuffix

	UserIDCol            = "id"
	UserCreationDateCol  = "creation_date"
	UserChangeDateCol    = "change_date"
	UserSequenceCol      = "sequence"
	UserStateCol         = "state"
	UserResourceOwnerCol = "resource_owner"
	UserInstanceIDCol    = "instance_id"
	UserUsernameCol      = "username"
	UserTypeCol          = "type"
	UserOwnerRemovedCol  = "owner_removed"

	UserHumanSuffix        = "humans"
	HumanUserIDCol         = "user_id"
	HumanUserInstanceIDCol = "instance_id"

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

	// machine
	UserMachineSuffix         = "machines"
	MachineUserIDCol          = "user_id"
	MachineUserInstanceIDCol  = "instance_id"
	MachineNameCol            = "name"
	MachineDescriptionCol     = "description"
	MachineHasSecretCol       = "has_secret"
	MachineAccessTokenTypeCol = "access_token_type"

	// notify
	UserNotifySuffix       = "notifications"
	NotifyUserIDCol        = "user_id"
	NotifyInstanceIDCol    = "instance_id"
	NotifyLastEmailCol     = "last_email"
	NotifyVerifiedEmailCol = "verified_email"
	NotifyLastPhoneCol     = "last_phone"
	NotifyVerifiedPhoneCol = "verified_phone"
	NotifyPasswordSetCol   = "password_set"
)

func newUserProjection(ctx context.Context, config crdb.StatementHandlerConfig) *userProjection {
	p := new(userProjection)
	config.ProjectionName = UserTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewMultiTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(UserIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(UserCreationDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(UserChangeDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(UserSequenceCol, crdb.ColumnTypeInt64),
			crdb.NewColumn(UserStateCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(UserResourceOwnerCol, crdb.ColumnTypeText),
			crdb.NewColumn(UserInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(UserUsernameCol, crdb.ColumnTypeText),
			crdb.NewColumn(UserTypeCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(UserOwnerRemovedCol, crdb.ColumnTypeBool, crdb.Default(false)),
		},
			crdb.NewPrimaryKey(UserInstanceIDCol, UserIDCol),
			crdb.WithIndex(crdb.NewIndex("username", []string{UserUsernameCol})),
			crdb.WithIndex(crdb.NewIndex("resource_owner", []string{UserResourceOwnerCol})),
			crdb.WithIndex(crdb.NewIndex("owner_removed", []string{UserOwnerRemovedCol})),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(HumanUserIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(HumanUserInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(HumanFirstNameCol, crdb.ColumnTypeText),
			crdb.NewColumn(HumanLastNameCol, crdb.ColumnTypeText),
			crdb.NewColumn(HumanNickNameCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(HumanDisplayNameCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(HumanPreferredLanguageCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(HumanGenderCol, crdb.ColumnTypeEnum, crdb.Nullable()),
			crdb.NewColumn(HumanAvatarURLCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(HumanEmailCol, crdb.ColumnTypeText),
			crdb.NewColumn(HumanIsEmailVerifiedCol, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(HumanPhoneCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(HumanIsPhoneVerifiedCol, crdb.ColumnTypeBool, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(HumanUserInstanceIDCol, HumanUserIDCol),
			UserHumanSuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys()),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(MachineUserIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(MachineUserInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(MachineNameCol, crdb.ColumnTypeText),
			crdb.NewColumn(MachineDescriptionCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(MachineHasSecretCol, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(MachineAccessTokenTypeCol, crdb.ColumnTypeEnum, crdb.Default(0)),
		},
			crdb.NewPrimaryKey(MachineUserInstanceIDCol, MachineUserIDCol),
			UserMachineSuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys()),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(NotifyUserIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(NotifyInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(NotifyLastEmailCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(NotifyVerifiedEmailCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(NotifyLastPhoneCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(NotifyVerifiedPhoneCol, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(NotifyPasswordSetCol, crdb.ColumnTypeBool, crdb.Default(false)),
		},
			crdb.NewPrimaryKey(NotifyInstanceIDCol, NotifyUserIDCol),
			UserNotifySuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys()),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

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
					Event:  user.UserDomainClaimedType,
					Reduce: p.reduceDomainClaimed,
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
				{
					Event:  user.HumanPasswordChangedType,
					Reduce: p.reduceHumanPasswordChanged,
				},
				{
					Event:  user.MachineSecretSetType,
					Reduce: p.reduceMachineSecretSet,
				},
				{
					Event:  user.MachineSecretRemovedType,
					Reduce: p.reduceMachineSecretRemoved,
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOwnerRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(UserInstanceIDCol),
				},
			},
		},
	}
}

func (p *userProjection) reduceHumanAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Ebynp", "reduce.wrong.event.type %s", user.HumanAddedType)
	}
	return crdb.NewMultiStatement(
		e,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(UserIDCol, e.Aggregate().ID),
				handler.NewCol(UserCreationDateCol, e.CreationDate()),
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserResourceOwnerCol, e.Aggregate().ResourceOwner),
				handler.NewCol(UserInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCol(UserStateCol, domain.UserStateActive),
				handler.NewCol(UserSequenceCol, e.Sequence()),
				handler.NewCol(UserUsernameCol, e.UserName),
				handler.NewCol(UserTypeCol, domain.UserTypeHuman),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(HumanUserIDCol, e.Aggregate().ID),
				handler.NewCol(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
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
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(NotifyUserIDCol, e.Aggregate().ID),
				handler.NewCol(NotifyInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCol(NotifyLastEmailCol, e.EmailAddress),
				handler.NewCol(NotifyLastPhoneCol, &sql.NullString{String: e.PhoneNumber, Valid: e.PhoneNumber != ""}),
				handler.NewCol(NotifyPasswordSetCol, e.Secret != nil),
			},
			crdb.WithTableSuffix(UserNotifySuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanRegistered(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanRegisteredEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-xE53M", "reduce.wrong.event.type %s", user.HumanRegisteredType)
	}
	return crdb.NewMultiStatement(
		e,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(UserIDCol, e.Aggregate().ID),
				handler.NewCol(UserCreationDateCol, e.CreationDate()),
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserResourceOwnerCol, e.Aggregate().ResourceOwner),
				handler.NewCol(UserInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCol(UserStateCol, domain.UserStateActive),
				handler.NewCol(UserSequenceCol, e.Sequence()),
				handler.NewCol(UserUsernameCol, e.UserName),
				handler.NewCol(UserTypeCol, domain.UserTypeHuman),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(HumanUserIDCol, e.Aggregate().ID),
				handler.NewCol(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
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
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(NotifyUserIDCol, e.Aggregate().ID),
				handler.NewCol(NotifyInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCol(NotifyLastEmailCol, e.EmailAddress),
				handler.NewCol(NotifyLastPhoneCol, &sql.NullString{String: e.PhoneNumber, Valid: e.PhoneNumber != ""}),
				handler.NewCol(NotifyPasswordSetCol, e.Secret != nil),
			},
			crdb.WithTableSuffix(UserNotifySuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanInitCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanInitialCodeAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Dvgws", "reduce.wrong.event.type %s", user.HumanInitialCodeAddedType)
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(UserStateCol, domain.UserStateInitial),
		},
		[]handler.Condition{
			handler.NewCond(UserIDCol, e.Aggregate().ID),
			handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *userProjection) reduceHumanInitCodeSucceeded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanInitializedCheckSucceededEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Dfvwq", "reduce.wrong.event.type %s", user.HumanInitializedCheckSucceededType)
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(UserStateCol, domain.UserStateActive),
		},
		[]handler.Condition{
			handler.NewCond(UserIDCol, e.Aggregate().ID),
			handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *userProjection) reduceUserLocked(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserLockedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-exyBF", "reduce.wrong.event.type %s", user.UserLockedType)
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
			handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *userProjection) reduceUserUnlocked(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserUnlockedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-JIyRl", "reduce.wrong.event.type %s", user.UserUnlockedType)
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
			handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *userProjection) reduceUserDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserDeactivatedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-6BNjj", "reduce.wrong.event.type %s", user.UserDeactivatedType)
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
			handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *userProjection) reduceUserReactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserReactivatedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-IoF6j", "reduce.wrong.event.type %s", user.UserReactivatedType)
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
			handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *userProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-BQB2t", "reduce.wrong.event.type %s", user.UserRemovedType)
	}

	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(UserIDCol, e.Aggregate().ID),
			handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *userProjection) reduceUserNameChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UsernameChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-QNKyV", "reduce.wrong.event.type %s", user.UserUserNameChangedType)
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
			handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *userProjection) reduceDomainClaimed(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.DomainClaimedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-ASwf3", "reduce.wrong.event.type %s", user.UserDomainClaimedType)
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
			handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *userProjection) reduceHumanProfileChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanProfileChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-769v4", "reduce.wrong.event.type %s", user.HumanProfileChangedType)
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
				handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
		crdb.AddUpdateStatement(
			cols,
			[]handler.Condition{
				handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
				handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
			},
			crdb.WithTableSuffix(UserHumanSuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanPhoneChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPhoneChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-xOGIA", "reduce.wrong.event.type %s", user.HumanPhoneChangedType)
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
				handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(HumanPhoneCol, e.PhoneNumber),
				handler.NewCol(HumanIsPhoneVerifiedCol, false),
			},
			[]handler.Condition{
				handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
				handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
			},
			crdb.WithTableSuffix(UserHumanSuffix),
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(NotifyLastPhoneCol, &sql.NullString{String: e.PhoneNumber, Valid: e.PhoneNumber != ""}),
			},
			[]handler.Condition{
				handler.NewCond(NotifyUserIDCol, e.Aggregate().ID),
				handler.NewCond(NotifyInstanceIDCol, e.Aggregate().InstanceID),
			},
			crdb.WithTableSuffix(UserNotifySuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanPhoneRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPhoneRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-JI4S1", "reduce.wrong.event.type %s", user.HumanPhoneRemovedType)
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
				handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(HumanPhoneCol, nil),
				handler.NewCol(HumanIsPhoneVerifiedCol, nil),
			},
			[]handler.Condition{
				handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
				handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
			},
			crdb.WithTableSuffix(UserHumanSuffix),
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(NotifyLastPhoneCol, nil),
				handler.NewCol(NotifyVerifiedPhoneCol, nil),
			},
			[]handler.Condition{
				handler.NewCond(NotifyUserIDCol, e.Aggregate().ID),
				handler.NewCond(NotifyInstanceIDCol, e.Aggregate().InstanceID),
			},
			crdb.WithTableSuffix(UserNotifySuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanPhoneVerified(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPhoneVerifiedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-LBnqG", "reduce.wrong.event.type %s", user.HumanPhoneVerifiedType)
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
				handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(HumanIsPhoneVerifiedCol, true),
			},
			[]handler.Condition{
				handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
				handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
			},
			crdb.WithTableSuffix(UserHumanSuffix),
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				crdb.NewCopyCol(NotifyVerifiedPhoneCol, NotifyLastPhoneCol),
			},
			[]handler.Condition{
				handler.NewCond(NotifyUserIDCol, e.Aggregate().ID),
				handler.NewCond(NotifyInstanceIDCol, e.Aggregate().InstanceID),
			},
			crdb.WithTableSuffix(UserNotifySuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanEmailChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanEmailChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-KwiHa", "reduce.wrong.event.type %s", user.HumanEmailChangedType)
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
				handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(HumanEmailCol, e.EmailAddress),
				handler.NewCol(HumanIsEmailVerifiedCol, false),
			},
			[]handler.Condition{
				handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
				handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
			},
			crdb.WithTableSuffix(UserHumanSuffix),
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(NotifyLastEmailCol, &sql.NullString{String: e.EmailAddress, Valid: e.EmailAddress != ""}),
			},
			[]handler.Condition{
				handler.NewCond(NotifyUserIDCol, e.Aggregate().ID),
				handler.NewCond(NotifyInstanceIDCol, e.Aggregate().InstanceID),
			},
			crdb.WithTableSuffix(UserNotifySuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanEmailVerified(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanEmailVerifiedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-JzcDq", "reduce.wrong.event.type %s", user.HumanEmailVerifiedType)
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
				handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(HumanIsEmailVerifiedCol, true),
			},
			[]handler.Condition{
				handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
				handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
			},
			crdb.WithTableSuffix(UserHumanSuffix),
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				crdb.NewCopyCol(NotifyVerifiedEmailCol, NotifyLastEmailCol),
			},
			[]handler.Condition{
				handler.NewCond(NotifyUserIDCol, e.Aggregate().ID),
				handler.NewCond(NotifyInstanceIDCol, e.Aggregate().InstanceID),
			},
			crdb.WithTableSuffix(UserNotifySuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanAvatarAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanAvatarAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-eDEdt", "reduce.wrong.event.type %s", user.HumanAvatarAddedType)
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
				handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(HumanAvatarURLCol, e.StoreKey),
			},
			[]handler.Condition{
				handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
				handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
			},
			crdb.WithTableSuffix(UserHumanSuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanAvatarRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanAvatarRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-KhETX", "reduce.wrong.event.type %s", user.HumanAvatarRemovedType)
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
				handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(HumanAvatarURLCol, nil),
			},
			[]handler.Condition{
				handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
				handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
			},
			crdb.WithTableSuffix(UserHumanSuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanPasswordChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPasswordChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-jqXUY", "reduce.wrong.event.type %s", user.HumanPasswordChangedType)
	}

	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(NotifyPasswordSetCol, true),
		},
		[]handler.Condition{
			handler.NewCond(NotifyUserIDCol, e.Aggregate().ID),
			handler.NewCond(NotifyInstanceIDCol, e.Aggregate().InstanceID),
		},
		crdb.WithTableSuffix(UserNotifySuffix),
	), nil
}

func (p *userProjection) reduceMachineSecretSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MachineSecretSetEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-x0p1n1i", "reduce.wrong.event.type %s", user.MachineSecretSetType)
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
				handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(MachineHasSecretCol, true),
			},
			[]handler.Condition{
				handler.NewCond(MachineUserIDCol, e.Aggregate().ID),
				handler.NewCond(MachineUserInstanceIDCol, e.Aggregate().InstanceID),
			},
			crdb.WithTableSuffix(UserMachineSuffix),
		),
	), nil
}

func (p *userProjection) reduceMachineSecretRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MachineSecretRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-x0p6n1i", "reduce.wrong.event.type %s", user.MachineSecretRemovedType)
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
				handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(MachineHasSecretCol, false),
			},
			[]handler.Condition{
				handler.NewCond(MachineUserIDCol, e.Aggregate().ID),
				handler.NewCond(MachineUserInstanceIDCol, e.Aggregate().InstanceID),
			},
			crdb.WithTableSuffix(UserMachineSuffix),
		),
	), nil
}

func (p *userProjection) reduceMachineAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MachineAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-q7ier", "reduce.wrong.event.type %s", user.MachineAddedEventType)
	}

	return crdb.NewMultiStatement(
		e,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(UserIDCol, e.Aggregate().ID),
				handler.NewCol(UserCreationDateCol, e.CreationDate()),
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserResourceOwnerCol, e.Aggregate().ResourceOwner),
				handler.NewCol(UserInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCol(UserStateCol, domain.UserStateActive),
				handler.NewCol(UserSequenceCol, e.Sequence()),
				handler.NewCol(UserUsernameCol, e.UserName),
				handler.NewCol(UserTypeCol, domain.UserTypeMachine),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(MachineUserIDCol, e.Aggregate().ID),
				handler.NewCol(MachineUserInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCol(MachineNameCol, e.Name),
				handler.NewCol(MachineDescriptionCol, &sql.NullString{String: e.Description, Valid: e.Description != ""}),
				handler.NewCol(MachineAccessTokenTypeCol, e.AccessTokenType),
			},
			crdb.WithTableSuffix(UserMachineSuffix),
		),
	), nil
}

func (p *userProjection) reduceMachineChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MachineChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-qYHvj", "reduce.wrong.event.type %s", user.MachineChangedEventType)
	}

	cols := make([]handler.Column, 0, 2)
	if e.Name != nil {
		cols = append(cols, handler.NewCol(MachineNameCol, *e.Name))
	}
	if e.Description != nil {
		cols = append(cols, handler.NewCol(MachineDescriptionCol, *e.Description))
	}
	if e.AccessTokenType != nil {
		cols = append(cols, handler.NewCol(MachineAccessTokenTypeCol, e.AccessTokenType))
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
				handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
		crdb.AddUpdateStatement(
			cols,
			[]handler.Condition{
				handler.NewCond(MachineUserIDCol, e.Aggregate().ID),
				handler.NewCond(MachineUserInstanceIDCol, e.Aggregate().InstanceID),
			},
			crdb.WithTableSuffix(UserMachineSuffix),
		),
	), nil

}

func (p *userProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-NCsdV", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(UserChangeDateCol, e.CreationDate()),
			handler.NewCol(UserSequenceCol, e.Sequence()),
			handler.NewCol(UserOwnerRemovedCol, true),
		},
		[]handler.Condition{
			handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(UserResourceOwnerCol, e.Aggregate().ID),
		},
	), nil
}

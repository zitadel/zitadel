package projection

import (
	"context"
	"database/sql"
	errs "errors"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	v3 "github.com/zitadel/zitadel/internal/eventstore/handler/v3"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type userProjection struct {
	crdb.StatementHandler
}

const (
	UserTable        = "projections.users5"
	UserHumanTable   = UserTable + "_" + UserHumanSuffix
	UserMachineTable = UserTable + "_" + UserMachineSuffix
	UserNotifyTable  = UserTable + "_" + UserNotifySuffix

	UserIDCol            = "id"
	UserCreationDateCol  = "creation_date"
	UserChangeDateCol    = "change_date"
	UserStateCol         = "state"
	UserResourceOwnerCol = "resource_owner"
	UserInstanceIDCol    = "instance_id"
	UserUsernameCol      = "username"
	UserTypeCol          = "type"

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
	UserMachineSuffix        = "machines"
	MachineUserIDCol         = "user_id"
	MachineUserInstanceIDCol = "instance_id"
	MachineNameCol           = "name"
	MachineDescriptionCol    = "description"

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

func newUserProjection(ctx context.Context, config v3.Config) *v3.IDProjection {
	p := new(userProjection)
	config.Check = crdb.NewMultiTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(UserIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(UserCreationDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(UserChangeDateCol, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(UserStateCol, crdb.ColumnTypeEnum),
			crdb.NewColumn(UserResourceOwnerCol, crdb.ColumnTypeText),
			crdb.NewColumn(UserInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(UserUsernameCol, crdb.ColumnTypeText),
			crdb.NewColumn(UserTypeCol, crdb.ColumnTypeEnum),
		},
			crdb.NewPrimaryKey(UserIDCol, UserInstanceIDCol),
			crdb.WithIndex(crdb.NewIndex("username_idx5", []string{UserUsernameCol})),
			crdb.WithIndex(crdb.NewIndex("user_ro_idx5", []string{UserResourceOwnerCol})),
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
			crdb.NewPrimaryKey(HumanUserIDCol, HumanUserInstanceIDCol),
			UserHumanSuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys("fk_human_ref_user5")),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(MachineUserIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(MachineUserInstanceIDCol, crdb.ColumnTypeText),
			crdb.NewColumn(MachineNameCol, crdb.ColumnTypeText),
			crdb.NewColumn(MachineDescriptionCol, crdb.ColumnTypeText, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(MachineUserIDCol, MachineUserInstanceIDCol),
			UserMachineSuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys("fk_machine_ref_user5")),
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
			crdb.NewPrimaryKey(NotifyUserIDCol, NotifyInstanceIDCol),
			UserNotifySuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys("fk_notify_ref_user5")),
		),
	)

	config.Reduces = map[eventstore.AggregateType][]v3.Reducer{
		user.AggregateType: {
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
				Event:          user.HumanInitialCodeAddedType,
				Reduce:         p.reduceHumanInitCodeAdded,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.UserV1InitialCodeAddedType,
				Reduce:         p.reduceHumanInitCodeAdded,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.HumanInitializedCheckSucceededType,
				Reduce:         p.reduceHumanInitCodeSucceeded,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.UserV1InitializedCheckSucceededType,
				Reduce:         p.reduceHumanInitCodeSucceeded,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.UserLockedType,
				Reduce:         p.reduceUserLocked,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.UserUnlockedType,
				Reduce:         p.reduceUserUnlocked,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.UserDeactivatedType,
				Reduce:         p.reduceUserDeactivated,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.UserReactivatedType,
				Reduce:         p.reduceUserReactivated,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.UserRemovedType,
				Reduce:         p.reduceUserRemoved,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.UserUserNameChangedType,
				Reduce:         p.reduceUserNameChanged,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.UserDomainClaimedType,
				Reduce:         p.reduceDomainClaimed,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.HumanProfileChangedType,
				Reduce:         p.reduceHumanProfileChanged,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.UserV1ProfileChangedType,
				Reduce:         p.reduceHumanProfileChanged,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.HumanPhoneChangedType,
				Reduce:         p.reduceHumanPhoneChanged,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.UserV1PhoneChangedType,
				Reduce:         p.reduceHumanPhoneChanged,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.HumanPhoneRemovedType,
				Reduce:         p.reduceHumanPhoneRemoved,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.UserV1PhoneRemovedType,
				Reduce:         p.reduceHumanPhoneRemoved,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.HumanPhoneVerifiedType,
				Reduce:         p.reduceHumanPhoneVerified,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.UserV1PhoneVerifiedType,
				Reduce:         p.reduceHumanPhoneVerified,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.HumanEmailChangedType,
				Reduce:         p.reduceHumanEmailChanged,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.UserV1EmailChangedType,
				Reduce:         p.reduceHumanEmailChanged,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.HumanEmailVerifiedType,
				Reduce:         p.reduceHumanEmailVerified,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.UserV1EmailVerifiedType,
				Reduce:         p.reduceHumanEmailVerified,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.HumanAvatarAddedType,
				Reduce:         p.reduceHumanAvatarAdded,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.HumanAvatarRemovedType,
				Reduce:         p.reduceHumanAvatarRemoved,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.MachineAddedEventType,
				Reduce:         p.reduceMachineAdded,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.MachineChangedEventType,
				Reduce:         p.reduceMachineChanged,
				PreviousEvents: p.previousEvents,
			},
			{
				Event:          user.HumanPasswordChangedType,
				Reduce:         p.reduceHumanPasswordChanged,
				PreviousEvents: p.previousEvents,
			},
		},
		instance.AggregateType: {
			{
				Event:  instance.InstanceRemovedEventType,
				Reduce: reduceInstanceRemovedHelper(UserInstanceIDCol),
			},
		},
	}

	return v3.New(UserTable, config)
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
	if len(cols) == 0 {
		return crdb.NewNoOpStatement(e), nil
	}

	return crdb.NewMultiStatement(
		e,
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
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

func (p *userProjection) previousEvents(tx *sql.Tx, event eventstore.Event) (*eventstore.SearchQueryBuilder, error) {
	row := tx.QueryRow("SELECT "+UserChangeDateCol+" FROM "+UserTable+" WHERE "+UserIDCol+" = $1 AND "+UserInstanceIDCol+" = $2 FOR UPDATE", event.Aggregate().ID, event.Aggregate().InstanceID)

	var changeDate time.Time

	if err := row.Scan(&changeDate); err != nil && !errs.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		// SetTx(tx).
		InstanceID(event.Aggregate().InstanceID).
		SystemTime(event.CreationDate()).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(event.Aggregate().ID).
		EventTypes(
			user.UserV1AddedType,
			user.HumanAddedType,
			user.UserV1RegisteredType,
			user.HumanRegisteredType,
			user.HumanInitialCodeAddedType,
			user.UserV1InitialCodeAddedType,
			user.HumanInitializedCheckSucceededType,
			user.UserV1InitializedCheckSucceededType,
			user.UserLockedType,
			user.UserUnlockedType,
			user.UserDeactivatedType,
			user.UserReactivatedType,
			user.UserRemovedType,
			user.UserUserNameChangedType,
			user.UserDomainClaimedType,
			user.HumanProfileChangedType,
			user.UserV1ProfileChangedType,
			user.HumanPhoneChangedType,
			user.UserV1PhoneChangedType,
			user.HumanPhoneRemovedType,
			user.UserV1PhoneRemovedType,
			user.HumanPhoneVerifiedType,
			user.UserV1PhoneVerifiedType,
			user.HumanEmailChangedType,
			user.UserV1EmailChangedType,
			user.HumanEmailVerifiedType,
			user.UserV1EmailVerifiedType,
			user.HumanAvatarAddedType,
			user.HumanAvatarRemovedType,
			user.MachineAddedEventType,
			user.MachineChangedEventType,
			user.HumanPasswordChangedType,
		).
		CreationDateAfter(changeDate).
		Builder(), nil
}

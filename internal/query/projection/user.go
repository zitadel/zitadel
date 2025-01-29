package projection

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	UserTable        = "projections.users14"
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

	UserHumanSuffix             = "humans"
	HumanUserIDCol              = "user_id"
	HumanUserInstanceIDCol      = "instance_id"
	HumanPasswordChangeRequired = "password_change_required"
	HumanPasswordChanged        = "password_changed"
	HumanMFAInitSkipped         = "mfa_init_skipped"

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
	MachineSecretCol          = "secret"
	MachineAccessTokenTypeCol = "access_token_type"

	// notify
	UserNotifySuffix            = "notifications"
	NotifyUserIDCol             = "user_id"
	NotifyInstanceIDCol         = "instance_id"
	NotifyLastEmailCol          = "last_email"
	NotifyVerifiedEmailCol      = "verified_email"
	NotifyVerifiedEmailLowerCol = "verified_email_lower"
	NotifyLastPhoneCol          = "last_phone"
	NotifyVerifiedPhoneCol      = "verified_phone"
	NotifyPasswordSetCol        = "password_set"
)

type userProjection struct{}

func newUserProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(userProjection))
}

func (*userProjection) Name() string {
	return UserTable
}

func (*userProjection) Init() *old_handler.Check {
	return handler.NewMultiTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(UserIDCol, handler.ColumnTypeText),
			handler.NewColumn(UserCreationDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(UserChangeDateCol, handler.ColumnTypeTimestamp),
			handler.NewColumn(UserSequenceCol, handler.ColumnTypeInt64),
			handler.NewColumn(UserStateCol, handler.ColumnTypeEnum),
			handler.NewColumn(UserResourceOwnerCol, handler.ColumnTypeText),
			handler.NewColumn(UserInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(UserUsernameCol, handler.ColumnTypeText),
			handler.NewColumn(UserTypeCol, handler.ColumnTypeEnum),
		},
			handler.NewPrimaryKey(UserInstanceIDCol, UserIDCol),
			handler.WithIndex(handler.NewIndex("username", []string{UserUsernameCol})),
			handler.WithIndex(handler.NewIndex("resource_owner", []string{UserResourceOwnerCol})),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(HumanUserIDCol, handler.ColumnTypeText),
			handler.NewColumn(HumanUserInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(HumanFirstNameCol, handler.ColumnTypeText),
			handler.NewColumn(HumanLastNameCol, handler.ColumnTypeText),
			handler.NewColumn(HumanNickNameCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(HumanDisplayNameCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(HumanPreferredLanguageCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(HumanGenderCol, handler.ColumnTypeEnum, handler.Nullable()),
			handler.NewColumn(HumanAvatarURLCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(HumanEmailCol, handler.ColumnTypeText),
			handler.NewColumn(HumanIsEmailVerifiedCol, handler.ColumnTypeBool, handler.Default(false)),
			handler.NewColumn(HumanPhoneCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(HumanIsPhoneVerifiedCol, handler.ColumnTypeBool, handler.Nullable()),
			handler.NewColumn(HumanPasswordChangeRequired, handler.ColumnTypeBool),
			handler.NewColumn(HumanPasswordChanged, handler.ColumnTypeTimestamp, handler.Nullable()),
			handler.NewColumn(HumanMFAInitSkipped, handler.ColumnTypeTimestamp, handler.Nullable()),
		},
			handler.NewPrimaryKey(HumanUserInstanceIDCol, HumanUserIDCol),
			UserHumanSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(MachineUserIDCol, handler.ColumnTypeText),
			handler.NewColumn(MachineUserInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(MachineNameCol, handler.ColumnTypeText),
			handler.NewColumn(MachineDescriptionCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(MachineSecretCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(MachineAccessTokenTypeCol, handler.ColumnTypeEnum, handler.Default(0)),
		},
			handler.NewPrimaryKey(MachineUserInstanceIDCol, MachineUserIDCol),
			UserMachineSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(NotifyUserIDCol, handler.ColumnTypeText),
			handler.NewColumn(NotifyInstanceIDCol, handler.ColumnTypeText),
			handler.NewColumn(NotifyLastEmailCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(NotifyVerifiedEmailCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(NotifyLastPhoneCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(NotifyVerifiedPhoneCol, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(NotifyPasswordSetCol, handler.ColumnTypeBool, handler.Default(false)),
		},
			handler.NewPrimaryKey(NotifyInstanceIDCol, NotifyUserIDCol),
			UserNotifySuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
		),
	)
}

func (p *userProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: user.AggregateType,
			EventReducers: []handler.EventReducer{
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
					Event:  user.MachineSecretHashUpdatedType,
					Reduce: p.reduceMachineSecretHashUpdated,
				},
				{
					Event:  user.MachineSecretRemovedType,
					Reduce: p.reduceMachineSecretRemoved,
				},
				{
					Event:  user.UserV1MFAOTPVerifiedType,
					Reduce: p.reduceUnsetMFAInitSkipped,
				},
				{
					Event:  user.HumanMFAOTPVerifiedType,
					Reduce: p.reduceUnsetMFAInitSkipped,
				},
				{
					Event:  user.HumanOTPSMSAddedType,
					Reduce: p.reduceUnsetMFAInitSkipped,
				},
				{
					Event:  user.HumanOTPEmailAddedType,
					Reduce: p.reduceUnsetMFAInitSkipped,
				},
				{
					Event:  user.HumanU2FTokenVerifiedType,
					Reduce: p.reduceUnsetMFAInitSkipped,
				},
				{
					Event:  user.HumanPasswordlessTokenVerifiedType,
					Reduce: p.reduceUnsetMFAInitSkipped,
				},
				{
					Event:  user.UserV1MFAInitSkippedType,
					Reduce: p.reduceMFAInitSkipped,
				},
				{
					Event:  user.HumanMFAInitSkippedType,
					Reduce: p.reduceMFAInitSkipped,
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOwnerRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Ebynp", "reduce.wrong.event.type %s", user.HumanAddedType)
	}
	passwordSet := crypto.SecretOrEncodedHash(e.Secret, e.EncodedHash) != ""
	return handler.NewMultiStatement(
		e,
		handler.AddCreateStatement(
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
		handler.AddCreateStatement(
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
				handler.NewCol(HumanPhoneCol, &sql.NullString{String: string(e.PhoneNumber), Valid: e.PhoneNumber != ""}),
				handler.NewCol(HumanPasswordChangeRequired, e.ChangeRequired),
				handler.NewCol(HumanPasswordChanged, &sql.NullTime{Time: e.CreatedAt(), Valid: passwordSet}),
			},
			handler.WithTableSuffix(UserHumanSuffix),
		),
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(NotifyUserIDCol, e.Aggregate().ID),
				handler.NewCol(NotifyInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCol(NotifyLastEmailCol, e.EmailAddress),
				handler.NewCol(NotifyLastPhoneCol, &sql.NullString{String: string(e.PhoneNumber), Valid: e.PhoneNumber != ""}),
				handler.NewCol(NotifyPasswordSetCol, passwordSet),
			},
			handler.WithTableSuffix(UserNotifySuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanRegistered(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanRegisteredEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-xE53M", "reduce.wrong.event.type %s", user.HumanRegisteredType)
	}
	passwordSet := crypto.SecretOrEncodedHash(e.Secret, e.EncodedHash) != ""
	return handler.NewMultiStatement(
		e,
		handler.AddCreateStatement(
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
		handler.AddCreateStatement(
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
				handler.NewCol(HumanPhoneCol, &sql.NullString{String: string(e.PhoneNumber), Valid: e.PhoneNumber != ""}),
				handler.NewCol(HumanPasswordChangeRequired, e.ChangeRequired),
				handler.NewCol(HumanPasswordChanged, &sql.NullTime{Time: e.CreatedAt(), Valid: passwordSet}),
			},
			handler.WithTableSuffix(UserHumanSuffix),
		),
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(NotifyUserIDCol, e.Aggregate().ID),
				handler.NewCol(NotifyInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCol(NotifyLastEmailCol, e.EmailAddress),
				handler.NewCol(NotifyLastPhoneCol, &sql.NullString{String: string(e.PhoneNumber), Valid: e.PhoneNumber != ""}),
				handler.NewCol(NotifyPasswordSetCol, passwordSet),
			},
			handler.WithTableSuffix(UserNotifySuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanInitCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanInitialCodeAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Dvgws", "reduce.wrong.event.type %s", user.HumanInitialCodeAddedType)
	}
	return handler.NewUpdateStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Dfvwq", "reduce.wrong.event.type %s", user.HumanInitializedCheckSucceededType)
	}
	return handler.NewUpdateStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-exyBF", "reduce.wrong.event.type %s", user.UserLockedType)
	}

	return handler.NewUpdateStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-JIyRl", "reduce.wrong.event.type %s", user.UserUnlockedType)
	}

	return handler.NewUpdateStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-6BNjj", "reduce.wrong.event.type %s", user.UserDeactivatedType)
	}

	return handler.NewUpdateStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-IoF6j", "reduce.wrong.event.type %s", user.UserReactivatedType)
	}

	return handler.NewUpdateStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-BQB2t", "reduce.wrong.event.type %s", user.UserRemovedType)
	}

	return handler.NewDeleteStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-QNKyV", "reduce.wrong.event.type %s", user.UserUserNameChangedType)
	}

	return handler.NewUpdateStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-ASwf3", "reduce.wrong.event.type %s", user.UserDomainClaimedType)
	}

	return handler.NewUpdateStatement(
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
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-769v4", "reduce.wrong.event.type %s", user.HumanProfileChangedType)
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

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserSequenceCol, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(UserIDCol, e.Aggregate().ID),
				handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
		handler.AddUpdateStatement(
			cols,
			[]handler.Condition{
				handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
				handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(UserHumanSuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanPhoneChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPhoneChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-xOGIA", "reduce.wrong.event.type %s", user.HumanPhoneChangedType)
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserSequenceCol, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(UserIDCol, e.Aggregate().ID),
				handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(HumanPhoneCol, e.PhoneNumber),
				handler.NewCol(HumanIsPhoneVerifiedCol, false),
			},
			[]handler.Condition{
				handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
				handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(UserHumanSuffix),
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(NotifyLastPhoneCol, &sql.NullString{String: string(e.PhoneNumber), Valid: e.PhoneNumber != ""}),
			},
			[]handler.Condition{
				handler.NewCond(NotifyUserIDCol, e.Aggregate().ID),
				handler.NewCond(NotifyInstanceIDCol, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(UserNotifySuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanPhoneRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPhoneRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-JI4S1", "reduce.wrong.event.type %s", user.HumanPhoneRemovedType)
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserSequenceCol, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(UserIDCol, e.Aggregate().ID),
				handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(HumanPhoneCol, nil),
				handler.NewCol(HumanIsPhoneVerifiedCol, nil),
			},
			[]handler.Condition{
				handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
				handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(UserHumanSuffix),
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(NotifyLastPhoneCol, nil),
				handler.NewCol(NotifyVerifiedPhoneCol, nil),
			},
			[]handler.Condition{
				handler.NewCond(NotifyUserIDCol, e.Aggregate().ID),
				handler.NewCond(NotifyInstanceIDCol, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(UserNotifySuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanPhoneVerified(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPhoneVerifiedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-LBnqG", "reduce.wrong.event.type %s", user.HumanPhoneVerifiedType)
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserSequenceCol, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(UserIDCol, e.Aggregate().ID),
				handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(HumanIsPhoneVerifiedCol, true),
			},
			[]handler.Condition{
				handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
				handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(UserHumanSuffix),
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCopyCol(NotifyVerifiedPhoneCol, NotifyLastPhoneCol),
			},
			[]handler.Condition{
				handler.NewCond(NotifyUserIDCol, e.Aggregate().ID),
				handler.NewCond(NotifyInstanceIDCol, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(UserNotifySuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanEmailChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanEmailChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-KwiHa", "reduce.wrong.event.type %s", user.HumanEmailChangedType)
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserSequenceCol, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(UserIDCol, e.Aggregate().ID),
				handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(HumanEmailCol, e.EmailAddress),
				handler.NewCol(HumanIsEmailVerifiedCol, false),
			},
			[]handler.Condition{
				handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
				handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(UserHumanSuffix),
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(NotifyLastEmailCol, &sql.NullString{String: string(e.EmailAddress), Valid: e.EmailAddress != ""}),
			},
			[]handler.Condition{
				handler.NewCond(NotifyUserIDCol, e.Aggregate().ID),
				handler.NewCond(NotifyInstanceIDCol, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(UserNotifySuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanEmailVerified(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanEmailVerifiedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-JzcDq", "reduce.wrong.event.type %s", user.HumanEmailVerifiedType)
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserSequenceCol, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(UserIDCol, e.Aggregate().ID),
				handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(HumanIsEmailVerifiedCol, true),
			},
			[]handler.Condition{
				handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
				handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(UserHumanSuffix),
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCopyCol(NotifyVerifiedEmailCol, NotifyLastEmailCol),
			},
			[]handler.Condition{
				handler.NewCond(NotifyUserIDCol, e.Aggregate().ID),
				handler.NewCond(NotifyInstanceIDCol, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(UserNotifySuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanAvatarAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanAvatarAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-eDEdt", "reduce.wrong.event.type %s", user.HumanAvatarAddedType)
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserSequenceCol, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(UserIDCol, e.Aggregate().ID),
				handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(HumanAvatarURLCol, e.StoreKey),
			},
			[]handler.Condition{
				handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
				handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(UserHumanSuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanAvatarRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanAvatarRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-KhETX", "reduce.wrong.event.type %s", user.HumanAvatarRemovedType)
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserSequenceCol, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(UserIDCol, e.Aggregate().ID),
				handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(HumanAvatarURLCol, nil),
			},
			[]handler.Condition{
				handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
				handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(UserHumanSuffix),
		),
	), nil
}

func (p *userProjection) reduceHumanPasswordChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPasswordChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-jqXUY", "reduce.wrong.event.type %s", user.HumanPasswordChangedType)
	}
	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(HumanPasswordChangeRequired, e.ChangeRequired),
				handler.NewCol(HumanPasswordChanged, &sql.NullTime{Time: e.CreatedAt(), Valid: true}),
			},
			[]handler.Condition{
				handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
				handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(UserHumanSuffix),
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(NotifyPasswordSetCol, true),
			},
			[]handler.Condition{
				handler.NewCond(NotifyUserIDCol, e.Aggregate().ID),
				handler.NewCond(NotifyInstanceIDCol, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(UserNotifySuffix),
		),
	), nil
}

func (p *userProjection) reduceMachineSecretSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MachineSecretSetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-x0p1n1i", "reduce.wrong.event.type %s", user.MachineSecretSetType)
	}
	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserSequenceCol, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(UserIDCol, e.Aggregate().ID),
				handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(MachineSecretCol, crypto.SecretOrEncodedHash(e.ClientSecret, e.HashedSecret)),
			},
			[]handler.Condition{
				handler.NewCond(MachineUserIDCol, e.Aggregate().ID),
				handler.NewCond(MachineUserInstanceIDCol, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(UserMachineSuffix),
		),
	), nil
}

func (p *userProjection) reduceMachineSecretHashUpdated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MachineSecretHashUpdatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Wieng4u", "reduce.wrong.event.type %s", user.MachineSecretHashUpdatedType)
	}
	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserSequenceCol, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(UserIDCol, e.Aggregate().ID),
				handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(MachineSecretCol, e.HashedSecret),
			},
			[]handler.Condition{
				handler.NewCond(MachineUserIDCol, e.Aggregate().ID),
				handler.NewCond(MachineUserInstanceIDCol, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(UserMachineSuffix),
		),
	), nil
}

func (p *userProjection) reduceMachineSecretRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MachineSecretRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-x0p6n1i", "reduce.wrong.event.type %s", user.MachineSecretRemovedType)
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserSequenceCol, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(UserIDCol, e.Aggregate().ID),
				handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(MachineSecretCol, nil),
			},
			[]handler.Condition{
				handler.NewCond(MachineUserIDCol, e.Aggregate().ID),
				handler.NewCond(MachineUserInstanceIDCol, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(UserMachineSuffix),
		),
	), nil
}

func (p *userProjection) reduceMachineAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MachineAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-q7ier", "reduce.wrong.event.type %s", user.MachineAddedEventType)
	}

	return handler.NewMultiStatement(
		e,
		handler.AddCreateStatement(
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
		handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(MachineUserIDCol, e.Aggregate().ID),
				handler.NewCol(MachineUserInstanceIDCol, e.Aggregate().InstanceID),
				handler.NewCol(MachineNameCol, e.Name),
				handler.NewCol(MachineDescriptionCol, &sql.NullString{String: e.Description, Valid: e.Description != ""}),
				handler.NewCol(MachineAccessTokenTypeCol, e.AccessTokenType),
			},
			handler.WithTableSuffix(UserMachineSuffix),
		),
	), nil
}

func (p *userProjection) reduceMachineChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MachineChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qYHvj", "reduce.wrong.event.type %s", user.MachineChangedEventType)
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
		return handler.NewNoOpStatement(e), nil
	}

	return handler.NewMultiStatement(
		e,
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(UserChangeDateCol, e.CreationDate()),
				handler.NewCol(UserSequenceCol, e.Sequence()),
			},
			[]handler.Condition{
				handler.NewCond(UserIDCol, e.Aggregate().ID),
				handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
			},
		),
		handler.AddUpdateStatement(
			cols,
			[]handler.Condition{
				handler.NewCond(MachineUserIDCol, e.Aggregate().ID),
				handler.NewCond(MachineUserInstanceIDCol, e.Aggregate().InstanceID),
			},
			handler.WithTableSuffix(UserMachineSuffix),
		),
	), nil

}

func (p *userProjection) reduceUnsetMFAInitSkipped(e eventstore.Event) (*handler.Statement, error) {
	switch e.(type) {
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-ojrf6", "reduce.wrong.event.type %s", e.Type())
	case *user.HumanOTPVerifiedEvent,
		*user.HumanOTPSMSAddedEvent,
		*user.HumanOTPEmailAddedEvent,
		*user.HumanU2FVerifiedEvent,
		*user.HumanPasswordlessVerifiedEvent:
	}

	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(HumanMFAInitSkipped, sql.NullTime{}),
		},
		[]handler.Condition{
			handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
			handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
		},
		handler.WithTableSuffix(UserHumanSuffix),
	), nil
}

func (p *userProjection) reduceMFAInitSkipped(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanMFAInitSkippedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qYHvj", "reduce.wrong.event.type %s", user.MachineChangedEventType)
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(HumanMFAInitSkipped, sql.NullTime{
				Time:  e.CreatedAt(),
				Valid: true,
			}),
		},
		[]handler.Condition{
			handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
			handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
		},
		handler.WithTableSuffix(UserHumanSuffix),
	), nil
}

func (p *userProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-NCsdV", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
			handler.NewCond(UserResourceOwnerCol, e.Aggregate().ID),
		},
	), nil
}

package projection

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/backend/v3/domain"
	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	UserRelationProjectionTable = "zitadel.users"
)

type userRelationalProjection struct{}

func (*userRelationalProjection) Name() string {
	return UserRelationProjectionTable
}

func newUserRelationalProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(userRelationalProjection))
}

func (p *userRelationalProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: user.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  user.UserV1AddedType,
					Reduce: p.reduceHumanAdded,
				},
				// 		{
				// 			Event:  user.HumanAddedType,
				// 			Reduce: p.reduceHumanAdded,
				// 		},
				// 		{
				// 			Event:  user.UserV1RegisteredType,
				// 			Reduce: p.reduceHumanRegistered,
				// 		},
				// 		{
				// 			Event:  user.HumanRegisteredType,
				// 			Reduce: p.reduceHumanRegistered,
				// 		},
				// 		{
				// 			Event:  user.HumanInitialCodeAddedType,
				// 			Reduce: p.reduceHumanInitCodeAdded,
				// 		},
				// 		{
				// 			Event:  user.UserV1InitialCodeAddedType,
				// 			Reduce: p.reduceHumanInitCodeAdded,
				// 		},
				// 		{
				// 			Event:  user.HumanInitializedCheckSucceededType,
				// 			Reduce: p.reduceHumanInitCodeSucceeded,
				// 		},
				// 		{
				// 			Event:  user.UserV1InitializedCheckSucceededType,
				// 			Reduce: p.reduceHumanInitCodeSucceeded,
				// 		},
				// 		{
				// 			Event:  user.UserLockedType,
				// 			Reduce: p.reduceUserLocked,
				// 		},
				// 		{
				// 			Event:  user.UserUnlockedType,
				// 			Reduce: p.reduceUserUnlocked,
				// 		},
				// 		{
				// 			Event:  user.UserDeactivatedType,
				// 			Reduce: p.reduceUserDeactivated,
				// 		},
				// 		{
				// 			Event:  user.UserReactivatedType,
				// 			Reduce: p.reduceUserReactivated,
				// 		},
				// 		{
				// 			Event:  user.UserRemovedType,
				// 			Reduce: p.reduceUserRemoved,
				// 		},
				// 		{
				// 			Event:  user.UserUserNameChangedType,
				// 			Reduce: p.reduceUserNameChanged,
				// 		},
				// 		{
				// 			Event:  user.UserDomainClaimedType,
				// 			Reduce: p.reduceDomainClaimed,
				// 		},
				// 		{
				// 			Event:  user.HumanProfileChangedType,
				// 			Reduce: p.reduceHumanProfileChanged,
				// 		},
				// 		{
				// 			Event:  user.UserV1ProfileChangedType,
				// 			Reduce: p.reduceHumanProfileChanged,
				// 		},
				// 		{
				// 			Event:  user.HumanPhoneChangedType,
				// 			Reduce: p.reduceHumanPhoneChanged,
				// 		},
				// 		{
				// 			Event:  user.UserV1PhoneChangedType,
				// 			Reduce: p.reduceHumanPhoneChanged,
				// 		},
				// 		{
				// 			Event:  user.HumanPhoneRemovedType,
				// 			Reduce: p.reduceHumanPhoneRemoved,
				// 		},
				// 		{
				// 			Event:  user.UserV1PhoneRemovedType,
				// 			Reduce: p.reduceHumanPhoneRemoved,
				// 		},
				// 		{
				// 			Event:  user.HumanPhoneVerifiedType,
				// 			Reduce: p.reduceHumanPhoneVerified,
				// 		},
				// 		{
				// 			Event:  user.UserV1PhoneVerifiedType,
				// 			Reduce: p.reduceHumanPhoneVerified,
				// 		},
				// 		{
				// 			Event:  user.HumanEmailChangedType,
				// 			Reduce: p.reduceHumanEmailChanged,
				// 		},
				// 		{
				// 			Event:  user.UserV1EmailChangedType,
				// 			Reduce: p.reduceHumanEmailChanged,
				// 		},
				// 		{
				// 			Event:  user.HumanEmailVerifiedType,
				// 			Reduce: p.reduceHumanEmailVerified,
				// 		},
				// 		{
				// 			Event:  user.UserV1EmailVerifiedType,
				// 			Reduce: p.reduceHumanEmailVerified,
				// 		},
				// 		{
				// 			Event:  user.HumanAvatarAddedType,
				// 			Reduce: p.reduceHumanAvatarAdded,
				// 		},
				// 		{
				// 			Event:  user.HumanAvatarRemovedType,
				// 			Reduce: p.reduceHumanAvatarRemoved,
				// 		},
				// 		{
				// 			Event:  user.MachineAddedEventType,
				// 			Reduce: p.reduceMachineAdded,
				// 		},
				// 		{
				// 			Event:  user.MachineChangedEventType,
				// 			Reduce: p.reduceMachineChanged,
				// 		},
				// 		{
				// 			Event:  user.HumanPasswordChangedType,
				// 			Reduce: p.reduceHumanPasswordChanged,
				// 		},
				// 		{
				// 			Event:  user.MachineSecretSetType,
				// 			Reduce: p.reduceMachineSecretSet,
				// 		},
				// 		{
				// 			Event:  user.MachineSecretHashUpdatedType,
				// 			Reduce: p.reduceMachineSecretHashUpdated,
				// 		},
				// 		{
				// 			Event:  user.MachineSecretRemovedType,
				// 			Reduce: p.reduceMachineSecretRemoved,
				// 		},
				// 		{
				// 			Event:  user.UserV1MFAOTPVerifiedType,
				// 			Reduce: p.reduceUnsetMFAInitSkipped,
				// 		},
				// 		{
				// 			Event:  user.HumanMFAOTPVerifiedType,
				// 			Reduce: p.reduceUnsetMFAInitSkipped,
				// 		},
				// 		{
				// 			Event:  user.HumanOTPSMSAddedType,
				// 			Reduce: p.reduceUnsetMFAInitSkipped,
				// 		},
				// 		{
				// 			Event:  user.HumanOTPEmailAddedType,
				// 			Reduce: p.reduceUnsetMFAInitSkipped,
				// 		},
				// 		{
				// 			Event:  user.HumanU2FTokenVerifiedType,
				// 			Reduce: p.reduceUnsetMFAInitSkipped,
				// 		},
				// 		{
				// 			Event:  user.HumanPasswordlessTokenVerifiedType,
				// 			Reduce: p.reduceUnsetMFAInitSkipped,
				// 		},
				// 		{
				// 			Event:  user.UserV1MFAInitSkippedType,
				// 			Reduce: p.reduceMFAInitSkipped,
				// 		},
				// 		{
				// 			Event:  user.HumanMFAInitSkippedType,
				// 			Reduce: p.reduceMFAInitSkipped,
				// 		},
				// 	},
				// },
				// {
				// 	Aggregate: org.AggregateType,
				// 	EventReducers: []handler.EventReducer{
				// 		{
				// 			Event:  org.OrgRemovedEventType,
				// 			Reduce: p.reduceOwnerRemoved,
				// 		},
				// 	},
				// },
				// {
				// 	Aggregate: instance.AggregateType,
				// 	EventReducers: []handler.EventReducer{
				// 		{
				// 			Event:  instance.InstanceRemovedEventType,
				// 			Reduce: reduceInstanceRemovedHelper(UserInstanceIDCol),
				// 		},
			},
		},
	}
}

func (u *userRelationalProjection) reduceHumanAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-HbYn4", "reduce.wrong.event.type %s", user.HumanAddedType)
	}
	// passwordSet := crypto.SecretOrEncodedHash(e.Secret, e.EncodedHash) != ""
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		userRepo := repository.UserRepository()

		_, err := userRepo.CreateHuman(ctx, v3_sql.SQLTx(tx),
			&domain.Human{
				User: domain.User{
					ID:         e.Aggregate().ID,
					InstanceID: e.Aggregate().InstanceID,
					OrgID:      e.Aggregate().ResourceOwner,
					Username:   e.UserName,
					State:      domain.UserStateActive,
					CreatedAt:  e.CreationDate(),
					UpdatedAt:  e.CreationDate(),
				},
				FirstName:         e.FirstName,
				LastName:          e.LastName,
				NickName:          e.NickName,
				DisplayName:       e.DisplayName,
				PreferredLanguage: e.PreferredLanguage.String(),
				Gender:            uint8(e.Gender),
			},
		)
		return err
	}), nil
	// return handler.NewMultiStatement(
	// 	e,
	// 	handler.AddCreateStatement(
	// 		[]handler.Column{
	// 			handler.NewCol(UserIDCol, e.Aggregate().ID),
	// 			handler.NewCol(UserCreationDateCol, e.CreationDate()),
	// 			handler.NewCol(UserChangeDateCol, e.CreationDate()),
	// 			handler.NewCol(UserResourceOwnerCol, e.Aggregate().ResourceOwner),
	// 			handler.NewCol(UserInstanceIDCol, e.Aggregate().InstanceID),
	// 			handler.NewCol(UserStateCol, domain.UserStateActive),
	// 			handler.NewCol(UserSequenceCol, e.Sequence()),
	// 			handler.NewCol(UserUsernameCol, e.UserName),
	// 			handler.NewCol(UserTypeCol, domain.UserTypeHuman),
	// 		},
	// 	),
	// 	handler.AddCreateStatement(
	// 		[]handler.Column{
	// 			handler.NewCol(HumanUserIDCol, e.Aggregate().ID),
	// 			handler.NewCol(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
	// 			handler.NewCol(HumanFirstNameCol, e.FirstName),
	// 			handler.NewCol(HumanLastNameCol, e.LastName),
	// 			handler.NewCol(HumanNickNameCol, &sql.NullString{String: e.NickName, Valid: e.NickName != ""}),
	// 			handler.NewCol(HumanDisplayNameCol, &sql.NullString{String: e.DisplayName, Valid: e.DisplayName != ""}),
	// 			handler.NewCol(HumanPreferredLanguageCol, &sql.NullString{String: e.PreferredLanguage.String(), Valid: !e.PreferredLanguage.IsRoot()}),
	// 			handler.NewCol(HumanGenderCol, &sql.NullInt16{Int16: int16(e.Gender), Valid: e.Gender.Specified()}),
	// 			handler.NewCol(HumanEmailCol, e.EmailAddress),
	// 			handler.NewCol(HumanPhoneCol, &sql.NullString{String: string(e.PhoneNumber), Valid: e.PhoneNumber != ""}),
	// 			handler.NewCol(HumanPasswordChangeRequired, e.ChangeRequired),
	// 			handler.NewCol(HumanPasswordChanged, &sql.NullTime{Time: e.CreatedAt(), Valid: passwordSet}),
	// 		},
	// 		handler.WithTableSuffix(UserHumanSuffix),
	// 	),
	// 	handler.AddCreateStatement(
	// 		[]handler.Column{
	// 			handler.NewCol(NotifyUserIDCol, e.Aggregate().ID),
	// 			handler.NewCol(NotifyInstanceIDCol, e.Aggregate().InstanceID),
	// 			handler.NewCol(NotifyLastEmailCol, e.EmailAddress),
	// 			handler.NewCol(NotifyLastPhoneCol, &sql.NullString{String: string(e.PhoneNumber), Valid: e.PhoneNumber != ""}),
	// 			handler.NewCol(NotifyPasswordSetCol, passwordSet),
	// 		},
	// 		handler.WithTableSuffix(UserNotifySuffix),
	// 	),
	// ), nil
}

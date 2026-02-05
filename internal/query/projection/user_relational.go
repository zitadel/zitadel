package projection

import (
	"context"
	"database/sql"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/crypto"
	old_domain "github.com/zitadel/zitadel/internal/domain"
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

				// {
				// 	Event:  user.HumanInitialCodeAddedType,
				// 	Reduce: p.reduceHumanInitCodeAdded,
				// },
				// {
				// 	Event:  user.UserV1InitialCodeAddedType,
				// 	Reduce: p.reduceHumanInitCodeAdded,
				// },
				// {
				// 	Event:  user.HumanInitializedCheckSucceededType,
				// 	Reduce: p.reduceHumanInitCodeSucceeded,
				// },
				// {
				// 	Event:  user.UserV1InitializedCheckSucceededType,
				// 	Reduce: p.reduceHumanInitCodeSucceeded,
				// },
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
					Reduce: p.reduceUsernameChanged,
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
					Event:  user.HumanEmailCodeAddedType,
					Reduce: p.reduceHumanEmailCodeAdded,
				},
				{
					Event:  user.UserV1EmailCodeAddedType,
					Reduce: p.reduceHumanEmailCodeAdded,
				},
				{
					Event:  user.HumanEmailVerificationFailedType,
					Reduce: p.reduceHumanEmailVerificationFailed,
				},
				{
					Event:  user.UserV1EmailVerificationFailedType,
					Reduce: p.reduceHumanEmailVerificationFailed,
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
					Event:  user.HumanPhoneCodeAddedType,
					Reduce: p.reduceHumanPhoneCodeAdded,
				},
				{
					Event:  user.UserV1PhoneCodeAddedType,
					Reduce: p.reduceHumanPhoneCodeAdded,
				},
				{
					Event:  user.HumanPhoneVerificationFailedType,
					Reduce: p.reduceHumanPhoneVerificationFailed,
				},
				{
					Event:  user.UserV1PhoneVerificationFailedType,
					Reduce: p.reduceHumanPhoneVerificationFailed,
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
					Event:  user.HumanPasswordChangedType,
					Reduce: p.reduceHumanPasswordChanged,
				},
				{
					Event:  user.HumanPasswordCodeAddedType,
					Reduce: p.reduceHumanPasswordCodeAdded,
				},
				{
					Event:  user.HumanPasswordCheckSucceededType,
					Reduce: p.reduceHumanPasswordCheckSucceeded,
				},
				{
					Event:  user.HumanPasswordCheckFailedType,
					Reduce: p.reduceHumanPasswordCheckFailed,
				},
				{
					Event:  user.HumanPasswordHashUpdatedType,
					Reduce: p.reduceHumanPasswordHashUpdated,
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
					Event:  user.MachineKeyAddedEventType,
					Reduce: p.reduceMachineKeyAdded,
				},
				{
					Event:  user.MachineKeyRemovedEventType,
					Reduce: p.reduceMachineKeyRemoved,
				},
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
				// {
				// 	Event:  user.HumanPasswordlessTokenVerifiedType,
				// 	Reduce: p.reduceUnsetMFAInitSkipped,
				// },

				// Pats only on machines

				{
					Event:  user.UserV1MFAInitSkippedType,
					Reduce: p.reduceMFAInitSkipped,
				},
				{
					Event:  user.HumanMFAInitSkippedType,
					Reduce: p.reduceMFAInitSkipped,
				},

				{
					Event: user.HumanRefreshTokenAddedType,
					// TODO: needed?
				},
				{
					Event: user.HumanRefreshTokenRenewedType,
					// TODO: needed?
				},
				{
					Event: user.HumanRefreshTokenRemovedType,
					// TODO: needed?
				},

				{
					Event:  user.PersonalAccessTokenAddedType,
					Reduce: p.reducePersonalAccessTokenAdded,
				},
				{
					Event:  user.PersonalAccessTokenRemovedType,
					Reduce: p.reducePersonalAccessTokenRemoved,
				},

				{
					Event:  user.MetadataSetType,
					Reduce: p.reduceMetadataSet,
				},
				{
					Event:  user.MetadataRemovedType,
					Reduce: p.reduceMetadataRemoved,
				},
				{
					Event:  user.MetadataRemovedAllType,
					Reduce: p.reduceMetadataRemovedAll,
				},

				{
					Event:  user.HumanPasswordlessTokenAddedType,
					Reduce: p.reducePasskeyAdded,
				},
				{
					Event:  user.HumanPasswordlessTokenVerifiedType,
					Reduce: p.reducePasskeyVerified,
				},
				{
					Event:  user.HumanPasswordlessTokenSignCountChangedType,
					Reduce: p.reducePasskeySignCountSet,
				},
				{
					Event:  user.HumanPasswordlessTokenRemovedType,
					Reduce: p.reducePasskeyRemoved,
				},

				{
					Event:  user.HumanU2FTokenAddedType,
					Reduce: p.reducePasskeyAdded,
				},
				{
					Event:  user.HumanU2FTokenVerifiedType,
					Reduce: p.reducePasskeyVerified,
				},
				{
					Event:  user.HumanU2FTokenSignCountChangedType,
					Reduce: p.reducePasskeySignCountSet,
				},
				{
					Event:  user.HumanU2FTokenRemovedType,
					Reduce: p.reducePasskeyRemoved,
				},

				{
					Event:  user.HumanPasswordlessInitCodeAddedType,
					Reduce: p.reducePasskeyInitCodeAdded,
				},
				{
					Event:  user.HumanPasswordlessInitCodeCheckFailedType,
					Reduce: p.reducePasskeyInitCodeCheckFailed,
				},
				{
					Event:  user.HumanPasswordlessInitCodeCheckSucceededType,
					Reduce: p.reducePasskeyInitCodeCheckSucceeded,
				},
				{
					Event:  user.HumanPasswordlessInitCodeRequestedType,
					Reduce: p.reducePasskeyInitCodeRequested,
				},

				{
					Event:  user.UserIDPLinkAddedType,
					Reduce: p.reduceIDPLinkAdded,
				},
				{
					Event:  user.UserIDPLinkCascadeRemovedType,
					Reduce: p.reduceIDPLinkCascadeRemoved,
				},
				{
					Event:  user.UserIDPLinkRemovedType,
					Reduce: p.reduceIDPLinkRemoved,
				},
				{
					Event:  user.UserIDPExternalIDMigratedType,
					Reduce: p.reduceIDPLinkUserIDMigrated,
				},
				{
					Event:  user.UserIDPExternalUsernameChangedType,
					Reduce: p.reduceIDPLinkUsernameChanged,
				},

				{
					Event:  user.HumanMFAOTPAddedType,
					Reduce: p.reduceTOTPAdded,
				},
				{
					Event:  user.HumanMFAOTPVerifiedType,
					Reduce: p.reduceTOTPVerified,
				},
				{
					Event:  user.HumanMFAOTPRemovedType,
					Reduce: p.reduceTOTPRemoved,
				},
				{
					Event:  user.HumanMFAOTPCheckSucceededType,
					Reduce: p.reduceTOTPCheckSucceeded,
				},
				{
					Event:  user.HumanMFAOTPCheckFailedType,
					Reduce: p.reduceTOTPCheckFailed,
				},
				{
					Event:  user.HumanOTPSMSAddedType,
					Reduce: p.reduceOTPSMSEnabled,
				},
				{
					Event:  user.HumanOTPSMSRemovedType,
					Reduce: p.reduceOTPSMSDisabled,
				},
				{
					Event:  user.HumanOTPSMSCodeAddedType,
					Reduce: p.reduceOTPSMSCodeAdded,
				},
				{
					Event:  user.HumanOTPSMSCheckSucceededType,
					Reduce: p.reduceOTPSMSCheckSucceeded,
				},
				{
					Event:  user.HumanOTPSMSCheckFailedType,
					Reduce: p.reduceOTPSMSCheckFailed,
				},
				{
					Event:  user.HumanOTPEmailAddedType,
					Reduce: p.reduceOTPEmailEnabled,
				},
				{
					Event:  user.HumanOTPEmailRemovedType,
					Reduce: p.reduceOTPEmailDisabled,
				},
				{
					Event:  user.HumanOTPEmailCodeAddedType,
					Reduce: p.reduceOTPEmailCodeAdded,
				},
				{
					Event:  user.HumanOTPEmailCheckSucceededType,
					Reduce: p.reduceOTPEmailCheckSucceeded,
				},
				{
					Event:  user.HumanOTPEmailCheckFailedType,
					Reduce: p.reduceOTPEmailCheckFailed,
				},
			},
		},
	}
}

func (u *userRelationalProjection) reduceHumanAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-HbYn4", "reduce.wrong.event.type %s", user.HumanAddedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		v3Tx := v3_sql.SQLTx(tx)

		userRepo := repository.UserRepository()
		humanRepo := userRepo.Human()
		condition := userRepo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID)

		err := userRepo.Create(ctx, v3Tx, &domain.User{
			InstanceID:     e.Aggregate().InstanceID,
			OrganizationID: e.Aggregate().ResourceOwner,
			ID:             e.Aggregate().ID,
			Username:       e.UserName,
			State:          domain.UserStateActive,
			CreatedAt:      e.CreatedAt(),
			UpdatedAt:      e.CreatedAt(),
			Human: &domain.HumanUser{
				FirstName:         e.FirstName,
				LastName:          e.LastName,
				Nickname:          e.NickName,
				DisplayName:       e.DisplayName,
				PreferredLanguage: e.PreferredLanguage,
				Gender:            mapHumanGender(e.Gender),
			},
		})
		if err != nil {
			return err
		}

		if password := crypto.SecretOrEncodedHash(e.Secret, e.EncodedHash); password != "" {
			_, err = userRepo.Update(ctx, v3Tx,
				condition,
				humanRepo.SetPassword(&domain.VerificationTypeSkipped{
					SkippedAt: e.CreatedAt(),
					Value:     &password,
				}),
			)
			if err != nil {
				return err
			}
		}
		if e.ChangeRequired {
			_, err = userRepo.Update(ctx, v3Tx,
				condition,
				humanRepo.SetPasswordChangeRequired(e.ChangeRequired),
				humanRepo.SetUpdatedAt(e.CreatedAt()),
			)
		}
		return err
		// &domain.Human{
		// 	FirstName:         e.FirstName,
		// 	LastName:          e.LastName,
		// 	Nickname:          e.NickName,
		// 	DisplayName:       e.DisplayName,
		// 	PreferredLanguage: e.PreferredLanguage.String(),
		// 	Gender:            uint8(e.Gender),
		// 	User: domain.User{
		// 		ID:         e.Aggregate().ID,
		// 		InstanceID: e.Aggregate().InstanceID,
		// 		OrgID:      e.Aggregate().ResourceOwner,
		// 		Username:   e.UserName,
		// 		// TODO check when to set username unique
		// 		// UsernameOrgUnique: false,
		// 		State:     domain.UserStateActive,
		// 		CreatedAt: e.CreationDate(),
		// 		UpdatedAt: e.CreationDate(),
		// 	},
		// 	HumanEmailContact: domain.HumanContact{
		// 		Type:       gu.Ptr(domain.ContactTypeEmail),
		// 		Value:      gu.Ptr(string(e.EmailAddress.Normalize())),
		// 		IsVerified: gu.Ptr(false),
		// 	},
		// 	HumanPhoneContact: func() *domain.HumanContact {
		// 		if e.PhoneNumber == "" {
		// 			return nil
		// 		}
		// 		return &domain.HumanContact{
		// 			Type:       gu.Ptr(domain.ContactTypePhone),
		// 			Value:      gu.Ptr(string(e.PhoneNumber)),
		// 			IsVerified: gu.Ptr(false),
		// 		}
		// 	}(),
		// 	HumanSecurity: domain.HumanSecurity{
		// 		PasswordChangeRequired: e.ChangeRequired,
		// 		PasswordChange: func() *time.Time {
		// 			if !passwordSet {
		// 				return nil
		// 			}
		// 			passwordChange := e.CreatedAt()
		// 			return &passwordChange
		// 		}(),
		// 	},
		// },
	}), nil
}

func (p *userRelationalProjection) reduceHumanRegistered(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanRegisteredEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-xhD3q", "reduce.wrong.event.type %s", user.HumanRegisteredType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		v3Tx := v3_sql.SQLTx(tx)

		userRepo := repository.UserRepository()
		humanRepo := userRepo.Human()
		condition := userRepo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID)

		err := userRepo.Create(ctx, v3Tx, &domain.User{
			InstanceID:     e.Aggregate().InstanceID,
			OrganizationID: e.Aggregate().ResourceOwner,
			ID:             e.Aggregate().ID,
			Username:       e.UserName,
			State:          domain.UserStateActive,
			CreatedAt:      e.CreatedAt(),
			UpdatedAt:      e.CreatedAt(),
			Human: &domain.HumanUser{
				FirstName:         e.FirstName,
				LastName:          e.LastName,
				Nickname:          e.NickName,
				DisplayName:       e.DisplayName,
				PreferredLanguage: e.PreferredLanguage,
				Gender:            mapHumanGender(e.Gender),
			},
		})
		if err != nil {
			return err
		}

		if password := crypto.SecretOrEncodedHash(e.Secret, e.EncodedHash); password != "" {
			_, err = userRepo.Update(ctx, v3Tx,
				condition,
				humanRepo.SetPassword(&domain.VerificationTypeSkipped{
					SkippedAt: e.CreatedAt(),
					Value:     &password,
				}),
			)
			if err != nil {
				return err
			}
		}
		if e.ChangeRequired {
			_, err = userRepo.Update(ctx, v3Tx,
				condition,
				humanRepo.SetPasswordChangeRequired(e.ChangeRequired),
				humanRepo.SetUpdatedAt(e.CreatedAt()),
			)
		}
		return err
		// _, err := userRepo.CreateHuman(ctx, v3_sql.SQLTx(tx),
		// 	&domain.Human{
		// 		FirstName:         e.FirstName,
		// 		LastName:          e.LastName,
		// 		NickName:          e.NickName,
		// 		DisplayName:       e.DisplayName,
		// 		PreferredLanguage: e.PreferredLanguage.String(),
		// 		Gender:            uint8(e.Gender),
		// 		User: domain.User{
		// 			ID:         e.Aggregate().ID,
		// 			InstanceID: e.Aggregate().InstanceID,
		// 			OrgID:      e.Aggregate().ResourceOwner,
		// 			Username:   e.UserName,
		// 			// TODO check when to set username unique
		// 			// UsernameOrgUnique: false,
		// 			State:     domain.UserStateActive,
		// 			CreatedAt: e.CreationDate(),
		// 			UpdatedAt: e.CreationDate(),
		// 		},
		// 		HumanEmailContact: domain.HumanContact{
		// 			Type:       gu.Ptr(domain.ContactTypeEmail),
		// 			Value:      gu.Ptr(string(e.EmailAddress.Normalize())),
		// 			IsVerified: gu.Ptr(false),
		// 		},
		// 		HumanPhoneContact: func() *domain.HumanContact {
		// 			if e.PhoneNumber == "" {
		// 				return nil
		// 			}
		// 			return &domain.HumanContact{
		// 				Type:       gu.Ptr(domain.ContactTypePhone),
		// 				Value:      gu.Ptr(string(e.PhoneNumber)),
		// 				IsVerified: gu.Ptr(false),
		// 			}
		// 		}(),
		// 		HumanSecurity: domain.HumanSecurity{
		// 			PasswordChangeRequired: e.ChangeRequired,
		// 			PasswordChange: func() *time.Time {
		// 				if !passwordSet {
		// 					return nil
		// 				}
		// 				passwordChange := e.CreatedAt()
		// 				return &passwordChange
		// 			}(),
		// 		},
		// 	},
		// )
	}), nil
}

func (p *userRelationalProjection) reduceUserLocked(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserLockedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-eUn8f", "reduce.wrong.event.type %s", user.UserLockedType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository()
		_, err := repo.Update(
			ctx,
			v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Agg.InstanceID, e.Aggregate().ID),
			repo.SetState(domain.UserStateLocked),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceUserUnlocked(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserUnlockedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-JIyRl", "reduce.wrong.event.type %s", user.UserUnlockedType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository()
		_, err := repo.Update(
			ctx,
			v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Agg.InstanceID, e.Aggregate().ID),
			repo.SetState(domain.UserStateActive),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceUserDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserDeactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-6BNjj", "reduce.wrong.event.type %s", user.UserDeactivatedType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository()
		_, err := repo.Update(
			ctx,
			v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Agg.InstanceID, e.Aggregate().ID),
			repo.SetState(domain.UserStateInactive),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceUserReactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserReactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-IoF6j", "reduce.wrong.event.type %s", user.UserReactivatedType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository()
		_, err := repo.Update(
			ctx,
			v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Agg.InstanceID, e.Aggregate().ID),
			repo.SetState(domain.UserStateActive),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-BQB2t", "reduce.wrong.event.type %s", user.UserRemovedType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository()
		_, err := repo.Delete(
			ctx,
			v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Agg.InstanceID, e.Aggregate().ID),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceUsernameChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UsernameChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-QNKyV", "reduce.wrong.event.type %s", user.UserUserNameChangedType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository()
		_, err := repo.Update(
			ctx,
			v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Agg.InstanceID, e.Aggregate().ID),
			repo.SetUsername(e.UserName),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceDomainClaimed(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.DomainClaimedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-ASwf3", "reduce.wrong.event.type %s", user.UserDomainClaimedType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		userRepo := repository.UserRepository()

		_, err := userRepo.Update(ctx, v3_sql.SQLTx(tx),
			userRepo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			userRepo.SetUsername(e.UserName),
			userRepo.SetUpdatedAt(e.CreatedAt()),
		)

		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanProfileChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanProfileChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-769v4", "reduce.wrong.event.type %s", user.HumanProfileChangedType)
	}

	repo := repository.UserRepository().Human()
	changes := make([]database.Change, 0, 7)

	if e.FirstName != "" {
		changes = append(changes, repo.SetFirstName(e.FirstName))
	}

	if e.LastName != "" {
		changes = append(changes, repo.SetLastName(e.LastName))
	}

	if e.NickName != nil {
		changes = append(changes, repo.SetNickname(*e.NickName))
	}

	if e.DisplayName != nil {
		changes = append(changes, repo.SetDisplayName(*e.DisplayName))
	}

	if e.PreferredLanguage != nil {
		changes = append(changes, repo.SetPreferredLanguage(*e.PreferredLanguage))
	}

	if e.Gender != nil {
		changes = append(changes, repo.SetGender(mapHumanGender(*e.Gender)))
	}
	changes = append(changes, repo.SetUpdatedAt(e.CreatedAt()))

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Agg.InstanceID, e.Aggregate().ID),
			changes...,
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanPhoneChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPhoneChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-xOGIA", "reduce.wrong.event.type %s", user.HumanPhoneChangedType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetPhone(&domain.VerificationTypeInit{
				Value:     (*string)(&e.PhoneNumber),
				CreatedAt: e.CreatedAt(),
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanPhoneRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPhoneRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-JI4S1", "reduce.wrong.event.type %s", user.HumanPhoneRemovedType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.RemovePhone(),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanPhoneVerified(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPhoneVerifiedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-LBnqG", "reduce.wrong.event.type %s", user.HumanPhoneVerifiedType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetPhone(&domain.VerificationTypeVerified{
				VerifiedAt: e.CreatedAt(),
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanPhoneCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPhoneCodeAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-JzcDq", "reduce.wrong.event.type %s", user.HumanPhoneCodeAddedType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository().Human()

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetPhone(&domain.VerificationTypeUpdate{
				Code:   e.Code,
				Expiry: &e.Expiry,
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanPhoneVerificationFailed(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPhoneVerificationFailedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-JzcDq", "reduce.wrong.event.type %s", user.HumanPhoneVerificationFailedType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository().Human()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetPhone(&domain.VerificationTypeFailed{
				FailedAt: e.CreatedAt(),
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanEmailChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanEmailChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-KwiHa", "reduce.wrong.event.type %s", user.HumanEmailChangedType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository().Human()

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetEmail(&domain.VerificationTypeInit{
				Value:     (*string)(&e.EmailAddress),
				CreatedAt: e.CreatedAt(),
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanEmailVerified(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanEmailVerifiedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-JzcDq", "reduce.wrong.event.type %s", user.HumanEmailVerifiedType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository().Human()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetEmail(&domain.VerificationTypeVerified{
				VerifiedAt: e.CreatedAt(),
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanEmailCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanEmailCodeAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-JzcDq", "reduce.wrong.event.type %s", user.HumanEmailCodeAddedType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository().Human()

		var expiry *time.Duration
		if e.Expiry > 0 {
			expiry = &e.Expiry
		}

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetEmail(&domain.VerificationTypeUpdate{
				Code:   e.Code,
				Expiry: expiry,
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanEmailVerificationFailed(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanEmailVerificationFailedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-JzcDq", "reduce.wrong.event.type %s", user.HumanEmailVerificationFailedType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository().Human()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetEmail(&domain.VerificationTypeFailed{
				FailedAt: e.CreatedAt(),
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanAvatarAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanAvatarAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-eDEdt", "reduce.wrong.event.type %s", user.HumanAvatarAddedType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository().Human()

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetAvatarKey(&e.StoreKey),
			repo.SetUpdatedAt(e.CreationDate()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanAvatarRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanAvatarRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-KhETX", "reduce.wrong.event.type %s", user.HumanAvatarRemovedType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository().Human()

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetAvatarKey(nil),
			repo.SetUpdatedAt(e.CreationDate()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanPasswordChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPasswordChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-jqXUY", "reduce.wrong.event.type %s", user.HumanPasswordChangedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		password := crypto.SecretOrEncodedHash(e.Secret, e.EncodedHash)

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetPasswordChangeRequired(e.ChangeRequired),
			repo.SetPassword(&domain.VerificationTypeInit{
				Value:     &password,
				CreatedAt: e.CreatedAt(),
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanPasswordCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPasswordCodeAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-nd7f3", "reduce.wrong.event.type %s", user.HumanPasswordCodeAddedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()

		var expiry *time.Duration
		if e.Expiry > 0 {
			expiry = &e.Expiry
		}

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetPassword(&domain.VerificationTypeUpdate{
				Code:   e.Code,
				Expiry: expiry,
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanPasswordCheckSucceeded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPasswordCheckSucceededEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-nd7f3", "reduce.wrong.event.type %s", user.HumanPasswordCheckSucceededType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetPassword(&domain.VerificationTypeVerified{
				VerifiedAt: e.CreatedAt(),
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanPasswordCheckFailed(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPasswordCheckFailedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-nd7f3", "reduce.wrong.event.type %s", user.HumanPasswordCheckFailedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetPassword(&domain.VerificationTypeFailed{
				FailedAt: e.CreatedAt(),
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanPasswordHashUpdated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPasswordHashUpdatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-nd7f3", "reduce.wrong.event.type %s", user.HumanPasswordHashUpdatedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetPassword(&domain.VerificationTypeSkipped{
				Value:     &e.EncodedHash,
				SkippedAt: e.CreatedAt(),
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceMachineSecretSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MachineSecretSetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-x0p1n1i", "reduce.wrong.event.type %s", user.MachineSecretSetType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.MachineUserRepository()

		secret := crypto.SecretOrEncodedHash(e.ClientSecret, e.HashedSecret)

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetSecret(&secret),
			repo.SetUpdatedAt(e.CreationDate()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceMachineSecretHashUpdated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MachineSecretHashUpdatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Wieng4u", "reduce.wrong.event.type %s", user.MachineSecretHashUpdatedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository().Machine()

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetSecret(&e.HashedSecret),
			repo.SetUpdatedAt(e.CreationDate()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceMachineSecretRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MachineSecretRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-x0p6n1i", "reduce.wrong.event.type %s", user.MachineSecretRemovedType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository().Machine()

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetSecret(nil),
			repo.SetUpdatedAt(e.CreationDate()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceMachineKeyAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MachineKeyAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-nd7f3", "reduce.wrong.event.type %s", user.MachineKeyAddedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.MachineUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.AddKey(&domain.MachineKey{
				ID:        e.KeyID,
				Type:      mapMachineKeyType(e.KeyType),
				PublicKey: e.PublicKey,
				CreatedAt: e.CreatedAt(),
				ExpiresAt: e.ExpirationDate,
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceMachineKeyRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MachineKeyRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-nd7f3", "reduce.wrong.event.type %s", user.MachineKeyRemovedEventType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.MachineUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.RemoveKey(e.KeyID),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceMachineAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MachineAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-q7ier", "reduce.wrong.event.type %s", user.MachineAddedEventType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository()
		return repo.Create(ctx, v3_sql.SQLTx(tx),
			&domain.User{
				ID:             e.Aggregate().ID,
				InstanceID:     e.Aggregate().InstanceID,
				OrganizationID: e.Aggregate().ResourceOwner,
				Username:       e.UserName,
				// TODO check when to set username unique
				// IsUsernameOrgUnique: ,
				State:     domain.UserStateActive,
				CreatedAt: e.CreatedAt(),
				UpdatedAt: e.CreatedAt(),
				Machine: &domain.MachineUser{
					Name:            e.Name,
					Description:     e.Description,
					AccessTokenType: mapMachineAccessTokenType(e.AccessTokenType),
				},
			},
		)
	}), nil
}

func (p *userRelationalProjection) reduceMachineChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MachineChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qYHvj", "reduce.wrong.event.type %s", user.MachineChangedEventType)
	}
	repo := repository.MachineUserRepository()

	changes := make([]database.Change, 0, 4)
	if e.Name != nil {
		changes = append(changes, repo.SetName(*e.Name))
	}
	if e.Description != nil {
		changes = append(changes, repo.SetDescription(*e.Description))
	}
	if e.AccessTokenType != nil {
		changes = append(changes, repo.SetAccessTokenType(mapMachineAccessTokenType(*e.AccessTokenType)))
	}
	changes = append(changes, repo.SetUpdatedAt(e.CreatedAt()))

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			database.And(
				repo.InstanceIDCondition(e.Aggregate().InstanceID),
				repo.IDCondition(e.Aggregate().ID),
			),
			changes...,
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reducePersonalAccessTokenAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.PersonalAccessTokenAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-nd7f3", "reduce.wrong.event.type %s", user.PersonalAccessTokenAddedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.MachineUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.AddPersonalAccessToken(&domain.PersonalAccessToken{
				ID:        e.TokenID,
				Scopes:    e.Scopes,
				ExpiresAt: e.Expiration,
				CreatedAt: e.CreatedAt(),
			}),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reducePersonalAccessTokenRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.PersonalAccessTokenRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-nd7f3", "reduce.wrong.event.type %s", user.PersonalAccessTokenRemovedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.MachineUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.RemovePersonalAccessToken(e.TokenID),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceMetadataSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MetadataSetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-xOO4e", "reduce.wrong.event.type %s", user.MetadataSetType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-xg4IJ", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.AddMetadata(&domain.Metadata{
				Key:       e.Key,
				Value:     e.Value,
				CreatedAt: e.CreationDate(),
				UpdatedAt: e.CreationDate(),
			}),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceMetadataRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MetadataRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-xOO4e", "reduce.wrong.event.type %s", user.MetadataRemovedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-xg4IJ", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.RemoveMetadata(repo.MetadataConditions().MetadataKeyCondition(database.TextOperationEqual, e.Key)),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceMetadataRemovedAll(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MetadataRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-xOO4e", "reduce.wrong.event.type %s", user.MetadataRemovedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-xg4IJ", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.RemoveMetadata(nil),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reducePasskeyAdded(event eventstore.Event) (*handler.Statement, error) {
	var (
		e   user.HumanWebAuthNAddedEvent
		typ domain.PasskeyType
	)
	switch typed := event.(type) {
	case *user.HumanPasswordlessAddedEvent:
		e = typed.HumanWebAuthNAddedEvent
		typ = domain.PasskeyTypePasswordless
	case *user.HumanU2FAddedEvent:
		e = typed.HumanWebAuthNAddedEvent
		typ = domain.PasskeyTypeU2F
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-nd7f3", "reduce.wrong.event.type for passkey %s", event.Type())
	}
	return handler.NewStatement(&e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			database.And(
				repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
				repo.PasskeyIDCondition(e.WebAuthNTokenID),
			),
			repo.AddPasskey(&domain.Passkey{
				ID:             e.WebAuthNTokenID,
				Challenge:      []byte(e.Challenge),
				RelyingPartyID: e.RPID,
				CreatedAt:      e.CreatedAt(),
				UpdatedAt:      e.CreatedAt(),
				Type:           typ,
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reducePasskeyVerified(event eventstore.Event) (*handler.Statement, error) {
	var (
		e user.HumanWebAuthNVerifiedEvent
	)
	switch typed := event.(type) {
	case *user.HumanPasswordlessVerifiedEvent:
		e = typed.HumanWebAuthNVerifiedEvent
	case *user.HumanU2FVerifiedEvent:
		e = typed.HumanWebAuthNVerifiedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-nd7f3", "reduce.wrong.event.type for passkey %s", event.Type())
	}
	return handler.NewStatement(&e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		// TODO: remaining fields to be updated?
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.UpdatePasskey(
				repo.PasskeyIDCondition(e.WebAuthNTokenID),
				repo.SetPasskeyKeyID(e.KeyID),
				repo.SetPasskeyPublicKey(e.PublicKey),
				repo.SetPasskeyAttestationType(e.AttestationType),
				repo.SetPasskeyAuthenticatorAttestationGUID(e.AAGUID),
				repo.SetPasskeySignCount(e.SignCount),
				repo.SetPasskeyName(e.WebAuthNTokenName),
				repo.SetPasskeyUpdatedAt(e.CreatedAt()),
			),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reducePasskeySignCountSet(event eventstore.Event) (*handler.Statement, error) {
	var (
		e user.HumanWebAuthNSignCountChangedEvent
	)
	switch typed := event.(type) {
	case *user.HumanPasswordlessSignCountChangedEvent:
		e = typed.HumanWebAuthNSignCountChangedEvent
	case *user.HumanU2FSignCountChangedEvent:
		e = typed.HumanWebAuthNSignCountChangedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-nd7f3", "reduce.wrong.event.type for passkey %s", event.Type())
	}
	return handler.NewStatement(&e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.UpdatePasskey(
				repo.PasskeyIDCondition(e.WebAuthNTokenID),
				repo.SetPasskeySignCount(e.SignCount),
			),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reducePasskeyRemoved(event eventstore.Event) (*handler.Statement, error) {
	var (
		e user.HumanWebAuthNRemovedEvent
	)
	switch typed := event.(type) {
	case *user.HumanPasswordlessRemovedEvent:
		e = typed.HumanWebAuthNRemovedEvent
	case *user.HumanU2FRemovedEvent:
		e = typed.HumanWebAuthNRemovedEvent
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-nd7f3", "reduce.wrong.event.type for passkey %s", event.Type())
	}
	return handler.NewStatement(&e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		userRepo := repository.UserRepository()
		humanRepo := userRepo.Human()
		_, err := userRepo.Update(ctx, v3_sql.SQLTx(tx),
			userRepo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.WebAuthNTokenID),
			humanRepo.RemovePasskey(humanRepo.PasskeyIDCondition(e.WebAuthNTokenID)),
			humanRepo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reducePasskeyInitCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPasswordlessInitCodeAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-nd7f3", "reduce.wrong.event.type %s", user.HumanPasswordlessInitCodeAddedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetVerification(&domain.VerificationTypeInit{
				ID:        &e.ID,
				Code:      e.Code,
				CreatedAt: e.CreatedAt(),
				Expiry:    &e.Expiry,
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reducePasskeyInitCodeCheckFailed(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPasswordlessInitCodeCheckFailedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-nd7f3", "reduce.wrong.event.type %s", user.HumanPasswordlessInitCodeCheckFailedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetVerification(&domain.VerificationTypeFailed{
				ID:       &e.ID,
				FailedAt: e.CreatedAt(),
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reducePasskeyInitCodeCheckSucceeded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPasswordlessInitCodeCheckSucceededEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-nd7f3", "reduce.wrong.event.type %s", user.HumanPasswordlessInitCodeCheckSucceededType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetVerification(&domain.VerificationTypeVerified{
				ID:         &e.ID,
				VerifiedAt: e.CreatedAt(),
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reducePasskeyInitCodeRequested(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanPasswordlessInitCodeRequestedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-nd7f3", "reduce.wrong.event.type %s", user.HumanPasswordlessInitCodeRequestedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetVerification(&domain.VerificationTypeInit{
				ID:        &e.ID,
				Code:      e.Code,
				CreatedAt: e.CreatedAt(),
				Expiry:    &e.Expiry,
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

// func (p *userRelationalProjection) reduceUnsetMFAInitSkipped(e eventstore.Event) (*handler.Statement, error) {
// 	switch e.(type) {
// 	default:
// 		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-ojrf6", "reduce.wrong.event.type %s", e.Type())
// 	case *user.HumanOTPVerifiedEvent,
// 		*user.HumanOTPSMSAddedEvent,
// 		*user.HumanOTPEmailAddedEvent,
// 		*user.HumanU2FVerifiedEvent,
// 		*user.HumanPasswordlessVerifiedEvent:
// 	}

// 	return handler.NewUpdateStatement(
// 		e,
// 		[]handler.Column{
// 			handler.NewCol(HumanMFAInitSkipped, sql.NullTime{}),
// 		},
// 		[]handler.Condition{
// 			handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
// 			handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
// 		},
// 		handler.WithTableSuffix(UserHumanSuffix),
// 	), nil
// }

func (p *userRelationalProjection) reduceMFAInitSkipped(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanMFAInitSkippedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qYHvj", "reduce.wrong.event.type %s", user.HumanMFAInitSkippedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SkipMultifactorInitializationAt(e.CreatedAt()),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceIDPLinkAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserIDPLinkAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qYHvj", "reduce.wrong.event.type %s", user.UserIDPLinkAddedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.AddIdentityProviderLink(&domain.IdentityProviderLink{
				ProviderID:       e.IDPConfigID,
				ProvidedUserID:   e.ExternalUserID,
				ProvidedUsername: e.DisplayName,
				CreatedAt:        e.CreatedAt(),
				UpdatedAt:        e.CreatedAt(),
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceIDPLinkCascadeRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserIDPLinkCascadeRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qYHvj", "reduce.wrong.event.type %s", user.UserIDPLinkCascadeRemovedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.RemoveIdentityProviderLink(e.IDPConfigID, e.ExternalUserID),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceIDPLinkRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserIDPLinkRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qYHvj", "reduce.wrong.event.type %s", user.UserIDPLinkRemovedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.RemoveIdentityProviderLink(e.IDPConfigID, e.ExternalUserID),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceIDPLinkUserIDMigrated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserIDPExternalIDMigratedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qYHvj", "reduce.wrong.event.type %s", user.UserIDPExternalIDMigratedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.UpdateIdentityProviderLink(
				nil, // TODO(adlerhurst): set correct condition
				repo.SetIdentityProviderLinkProvidedID(e.IDPConfigID, e.PreviousID, e.NewID),
			),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceIDPLinkUsernameChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserIDPExternalUsernameEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qYHvj", "reduce.wrong.event.type %s", user.UserIDPExternalUsernameChangedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.UpdateIdentityProviderLink(
				database.And(
					repo.LinkedIdentityProviderIDCondition(e.IDPConfigID),
					nil, // TODO(adlerhurst): add e.ExternalUserID as condition
				),
				repo.SetIdentityProviderLinkUsername(e.ExternalUsername),
			),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceTOTPAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanOTPAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qYHvj", "reduce.wrong.event.type %s", user.HumanMFAOTPAddedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			// repo.SetTOTP(&domain.VerificationTypeInit{
			// 	CreatedAt: e.CreatedAt(),
			// 	Code:      e.Secret,
			// 	Value:     gu.Ptr(string(e.Secret.Crypted)),
			// }),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceTOTPVerified(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanOTPVerifiedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qYHvj", "reduce.wrong.event.type %s", user.HumanMFAOTPVerifiedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			// repo.SetTOTP(&domain.VerificationTypeVerified{
			// 	VerifiedAt: e.CreatedAt(),
			// }),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceTOTPRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanOTPRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qYHvj", "reduce.wrong.event.type %s", user.HumanMFAOTPRemovedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.RemoveTOTP(),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceTOTPCheckSucceeded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanOTPCheckSucceededEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qYHvj", "reduce.wrong.event.type %s", user.HumanMFAOTPCheckSucceededType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			// repo.CheckTOTP(&domain.CheckTypeSucceeded{
			// 	SucceededAt: e.CreatedAt(),
			// }),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceTOTPCheckFailed(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanOTPCheckFailedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qYHvj", "reduce.wrong.event.type %s", user.HumanMFAOTPCheckFailedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			// repo.CheckTOTP(&domain.CheckTypeFailed{
			// 	FailedAt: e.CreatedAt(),
			// }),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceOTPSMSEnabled(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanOTPSMSAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qYHvj", "reduce.wrong.event.type %s", user.HumanOTPSMSAddedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.EnableSMSOTPAt(e.CreatedAt()),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceOTPSMSDisabled(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanOTPSMSRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qYHvj", "reduce.wrong.event.type %s", user.HumanOTPSMSRemovedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.DisableSMSOTP(),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceOTPSMSCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanOTPSMSCodeAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qYHvj", "reduce.wrong.event.type %s", user.HumanOTPSMSCodeAddedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			// repo.CheckSMSOTP(&domain.CheckTypeInit{
			// 	Code:      e.Code,
			// 	CreatedAt: e.CreatedAt(),
			// 	Expiry:    &e.Expiry,
			// }),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceOTPSMSCheckSucceeded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanOTPSMSCheckSucceededEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qYHvj", "reduce.wrong.event.type %s", user.HumanOTPSMSCheckSucceededType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			// repo.CheckSMSOTP(&domain.CheckTypeSucceeded{
			// 	SucceededAt: e.CreatedAt(),
			// }),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceOTPSMSCheckFailed(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanOTPSMSCheckFailedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qYHvj", "reduce.wrong.event.type %s", user.HumanOTPSMSCheckFailedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			// repo.CheckSMSOTP(&domain.CheckTypeFailed{
			// 	FailedAt: e.CreatedAt(),
			// }),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceOTPEmailEnabled(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanOTPEmailAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qYHvj", "reduce.wrong.event.type %s", user.HumanOTPEmailAddedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.EnableEmailOTPAt(e.CreatedAt()),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceOTPEmailDisabled(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanOTPEmailRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qYHvj", "reduce.wrong.event.type %s", user.HumanOTPEmailRemovedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.DisableEmailOTP(),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceOTPEmailCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanOTPEmailCodeAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qYHvj", "reduce.wrong.event.type %s", user.HumanOTPEmailCodeAddedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			// repo.CheckEmailOTP(&domain.CheckTypeInit{
			// 	Code:      e.Code,
			// 	CreatedAt: e.CreatedAt(),
			// 	Expiry:    &e.Expiry,
			// }),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceOTPEmailCheckSucceeded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanOTPEmailCheckSucceededEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qYHvj", "reduce.wrong.event.type %s", user.HumanOTPEmailCheckSucceededType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			// repo.CheckEmailOTP(&domain.CheckTypeSucceeded{
			// 	SucceededAt: e.CreatedAt(),
			// }),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceOTPEmailCheckFailed(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanOTPEmailCheckFailedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-qYHvj", "reduce.wrong.event.type %s", user.HumanOTPEmailCheckFailedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			// repo.CheckEmailOTP(&domain.CheckTypeFailed{
			// 	FailedAt: e.CreatedAt(),
			// }),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func mapHumanGender(gender old_domain.Gender) domain.HumanGender {
	switch gender {
	case old_domain.GenderFemale:
		return domain.HumanGenderFemale
	case old_domain.GenderMale:
		return domain.HumanGenderMale
	case old_domain.GenderDiverse:
		return domain.HumanGenderDiverse
	default:
		return domain.HumanGenderUnspecified
	}
}

func mapMachineAccessTokenType(tokenType old_domain.OIDCTokenType) domain.AccessTokenType {
	switch tokenType {
	case old_domain.OIDCTokenTypeBearer:
		return domain.AccessTokenTypeBearer
	case old_domain.OIDCTokenTypeJWT:
		return domain.AccessTokenTypeJWT
	default:
		return domain.AccessTokenTypeUnspecified
	}
}

func mapMachineKeyType(keyType old_domain.AuthNKeyType) domain.MachineKeyType {
	switch keyType {
	case old_domain.AuthNKeyTypeJSON:
		return domain.MachineKeyTypeJSON
	case old_domain.AuthNKeyTypeNONE:
		fallthrough
	default:
		return domain.MachineKeyTypeNone
	}
}

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
					Event:  user.UserV1PasswordChangedType,
					Reduce: p.reduceHumanPasswordChanged,
				},
				{
					Event:  user.HumanPasswordCodeAddedType,
					Reduce: p.reduceHumanPasswordCodeAdded,
				},
				{
					Event:  user.UserV1PasswordCodeAddedType,
					Reduce: p.reduceHumanPasswordCodeAdded,
				},
				{
					Event:  user.HumanPasswordCheckSucceededType,
					Reduce: p.reduceHumanPasswordCheckSucceeded,
				},
				{
					Event:  user.UserV1PasswordCheckSucceededType,
					Reduce: p.reduceHumanPasswordCheckSucceeded,
				},
				{
					Event:  user.HumanPasswordCheckFailedType,
					Reduce: p.reduceHumanPasswordCheckFailed,
				},
				{
					Event:  user.UserV1PasswordCheckFailedType,
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
				{
					Event:  user.UserV1MFAInitSkippedType,
					Reduce: p.reduceMFAInitSkipped,
				},
				{
					Event:  user.HumanMFAInitSkippedType,
					Reduce: p.reduceMFAInitSkipped,
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
					Event:  user.UserV1MFAOTPAddedType,
					Reduce: p.reduceTOTPAdded,
				},
				{
					Event:  user.HumanMFAOTPVerifiedType,
					Reduce: p.reduceTOTPVerified,
				},
				{
					Event:  user.UserV1MFAOTPVerifiedType,
					Reduce: p.reduceTOTPVerified,
				},
				{
					Event:  user.HumanMFAOTPRemovedType,
					Reduce: p.reduceTOTPRemoved,
				},
				{
					Event:  user.UserV1MFAOTPRemovedType,
					Reduce: p.reduceTOTPRemoved,
				},
				{
					Event:  user.HumanMFAOTPCheckSucceededType,
					Reduce: p.reduceTOTPCheckSucceeded,
				},
				{
					Event:  user.UserV1MFAOTPCheckSucceededType,
					Reduce: p.reduceTOTPCheckSucceeded,
				},
				{
					Event:  user.HumanMFAOTPCheckFailedType,
					Reduce: p.reduceTOTPCheckFailed,
				},
				{
					Event:  user.UserV1MFAOTPCheckFailedType,
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
					Event:  user.HumanOTPEmailCheckSucceededType,
					Reduce: p.reduceOTPEmailCheckSucceeded,
				},
				{
					Event:  user.HumanOTPEmailCheckFailedType,
					Reduce: p.reduceOTPEmailCheckFailed,
				},
				{
					Event:  user.HumanInviteCodeAddedType,
					Reduce: p.reduceInviteCodeAdded,
				},
				{
					Event:  user.HumanInviteCheckSucceededType,
					Reduce: p.reduceInviteCheckSucceeded,
				},
				{
					Event:  user.HumanInviteCheckFailedType,
					Reduce: p.reduceInviteCheckFailed,
				},
			},
		},
	}
}

func (u *userRelationalProjection) reduceHumanAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanAddedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		v3Tx := v3_sql.SQLTx(tx)

		userRepo := repository.UserRepository()

		var phone *domain.HumanPhone
		if e.PhoneNumber != "" {
			phone = &domain.HumanPhone{
				UnverifiedNumber: string(e.PhoneNumber),
			}
		}

		password := domain.HumanPassword{
			IsChangeRequired: e.ChangeRequired,
		}
		if hash := crypto.SecretOrEncodedHash(e.Secret, e.EncodedHash); hash != "" {
			password.Hash = hash
			password.ChangedAt = e.CreatedAt()
		}

		return userRepo.Create(ctx, v3Tx, &domain.User{
			InstanceID:     e.Aggregate().InstanceID,
			OrganizationID: e.Aggregate().ResourceOwner,
			ID:             e.Aggregate().ID,
			Username:       e.UserName,
			State:          domain.UserStateActive,
			CreatedAt:      e.CreatedAt(),
			UpdatedAt:      e.CreatedAt(),
			// TODO check when to set username unique
			// IsUsernameOrgUnique: ,
			Human: &domain.HumanUser{
				FirstName:         e.FirstName,
				LastName:          e.LastName,
				Nickname:          e.NickName,
				DisplayName:       e.DisplayName,
				PreferredLanguage: e.PreferredLanguage,
				Gender:            mapHumanGender(e.Gender),
				Email: domain.HumanEmail{
					UnverifiedAddress: string(e.EmailAddress),
				},
				Phone:    phone,
				Password: password,
			},
		})
	}), nil
}

func (p *userRelationalProjection) reduceHumanRegistered(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanRegisteredEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		v3Tx := v3_sql.SQLTx(tx)

		userRepo := repository.UserRepository()

		var phone *domain.HumanPhone
		if e.PhoneNumber != "" {
			phone = &domain.HumanPhone{
				UnverifiedNumber: string(e.PhoneNumber),
			}
		}

		password := domain.HumanPassword{
			IsChangeRequired: e.ChangeRequired,
		}
		if hash := crypto.SecretOrEncodedHash(e.Secret, e.EncodedHash); hash != "" {
			password.Hash = hash
			password.ChangedAt = e.CreatedAt()
		}

		return userRepo.Create(ctx, v3Tx, &domain.User{
			InstanceID:     e.Aggregate().InstanceID,
			OrganizationID: e.Aggregate().ResourceOwner,
			ID:             e.Aggregate().ID,
			Username:       e.UserName,
			State:          domain.UserStateActive,
			CreatedAt:      e.CreatedAt(),
			UpdatedAt:      e.CreatedAt(),
			// TODO check when to set username unique
			// IsUsernameOrgUnique: ,
			Human: &domain.HumanUser{
				FirstName:         e.FirstName,
				LastName:          e.LastName,
				Nickname:          e.NickName,
				DisplayName:       e.DisplayName,
				PreferredLanguage: e.PreferredLanguage,
				Gender:            mapHumanGender(e.Gender),
				Email: domain.HumanEmail{
					UnverifiedAddress: string(e.EmailAddress),
				},
				Phone:    phone,
				Password: password,
			},
		})
	}), nil
}

func (p *userRelationalProjection) reduceUserLocked(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.UserLockedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.UserUnlockedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.UserDeactivatedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.UserReactivatedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.UserRemovedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.UsernameChangedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.DomainClaimedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.HumanProfileChangedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.HumanPhoneChangedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetUnverifiedPhone(string(e.PhoneNumber)),
			repo.SetPhoneVerification(&domain.VerificationTypeInit{
				CreatedAt: e.CreatedAt(),
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanPhoneRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanPhoneRemovedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.HumanPhoneVerifiedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetPhoneVerification(&domain.VerificationTypeSucceeded{
				VerifiedAt: e.CreatedAt(),
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanPhoneCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanPhoneCodeAddedEvent](event)
	if err != nil {
		return nil, err
	}

	var expiry *time.Duration
	if e.Expiry > 0 {
		expiry = &e.Expiry
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository().Human()

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetPhoneVerification(&domain.VerificationTypeUpdate{
				Code:   e.Code,
				Expiry: expiry,
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanPhoneVerificationFailed(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanPhoneVerificationFailedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository().Human()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetPhoneVerification(&domain.VerificationTypeFailed{
				FailedAt: e.CreatedAt(),
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanEmailChanged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanEmailChangedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository().Human()

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetUnverifiedEmail(string(e.EmailAddress)),
			repo.SetEmailVerification(&domain.VerificationTypeInit{
				CreatedAt: e.CreatedAt(),
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanEmailVerified(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanEmailVerifiedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository().Human()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetEmailVerification(&domain.VerificationTypeSucceeded{
				VerifiedAt: e.CreatedAt(),
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanEmailCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanEmailCodeAddedEvent](event)
	if err != nil {
		return nil, err
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
			repo.SetEmailVerification(&domain.VerificationTypeUpdate{
				Code:   e.Code,
				Expiry: expiry,
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanEmailVerificationFailed(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanEmailVerificationFailedEvent](event)
	if err != nil {
		return nil, err
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository().Human()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetEmailVerification(&domain.VerificationTypeFailed{
				FailedAt: e.CreatedAt(),
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanAvatarAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanAvatarAddedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.HumanAvatarRemovedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.HumanPasswordChangedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()

		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetPasswordChangeRequired(e.ChangeRequired),
			repo.SetPassword(crypto.SecretOrEncodedHash(e.Secret, e.EncodedHash)),
			repo.SetResetPasswordVerification(&domain.VerificationTypeSucceeded{VerifiedAt: e.CreatedAt()}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanPasswordCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanPasswordCodeAddedEvent](event)
	if err != nil {
		return nil, err
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
			repo.SetResetPasswordVerification(&domain.VerificationTypeUpdate{
				Code:   e.Code,
				Expiry: expiry,
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanPasswordCheckSucceeded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanPasswordCheckSucceededEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetLastSuccessfulPasswordCheck(e.CreatedAt()),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanPasswordCheckFailed(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanPasswordCheckFailedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.IncrementPasswordFailedAttempts(),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanPasswordHashUpdated(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanPasswordHashUpdatedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetPassword(e.EncodedHash),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceMachineSecretSet(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.MachineSecretSetEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.MachineSecretHashUpdatedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.MachineSecretRemovedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.MachineKeyAddedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.MachineKeyRemovedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.MachineAddedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.MachineChangedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.PersonalAccessTokenAddedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.PersonalAccessTokenRemovedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.MetadataSetEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-xg4IJ", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.UserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetMetadata(&domain.Metadata{
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
	e, err := assertEvent[*user.MetadataRemovedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.MetadataRemovedEvent](event)
	if err != nil {
		return nil, err
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

		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.UpdatePasskey(
				repo.PasskeyConditions().IDCondition(e.WebAuthNTokenID),
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
				repo.PasskeyConditions().IDCondition(e.WebAuthNTokenID),
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
			userRepo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			humanRepo.RemovePasskey(humanRepo.PasskeyConditions().IDCondition(e.WebAuthNTokenID)),
			humanRepo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reducePasskeyInitCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanPasswordlessInitCodeAddedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.HumanPasswordlessInitCodeCheckFailedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.HumanPasswordlessInitCodeCheckSucceededEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetVerification(&domain.VerificationTypeSucceeded{
				ID:         &e.ID,
				VerifiedAt: e.CreatedAt(),
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reducePasskeyInitCodeRequested(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanPasswordlessInitCodeRequestedEvent](event)
	if err != nil {
		return nil, err
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

func (p *userRelationalProjection) reduceMFAInitSkipped(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanMFAInitSkippedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.UserIDPLinkAddedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.UserIDPLinkCascadeRemovedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.UserIDPLinkRemovedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.UserIDPExternalIDMigratedEvent](event)
	if err != nil {
		return nil, err
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
				database.And(repo.IdentityProviderLinkConditions().ProviderIDCondition(e.IDPConfigID), repo.IdentityProviderLinkConditions().ProvidedUserIDCondition(e.PreviousID)),
				repo.SetIdentityProviderLinkProvidedID(e.NewID),
			),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceIDPLinkUsernameChanged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.UserIDPExternalUsernameEvent](event)
	if err != nil {
		return nil, err
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
				database.And(repo.IdentityProviderLinkConditions().ProviderIDCondition(e.IDPConfigID), repo.IdentityProviderLinkConditions().ProvidedUserIDCondition(e.ExternalUserID)),
				repo.SetIdentityProviderLinkUsername(e.ExternalUsername),
			),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceTOTPAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanOTPAddedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetTOTPSecret(e.Secret),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceTOTPVerified(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanOTPVerifiedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetTOTPVerifiedAt(e.CreatedAt()),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceTOTPRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanOTPRemovedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.HumanOTPCheckSucceededEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetLastSuccessfulTOTPCheck(e.CreatedAt()),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceTOTPCheckFailed(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanOTPCheckFailedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.IncrementTOTPFailedAttempts(),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceOTPSMSEnabled(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanOTPSMSAddedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.HumanOTPSMSRemovedEvent](event)
	if err != nil {
		return nil, err
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

func (p *userRelationalProjection) reduceOTPSMSCheckSucceeded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanOTPSMSCheckSucceededEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetLastSuccessfulSMSOTPCheck(e.CreatedAt()),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceOTPSMSCheckFailed(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanOTPSMSCheckFailedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.IncrementSMSOTPFailedAttempts(),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceOTPEmailEnabled(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanOTPEmailAddedEvent](event)
	if err != nil {
		return nil, err
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
	e, err := assertEvent[*user.HumanOTPEmailRemovedEvent](event)
	if err != nil {
		return nil, err
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

func (p *userRelationalProjection) reduceOTPEmailCheckSucceeded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanOTPEmailCheckSucceededEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetLastSuccessfulEmailOTPCheck(e.CreatedAt()),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceOTPEmailCheckFailed(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanOTPEmailCheckFailedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.IncrementEmailOTPFailedAttempts(),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceInviteCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanInviteCodeAddedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}

		var expiry *time.Duration
		if e.Expiry > 0 {
			expiry = &e.Expiry
		}

		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetInviteVerification(&domain.VerificationTypeInit{
				CreatedAt: e.CreatedAt(),
				Code:      e.Code,
				Expiry:    expiry,
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceInviteCheckSucceeded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanInviteCheckSucceededEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetInviteVerification(&domain.VerificationTypeSucceeded{
				VerifiedAt: e.CreatedAt(),
			}),
			repo.SetUpdatedAt(e.CreatedAt()),
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceInviteCheckFailed(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.HumanInviteCheckFailedEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}

		repo := repository.HumanUserRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(e.Aggregate().InstanceID, e.Aggregate().ID),
			repo.SetInviteVerification(&domain.VerificationTypeFailed{
				FailedAt: e.CreatedAt(),
			}),
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
	case old_domain.GenderUnspecified:
		fallthrough
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

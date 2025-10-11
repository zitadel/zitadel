package projection

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/org"
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
				// TODO
				// {
				// 	Event:  user.HumanPhoneChangedType,
				// 	Reduce: p.reduceHumanPhoneChanged,
				// },
				// {
				// 	Event:  user.UserV1PhoneChangedType,
				// 	Reduce: p.reduceHumanPhoneChanged,
				// },
				// TODO
				// 		{
				// 			Event:  user.HumanPhoneRemovedType,
				// 			Reduce: p.reduceHumanPhoneRemoved,
				// 		},
				// 		{
				// 			Event:  user.UserV1PhoneRemovedType,
				// 			Reduce: p.reduceHumanPhoneRemoved,
				// 		},
				// TODO
				// 		{
				// 			Event:  user.HumanPhoneVerifiedType,
				// 			Reduce: p.reduceHumanPhoneVerified,
				// 		},
				// 		{
				// 			Event:  user.UserV1PhoneVerifiedType,
				// 			Reduce: p.reduceHumanPhoneVerified,
				// 		},
				// TODO
				// 		{
				// 			Event:  user.HumanEmailChangedType,
				// 			Reduce: p.reduceHumanEmailChanged,
				// 		},
				// 		{
				// 			Event:  user.UserV1EmailChangedType,
				// 			Reduce: p.reduceHumanEmailChanged,
				// 		},
				// TODO
				// 		{
				// 			Event:  user.HumanEmailVerifiedType,
				// 			Reduce: p.reduceHumanEmailVerified,
				// 		},
				// 		{
				// 			Event:  user.UserV1EmailVerifiedType,
				// 			Reduce: p.reduceHumanEmailVerified,
				// 		},
				///////////////////////////////////////////////////////////////////////////////// TARGET
				// TODO
				// {
				// 	Event:  user.HumanAvatarAddedType,
				// 	Reduce: p.reduceHumanAvatarAdded,
				// },
				// 		{
				// 			Event:  user.HumanAvatarRemovedType,
				// 			Reduce: p.reduceHumanAvatarRemoved,
				// 		},
						{
							Event:  user.MachineAddedEventType,
							Reduce: p.reduceMachineAdded,
						},
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
	passwordSet := crypto.SecretOrEncodedHash(e.Secret, e.EncodedHash) != ""
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		userRepo := repository.UserRepository()

		_, err := userRepo.CreateHuman(ctx, v3_sql.SQLTx(tx),
			&domain.Human{
				FirstName:         e.FirstName,
				LastName:          e.LastName,
				NickName:          e.NickName,
				DisplayName:       e.DisplayName,
				PreferredLanguage: e.PreferredLanguage.String(),
				Gender:            uint8(e.Gender),
				User: domain.User{
					ID:         e.Aggregate().ID,
					InstanceID: e.Aggregate().InstanceID,
					OrgID:      e.Aggregate().ResourceOwner,
					Username:   e.UserName,
					// TODO check when to set username unique
					// UsernameOrgUnique: false,
					State:     domain.UserStateActive,
					CreatedAt: e.CreationDate(),
					UpdatedAt: e.CreationDate(),
				},
				HumanEmailContact: domain.HumanContact{
					Type:       gu.Ptr(domain.ContactTypeEmail),
					Value:      gu.Ptr(string(e.EmailAddress.Normalize())),
					IsVerified: gu.Ptr(false),
				},
				HumanPhoneContact: func() *domain.HumanContact {
					if e.PhoneNumber == "" {
						return nil
					}
					return &domain.HumanContact{
						Type:       gu.Ptr(domain.ContactTypePhone),
						Value:      gu.Ptr(string(e.PhoneNumber)),
						IsVerified: gu.Ptr(false),
					}
				}(),
				HumanSecurity: domain.HumanSecurity{
					PasswordChangeRequired: e.ChangeRequired,
					PasswordChange: func() *time.Time {
						if !passwordSet {
							return nil
						}
						passwordChange := e.CreatedAt()
						return &passwordChange
					}(),
				},
			},
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanRegistered(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanRegisteredEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-xhD3q", "reduce.wrong.event.type %s", user.HumanRegisteredType)
	}
	fmt.Println("[DEBUGPRINT] [users_test.go:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> HUMAN REGISTERED REDUCE")

	passwordSet := crypto.SecretOrEncodedHash(e.Secret, e.EncodedHash) != ""
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		userRepo := repository.UserRepository()

		_, err := userRepo.CreateHuman(ctx, v3_sql.SQLTx(tx),
			&domain.Human{
				FirstName:         e.FirstName,
				LastName:          e.LastName,
				NickName:          e.NickName,
				DisplayName:       e.DisplayName,
				PreferredLanguage: e.PreferredLanguage.String(),
				Gender:            uint8(e.Gender),
				User: domain.User{
					ID:         e.Aggregate().ID,
					InstanceID: e.Aggregate().InstanceID,
					OrgID:      e.Aggregate().ResourceOwner,
					Username:   e.UserName,
					// TODO check when to set username unique
					// UsernameOrgUnique: false,
					State:     domain.UserStateActive,
					CreatedAt: e.CreationDate(),
					UpdatedAt: e.CreationDate(),
				},
				HumanEmailContact: domain.HumanContact{
					Type:       gu.Ptr(domain.ContactTypeEmail),
					Value:      gu.Ptr(string(e.EmailAddress.Normalize())),
					IsVerified: gu.Ptr(false),
				},
				HumanPhoneContact: func() *domain.HumanContact {
					if e.PhoneNumber == "" {
						return nil
					}
					return &domain.HumanContact{
						Type:       gu.Ptr(domain.ContactTypePhone),
						Value:      gu.Ptr(string(e.PhoneNumber)),
						IsVerified: gu.Ptr(false),
					}
				}(),
				HumanSecurity: domain.HumanSecurity{
					PasswordChangeRequired: e.ChangeRequired,
					PasswordChange: func() *time.Time {
						if !passwordSet {
							return nil
						}
						passwordChange := e.CreatedAt()
						return &passwordChange
					}(),
				},
			},
		)
		return err
	}), nil
}

func (p *userRelationalProjection) reduceHumanInitCodeAdded(event eventstore.Event) (*handler.Statement, error) {
	fmt.Println("[DEBUGPRINT] [users_test.go:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> INITAL")
	e, ok := event.(*user.HumanInitialCodeAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-vv7Qs", "reduce.wrong.event.type %s", user.HumanInitialCodeAddedType)
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(UserStateCol, domain.UserStateInital),
		},
		[]handler.Condition{
			handler.NewCond(UserIDCol, e.Aggregate().ID),
			handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *userRelationalProjection) reduceHumanInitCodeSucceeded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanInitializedCheckSucceededEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-ifH8N", "reduce.wrong.event.type %s", user.HumanInitializedCheckSucceededType)
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

func (p *userRelationalProjection) reduceUserLocked(event eventstore.Event) (*handler.Statement, error) {
	fmt.Println("[DEBUGPRINT] [users_test.go:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> LOCKED")
	e, ok := event.(*user.UserLockedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-eUn8f", "reduce.wrong.event.type %s", user.UserLockedType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		userRepo := repository.UserRepository()

		noOfRecordsUpdated, err := userRepo.UpdateHuman(ctx, v3_sql.SQLTx(tx),
			database.And(
				userRepo.Human().InstanceIDCondition(e.Aggregate().InstanceID),
				userRepo.Human().IDCondition(e.Aggregate().ID),
			),
			userRepo.Human().SetState(domain.UserStateLocked),
			userRepo.Human().SetUpdatedAt(e.CreationDate()),
		)
		if err != nil {
			return err
		}
		if noOfRecordsUpdated == 0 {
			return zerrors.ThrowNotFound(nil, "HANDL-SD3fs", "Errors.User.NotFound")
		} else if noOfRecordsUpdated > 1 {
			tx.Rollback()
			// TODO add "Errors.User.TooManyEntries"
			return zerrors.ThrowInternal(nil, "HANDL-Df3fs", "Errors.User.TooManyEntries")
		}
		return nil
	}), nil
	// return handler.NewUpdateStatement(
	// 	e,
	// 	[]handler.Column{
	// 		handler.NewCol(UserChangeDateCol, e.CreationDate()),
	// 		handler.NewCol(UserStateCol, domain.UserStateLocked),
	// 		handler.NewCol(UserSequenceCol, e.Sequence()),
	// 	},
	// 	[]handler.Condition{
	// 		handler.NewCond(UserIDCol, e.Aggregate().ID),
	// 		handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
	// 	},
	// ), nil
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
		userRepo := repository.UserRepository()

		noOfRecordsUpdated, err := userRepo.UpdateHuman(ctx, v3_sql.SQLTx(tx),
			database.And(
				userRepo.Human().InstanceIDCondition(e.Aggregate().InstanceID),
				userRepo.Human().IDCondition(e.Aggregate().ID),
			),
			userRepo.Human().SetState(domain.UserStateActive),
			userRepo.Human().SetUpdatedAt(e.CreationDate()),
		)
		if err != nil {
			return err
		}
		if noOfRecordsUpdated == 0 {
			return zerrors.ThrowNotFound(nil, "HANDL-SD3fs", "Errors.User.NotFound")
		} else if noOfRecordsUpdated > 1 {
			tx.Rollback()
			// TODO add "Errors.User.TooManyEntries"
			return zerrors.ThrowInternal(nil, "HANDL-Df3fs", "Errors.User.TooManyEntries")
		}
		return nil
	}), nil
	// return handler.NewUpdateStatement(
	// 	e,
	// 	[]handler.Column{
	// 		handler.NewCol(UserChangeDateCol, e.CreationDate()),
	// 		handler.NewCol(UserStateCol, domain.UserStateActive),
	// 		handler.NewCol(UserSequenceCol, e.Sequence()),
	// 	},
	// 	[]handler.Condition{
	// 		handler.NewCond(UserIDCol, e.Aggregate().ID),
	// 		handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
	// 	},
	// ), nil
}

func (p *userRelationalProjection) reduceUserDeactivated(event eventstore.Event) (*handler.Statement, error) {
	fmt.Println("[DEBUGPRINT] [users_test.go:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> DEACTIVATED")
	e, ok := event.(*user.UserDeactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-6BNjj", "reduce.wrong.event.type %s", user.UserDeactivatedType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		userRepo := repository.UserRepository()

		noOfRecordsUpdated, err := userRepo.UpdateHuman(ctx, v3_sql.SQLTx(tx),
			database.And(
				userRepo.Human().InstanceIDCondition(e.Aggregate().InstanceID),
				userRepo.Human().IDCondition(e.Aggregate().ID),
			),
			userRepo.Human().SetState(domain.UserStateInactive),
			userRepo.Human().SetUpdatedAt(e.CreationDate()),
		)
		if err == nil && noOfRecordsUpdated == 1 {
			return nil
		}
		if err != nil {
			return err
		}
		if noOfRecordsUpdated == 0 {
			return zerrors.ThrowNotFound(nil, "HANDL-SD3fs", "Errors.User.NotFound")
		} else if noOfRecordsUpdated > 1 {
			tx.Rollback()
			// TODO add "Errors.User.TooManyEntries"
			return zerrors.ThrowInternal(nil, "HANDL-Df3fs", "Errors.User.TooManyEntries")
		}
		return nil
	}), nil
	// return handler.NewUpdateStatement(
	// 	e,
	// 	[]handler.Column{
	// 		handler.NewCol(UserChangeDateCol, e.CreationDate()),
	// 		handler.NewCol(UserStateCol, domain.UserStateInactive),
	// 		handler.NewCol(UserSequenceCol, e.Sequence()),
	// 	},
	// 	[]handler.Condition{
	// 		handler.NewCond(UserIDCol, e.Aggregate().ID),
	// 		handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
	// 	},
	// ), nil
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
		userRepo := repository.UserRepository()

		noOfRecordsUpdated, err := userRepo.UpdateHuman(ctx, v3_sql.SQLTx(tx),
			database.And(
				userRepo.Human().InstanceIDCondition(e.Aggregate().InstanceID),
				userRepo.Human().IDCondition(e.Aggregate().ID),
			),
			userRepo.Human().SetState(domain.UserStateActive),
			userRepo.Human().SetUpdatedAt(e.CreationDate()),
		)
		if err != nil {
			return err
		}
		if noOfRecordsUpdated == 0 {
			return zerrors.ThrowNotFound(nil, "HANDL-SD3fs", "Errors.User.NotFound")
		} else if noOfRecordsUpdated > 1 {
			tx.Rollback()
			// TODO add "Errors.User.TooManyEntries"
			return zerrors.ThrowInternal(nil, "HANDL-Df3fs", "Errors.User.TooManyEntries")
		}
		return nil
	}), nil

	// return handler.NewUpdateStatement(
	// 	e,
	// 	[]handler.Column{
	// 		handler.NewCol(UserChangeDateCol, e.CreationDate()),
	// 		handler.NewCol(UserStateCol, domain.UserStateActive),
	// 		handler.NewCol(UserSequenceCol, e.Sequence()),
	// 	},
	// 	[]handler.Condition{
	// 		handler.NewCond(UserIDCol, e.Aggregate().ID),
	// 		handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
	// 	},
	// ), nil
}

func (p *userRelationalProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	fmt.Println("[DEBUGPRINT] [users_test.go:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> USER DELETED")
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-BQB2t", "reduce.wrong.event.type %s", user.UserRemovedType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		userRepo := repository.UserRepository()

		noOfRecordsDeleted, err := userRepo.DeleteHuman(ctx, v3_sql.SQLTx(tx),
			database.And(
				userRepo.Human().InstanceIDCondition(e.Aggregate().InstanceID),
				userRepo.Human().IDCondition(e.Aggregate().ID),
			),
		)
		if err != nil {
			return err
		}
		if noOfRecordsDeleted == 0 {
			return zerrors.ThrowNotFound(nil, "HANDL-SD3fs", "Errors.User.NotFound")
		} else if noOfRecordsDeleted > 1 {
			tx.Rollback()
			// TODO add "Errors.User.TooManyEntries"
			return zerrors.ThrowInternal(nil, "HANDL-Df3fs", "Errors.User.TooManyEntries")
		}
		return nil
	}), nil

	// return handler.NewDeleteStatement(
	// 	e,
	// 	[]handler.Condition{
	// 		handler.NewCond(UserIDCol, e.Aggregate().ID),
	// 		handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
	// 	},
	// ), nil
}

func (p *userRelationalProjection) reduceUserNameChanged(event eventstore.Event) (*handler.Statement, error) {
	fmt.Println("[DEBUGPRINT] [users_test.go:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> USERNAME CHANGE")
	e, ok := event.(*user.UsernameChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-QNKyV", "reduce.wrong.event.type %s", user.UserUserNameChangedType)
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		userRepo := repository.UserRepository()

		noOfRecordsUpdated, err := userRepo.UpdateHuman(ctx, v3_sql.SQLTx(tx),
			database.And(
				userRepo.Human().InstanceIDCondition(e.Aggregate().InstanceID),
				userRepo.Human().IDCondition(e.Aggregate().ID),
			),
			userRepo.Human().SetUsername(e.UserName),
			userRepo.Human().SetUpdatedAt(e.CreationDate()),
		)
		if err != nil {
			return err
		}
		if noOfRecordsUpdated == 0 {
			return zerrors.ThrowNotFound(nil, "HANDL-SD3fs", "Errors.User.NotFound")
		} else if noOfRecordsUpdated > 1 {
			tx.Rollback()
			// TODO add "Errors.User.TooManyEntries"
			return zerrors.ThrowInternal(nil, "HANDL-Df3fs", "Errors.User.TooManyEntries")
		}
		return nil
	}), nil

	// return handler.NewUpdateStatement(
	// 	e,
	// 	[]handler.Column{
	// 		handler.NewCol(UserChangeDateCol, e.CreationDate()),
	// 		handler.NewCol(UserUsernameCol, e.UserName),
	// 		handler.NewCol(UserSequenceCol, e.Sequence()),
	// 	},
	// 	[]handler.Condition{
	// 		handler.NewCond(UserIDCol, e.Aggregate().ID),
	// 		handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
	// 	},
	// ), nil
}

func (p *userRelationalProjection) reduceDomainClaimed(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.DomainClaimedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-ASwf3", "reduce.wrong.event.type %s", user.UserDomainClaimedType)
	}
	fmt.Println("[DEBUGPRINT] [users_test.go:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> DOMAIN CLAIMED")
	fmt.Printf("[DEBUGPRINT] [users_test.go:1] >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> e.UserName = %+v\n", e.UserName)

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}
		userRepo := repository.UserRepository()

		noOfRecordsUpdated, err := userRepo.UpdateHuman(ctx, v3_sql.SQLTx(tx),
			database.And(
				userRepo.Human().InstanceIDCondition(e.Aggregate().InstanceID),
				userRepo.Human().IDCondition(e.Aggregate().ID),
			),
			userRepo.Human().SetUsername(e.UserName),
			userRepo.Human().SetUpdatedAt(e.CreationDate()),
		)
		if err != nil {
			return err
		}
		if noOfRecordsUpdated == 0 {
			return zerrors.ThrowNotFound(nil, "HANDL-SD3fs", "Errors.User.NotFound")
		} else if noOfRecordsUpdated > 1 {
			tx.Rollback()
			// TODO add "Errors.User.TooManyEntries"
			return zerrors.ThrowInternal(nil, "HANDL-Df3fs", "Errors.User.TooManyEntries")
		}
		return nil
	}), nil

	// return handler.NewUpdateStatement(
	// 	e,
	// 	[]handler.Column{
	// 		handler.NewCol(UserChangeDateCol, e.CreationDate()),
	// 		handler.NewCol(UserUsernameCol, e.UserName),
	// 		handler.NewCol(UserSequenceCol, e.Sequence()),
	// 	},
	// 	[]handler.Condition{
	// 		handler.NewCond(UserIDCol, e.Aggregate().ID),
	// 		handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
	// 	},
	// ), nil
}

func (p *userRelationalProjection) reduceHumanProfileChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.HumanProfileChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-769v4", "reduce.wrong.event.type %s", user.HumanProfileChangedType)
	}

	userRepo := repository.UserRepository()
	changes := make([]database.Change, 0, 6)

	if e.FirstName != "" {
		changes = append(changes, userRepo.Human().SetFirstName(e.FirstName))
	}

	if e.LastName != "" {
		changes = append(changes, userRepo.Human().SetLastName(e.LastName))
	}

	if e.NickName != nil {
		changes = append(changes, userRepo.Human().SetNickName(*e.NickName))
	}

	if e.DisplayName != nil {
		changes = append(changes, userRepo.Human().SetDisplayName(*e.DisplayName))
	}

	if e.PreferredLanguage != nil {
		changes = append(changes, userRepo.Human().SetPreferredLanguage(e.PreferredLanguage.String()))
	}

	if e.Gender != nil {
		changes = append(changes, userRepo.Human().SetGender(uint8(*e.Gender)))
	}

	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, projectionName string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-iZGH3", "reduce.wrong.db.pool %T", ex)
		}

		noOfRecordsUpdated, err := userRepo.UpdateHuman(ctx, v3_sql.SQLTx(tx),
			database.And(
				userRepo.Human().InstanceIDCondition(e.Aggregate().InstanceID),
				userRepo.Human().IDCondition(e.Aggregate().ID),
			),
			changes...,
		)
		if err != nil {
			return err
		}
		if noOfRecordsUpdated == 0 {
			return zerrors.ThrowNotFound(nil, "HANDL-SD3fs", "Errors.User.NotFound")
		} else if noOfRecordsUpdated > 1 {
			tx.Rollback()
			// TODO add "Errors.User.TooManyEntries"
			return zerrors.ThrowInternal(nil, "HANDL-Df3fs", "Errors.User.TooManyEntries")
		}
		return nil
	}), nil

	// return handler.NewMultiStatement(
	// 	e,
	// 	handler.AddUpdateStatement(
	// 		[]handler.Column{
	// 			handler.NewCol(UserChangeDateCol, e.CreationDate()),
	// 			handler.NewCol(UserSequenceCol, e.Sequence()),
	// 		},
	// 		[]handler.Condition{
	// 			handler.NewCond(UserIDCol, e.Aggregate().ID),
	// 			handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
	// 		},
	// 	),
	// 	handler.AddUpdateStatement(
	// 		cols,
	// 		[]handler.Condition{
	// 			handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
	// 			handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
	// 		},
	// 		handler.WithTableSuffix(UserHumanSuffix),
	// 	),
	// ), nil
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
		userRepo := repository.UserRepository()

		noOfRecordsUpdated, err := userRepo.UpdateHuman(ctx, v3_sql.SQLTx(tx),
			database.And(
				userRepo.Human().InstanceIDCondition(e.Aggregate().InstanceID),
				userRepo.Human().IDCondition(e.Aggregate().ID),
			),
			userRepo.Human().SetUsername(e.UserName),
			userRepo.Human().SetUpdatedAt(e.CreationDate()),
		)
		if err != nil {
			return err
		}
		if noOfRecordsUpdated == 0 {
			return zerrors.ThrowNotFound(nil, "HANDL-SD3fs", "Errors.User.NotFound")
		} else if noOfRecordsUpdated > 1 {
			tx.Rollback()
			// TODO add "Errors.User.TooManyEntries"
			return zerrors.ThrowInternal(nil, "HANDL-Df3fs", "Errors.User.TooManyEntries")
		}
		return nil
	}), nil

	// return handler.NewMultiStatement(
	// 	e,
	// 	handler.AddUpdateStatement(
	// 		[]handler.Column{
	// 			handler.NewCol(UserChangeDateCol, e.CreationDate()),
	// 			handler.NewCol(UserSequenceCol, e.Sequence()),
	// 		},
	// 		[]handler.Condition{
	// 			handler.NewCond(UserIDCol, e.Aggregate().ID),
	// 			handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
	// 		},
	// 	),
	// 	handler.AddUpdateStatement(
	// 		[]handler.Column{
	// 			handler.NewCol(HumanPhoneCol, e.PhoneNumber),
	// 			handler.NewCol(HumanIsPhoneVerifiedCol, false),
	// 		},
	// 		[]handler.Condition{
	// 			handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
	// 			handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
	// 		},
	// 		handler.WithTableSuffix(UserHumanSuffix),
	// 	),
	// 	handler.AddUpdateStatement(
	// 		[]handler.Column{
	// 			handler.NewCol(NotifyLastPhoneCol, &sql.NullString{String: string(e.PhoneNumber), Valid: e.PhoneNumber != ""}),
	// 		},
	// 		[]handler.Condition{
	// 			handler.NewCond(NotifyUserIDCol, e.Aggregate().ID),
	// 			handler.NewCond(NotifyInstanceIDCol, e.Aggregate().InstanceID),
	// 		},
	// 		handler.WithTableSuffix(UserNotifySuffix),
	// 	),
	// ), nil
}

func (p *userRelationalProjection) reduceHumanPhoneRemoved(event eventstore.Event) (*handler.Statement, error) {
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

func (p *userRelationalProjection) reduceHumanPhoneVerified(event eventstore.Event) (*handler.Statement, error) {
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

func (p *userRelationalProjection) reduceHumanEmailChanged(event eventstore.Event) (*handler.Statement, error) {
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

func (p *userRelationalProjection) reduceHumanEmailVerified(event eventstore.Event) (*handler.Statement, error) {
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
		userRepo := repository.UserRepository()
		noOfRecordsUpdated, err := userRepo.UpdateHuman(ctx, v3_sql.SQLTx(tx),
			database.And(
				userRepo.Human().InstanceIDCondition(e.Aggregate().InstanceID),
				userRepo.Human().IDCondition(e.Aggregate().ID),
			),
			userRepo.Human().SetAvatarKey(e.StoreKey)
			userRepo.Human().SetUpdatedAt(e.CreationDate()),
		)
		if err != nil {
			return err
		}
		if noOfRecordsUpdated == 0 {
			return zerrors.ThrowNotFound(nil, "HANDL-SD3fs", "Errors.User.NotFound")
		} else if noOfRecordsUpdated > 1 {
			tx.Rollback()
			// TODO add "Errors.User.TooManyEntries"
			return zerrors.ThrowInternal(nil, "HANDL-Df3fs", "Errors.User.TooManyEntries")
		}
		return nil
	}), nil

	// return handler.NewMultiStatement(
	// 	e,
	// 	handler.AddUpdateStatement(
	// 		[]handler.Column{
	// 			handler.NewCol(UserChangeDateCol, e.CreationDate()),
	// 			handler.NewCol(UserSequenceCol, e.Sequence()),
	// 		},
	// 		[]handler.Condition{
	// 			handler.NewCond(UserIDCol, e.Aggregate().ID),
	// 			handler.NewCond(UserInstanceIDCol, e.Aggregate().InstanceID),
	// 		},
	// 	),
	// 	handler.AddUpdateStatement(
	// 		[]handler.Column{
	// 			handler.NewCol(HumanAvatarURLCol, e.StoreKey),
	// 		},
	// 		[]handler.Condition{
	// 			handler.NewCond(HumanUserIDCol, e.Aggregate().ID),
	// 			handler.NewCond(HumanUserInstanceIDCol, e.Aggregate().InstanceID),
	// 		},
	// 		handler.WithTableSuffix(UserHumanSuffix),
	// 	),
	// ), nil
}

func (p *userRelationalProjection) reduceHumanAvatarRemoved(event eventstore.Event) (*handler.Statement, error) {
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

func (p *userRelationalProjection) reduceHumanPasswordChanged(event eventstore.Event) (*handler.Statement, error) {
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

func (p *userRelationalProjection) reduceMachineSecretSet(event eventstore.Event) (*handler.Statement, error) {
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

func (p *userRelationalProjection) reduceMachineSecretHashUpdated(event eventstore.Event) (*handler.Statement, error) {
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

func (p *userRelationalProjection) reduceMachineSecretRemoved(event eventstore.Event) (*handler.Statement, error) {
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

func (p *userRelationalProjection) reduceMachineAdded(event eventstore.Event) (*handler.Statement, error) {
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
				// handler.NewCol(UserTypeCol, domain.UserTypeMachine),
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

func (p *userRelationalProjection) reduceMachineChanged(event eventstore.Event) (*handler.Statement, error) {
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

func (p *userRelationalProjection) reduceUnsetMFAInitSkipped(e eventstore.Event) (*handler.Statement, error) {
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

func (p *userRelationalProjection) reduceMFAInitSkipped(event eventstore.Event) (*handler.Statement, error) {
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

func (p *userRelationalProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
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

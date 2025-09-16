package command

import (
	"context"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) ImportHumanRecoveryCodes(ctx context.Context, userID, resourceOwner string, codes []string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if _, err = c.checkUserExists(ctx, userID, resourceOwner); err != nil {
		return err
	}

	recoveryCodeWriteModel := NewHumanRecoveryCodeWriteModel(userID, resourceOwner)
	if err := c.eventstore.FilterToQueryReducer(ctx, recoveryCodeWriteModel); err != nil {
		return err
	}

	if len(recoveryCodeWriteModel.Codes())+len(codes) > c.multifactors.RecoveryCodes.MaxCount {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-53cjw", "Errors.User.MFA.RecoveryCodes.MaxCountExceeded")
	}

	hashedCodes, err := domain.RecoveryCodesFromRaw(codes, c.secretHasher)
	if err != nil {
		return err
	}

	userAgg := UserAggregateFromWriteModelCtx(ctx, &recoveryCodeWriteModel.WriteModel)

	_, err = c.eventstore.Push(ctx,
		user.NewHumanRecoveryCodesAddedEvent(ctx, userAgg, hashedCodes, nil),
	)
	return err
}

type RecoveryCodesDetails struct {
	domain.ObjectDetails
	RawCodes []string
}

func (c *Commands) GenerateRecoveryCodes(ctx context.Context, userID string, count int, resourceOwner string, authRequest *domain.AuthRequest) (*RecoveryCodesDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-4kje7", "Errors.User.UserIDMissing")
	}

	resourceOwner, err := c.checkUserExists(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}

	if err := c.checkPermissionUpdateUserCredentials(ctx, resourceOwner, userID); err != nil {
		return nil, err
	}

	if count <= 0 || count > c.multifactors.RecoveryCodes.MaxCount {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-7c0nx", "Errors.User.RecoveryCodes.CountInvalid")
	}

	if err := c.multifactors.RecoveryCodes.Valid(); err != nil {
		return nil, err
	}

	recoveryCodeWriteModel := NewHumanRecoveryCodeWriteModel(userID, resourceOwner)
	if err := c.eventstore.FilterToQueryReducer(ctx, recoveryCodeWriteModel); err != nil {
		return nil, err
	}

	if len(recoveryCodeWriteModel.Codes())+count > c.multifactors.RecoveryCodes.MaxCount {
		return nil, zerrors.ThrowAlreadyExists(nil, "COMMAND-8f2k9", "Errors.User.MFA.RecoveryCodes.MaxCountExceeded")
	}

	hashedCodes, rawCodes, err := domain.GenerateRecoveryCodes(count, c.multifactors.RecoveryCodes, c.secretHasher)
	if err != nil {
		return nil, err
	}

	userAgg := UserAggregateFromWriteModelCtx(ctx, &recoveryCodeWriteModel.WriteModel)

	_, err = c.eventstore.Push(ctx,
		user.NewHumanRecoveryCodesAddedEvent(ctx, userAgg, hashedCodes, authRequestDomainToAuthRequestInfo(authRequest)),
	)
	if err != nil {
		return nil, err
	}

	return &RecoveryCodesDetails{
		ObjectDetails: domain.ObjectDetails{
			ResourceOwner: resourceOwner,
		},
		RawCodes: rawCodes,
	}, nil
}

func (c *Commands) RemoveRecoveryCodes(ctx context.Context, userID, resourceOwner string, authRequest *domain.AuthRequest) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-l2n9r", "Errors.User.UserIDMissing")
	}

	writeModel := NewHumanRecoveryCodeWriteModel(userID, resourceOwner)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}

	if err := c.checkPermissionUpdateUserCredentials(ctx, writeModel.ResourceOwner, userID); err != nil {
		return nil, err
	}

	if writeModel.UserLocked() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-d9u8q", "Errors.User.RecoveryCodes.Locked")
	}

	if writeModel.State != domain.MFAStateReady {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-84rgg", "Errors.User.RecoveryCodes.NotAdded")
	}

	userAgg := UserAggregateFromWriteModelCtx(ctx, &writeModel.WriteModel)

	_, err := c.eventstore.Push(ctx, user.NewHumanRecoveryCodeRemovedEvent(ctx, userAgg, authRequestDomainToAuthRequestInfo(authRequest)))
	if err != nil {
		return nil, err
	}

	return &domain.ObjectDetails{
		ResourceOwner: resourceOwner,
	}, nil
}

func (c *Commands) HumanCheckRecoveryCode(ctx context.Context, userID, code, resourceOwner string, authRequest *domain.AuthRequest) error {
	commands, err := checkRecoveryCode(ctx, userID, code, resourceOwner, authRequest, c.eventstore.FilterToQueryReducer, c.secretHasher)
	if len(commands) > 0 {
		_, err = c.eventstore.Push(ctx, commands...)
		logging.OnError(err).Error("failed to push recovery code check events")
	}
	return err
}

func checkRecoveryCode(
	ctx context.Context,
	userID, code, resourceOwner string,
	authRequest *domain.AuthRequest,
	queryReducer func(ctx context.Context, r eventstore.QueryReducer) error,
	secretHasher *crypto.Hasher,
) ([]eventstore.Command, error) {
	if code == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-u0b6c", "Errors.User.UserIDMissing")
	}

	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-4m9s2", "Errors.User.UserIDMissing")
	}

	recoveryCodeWm := NewHumanRecoveryCodeWriteModel(userID, resourceOwner)
	err := queryReducer(ctx, recoveryCodeWm)
	if err != nil {
		return nil, err
	}

	if recoveryCodeWm.UserLocked() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-2w6oa", "Errors.User.MFA.RecoveryCodes.Locked")
	}

	if recoveryCodeWm.State != domain.MFAStateReady {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-84rgg", "Errors.User.RecoveryCodes.NotReady")
	}

	hashedCode, err := domain.ValidateRecoveryCode(code, toHumanRecoveryCode(ctx, recoveryCodeWm), secretHasher)

	userAgg := UserAggregateFromWriteModelCtx(ctx, &recoveryCodeWm.WriteModel)
	commands := make([]eventstore.Command, 0, 2)

	authRequestInfo := authRequestDomainToAuthRequestInfo(authRequest)

	if err == nil {
		commands = append(commands, user.NewHumanRecoveryCodeCheckSucceededEvent(ctx, userAgg, hashedCode, authRequestInfo))
	} else {
		commands = append(commands, user.NewHumanRecoveryCodeCheckFailedEvent(ctx, userAgg, authRequestInfo))

		lockoutPolicy, lockoutErr := getLockoutPolicy(ctx, recoveryCodeWm.ResourceOwner, queryReducer)
		logging.OnError(lockoutErr).Error("failed to get lockout policy")

		if lockoutPolicy != nil && lockoutPolicy.MaxOTPAttempts > 0 && recoveryCodeWm.FailedAttempts+1 >= lockoutPolicy.MaxOTPAttempts {
			commands = append(commands, user.NewUserLockedEvent(ctx, userAgg))
		}
	}

	return commands, err
}

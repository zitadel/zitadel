package command

import (
	"context"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) ImportHumanRecoveryCodes(ctx context.Context, userID, resourceOwner string, codes []string) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	if err = c.checkUserExists(ctx, userID, resourceOwner); err != nil {
		return err
	}

	recoveryCodeWriteModel := NewHumanRecoveryCodeWriteModel(userID, resourceOwner)
	if err := c.eventstore.FilterToQueryReducer(ctx, recoveryCodeWriteModel); err != nil {
		return err
	}

	if recoveryCodeWriteModel.State == domain.MFAStateReady {
		return zerrors.ThrowAlreadyExists(nil, "COMMAND-x7k9p", "Errors.User.MFA.RecoveryCodes.AlreadyExists")
	}

	hashedCodes, err := domain.RecoveryCodesFromRaw(codes, c.secretHasher)
	if err != nil {
		return err
	}

	userAgg := UserAggregateFromWriteModel(&recoveryCodeWriteModel.WriteModel)

	_, err = c.eventstore.Push(ctx,
		user.NewHumanRecoveryCodesAddedEvent(ctx, userAgg, hashedCodes, nil),
	)
	return err
}

type RecoveryCodesDetails struct {
	domain.ObjectDetails
	RawCodes []string
}

func (c *Commands) GenerateRecoveryCodes(ctx context.Context, userID, resourceOwner string, authRequest *domain.AuthRequest) (*RecoveryCodesDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-8f2k9", "Errors.User.UserIDMissing")
	}

	recoveryCodeWriteModel := NewHumanRecoveryCodeWriteModel(userID, resourceOwner)
	if err := c.eventstore.FilterToQueryReducer(ctx, recoveryCodeWriteModel); err != nil {
		return nil, err
	}

	if recoveryCodeWriteModel.State == domain.MFAStateReady {
		return nil, zerrors.ThrowAlreadyExists(nil, "COMMAND-8f2k9", "Errors.User.MFA.RecoveryCodes.AlreadyExists")
	}

	hashedCodes, rawCodes, err := domain.GenerateRecoveryCodes(c.multifactors.RecoveryCodes.Count, c.secretHasher)
	if err != nil {
		return nil, err
	}

	userAgg := UserAggregateFromWriteModel(&recoveryCodeWriteModel.WriteModel)

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

func (c *Commands) VerifyRecoveryCode(ctx context.Context, userID, resourceOwner, code string, authRequest *domain.AuthRequest) error {
	if userID == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-S453v", "Errors.User.UserIDMissing")
	}
	if code == "" {
		return zerrors.ThrowInvalidArgument(nil, "COMMAND-SJl2g", "Errors.User.Code.Empty")
	}

	writeModel := NewHumanRecoveryCodeWriteModel(userID, resourceOwner)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return err
	}

	if writeModel.UserLocked() {
		return zerrors.ThrowNotFound(nil, "COMMAND-3f6gz", "Errors.User.RecoveryCodes.Locked")
	}

	userAgg := UserAggregateFromWriteModel(&writeModel.WriteModel)

	recoveryCodes := c.toHumanRecoveryCode(ctx, writeModel)

	valid, index, err := domain.ValidateRecoveryCode(code, recoveryCodes, c.secretHasher)
	if err != nil {
		return err
	}

	// TODO: handle invalid code check and lockout policy

	if !valid {
		return zerrors.ThrowInvalidArgument(err, "COMMAND-84rgg", "Errors.User.Code.Invalid")
	}

	events := []eventstore.Command{user.NewHumanRecoveryCodeCheckSucceededEvent(ctx, userAgg, index, authRequestDomainToAuthRequestInfo(authRequest))}

	_, err = c.eventstore.Push(ctx, events...)
	logging.OnError(err).Error("error creating recovery code check succeeded event")
	return err
}

func (c *Commands) RemoveRecoveryCodes(ctx context.Context, userID, resourceOwner string, authRequest *domain.AuthRequest) (*domain.ObjectDetails, error) {
	if userID == "" {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-l2n9r", "Errors.User.UserIDMissing")
	}

	writeModel := NewHumanRecoveryCodeWriteModel(userID, resourceOwner)
	if err := c.eventstore.FilterToQueryReducer(ctx, writeModel); err != nil {
		return nil, err
	}

	if writeModel.UserLocked() {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-d9u8q", "Errors.User.RecoveryCodes.Locked")
	}

	if writeModel.State != domain.MFAStateReady {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-84rgg", "Errors.User.RecoveryCodes.NotAdded")
	}

	userAgg := UserAggregateFromWriteModel(&writeModel.WriteModel)

	_, err := c.eventstore.Push(ctx, user.NewHumanRecoveryCodeRemovedEvent(ctx, userAgg, authRequestDomainToAuthRequestInfo(authRequest)))
	if err != nil {
		return nil, err
	}

	return &domain.ObjectDetails{
		ResourceOwner: resourceOwner,
	}, nil
}

func (c *Commands) toHumanRecoveryCode(ctx context.Context, recoveryCodeWriteModel *HumanRecoveryCodeWriteModel) *domain.HumanRecoveryCodes {
	return &domain.HumanRecoveryCodes{
		ObjectDetails: writeModelToObjectDetails(&recoveryCodeWriteModel.WriteModel),
		Codes:         recoveryCodeWriteModel.Codes(),
	}
}

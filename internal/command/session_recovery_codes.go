package command

import (
	"context"

	"github.com/zitadel/logging"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) CreateRecoveryCodeChallenge() SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) ([]eventstore.Command, error) {
		if cmd.sessionWriteModel.UserID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-4m9s2", "Errors.User.UserIDMissing")
		}

		recoveryCodeWm := NewHumanRecoveryCodeWriteModel(cmd.sessionWriteModel.UserID, "")
		if err := cmd.eventstore.FilterToQueryReducer(ctx, recoveryCodeWm); err != nil {
			return nil, err
		}

		if recoveryCodeWm.UserLocked() {
			return nil, zerrors.ThrowNotFound(nil, "COMMAND-2w6oa", "Errors.User.MFA.RecoveryCodes.Locked")
		}

		if recoveryCodeWm.State != domain.MFAStateReady {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-84rgg", "Errors.User.RecoveryCodes.NotReady")
		}

		cmd.RecoveryCodeChallenged(ctx)
		return nil, nil
	}
}

func CheckRecoveryCode(code string) SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) ([]eventstore.Command, error) {
		if code == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-u0b6c", "Errors.User.UserIDMissing")
		}

		if cmd.sessionWriteModel.UserID == "" {
			return nil, zerrors.ThrowInvalidArgument(nil, "COMMAND-4m9s2", "Errors.User.UserIDMissing")
		}

		recoveryCodeWm := NewHumanRecoveryCodeWriteModel(cmd.sessionWriteModel.UserID, "")

		queryReducer := cmd.eventstore.FilterToQueryReducer

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

		valid, index, err := domain.ValidateRecoveryCode(code, toHumanRecoveryCode(ctx, recoveryCodeWm), cmd.hasher)
		if err != nil {
			return nil, err
		}

		// TODO: is recheck of new events needed here like in CheckHumanOTP?

		userAgg := UserAggregateFromWriteModel(&recoveryCodeWm.WriteModel)

		commands := make([]eventstore.Command, 0, 2)

		if valid {
			commands = append(commands, user.NewHumanRecoveryCodeCheckSucceededEvent(ctx, userAgg, index, nil))
		} else {
			commands = append(commands, user.NewHumanRecoveryCodeCheckFailedEvent(ctx, userAgg, nil))

			lockoutPolicy, lockoutErr := getLockoutPolicy(ctx, recoveryCodeWm.ResourceOwner, queryReducer)
			logging.OnError(lockoutErr).Error("failed to get lockout policy")

			if lockoutPolicy != nil && lockoutPolicy.MaxOTPAttempts > 0 && recoveryCodeWm.FailedAttempts+1 >= lockoutPolicy.MaxOTPAttempts {
				commands = append(commands, user.NewUserLockedEvent(ctx, userAgg))
			}
		}

		cmd.eventCommands = append(cmd.eventCommands, commands...)
		cmd.RecoveryCodeChecked(ctx, cmd.now())
		return nil, nil
	}
}

func toHumanRecoveryCode(ctx context.Context, recoveryCodeWriteModel *HumanRecoveryCodeWriteModel) *domain.HumanRecoveryCodes {
	return &domain.HumanRecoveryCodes{
		ObjectDetails: writeModelToObjectDetails(&recoveryCodeWriteModel.WriteModel),
		Codes:         recoveryCodeWriteModel.Codes(),
	}
}

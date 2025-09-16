package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

func CheckRecoveryCode(code string) SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) ([]eventstore.Command, error) {
		commands, err := checkRecoveryCode(ctx, cmd.sessionWriteModel.UserID, code, cmd.sessionWriteModel.UserResourceOwner, nil, cmd.eventstore.FilterToQueryReducer, cmd.secretHasher)
		if err != nil {
			return commands, err
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

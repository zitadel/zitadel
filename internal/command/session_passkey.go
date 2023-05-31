package command

import (
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
)

func (s *SessionCommands) getHumanPasswordlessTokenReadModel(ctx context.Context) (*HumanPasswordlessTokensReadModel, error) {
	tokenReadModel := NewHumanPasswordlessTokensReadModel(s.sessionWriteModel.UserID, s.sessionWriteModel.ResourceOwner)
	err := s.eventstore.FilterToQueryReducer(ctx, tokenReadModel)
	if err != nil {
		return nil, err
	}
	return tokenReadModel, nil
}

func (c *Commands) CreatePasskeyChallenge(userVerification domain.UserVerificationRequirement, dst json.Unmarshaler) SessionCommand {
	return func(ctx context.Context, s *SessionCommands) error {
		human, err := s.gethumanWriteModel(ctx)
		if err != nil {
			return err
		}
		tokenReadModel, err := s.getHumanPasswordlessTokenReadModel(ctx)
		if err != nil {
			return err
		}
		webAuthNLogin, err := c.webauthnConfig.BeginLogin(ctx, writeModelToHuman(human), userVerification, readModelToPasswordlessTokens(tokenReadModel)...)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(webAuthNLogin.CredentialAssertionData, dst); err != nil {
			return caos_errs.ThrowInternal(err, "COMMAND-Yah6A", "Errors.Internal")
		}

		s.sessionWriteModel.PasskeyChallenged(ctx, webAuthNLogin.Challenge, webAuthNLogin.AllowedCredentialIDs, webAuthNLogin.UserVerification)
		return nil
	}
}

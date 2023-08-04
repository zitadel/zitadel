package command

import (
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
)

type humanPasskeys struct {
	human  *domain.Human
	tokens []*domain.WebAuthNToken
}

func (s *SessionCommands) getHumanPasskeys(ctx context.Context) (*humanPasskeys, error) {
	humanWritemodel, err := s.gethumanWriteModel(ctx)
	if err != nil {
		return nil, err
	}
	tokenReadModel, err := s.getHumanPasswordlessTokenReadModel(ctx)
	if err != nil {
		return nil, err
	}
	return &humanPasskeys{
		human:  writeModelToHuman(humanWritemodel),
		tokens: readModelToPasswordlessTokens(tokenReadModel),
	}, nil
}

func (s *SessionCommands) getHumanPasswordlessTokenReadModel(ctx context.Context) (*HumanPasswordlessTokensReadModel, error) {
	tokenReadModel := NewHumanPasswordlessTokensReadModel(s.sessionWriteModel.UserID, s.sessionWriteModel.ResourceOwner)
	err := s.eventstore.FilterToQueryReducer(ctx, tokenReadModel)
	if err != nil {
		return nil, err
	}
	return tokenReadModel, nil
}

func (c *Commands) CreatePasskeyChallenge(userVerification domain.UserVerificationRequirement, dst json.Unmarshaler) SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) error {
		humanPasskeys, err := cmd.getHumanPasskeys(ctx)
		if err != nil {
			return err
		}
		webAuthNLogin, err := c.webauthnConfig.BeginLogin(ctx, humanPasskeys.human, userVerification, cmd.sessionWriteModel.Domain, humanPasskeys.tokens...)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(webAuthNLogin.CredentialAssertionData, dst); err != nil {
			return caos_errs.ThrowInternal(err, "COMMAND-Yah6A", "Errors.Internal")
		}

		cmd.PasskeyChallenged(ctx, webAuthNLogin.Challenge, webAuthNLogin.AllowedCredentialIDs, webAuthNLogin.UserVerification)
		return nil
	}
}

func (c *Commands) CheckPasskey(credentialAssertionData json.Marshaler) SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) error {
		credentialAssertionData, err := json.Marshal(credentialAssertionData)
		if err != nil {
			return caos_errs.ThrowInvalidArgument(err, "COMMAND-ohG2o", "todo")
		}
		humanPasskeys, err := cmd.getHumanPasskeys(ctx)
		if err != nil {
			return err
		}
		webAuthN, err := cmd.sessionWriteModel.PasskeyChallenge.WebAuthNLogin(humanPasskeys.human, credentialAssertionData)
		if err != nil {
			return err
		}
		keyID, signCount, err := c.webauthnConfig.FinishLogin(ctx, humanPasskeys.human, webAuthN, credentialAssertionData, humanPasskeys.tokens...)
		if err != nil && keyID == nil {
			return err
		}
		_, token := domain.GetTokenByKeyID(humanPasskeys.tokens, keyID)
		if token == nil {
			return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Aej7i", "Errors.User.WebAuthN.NotFound")
		}
		cmd.PasskeyChecked(ctx, cmd.now(), token.WebAuthNTokenID, signCount)
		return nil
	}
}

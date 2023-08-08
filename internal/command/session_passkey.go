package command

import (
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
)

type humanWebAuthN struct {
	human  *domain.Human
	tokens []*domain.WebAuthNToken
}

func (s *SessionCommands) getHumanPasskeys(ctx context.Context) (*humanWebAuthN, error) {
	humanWritemodel, err := s.gethumanWriteModel(ctx)
	if err != nil {
		return nil, err
	}
	tokenReadModel, err := s.getHumanPasswordlessTokenReadModel(ctx)
	if err != nil {
		return nil, err
	}
	return &humanWebAuthN{
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
		humanPasskeys, err := cmd.getHumanPasskeys(ctx)
		if err != nil {
			return err
		}
		if cmd.sessionWriteModel.PasskeyChallenge == nil {
			return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Ioqu5", "Errors.Session.Passkey.NoChallenge")
		}
		tokenID, signCount, err := c.webAuthNLogin(ctx, humanPasskeys, cmd.sessionWriteModel.PasskeyChallenge, credentialAssertionData)
		if err != nil {
			return err
		}
		cmd.PasskeyChecked(ctx, cmd.now(), tokenID, signCount)
		return nil
	}
}

func (s *SessionCommands) getHumanU2F(ctx context.Context) (*humanWebAuthN, error) {
	humanWritemodel, err := s.gethumanWriteModel(ctx)
	if err != nil {
		return nil, err
	}
	tokenReadModel, err := s.getHumanU2FTokenReadModel(ctx)
	if err != nil {
		return nil, err
	}
	return &humanWebAuthN{
		human:  writeModelToHuman(humanWritemodel),
		tokens: readModelToU2FTokens(tokenReadModel),
	}, nil
}

func (s *SessionCommands) getHumanU2FTokenReadModel(ctx context.Context) (*HumanU2FTokensReadModel, error) {
	tokenReadModel := NewHumanU2FTokensReadModel(s.sessionWriteModel.UserID, s.sessionWriteModel.ResourceOwner)
	err := s.eventstore.FilterToQueryReducer(ctx, tokenReadModel)
	if err != nil {
		return nil, err
	}
	return tokenReadModel, nil
}

func (c *Commands) CreateU2FChallenge(userVerification domain.UserVerificationRequirement, dst json.Unmarshaler) SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) error {
		humanPasskeys, err := cmd.getHumanU2F(ctx)
		if err != nil {
			return err
		}
		webAuthNLogin, err := c.webauthnConfig.BeginLogin(ctx, humanPasskeys.human, userVerification, cmd.sessionWriteModel.Domain, humanPasskeys.tokens...)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(webAuthNLogin.CredentialAssertionData, dst); err != nil {
			return caos_errs.ThrowInternal(err, "COMMAND-Aiy5o", "Errors.Internal")
		}

		cmd.U2FChallenged(ctx, webAuthNLogin.Challenge, webAuthNLogin.AllowedCredentialIDs, webAuthNLogin.UserVerification)
		return nil
	}
}

func (c *Commands) CheckU2F(credentialAssertionData json.Marshaler) SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) error {
		humanU2F, err := cmd.getHumanU2F(ctx)
		if err != nil {
			return err
		}
		if cmd.sessionWriteModel.PasskeyChallenge == nil {
			return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Ioqu5", "Errors.Session.U2F.NoChallenge")
		}
		tokenID, signCount, err := c.webAuthNLogin(ctx, humanU2F, cmd.sessionWriteModel.PasskeyChallenge, credentialAssertionData)
		if err != nil {
			return err
		}
		cmd.U2FChecked(ctx, cmd.now(), tokenID, signCount)
		return nil
	}
}

func (c *Commands) webAuthNLogin(ctx context.Context, humanTokens *humanWebAuthN, challenge *WebAuthNChallengeModel, credentialAssertionData json.Marshaler) (tokenID string, signCount uint32, err error) {
	credentialAssertionDataB, err := json.Marshal(credentialAssertionData)
	if err != nil {
		return "", 0, caos_errs.ThrowInternal(err, "COMMAND-Eba6g", "Errors.Internal")
	}
	webAuthN := challenge.WebAuthNLogin(humanTokens.human, credentialAssertionDataB)
	keyID, signCount, err := c.webauthnConfig.FinishLogin(ctx, humanTokens.human, webAuthN, credentialAssertionDataB, humanTokens.tokens...)
	if err != nil && keyID == nil {
		return "", 0, err
	}
	_, token := domain.GetTokenByKeyID(humanTokens.tokens, keyID)
	if token == nil {
		return "", 0, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-ane2E", "Errors.User.WebAuthN.NotFound")
	}
	return token.WebAuthNTokenID, signCount, nil
}

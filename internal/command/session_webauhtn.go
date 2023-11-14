package command

import (
	"context"
	"encoding/json"

	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
)

type humanWebAuthNTokens struct {
	human  *domain.Human
	tokens []*domain.WebAuthNToken
}

func (s *SessionCommands) getHumanWebAuthNTokens(ctx context.Context, userVerification domain.UserVerificationRequirement) (*humanWebAuthNTokens, error) {
	humanWritemodel, err := s.gethumanWriteModel(ctx)
	if err != nil {
		return nil, err
	}
	tokenReadModel, err := s.getHumanWebAuthNTokenReadModel(ctx, userVerification)
	if err != nil {
		return nil, err
	}
	return &humanWebAuthNTokens{
		human:  writeModelToHuman(humanWritemodel),
		tokens: readModelToWebAuthNTokens(tokenReadModel),
	}, nil
}

func (s *SessionCommands) getHumanWebAuthNTokenReadModel(ctx context.Context, userVerification domain.UserVerificationRequirement) (readModel HumanWebAuthNTokensReadModel, err error) {
	readModel = NewHumanU2FTokensReadModel(s.sessionWriteModel.UserID, "")
	if userVerification == domain.UserVerificationRequirementRequired {
		readModel = NewHumanPasswordlessTokensReadModel(s.sessionWriteModel.UserID, "")
	}
	err = s.eventstore.FilterToQueryReducer(ctx, readModel)
	if err != nil {
		return nil, err
	}
	return readModel, nil
}

func (c *Commands) CreateWebAuthNChallenge(userVerification domain.UserVerificationRequirement, rpid string, dst json.Unmarshaler) SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) error {
		humanPasskeys, err := cmd.getHumanWebAuthNTokens(ctx, userVerification)
		if err != nil {
			return err
		}
		webAuthNLogin, err := c.webauthnConfig.BeginLogin(ctx, humanPasskeys.human, userVerification, rpid, humanPasskeys.tokens...)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(webAuthNLogin.CredentialAssertionData, dst); err != nil {
			return caos_errs.ThrowInternal(err, "COMMAND-Yah6A", "Errors.Internal")
		}

		cmd.WebAuthNChallenged(ctx, webAuthNLogin.Challenge, webAuthNLogin.AllowedCredentialIDs, webAuthNLogin.UserVerification, rpid)
		return nil
	}
}

func (c *Commands) CheckWebAuthN(credentialAssertionData json.Marshaler) SessionCommand {
	return func(ctx context.Context, cmd *SessionCommands) error {
		credentialAssertionData, err := json.Marshal(credentialAssertionData)
		if err != nil {
			return caos_errs.ThrowInternal(err, "COMMAND-ohG2o", "Errors.Internal")
		}
		challenge := cmd.sessionWriteModel.WebAuthNChallenge
		if challenge == nil {
			return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Ioqu5", "Errors.Session.WebAuthN.NoChallenge")
		}
		webAuthNTokens, err := cmd.getHumanWebAuthNTokens(ctx, challenge.UserVerification)
		if err != nil {
			return err
		}
		webAuthN := challenge.WebAuthNLogin(webAuthNTokens.human, credentialAssertionData)

		credential, err := c.webauthnConfig.FinishLogin(ctx, webAuthNTokens.human, webAuthN, credentialAssertionData, webAuthNTokens.tokens...)
		if err != nil && (credential == nil || credential.ID == nil) {
			return err
		}
		_, token := domain.GetTokenByKeyID(webAuthNTokens.tokens, credential.ID)
		if token == nil {
			return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-Aej7i", "Errors.User.WebAuthN.NotFound")
		}
		cmd.WebAuthNChecked(ctx, cmd.now(), token.WebAuthNTokenID, credential.Authenticator.SignCount, credential.Flags.UserVerified)
		return nil
	}
}

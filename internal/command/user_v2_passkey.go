package command

import (
	"context"
	"io"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/user"
	usr_repo "github.com/zitadel/zitadel/internal/repository/user"
)

func (c *Commands) RegisterUserPasskey(ctx context.Context, userID, resourceOwner string, authenticator domain.AuthenticatorAttachment) (*domain.PasskeyRegistrationDetails, error) {
	return c.registerUserPasskey(ctx, userID, resourceOwner, authenticator)
}

func (c *Commands) RegisterUserPasskeyWithCode(ctx context.Context, userID, resourceOwner string, authenticator domain.AuthenticatorAttachment, codeID, code string, alg crypto.EncryptionAlgorithm) (*domain.PasskeyRegistrationDetails, error) {
	event, err := c.verifyUserPasskeyCode(ctx, userID, resourceOwner, codeID, code, alg)
	if err != nil {
		return nil, err
	}

	return c.registerUserPasskey(ctx, userID, resourceOwner, authenticator, event)
}

type eventCallback func(context.Context, *eventstore.Aggregate) eventstore.Command

func (c *Commands) verifyUserPasskeyCode(ctx context.Context, userID, resourceOwner, codeID, code string, alg crypto.EncryptionAlgorithm) (eventCallback, error) {
	wm := NewHumanPasswordlessInitCodeWriteModel(userID, codeID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, wm)
	if err != nil {
		return nil, err
	}
	err = verifyCryptoCode(ctx, c.eventstore.Filter, domain.SecretGeneratorTypePasswordlessInitCode, alg, wm.ChangeDate, wm.Expiration, wm.CryptoCode, code)
	if err != nil || wm.State != domain.PasswordlessInitCodeStateActive {
		c.verifyUserPasskeyCodeFailed(ctx, wm)
		return nil, caos_errs.ThrowInvalidArgument(err, "COMMAND-Eeb2a", "Errors.User.Code.Invalid")
	}
	return func(ctx context.Context, userAgg *eventstore.Aggregate) eventstore.Command {
		return usr_repo.NewHumanPasswordlessInitCodeCheckSucceededEvent(ctx, userAgg, codeID)
	}, nil
}

func (c *Commands) verifyUserPasskeyCodeFailed(ctx context.Context, wm *HumanPasswordlessInitCodeWriteModel) {
	userAgg := UserAggregateFromWriteModel(&wm.WriteModel)
	_, err := c.eventstore.Push(ctx, usr_repo.NewHumanPasswordlessInitCodeCheckFailedEvent(ctx, userAgg, wm.CodeID))
	logging.WithFields("userID", userAgg.ID).OnError(err).Error("RegisterUserPasskeyWithCode push failed")
}

func (c *Commands) registerUserPasskey(ctx context.Context, userID, resourceOwner string, authenticator domain.AuthenticatorAttachment, events ...eventCallback) (*domain.PasskeyRegistrationDetails, error) {
	passwordlessTokens, err := c.getHumanPasswordlessTokens(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}
	wm, userAgg, webAuthN, err := c.addHumanWebAuthN(ctx, userID, resourceOwner, false, passwordlessTokens, authenticator, domain.UserVerificationRequirementRequired)
	if err != nil {
		return nil, err
	}

	cmds := make([]eventstore.Command, len(events)+1)
	cmds[0] = usr_repo.NewHumanPasswordlessAddedEvent(ctx, userAgg, wm.WebauthNTokenID, webAuthN.Challenge)
	for i, event := range events {
		cmds[i+1] = event(ctx, userAgg)
	}

	err = c.pushAppendAndReduce(ctx, wm, cmds...)
	if err != nil {
		return nil, err
	}
	return webAuthN.PasskeyRegistrationDetails(writeModelToObjectDetails(&wm.WriteModel)), nil
}

func (c *Commands) AddUserPasskeyCode(ctx context.Context, userID, resourceOwner string, alg crypto.EncryptionAlgorithm) (*domain.ObjectDetails, error) {
	details, err := c.addUserPasskeyCode(ctx, userID, resourceOwner, alg, "", false)
	return details.ObjectDetails, err
}

func (c *Commands) AddUserPasskeyCodeURLTemplate(ctx context.Context, userID, resourceOwner string, alg crypto.EncryptionAlgorithm, urlTmpl string) (*domain.ObjectDetails, error) {
	if err := domain.RenderPasskeyURLTemplate(io.Discard, urlTmpl, userID, resourceOwner, "codeID", "code"); err != nil {
		return nil, err
	}
	details, err := c.addUserPasskeyCode(ctx, userID, resourceOwner, alg, urlTmpl, false)
	return details.ObjectDetails, err
}

func (c *Commands) AddUserPasskeyCodeReturn(ctx context.Context, userID, resourceOwner string, alg crypto.EncryptionAlgorithm) (*domain.PasskeyCodeDetails, error) {
	return c.addUserPasskeyCode(ctx, userID, resourceOwner, alg, "", true)
}

func (c *Commands) addUserPasskeyCode(ctx context.Context, userID, resourceOwner string, alg crypto.EncryptionAlgorithm, urlTmpl string, returnCode bool) (*domain.PasskeyCodeDetails, error) {
	codeID, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	code, err := newCryptoCodeWithExpiry(ctx, c.eventstore.Filter, domain.SecretGeneratorTypePasswordlessInitCode, alg)
	if err != nil {
		return nil, err
	}
	wm := NewHumanPasswordlessInitCodeWriteModel(userID, codeID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, wm)
	if err != nil {
		return nil, err
	}
	agg := UserAggregateFromWriteModel(&wm.WriteModel)
	cmd := user.NewHumanPasswordlessInitCodeRequestedEvent(ctx, agg, codeID, code.Crypted, code.Expiry, urlTmpl, returnCode)
	err = c.pushAppendAndReduce(ctx, wm, cmd)
	if err != nil {
		return nil, err
	}
	return &domain.PasskeyCodeDetails{
		ObjectDetails: writeModelToObjectDetails(&wm.WriteModel),
		CodeID:        codeID,
		Code:          code.Plain,
	}, nil
}

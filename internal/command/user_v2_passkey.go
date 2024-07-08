package command

import (
	"context"
	"io"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/repository/user"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// RegisterUserPasskey creates a passkey registration for the current authenticated user.
// UserID, usually taken from the request is compared against the user ID in the context.
func (c *Commands) RegisterUserPasskey(ctx context.Context, userID, resourceOwner, rpID string, authenticator domain.AuthenticatorAttachment) (*domain.WebAuthNRegistrationDetails, error) {
	if err := authz.UserIDInCTX(ctx, userID); err != nil {
		return nil, err
	}
	return c.registerUserPasskey(ctx, userID, resourceOwner, rpID, authenticator)
}

// RegisterUserPasskeyWithCode registers a new passkey for a unauthenticated user id.
// The resource is protected by the code, identified by the codeID.
func (c *Commands) RegisterUserPasskeyWithCode(ctx context.Context, userID, resourceOwner string, authenticator domain.AuthenticatorAttachment, codeID, code, rpID string, alg crypto.EncryptionAlgorithm) (*domain.WebAuthNRegistrationDetails, error) {
	event, err := c.verifyUserPasskeyCode(ctx, userID, resourceOwner, codeID, code, alg)
	if err != nil {
		return nil, err
	}

	return c.registerUserPasskey(ctx, userID, resourceOwner, rpID, authenticator, event)
}

type eventCallback func(context.Context, *eventstore.Aggregate) eventstore.Command

// verifyUserPasskeyCode verifies a passkey code, identified by codeID and userID.
// A code can only be used once.
// Upon success an event callback is returned, which must be called after
// all other events for the current request are created.
// This prevents consuming a code when another error occurred after verification.
func (c *Commands) verifyUserPasskeyCode(ctx context.Context, userID, resourceOwner, codeID, code string, alg crypto.EncryptionAlgorithm) (eventCallback, error) {
	wm := NewHumanPasswordlessInitCodeWriteModel(userID, codeID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, wm)
	if err != nil {
		return nil, err
	}
	err = verifyEncryptedCode(ctx, c.eventstore.Filter, domain.SecretGeneratorTypePasswordlessInitCode, alg, wm.ChangeDate, wm.Expiration, wm.CryptoCode, code) //nolint:staticcheck
	if err != nil || wm.State != domain.PasswordlessInitCodeStateActive {
		c.verifyUserPasskeyCodeFailed(ctx, wm)
		return nil, zerrors.ThrowInvalidArgument(err, "COMMAND-Eeb2a", "Errors.User.Code.Invalid")
	}
	return func(ctx context.Context, userAgg *eventstore.Aggregate) eventstore.Command {
		return user.NewHumanPasswordlessInitCodeCheckSucceededEvent(ctx, userAgg, codeID)
	}, nil
}

func (c *Commands) verifyUserPasskeyCodeFailed(ctx context.Context, wm *HumanPasswordlessInitCodeWriteModel) {
	userAgg := UserAggregateFromWriteModel(&wm.WriteModel)
	_, err := c.eventstore.Push(ctx, user.NewHumanPasswordlessInitCodeCheckFailedEvent(ctx, userAgg, wm.CodeID))
	logging.WithFields("userID", userAgg.ID).OnError(err).Error("RegisterUserPasskeyWithCode push failed")
}

func (c *Commands) registerUserPasskey(ctx context.Context, userID, resourceOwner, rpID string, authenticator domain.AuthenticatorAttachment, events ...eventCallback) (*domain.WebAuthNRegistrationDetails, error) {
	wm, userAgg, webAuthN, err := c.createUserPasskey(ctx, userID, resourceOwner, rpID, authenticator)
	if err != nil {
		return nil, err
	}
	return c.pushUserPasskey(ctx, wm, userAgg, webAuthN, events...)
}

func (c *Commands) createUserPasskey(ctx context.Context, userID, resourceOwner, rpID string, authenticator domain.AuthenticatorAttachment) (*HumanWebAuthNWriteModel, *eventstore.Aggregate, *domain.WebAuthNToken, error) {
	passwordlessTokens, err := c.getHumanPasswordlessTokens(ctx, userID, resourceOwner)
	if err != nil {
		return nil, nil, nil, err
	}
	return c.addHumanWebAuthN(ctx, userID, resourceOwner, rpID, passwordlessTokens, authenticator, domain.UserVerificationRequirementRequired)
}

func (c *Commands) pushUserPasskey(ctx context.Context, wm *HumanWebAuthNWriteModel, userAgg *eventstore.Aggregate, webAuthN *domain.WebAuthNToken, events ...eventCallback) (*domain.WebAuthNRegistrationDetails, error) {
	cmds := make([]eventstore.Command, len(events)+1)
	cmds[0] = user.NewHumanPasswordlessAddedEvent(ctx, userAgg, wm.WebauthNTokenID, webAuthN.Challenge, webAuthN.RPID)
	for i, event := range events {
		cmds[i+1] = event(ctx, userAgg)
	}

	err := c.pushAppendAndReduce(ctx, wm, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.WebAuthNRegistrationDetails{
		ObjectDetails:                      writeModelToObjectDetails(&wm.WriteModel),
		ID:                                 wm.WebauthNTokenID,
		PublicKeyCredentialCreationOptions: webAuthN.CredentialCreationData,
	}, nil
}

// AddUserPasskeyCode generates a Passkey code and sends an email
// with the default generated URL (pointing to zitadel).
func (c *Commands) AddUserPasskeyCode(ctx context.Context, userID, resourceOwner string, alg crypto.EncryptionAlgorithm) (*domain.ObjectDetails, error) {
	details, err := c.addUserPasskeyCode(ctx, userID, resourceOwner, alg, "", false)
	if err != nil {
		return nil, err
	}
	return details.ObjectDetails, err
}

// AddUserPasskeyCodeURLTemplate generates a Passkey code and sends an email
// with the URL created from passed template string.
// The template is executed as a test, before pushing to the eventstore.
func (c *Commands) AddUserPasskeyCodeURLTemplate(ctx context.Context, userID, resourceOwner string, alg crypto.EncryptionAlgorithm, urlTmpl string) (*domain.ObjectDetails, error) {
	if err := domain.RenderPasskeyURLTemplate(io.Discard, urlTmpl, userID, resourceOwner, "codeID", "code"); err != nil {
		return nil, err
	}
	details, err := c.addUserPasskeyCode(ctx, userID, resourceOwner, alg, urlTmpl, false)
	if err != nil {
		return nil, err
	}
	return details.ObjectDetails, err
}

// AddUserPasskeyCodeReturn generates and returns a Passkey code.
// No email will be sent to the user.
func (c *Commands) AddUserPasskeyCodeReturn(ctx context.Context, userID, resourceOwner string, alg crypto.EncryptionAlgorithm) (*domain.PasskeyCodeDetails, error) {
	return c.addUserPasskeyCode(ctx, userID, resourceOwner, alg, "", true)
}

func (c *Commands) addUserPasskeyCode(ctx context.Context, userID, resourceOwner string, alg crypto.EncryptionAlgorithm, urlTmpl string, returnCode bool) (*domain.PasskeyCodeDetails, error) {
	codeID, err := id_generator.Next()
	if err != nil {
		return nil, err
	}
	code, err := c.newPasskeyCode(ctx, c.eventstore.Filter, alg)
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

func (c *Commands) newPasskeyCode(ctx context.Context, filter preparation.FilterToQueryReducer, alg crypto.EncryptionAlgorithm) (*EncryptedCode, error) {
	return c.newEncryptedCode(ctx, filter, domain.SecretGeneratorTypePasswordlessInitCode, alg)
}

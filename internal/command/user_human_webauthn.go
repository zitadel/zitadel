package command

import (
	"context"
	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/eventstore"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	usr_repo "github.com/caos/zitadel/internal/repository/user"
	"github.com/caos/zitadel/internal/telemetry/tracing"
)

func (c *Commands) getHumanU2FTokens(ctx context.Context, userID, resourceowner string) ([]*domain.WebAuthNToken, error) {
	tokenReadModel := NewHumanU2FTokensReadModel(userID, resourceowner)
	err := c.eventstore.FilterToQueryReducer(ctx, tokenReadModel)
	if err != nil {
		return nil, err
	}
	if tokenReadModel.UserState == domain.UserStateDeleted {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-4M0ds", "Errors.User.NotFound")
	}
	return readModelToU2FTokens(tokenReadModel), nil
}

func (c *Commands) getHumanPasswordlessTokens(ctx context.Context, userID, resourceowner string) ([]*domain.WebAuthNToken, error) {
	tokenReadModel := NewHumanPasswordlessTokensReadModel(userID, resourceowner)
	err := c.eventstore.FilterToQueryReducer(ctx, tokenReadModel)
	if err != nil {
		return nil, err
	}
	if tokenReadModel.UserState == domain.UserStateDeleted {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-Mv9sd", "Errors.User.NotFound")
	}
	return readModelToPasswordlessTokens(tokenReadModel), nil
}

func (c *Commands) getHumanU2FLogin(ctx context.Context, userID, authReqID, resourceowner string) (*domain.WebAuthNLogin, error) {
	tokenReadModel := NewHumanU2FLoginReadModel(userID, authReqID, resourceowner)
	err := c.eventstore.FilterToQueryReducer(ctx, tokenReadModel)
	if err != nil {
		return nil, err
	}
	if tokenReadModel.State == domain.UserStateDeleted {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-5m88U", "Errors.User.NotFound")
	}
	return &domain.WebAuthNLogin{
		Challenge: tokenReadModel.Challenge,
	}, nil
}

func (c *Commands) getHumanPasswordlessLogin(ctx context.Context, userID, authReqID, resourceowner string) (*domain.WebAuthNLogin, error) {
	tokenReadModel := NewHumanPasswordlessLoginReadModel(userID, authReqID, resourceowner)
	err := c.eventstore.FilterToQueryReducer(ctx, tokenReadModel)
	if err != nil {
		return nil, err
	}
	if tokenReadModel.State == domain.UserStateDeleted {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-fm84R", "Errors.User.NotFound")
	}
	return &domain.WebAuthNLogin{
		Challenge: tokenReadModel.Challenge,
	}, nil
}

func (c *Commands) HumanAddU2FSetup(ctx context.Context, userID, resourceowner string, isLoginUI bool) (*domain.WebAuthNToken, error) {
	u2fTokens, err := c.getHumanU2FTokens(ctx, userID, resourceowner)
	if err != nil {
		return nil, err
	}
	addWebAuthN, userAgg, webAuthN, err := c.addHumanWebAuthN(ctx, userID, resourceowner, isLoginUI, u2fTokens)
	if err != nil {
		return nil, err
	}

	events, err := c.eventstore.PushEvents(ctx, usr_repo.NewHumanU2FAddedEvent(ctx, userAgg, addWebAuthN.WebauthNTokenID, webAuthN.Challenge))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addWebAuthN, events...)
	if err != nil {
		return nil, err
	}

	createdWebAuthN := writeModelToWebAuthN(addWebAuthN)
	createdWebAuthN.CredentialCreationData = webAuthN.CredentialCreationData
	createdWebAuthN.AllowedCredentialIDs = webAuthN.AllowedCredentialIDs
	createdWebAuthN.UserVerification = webAuthN.UserVerification
	return createdWebAuthN, nil
}

func (c *Commands) HumanAddPasswordlessSetup(ctx context.Context, userID, resourceowner string, isLoginUI bool) (*domain.WebAuthNToken, error) {
	passwordlessTokens, err := c.getHumanPasswordlessTokens(ctx, userID, resourceowner)
	if err != nil {
		return nil, err
	}
	addWebAuthN, userAgg, webAuthN, err := c.addHumanWebAuthN(ctx, userID, resourceowner, isLoginUI, passwordlessTokens)
	if err != nil {
		return nil, err
	}

	events, err := c.eventstore.PushEvents(ctx, usr_repo.NewHumanPasswordlessAddedEvent(ctx, userAgg, addWebAuthN.WebauthNTokenID, webAuthN.Challenge))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(addWebAuthN, events...)
	if err != nil {
		return nil, err
	}

	createdWebAuthN := writeModelToWebAuthN(addWebAuthN)
	createdWebAuthN.CredentialCreationData = webAuthN.CredentialCreationData
	createdWebAuthN.AllowedCredentialIDs = webAuthN.AllowedCredentialIDs
	createdWebAuthN.UserVerification = webAuthN.UserVerification
	return createdWebAuthN, nil
}

func (c *Commands) addHumanWebAuthN(ctx context.Context, userID, resourceowner string, isLoginUI bool, tokens []*domain.WebAuthNToken) (*HumanWebAuthNWriteModel, *eventstore.Aggregate, *domain.WebAuthNToken, error) {
	if userID == "" || resourceowner == "" {
		return nil, nil, nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3M0od", "Errors.IDMissing")
	}
	user, err := c.getHuman(ctx, userID, resourceowner)
	if err != nil {
		return nil, nil, nil, err
	}
	org, err := c.getOrg(ctx, user.ResourceOwner)
	if err != nil {
		return nil, nil, nil, err
	}
	orgPolicy, err := c.getOrgIAMPolicy(ctx, org.AggregateID)
	if err != nil {
		return nil, nil, nil, err
	}
	accountName := domain.GenerateLoginName(user.GetUsername(), org.PrimaryDomain, orgPolicy.UserLoginMustBeDomain)
	if accountName == "" {
		accountName = user.EmailAddress
	}
	webAuthN, err := c.webauthn.BeginRegistration(user, accountName, domain.AuthenticatorAttachmentUnspecified, domain.UserVerificationRequirementDiscouraged, isLoginUI, tokens...)
	if err != nil {
		return nil, nil, nil, err
	}
	tokenID, err := c.idGenerator.Next()
	if err != nil {
		return nil, nil, nil, err
	}
	addWebAuthN, err := c.webauthNWriteModelByID(ctx, userID, tokenID, resourceowner)
	if err != nil {
		return nil, nil, nil, err
	}

	userAgg := UserAggregateFromWriteModel(&addWebAuthN.WriteModel)
	return addWebAuthN, userAgg, webAuthN, nil
}

func (c *Commands) HumanVerifyU2FSetup(ctx context.Context, userID, resourceowner, tokenName, userAgentID string, credentialData []byte) (*domain.ObjectDetails, error) {
	u2fTokens, err := c.getHumanU2FTokens(ctx, userID, resourceowner)
	if err != nil {
		return nil, err
	}
	userAgg, webAuthN, verifyWebAuthN, err := c.verifyHumanWebAuthN(ctx, userID, resourceowner, tokenName, userAgentID, credentialData, u2fTokens)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx,
		usr_repo.NewHumanU2FVerifiedEvent(
			ctx,
			userAgg,
			verifyWebAuthN.WebauthNTokenID,
			webAuthN.WebAuthNTokenName,
			webAuthN.AttestationType,
			webAuthN.KeyID,
			webAuthN.PublicKey,
			webAuthN.AAGUID,
			webAuthN.SignCount,
		),
	)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(verifyWebAuthN, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&verifyWebAuthN.WriteModel), nil
}

func (c *Commands) HumanHumanPasswordlessSetup(ctx context.Context, userID, resourceowner, tokenName, userAgentID string, credentialData []byte) (*domain.ObjectDetails, error) {
	u2fTokens, err := c.getHumanPasswordlessTokens(ctx, userID, resourceowner)
	if err != nil {
		return nil, err
	}
	userAgg, webAuthN, verifyWebAuthN, err := c.verifyHumanWebAuthN(ctx, userID, resourceowner, tokenName, userAgentID, credentialData, u2fTokens)
	if err != nil {
		return nil, err
	}

	pushedEvents, err := c.eventstore.PushEvents(ctx,
		usr_repo.NewHumanPasswordlessVerifiedEvent(
			ctx,
			userAgg,
			verifyWebAuthN.WebauthNTokenID,
			webAuthN.WebAuthNTokenName,
			webAuthN.AttestationType,
			webAuthN.KeyID,
			webAuthN.PublicKey,
			webAuthN.AAGUID,
			webAuthN.SignCount,
		),
	)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(verifyWebAuthN, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&verifyWebAuthN.WriteModel), nil
}

func (c *Commands) verifyHumanWebAuthN(ctx context.Context, userID, resourceowner, tokenName, userAgentID string, credentialData []byte, tokens []*domain.WebAuthNToken) (*eventstore.Aggregate, *domain.WebAuthNToken, *HumanWebAuthNWriteModel, error) {
	if userID == "" || resourceowner == "" {
		return nil, nil, nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3M0od", "Errors.IDMissing")
	}
	user, err := c.getHuman(ctx, userID, resourceowner)
	if err != nil {
		return nil, nil, nil, err
	}
	_, token := domain.GetTokenToVerify(tokens)
	webAuthN, err := c.webauthn.FinishRegistration(user, token, tokenName, credentialData, userAgentID != "")
	if err != nil {
		return nil, nil, nil, err
	}

	verifyWebAuthN, err := c.webauthNWriteModelByID(ctx, userID, token.WebAuthNTokenID, resourceowner)
	if err != nil {
		return nil, nil, nil, err
	}

	userAgg := UserAggregateFromWriteModel(&verifyWebAuthN.WriteModel)
	return userAgg, webAuthN, verifyWebAuthN, nil
}

func (c *Commands) HumanBeginU2FLogin(ctx context.Context, userID, resourceOwner string, authRequest *domain.AuthRequest, isLoginUI bool) (*domain.WebAuthNLogin, error) {
	u2fTokens, err := c.getHumanU2FTokens(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}

	userAgg, webAuthNLogin, err := c.beginWebAuthNLogin(ctx, userID, resourceOwner, u2fTokens, isLoginUI)
	if err != nil {
		return nil, err
	}

	_, err = c.eventstore.PushEvents(ctx,
		usr_repo.NewHumanU2FBeginLoginEvent(
			ctx,
			userAgg,
			webAuthNLogin.Challenge,
			authRequestDomainToAuthRequestInfo(authRequest),
		),
	)

	return webAuthNLogin, err
}

func (c *Commands) HumanBeginPasswordlessLogin(ctx context.Context, userID, resourceOwner string, authRequest *domain.AuthRequest, isLoginUI bool) (*domain.WebAuthNLogin, error) {
	u2fTokens, err := c.getHumanPasswordlessTokens(ctx, userID, resourceOwner)
	if err != nil {
		return nil, err
	}

	userAgg, webAuthNLogin, err := c.beginWebAuthNLogin(ctx, userID, resourceOwner, u2fTokens, isLoginUI)
	if err != nil {
		return nil, err
	}
	_, err = c.eventstore.PushEvents(ctx,
		usr_repo.NewHumanPasswordlessBeginLoginEvent(
			ctx,
			userAgg,
			webAuthNLogin.Challenge,
			authRequestDomainToAuthRequestInfo(authRequest),
		),
	)
	return webAuthNLogin, err
}

func (c *Commands) beginWebAuthNLogin(ctx context.Context, userID, resourceOwner string, tokens []*domain.WebAuthNToken, isLoginUI bool) (*eventstore.Aggregate, *domain.WebAuthNLogin, error) {
	if userID == "" {
		return nil, nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-hh8K9", "Errors.IDMissing")
	}

	human, err := c.getHuman(ctx, userID, resourceOwner)
	if err != nil {
		return nil, nil, err
	}
	webAuthNLogin, err := c.webauthn.BeginLogin(human, domain.UserVerificationRequirementDiscouraged, isLoginUI, tokens...)
	if err != nil {
		return nil, nil, err
	}

	writeModel, err := c.webauthNWriteModelByID(ctx, userID, "", resourceOwner)
	if err != nil {
		return nil, nil, err
	}
	userAgg := UserAggregateFromWriteModel(&writeModel.WriteModel)

	return userAgg, webAuthNLogin, nil
}

func (c *Commands) HumanFinishU2FLogin(ctx context.Context, userID, resourceOwner string, credentialData []byte, authRequest *domain.AuthRequest, isLoginUI bool) error {
	webAuthNLogin, err := c.getHumanU2FLogin(ctx, userID, authRequest.ID, resourceOwner)
	if err != nil {
		return err
	}
	u2fTokens, err := c.getHumanU2FTokens(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}

	userAgg, token, signCount, err := c.finishWebAuthNLogin(ctx, userID, resourceOwner, credentialData, webAuthNLogin, u2fTokens, isLoginUI)
	if err != nil {
		_, pushErr := c.eventstore.PushEvents(ctx,
			usr_repo.NewHumanU2FCheckFailedEvent(
				ctx,
				userAgg,
				authRequestDomainToAuthRequestInfo(authRequest),
			),
		)
		logging.Log("EVENT-33M9f").OnError(pushErr).WithField("userID", userID).Warn("could not push failed passwordless check event")
		return err
	}

	_, err = c.eventstore.PushEvents(ctx,
		usr_repo.NewHumanU2FCheckSucceededEvent(
			ctx,
			userAgg,
			authRequestDomainToAuthRequestInfo(authRequest),
		),
		usr_repo.NewHumanU2FSignCountChangedEvent(
			ctx,
			userAgg,
			token.WebAuthNTokenID,
			signCount,
		),
	)

	return err
}

func (c *Commands) HumanFinishPasswordlessLogin(ctx context.Context, userID, resourceOwner string, credentialData []byte, authRequest *domain.AuthRequest, isLoginUI bool) error {
	webAuthNLogin, err := c.getHumanPasswordlessLogin(ctx, userID, authRequest.ID, resourceOwner)
	if err != nil {
		return err
	}

	passwordlessTokens, err := c.getHumanPasswordlessTokens(ctx, userID, resourceOwner)
	if err != nil {
		return err
	}

	userAgg, token, signCount, err := c.finishWebAuthNLogin(ctx, userID, resourceOwner, credentialData, webAuthNLogin, passwordlessTokens, isLoginUI)
	if err != nil {
		_, pushErr := c.eventstore.PushEvents(ctx,
			usr_repo.NewHumanPasswordlessCheckFailedEvent(
				ctx,
				userAgg,
				authRequestDomainToAuthRequestInfo(authRequest),
			),
		)
		logging.Log("EVENT-33M9f").OnError(pushErr).WithField("userID", userID).Warn("could not push failed passwordless check event")
		return err
	}

	_, err = c.eventstore.PushEvents(ctx,
		usr_repo.NewHumanU2FCheckSucceededEvent(
			ctx,
			userAgg,
			authRequestDomainToAuthRequestInfo(authRequest),
		),
		usr_repo.NewHumanPasswordlessSignCountChangedEvent(
			ctx,
			userAgg,
			token.WebAuthNTokenID,
			signCount,
		),
	)
	return err
}

func (c *Commands) finishWebAuthNLogin(ctx context.Context, userID, resourceOwner string, credentialData []byte, webAuthN *domain.WebAuthNLogin, tokens []*domain.WebAuthNToken, isLoginUI bool) (*eventstore.Aggregate, *domain.WebAuthNToken, uint32, error) {
	if userID == "" {
		return nil, nil, 0, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-hh8K9", "Errors.IDMissing")
	}

	human, err := c.getHuman(ctx, userID, resourceOwner)
	if err != nil {
		return nil, nil, 0, err
	}
	keyID, signCount, err := c.webauthn.FinishLogin(human, webAuthN, credentialData, isLoginUI, tokens...)
	if err != nil && keyID == nil {
		return nil, nil, 0, err
	}

	_, token := domain.GetTokenByKeyID(tokens, keyID)
	if token == nil {
		return nil, nil, 0, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3b7zs", "Errors.User.WebAuthN.NotFound")
	}

	writeModel, err := c.webauthNWriteModelByID(ctx, userID, "", resourceOwner)
	if err != nil {
		return nil, nil, 0, err
	}
	userAgg := UserAggregateFromWriteModel(&writeModel.WriteModel)

	return userAgg, token, signCount, nil
}

func (c *Commands) HumanRemoveU2F(ctx context.Context, userID, webAuthNID, resourceOwner string) (*domain.ObjectDetails, error) {
	event := usr_repo.PrepareHumanU2FRemovedEvent(ctx, webAuthNID)
	return c.removeHumanWebAuthN(ctx, userID, webAuthNID, resourceOwner, event)
}

func (c *Commands) HumanRemovePasswordless(ctx context.Context, userID, webAuthNID, resourceOwner string) (*domain.ObjectDetails, error) {
	event := usr_repo.PrepareHumanPasswordlessRemovedEvent(ctx, webAuthNID)
	return c.removeHumanWebAuthN(ctx, userID, webAuthNID, resourceOwner, event)
}

func (c *Commands) removeHumanWebAuthN(ctx context.Context, userID, webAuthNID, resourceOwner string, preparedEvent func(*eventstore.Aggregate) eventstore.EventPusher) (*domain.ObjectDetails, error) {
	if userID == "" || webAuthNID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-6M9de", "Errors.IDMissing")
	}

	existingWebAuthN, err := c.webauthNWriteModelByID(ctx, userID, webAuthNID, resourceOwner)
	if err != nil {
		return nil, err
	}
	if existingWebAuthN.State == domain.MFAStateUnspecified || existingWebAuthN.State == domain.MFAStateRemoved {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-2M9ds", "Errors.User.ExternalIDP.NotFound")
	}

	userAgg := UserAggregateFromWriteModel(&existingWebAuthN.WriteModel)
	pushedEvents, err := c.eventstore.PushEvents(ctx, preparedEvent(userAgg))
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(existingWebAuthN, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&existingWebAuthN.WriteModel), nil
}

func (c *Commands) webauthNWriteModelByID(ctx context.Context, userID, webAuthNID, resourceOwner string) (writeModel *HumanWebAuthNWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanWebAuthNWriteModel(userID, webAuthNID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

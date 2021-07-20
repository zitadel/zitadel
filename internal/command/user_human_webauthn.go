package command

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
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
	if tokenReadModel.State == domain.UserStateUnspecified || tokenReadModel.State == domain.UserStateDeleted {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-5m88U", "Errors.User.NotFound")
	}
	return &domain.WebAuthNLogin{
		ObjectRoot: models.ObjectRoot{
			AggregateID: tokenReadModel.AggregateID,
		},
		Challenge:            tokenReadModel.Challenge,
		AllowedCredentialIDs: tokenReadModel.AllowedCredentialIDs,
		UserVerification:     tokenReadModel.UserVerification,
	}, nil
}

func (c *Commands) getHumanPasswordlessLogin(ctx context.Context, userID, authReqID, resourceowner string) (*domain.WebAuthNLogin, error) {
	tokenReadModel := NewHumanPasswordlessLoginReadModel(userID, authReqID, resourceowner)
	err := c.eventstore.FilterToQueryReducer(ctx, tokenReadModel)
	if err != nil {
		return nil, err
	}
	if tokenReadModel.State == domain.UserStateUnspecified || tokenReadModel.State == domain.UserStateDeleted {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-fm84R", "Errors.User.NotFound")
	}
	return &domain.WebAuthNLogin{
		ObjectRoot: models.ObjectRoot{
			AggregateID: tokenReadModel.AggregateID,
		},
		Challenge:            tokenReadModel.Challenge,
		AllowedCredentialIDs: tokenReadModel.AllowedCredentialIDs,
		UserVerification:     tokenReadModel.UserVerification,
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

func (c *Commands) HumanAddPasswordlessSetupInitCode(ctx context.Context, userID, resourceowner, codeID, verificationCode string) (*domain.WebAuthNToken, error) {
	err := c.humanVerifyPasswordlessInitCode(ctx, userID, resourceowner, codeID, verificationCode)
	if err != nil {
		return nil, err
	}
	return c.HumanAddPasswordlessSetup(ctx, userID, resourceowner, true)
}

func (c *Commands) addHumanWebAuthN(ctx context.Context, userID, resourceowner string, isLoginUI bool, tokens []*domain.WebAuthNToken) (*HumanWebAuthNWriteModel, *eventstore.Aggregate, *domain.WebAuthNToken, error) {
	if userID == "" {
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
			userAgentID,
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

func (c *Commands) HumanPasswordlessSetupInitCode(ctx context.Context, userID, resourceowner, tokenName, userAgentID, codeID, verificationCode string, credentialData []byte) (*domain.ObjectDetails, error) {
	err := c.humanVerifyPasswordlessInitCode(ctx, userID, resourceowner, codeID, verificationCode)
	if err != nil {
		return nil, err
	}
	succeededEvent := func(userAgg *eventstore.Aggregate) *usr_repo.HumanPasswordlessInitCodeCheckSucceededEvent {
		return usr_repo.NewHumanPasswordlessInitCodeCheckSucceededEvent(ctx, userAgg, codeID)
	}
	return c.humanHumanPasswordlessSetup(ctx, userID, resourceowner, tokenName, userAgentID, credentialData, succeededEvent)
}

func (c *Commands) HumanHumanPasswordlessSetup(ctx context.Context, userID, resourceowner, tokenName, userAgentID string, credentialData []byte) (*domain.ObjectDetails, error) {
	return c.humanHumanPasswordlessSetup(ctx, userID, resourceowner, tokenName, userAgentID, credentialData, nil)
}

func (c *Commands) humanHumanPasswordlessSetup(ctx context.Context, userID, resourceowner, tokenName, userAgentID string, credentialData []byte,
	codeCheckEvent func(*eventstore.Aggregate) *usr_repo.HumanPasswordlessInitCodeCheckSucceededEvent) (*domain.ObjectDetails, error) {

	u2fTokens, err := c.getHumanPasswordlessTokens(ctx, userID, resourceowner)
	if err != nil {
		return nil, err
	}
	userAgg, webAuthN, verifyWebAuthN, err := c.verifyHumanWebAuthN(ctx, userID, resourceowner, tokenName, userAgentID, credentialData, u2fTokens)
	if err != nil {
		return nil, err
	}

	events := []eventstore.EventPusher{
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
			userAgentID,
		),
	}
	if codeCheckEvent != nil {
		events = append(events, codeCheckEvent(userAgg))
	}
	pushedEvents, err := c.eventstore.PushEvents(ctx, events...)
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
	if userID == "" {
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
			webAuthNLogin.AllowedCredentialIDs,
			webAuthNLogin.UserVerification,
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
			webAuthNLogin.AllowedCredentialIDs,
			webAuthNLogin.UserVerification,
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
		if userAgg == nil {
			logging.LogWithFields("EVENT-Addqd", "userID", userID, "resourceOwner", resourceOwner).WithError(err).Warn("missing userAggregate for pushing failed u2f check event")
			return err
		}
		_, pushErr := c.eventstore.PushEvents(ctx,
			usr_repo.NewHumanU2FCheckFailedEvent(
				ctx,
				userAgg,
				authRequestDomainToAuthRequestInfo(authRequest),
			),
		)
		logging.LogWithFields("EVENT-Bdgd2", "userID", userID, "resourceOwner", resourceOwner).OnError(pushErr).Warn("could not push failed u2f check event")
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
		if userAgg == nil {
			logging.LogWithFields("EVENT-Dbbbw", "userID", userID, "resourceOwner", resourceOwner).WithError(err).Warn("missing userAggregate for pushing failed passwordless check event")
			return err
		}
		_, pushErr := c.eventstore.PushEvents(ctx,
			usr_repo.NewHumanPasswordlessCheckFailedEvent(
				ctx,
				userAgg,
				authRequestDomainToAuthRequestInfo(authRequest),
			),
		)
		logging.LogWithFields("EVENT-33M9f", "userID", userID, "resourceOwner", resourceOwner).OnError(pushErr).Warn("could not push failed passwordless check event")
		return err
	}

	_, err = c.eventstore.PushEvents(ctx,
		usr_repo.NewHumanPasswordlessCheckSucceededEvent(
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

func (c *Commands) HumanAddPasswordlessInitCode(ctx context.Context, userID, resourceOwner string) (*domain.PasswordlessInitCode, error) {
	return c.humanAddPasswordlessInitCode(ctx, userID, resourceOwner, false)
}

func (c *Commands) HumanSendPasswordlessInitCode(ctx context.Context, userID, resourceOwner string) (*domain.PasswordlessInitCode, error) {
	return c.humanAddPasswordlessInitCode(ctx, userID, resourceOwner, true)
}

func (c *Commands) humanAddPasswordlessInitCode(ctx context.Context, userID, resourceOwner string, send bool) (*domain.PasswordlessInitCode, error) {
	if userID == "" {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-GVfg3", "Errors.IDMissing")
	}

	codeID, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}
	initCode := NewHumanPasswordlessInitCodeWriteModel(userID, codeID, resourceOwner)
	err = c.eventstore.FilterToQueryReducer(ctx, initCode)
	if err != nil {
		return nil, err
	}

	cryptoCode, code, err := crypto.NewCode(c.passwordlessInitCode)
	if err != nil {
		return nil, err
	}
	codeEvent := usr_repo.NewHumanPasswordlessInitCodeAddedEvent(ctx,
		UserAggregateFromWriteModel(&initCode.WriteModel),
		codeID,
		cryptoCode,
		c.passwordlessInitCode.Expiry(),
		send,
	)
	pushedEvents, err := c.eventstore.PushEvents(ctx, codeEvent)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(initCode, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToPasswordlessInitCode(initCode, code), nil
}

func (c *Commands) HumanPasswordlessInitCodeSent(ctx context.Context, userID, resourceOwner, codeID string) error {
	if userID == "" || codeID == "" {
		return caos_errs.ThrowInvalidArgument(nil, "COMMAND-ADggh", "Errors.IDMissing")
	}
	initCode := NewHumanPasswordlessInitCodeWriteModel(userID, codeID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, initCode)
	if err != nil {
		return err
	}

	//if !isUserStateExists(existingUser.UserState) {
	//	return caos_errs.ThrowNotFound(nil, "COMMAND-5m9gK", "Errors.User.NotFound")
	//}

	_, err = c.eventstore.PushEvents(ctx,
		usr_repo.NewHumanPasswordlessInitCodeSentEvent(ctx, UserAggregateFromWriteModel(&initCode.WriteModel), codeID),
	)
	return err
}

func (c *Commands) humanVerifyPasswordlessInitCode(ctx context.Context, userID, resourceOwner, codeID, verificationCode string) error {
	if userID == "" || codeID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-GVfg3", "Errors.IDMissing")
	}
	initCode := NewHumanPasswordlessInitCodeWriteModel(userID, codeID, resourceOwner)
	err := c.eventstore.FilterToQueryReducer(ctx, initCode)
	if err != nil {
		return err
	}
	if err = crypto.VerifyCode(initCode.ChangeDate, initCode.Expiration, initCode.CryptoCode, verificationCode, c.passwordlessInitCode); err != nil || !initCode.Active {
		userAgg := UserAggregateFromWriteModel(&initCode.WriteModel)
		_, err = c.eventstore.PushEvents(ctx, usr_repo.NewHumanPasswordlessInitCodeCheckFailedEvent(ctx, userAgg, codeID))
		logging.LogWithFields("COMMAND-Gkuud", "userID", userAgg.ID).OnError(err).Error("NewHumanPasswordlessInitCodeCheckFailedEvent push failed")
		return caos_errs.ThrowInvalidArgument(err, "COMMAND-Dhz8i", "Errors.User.Code.Invalid")
	}
	return nil
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
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-DAfb2", "Errors.User.WebAuthN.NotFound")
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

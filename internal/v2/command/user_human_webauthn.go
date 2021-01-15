package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	"github.com/caos/zitadel/internal/v2/domain"
	usr_repo "github.com/caos/zitadel/internal/v2/repository/user"
)

func (r *CommandSide) getHumanU2FTokens(ctx context.Context, userID, resourceowner string) ([]*domain.WebAuthNToken, error) {
	tokenReadModel := NewHumanU2FTokensReadModel(userID, resourceowner)
	err := r.eventstore.FilterToQueryReducer(ctx, tokenReadModel)
	if err != nil {
		return nil, err
	}
	if tokenReadModel.UserState == domain.UserStateDeleted {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-4M0ds", "Errors.User.NotFound")
	}
	return readModelToU2FTokens(tokenReadModel), nil
}

func (r *CommandSide) getHumanPasswordlessTokens(ctx context.Context, userID, resourceowner string) ([]*domain.WebAuthNToken, error) {
	tokenReadModel := NewHumanPasswordlessTokensReadModel(userID, resourceowner)
	err := r.eventstore.FilterToQueryReducer(ctx, tokenReadModel)
	if err != nil {
		return nil, err
	}
	if tokenReadModel.UserState == domain.UserStateDeleted {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-Mv9sd", "Errors.User.NotFound")
	}
	return readModelToPasswordlessTokens(tokenReadModel), nil
}

func (r *CommandSide) AddHumanU2F(ctx context.Context, userID, resourceowner string, isLoginUI bool) (*domain.WebAuthNToken, error) {
	u2fTokens, err := r.getHumanU2FTokens(ctx, userID, resourceowner)
	if err != nil {
		return nil, err
	}
	addWebAuthN, userAgg, webAuthN, err := r.addHumanWebAuthN(ctx, userID, resourceowner, isLoginUI, u2fTokens)
	if err != nil {
		return nil, err
	}
	userAgg.PushEvents(usr_repo.NewHumanU2FAddedEvent(ctx, addWebAuthN.WebauthNTokenID, webAuthN.Challenge))

	err = r.eventstore.PushAggregate(ctx, addWebAuthN, userAgg)
	if err != nil {
		return nil, err
	}
	createdWebAuthN := writeModelToWebAuthN(addWebAuthN)
	createdWebAuthN.CredentialCreationData = webAuthN.CredentialCreationData
	createdWebAuthN.AllowedCredentialIDs = webAuthN.AllowedCredentialIDs
	createdWebAuthN.UserVerification = webAuthN.UserVerification
	return createdWebAuthN, nil
}

func (r *CommandSide) AddHumanPasswordless(ctx context.Context, userID, resourceowner string, isLoginUI bool) (*domain.WebAuthNToken, error) {
	passwordlessTokens, err := r.getHumanPasswordlessTokens(ctx, userID, resourceowner)
	if err != nil {
		return nil, err
	}
	addWebAuthN, userAgg, webAuthN, err := r.addHumanWebAuthN(ctx, userID, resourceowner, isLoginUI, passwordlessTokens)
	if err != nil {
		return nil, err
	}
	userAgg.PushEvents(usr_repo.NewHumanU2FAddedEvent(ctx, addWebAuthN.WebauthNTokenID, webAuthN.Challenge))

	err = r.eventstore.PushAggregate(ctx, addWebAuthN, userAgg)
	if err != nil {
		return nil, err
	}
	createdWebAuthN := writeModelToWebAuthN(addWebAuthN)
	createdWebAuthN.CredentialCreationData = webAuthN.CredentialCreationData
	createdWebAuthN.AllowedCredentialIDs = webAuthN.AllowedCredentialIDs
	createdWebAuthN.UserVerification = webAuthN.UserVerification
	return createdWebAuthN, nil
}

func (r *CommandSide) addHumanWebAuthN(ctx context.Context, userID, resourceowner string, isLoginUI bool, tokens []*domain.WebAuthNToken) (*HumanWebAuthNWriteModel, *usr_repo.Aggregate, *domain.WebAuthNToken, error) {
	if userID == "" || resourceowner == "" {
		return nil, nil, nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3M0od", "Errors.IDMissing")
	}
	user, err := r.getHuman(ctx, userID, resourceowner)
	if err != nil {
		return nil, nil, nil, err
	}
	org, err := r.getOrg(ctx, user.ResourceOwner)
	if err != nil {
		return nil, nil, nil, err
	}
	orgPolicy, err := r.getOrgIAMPolicy(ctx, org.AggregateID)
	if err != nil {
		return nil, nil, nil, err
	}
	accountName := domain.GenerateLoginName(user.GetUsername(), org.PrimaryDomain, orgPolicy.UserLoginMustBeDomain)
	if accountName == "" {
		accountName = user.EmailAddress
	}
	webAuthN, err := r.webauthn.BeginRegistration(user, accountName, domain.AuthenticatorAttachmentUnspecified, domain.UserVerificationRequirementDiscouraged, isLoginUI, tokens...)
	if err != nil {
		return nil, nil, nil, err
	}
	tokenID, err := r.idGenerator.Next()
	if err != nil {
		return nil, nil, nil, err
	}
	addWebAuthN, err := r.webauthNWriteModelByID(ctx, userID, tokenID, resourceowner)
	if err != nil {
		return nil, nil, nil, err
	}

	userAgg := UserAggregateFromWriteModel(&addWebAuthN.WriteModel)
	return addWebAuthN, userAgg, webAuthN, nil
}

func (r *CommandSide) VerifyHumanU2F(ctx context.Context, userID, resourceowner, tokenName, userAgentID string, credentialData []byte) error {
	u2fTokens, err := r.getHumanU2FTokens(ctx, userID, resourceowner)
	if err != nil {
		return err
	}
	verifyWebAuthN, userAgg, webAuthN, err := r.verifyHumanWebAuthN(ctx, userID, resourceowner, tokenName, userAgentID, credentialData, u2fTokens)
	if err != nil {
		return err
	}
	userAgg.PushEvents(
		usr_repo.NewHumanU2FVerifiedEvent(
			ctx,
			verifyWebAuthN.WebauthNTokenID,
			webAuthN.WebAuthNTokenName,
			webAuthN.AttestationType,
			webAuthN.KeyID,
			webAuthN.PublicKey,
			webAuthN.AAGUID,
			webAuthN.SignCount,
		),
	)

	return r.eventstore.PushAggregate(ctx, verifyWebAuthN, userAgg)
}

func (r *CommandSide) VerifyHumanPasswordless(ctx context.Context, userID, resourceowner, tokenName, userAgentID string, credentialData []byte) error {
	u2fTokens, err := r.getHumanPasswordlessTokens(ctx, userID, resourceowner)
	if err != nil {
		return err
	}
	verifyWebAuthN, userAgg, webAuthN, err := r.verifyHumanWebAuthN(ctx, userID, resourceowner, tokenName, userAgentID, credentialData, u2fTokens)
	if err != nil {
		return err
	}
	userAgg.PushEvents(
		usr_repo.NewHumanU2FVerifiedEvent(
			ctx,
			verifyWebAuthN.WebauthNTokenID,
			webAuthN.WebAuthNTokenName,
			webAuthN.AttestationType,
			webAuthN.KeyID,
			webAuthN.PublicKey,
			webAuthN.AAGUID,
			webAuthN.SignCount,
		),
	)
	return r.eventstore.PushAggregate(ctx, verifyWebAuthN, userAgg)
}

func (r *CommandSide) verifyHumanWebAuthN(ctx context.Context, userID, resourceowner, tokenName, userAgentID string, credentialData []byte, tokens []*domain.WebAuthNToken) (*HumanWebAuthNWriteModel, *usr_repo.Aggregate, *domain.WebAuthNToken, error) {
	if userID == "" || resourceowner == "" {
		return nil, nil, nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3M0od", "Errors.IDMissing")
	}
	user, err := r.getHuman(ctx, userID, resourceowner)
	if err != nil {
		return nil, nil, nil, err
	}
	_, token := domain.GetTokenToVerify(tokens)
	webAuthN, err := r.webauthn.FinishRegistration(user, token, tokenName, credentialData, userAgentID != "")
	if err != nil {
		return nil, nil, nil, err
	}

	verifyWebAuthN, err := r.webauthNWriteModelByID(ctx, userID, token.WebAuthNTokenID, resourceowner)
	if err != nil {
		return nil, nil, nil, err
	}

	userAgg := UserAggregateFromWriteModel(&verifyWebAuthN.WriteModel)
	return verifyWebAuthN, userAgg, webAuthN, nil
}

func (r *CommandSide) RemoveHumanU2F(ctx context.Context, userID, webAuthNID, resourceOwner string) error {
	event := usr_repo.NewHumanU2FRemovedEvent(ctx, webAuthNID)
	return r.removeHumanWebAuthN(ctx, userID, webAuthNID, resourceOwner, event)
}

func (r *CommandSide) RemoveHumanPasswordless(ctx context.Context, userID, webAuthNID, resourceOwner string) error {
	event := usr_repo.NewHumanPasswordlessRemovedEvent(ctx, webAuthNID)
	return r.removeHumanWebAuthN(ctx, userID, webAuthNID, resourceOwner, event)
}

func (r *CommandSide) removeHumanWebAuthN(ctx context.Context, userID, webAuthNID, resourceOwner string, event eventstore.EventPusher) error {
	if userID == "" || webAuthNID == "" {
		return caos_errs.ThrowPreconditionFailed(nil, "COMMAND-6M9de", "Errors.IDMissing")
	}

	existingWebAuthN, err := r.webauthNWriteModelByID(ctx, userID, webAuthNID, resourceOwner)
	if err != nil {
		return err
	}
	if existingWebAuthN.State == domain.MFAStateUnspecified || existingWebAuthN.State == domain.MFAStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-2M9ds", "Errors.User.ExternalIDP.NotFound")
	}
	userAgg := UserAggregateFromWriteModel(&existingWebAuthN.WriteModel)
	userAgg.PushEvents(event)

	return r.eventstore.PushAggregate(ctx, existingWebAuthN, userAgg)
}

func (r *CommandSide) webauthNWriteModelByID(ctx context.Context, userID, webAuthNID, resourceOwner string) (writeModel *HumanWebAuthNWriteModel, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	writeModel = NewHumanWebAuthNWriteModel(userID, webAuthNID, resourceOwner)
	err = r.eventstore.FilterToQueryReducer(ctx, writeModel)
	if err != nil {
		return nil, err
	}
	return writeModel, nil
}

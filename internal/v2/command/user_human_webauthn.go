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
	return nil, nil
}

func (r *CommandSide) AddHumanU2F(ctx context.Context, userID, resourceowner string, isLoginUI bool) (*domain.WebAuthNToken, error) {

	addWebAuthN, userAgg, webAuthN, err := r.addHumanWebAuthN(ctx, userID, resourceowner, isLoginUI)
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
	addWebAuthN, userAgg, webAuthN, err := r.addHumanWebAuthN(ctx, userID, resourceowner, isLoginUI)
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

func (r *CommandSide) addHumanWebAuthN(ctx context.Context, userID, resourceowner string, isLoginUI bool) (*HumanWebAuthNWriteModel, *usr_repo.Aggregate, *domain.WebAuthNToken, error) {
	if userID == "" || resourceowner == "" {
		return nil, nil, nil, caos_errs.ThrowPreconditionFailed(nil, "COMMAND-3M0od", "Errors.IDMissing")
	}

	user, err := r.getUser(ctx, userID, resourceowner)
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
	accountName := domain.GenerateLoginName(user.UserName, org.PrimaryDomain, orgPolicy.UserLoginMustBeDomain)
	if accountName == "" {
		accountName = user.EmailAddress
	}
	webAuthN, err := r.webauthn.BeginRegistration(user, accountName, domain.AuthenticatorAttachmentUnspecified, domain.UserVerificationRequirementDiscouraged, isLoginUI, user.U2FTokens...)
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
	if existingWebAuthN.State == domain.WebAuthNStateUnspecified || existingWebAuthN.State == domain.WebAuthNStateRemoved {
		return caos_errs.ThrowNotFound(nil, "COMMAND-5M0ds", "Errors.User.ExternalIDP.NotFound")
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

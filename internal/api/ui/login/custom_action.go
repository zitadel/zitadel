package login

import (
	"context"

	"github.com/zitadel/oidc/v2/pkg/oidc"

	"github.com/zitadel/zitadel/internal/actions"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	iam_model "github.com/zitadel/zitadel/internal/iam/model"
)

func (l *Login) customExternalUserMapping(ctx context.Context, user *domain.ExternalUser, tokens *oidc.Tokens, req *domain.AuthRequest, config *iam_model.IDPConfigView) (*domain.ExternalUser, error) {
	resourceOwner := req.RequestedOrgID
	if resourceOwner == "" {
		resourceOwner = config.AggregateID
	}
	instance := authz.GetInstance(ctx)
	if resourceOwner == instance.InstanceID() {
		resourceOwner = instance.DefaultOrganisationID()
	}
	triggerActions, err := l.query.GetActiveActionsByFlowAndTriggerType(ctx, domain.FlowTypeExternalAuthentication, domain.TriggerTypePostAuthentication, resourceOwner)
	if err != nil {
		return nil, err
	}
	actionCtx := (&actions.Context{}).SetToken(tokens)
	api := (&actions.API{}).SetExternalUser(user).SetMetadata(&user.Metadatas)
	for _, a := range triggerActions {
		err = actions.Run(actionCtx, api, a.Script, a.Name, a.Options()...)
		if err != nil {
			return nil, err
		}
	}
	return user, err
}

func (l *Login) customExternalUserToLoginUserMapping(ctx context.Context, user *domain.Human, tokens *oidc.Tokens, req *domain.AuthRequest, config *iam_model.IDPConfigView, metadata []*domain.Metadata, resourceOwner string) (*domain.Human, []*domain.Metadata, error) {
	triggerActions, err := l.query.GetActiveActionsByFlowAndTriggerType(ctx, domain.FlowTypeExternalAuthentication, domain.TriggerTypePreCreation, resourceOwner)
	if err != nil {
		return nil, nil, err
	}
	actionCtx := (&actions.Context{}).SetToken(tokens)
	api := (&actions.API{}).SetHuman(user).SetMetadata(&metadata)
	for _, a := range triggerActions {
		err = actions.Run(actionCtx, api, a.Script, a.Name, a.Options()...)
		if err != nil {
			return nil, nil, err
		}
	}
	return user, metadata, err
}

func (l *Login) customGrants(ctx context.Context, userID string, tokens *oidc.Tokens, req *domain.AuthRequest, config *iam_model.IDPConfigView, resourceOwner string) ([]*domain.UserGrant, error) {
	triggerActions, err := l.query.GetActiveActionsByFlowAndTriggerType(ctx, domain.FlowTypeExternalAuthentication, domain.TriggerTypePostCreation, resourceOwner)
	if err != nil {
		return nil, err
	}
	actionCtx := (&actions.Context{}).SetToken(tokens)
	actionUserGrants := make([]actions.UserGrant, 0)
	api := (&actions.API{}).SetUserGrants(&actionUserGrants)
	for _, a := range triggerActions {
		err = actions.Run(actionCtx, api, a.Script, a.Name, a.Options()...)
		if err != nil {
			return nil, err
		}
	}
	return actionUserGrantsToDomain(userID, actionUserGrants), err
}

func actionUserGrantsToDomain(userID string, actionUserGrants []actions.UserGrant) []*domain.UserGrant {
	if actionUserGrants == nil {
		return nil
	}
	userGrants := make([]*domain.UserGrant, len(actionUserGrants))
	for i, grant := range actionUserGrants {
		userGrants[i] = &domain.UserGrant{
			UserID:         userID,
			ProjectID:      grant.ProjectID,
			ProjectGrantID: grant.ProjectGrantID,
			RoleKeys:       grant.Roles,
		}
	}
	return userGrants
}

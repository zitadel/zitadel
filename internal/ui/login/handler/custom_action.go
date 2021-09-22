package handler

import (
	"context"

	"github.com/caos/oidc/pkg/oidc"
	"github.com/caos/zitadel/internal/actions"
	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

func (l *Login) customExternalUserMapping(user *domain.ExternalUser, tokens *oidc.Tokens, req *domain.AuthRequest, config *iam_model.IDPConfigView) (*domain.ExternalUser, error) {
	triggerActions, err := l.query.GetActionsByFlowAndTriggerType(context.TODO(), domain.FlowTypeExternalAuthentication, domain.TriggerTypePostAuthentication)
	if err != nil {
		return nil, err
	}
	ctx := (&actions.Context{}).SetToken(tokens)
	api := (&actions.API{}).SetExternalUser(user).SetMetadata(&user.Metadatas)
	for _, a := range triggerActions {
		err = actions.Run(ctx, api, a.Script, a.Name, a.Timeout, a.AllowedToFail)
		if err != nil {
			return nil, err
		}
	}
	return user, err
}

func (l *Login) customExternalUserToLoginUserMapping(user *domain.Human, tokens *oidc.Tokens, req *domain.AuthRequest, config *iam_model.IDPConfigView, metadata []*domain.Metadata) (*domain.Human, []*domain.Metadata, error) {
	triggerActions, err := l.query.GetActionsByFlowAndTriggerType(context.TODO(), domain.FlowTypeExternalAuthentication, domain.TriggerTypePreCreation)
	if err != nil {
		return nil, nil, err
	}
	ctx := (&actions.Context{}).SetToken(tokens)
	api := (&actions.API{}).SetHuman(user).SetMetadata(&metadata)
	for _, a := range triggerActions {
		err = actions.Run(ctx, api, a.Script, a.Name, a.Timeout, a.AllowedToFail)
		if err != nil {
			return nil, nil, err
		}
	}
	return user, metadata, err
}

func (l *Login) customGrants(userID string, tokens *oidc.Tokens, req *domain.AuthRequest, config *iam_model.IDPConfigView) ([]*domain.UserGrant, error) {
	triggerActions, err := l.query.GetActionsByFlowAndTriggerType(context.TODO(), domain.FlowTypeExternalAuthentication, domain.TriggerTypePostCreation)
	if err != nil {
		return nil, err
	}
	ctx := (&actions.Context{}).SetToken(tokens)
	actionUserGrants := make([]actions.UserGrant, 0)
	api := (&actions.API{}).SetUserGrants(&actionUserGrants)
	for _, a := range triggerActions {
		err = actions.Run(ctx, api, a.Script, a.Name, a.Timeout, a.AllowedToFail)
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

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

//func (l *Login) preRedirectAction(ctx context.Context, request *domain.AuthRequest, user *user_model.UserView) error {
//	triggerActions, err := l.query.GetActionsByTriggerType(ctx, domain.TriggerTypePreRedirect)
//	if err != nil {
//		return err
//	}
//	c := &actions.Context{User: user}
//	list := make([]actions.UserGrant, 0)
//	for _, a := range triggerActions {
//		err = actions.Run(c, a.Script, a.Name, a.Timeout, a.AllowedToFail, actions.SetUser(user), actions.Appender(&list))
//		if err != nil {
//			return err
//		}
//	}
//	for _, grant := range list {
//		_, err = l.command.AddUserGrant(ctx, &domain.UserGrant{
//			UserID:         request.UserID,
//			ProjectID:      grant.ProjectID,
//			ProjectGrantID: grant.ProjectGrantID,
//			RoleKeys:       grant.Roles,
//		}, request.UserOrgID)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

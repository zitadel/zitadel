package handler

import (
	"context"

	"github.com/caos/oidc/pkg/oidc"
	"github.com/caos/zitadel/internal/actions"
	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

func (l *Login) customMapping(user *domain.ExternalUser, tokens *oidc.Tokens, req *domain.AuthRequest, config *iam_model.IDPConfigView) (*domain.ExternalUser, error) {
	triggerActions, err := l.query.GetActionsByFlowAndTriggerType(context.TODO(), domain.FlowTypeExternalAuthentication, domain.TriggerTypePostAuthentication)
	if err != nil {
		return nil, err
	}
	ctx := &actions.Context{ExternalUser: user, Tokens: tokens}
	apiuser := *user
	for _, a := range triggerActions {
		err = actions.Run(ctx, a.Script, a.Name, a.Timeout, a.AllowedToFail, actions.SetExternalUser(&apiuser))
		if err != nil {
			return nil, err
		}
	}
	return &apiuser, err
}

func (l *Login) customRegistrationMapping(user *domain.Human, tokens *oidc.Tokens, req *domain.AuthRequest, config *iam_model.IDPConfigView) (*domain.Human, []*domain.Metadata, error) {
	triggerActions, err := l.query.GetActionsByFlowAndTriggerType(context.TODO(), domain.FlowTypeExternalRegistration, domain.TriggerTypePostAuthentication)
	if err != nil {
		return nil, nil, err
	}
	t := *tokens
	ctx := &actions.Context{User: user, Tokens: &t}

	//apiuser := actions.NewUser(user.FirstName, user.LastName)
	apiuser := func(firstname string) {
		user.FirstName = firstname
	}
	metadata := make([]*domain.Metadata, 0)
	for _, a := range triggerActions {
		err = actions.Run(ctx, a.Script, a.Name, a.Timeout, a.AllowedToFail, actions.SetUser(apiuser), actions.SetMetadata(metadata))
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

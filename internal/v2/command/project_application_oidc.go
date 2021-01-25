package command

import (
	"context"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/project"
)

func (r *CommandSide) AddOIDCApplication(ctx context.Context, application *domain.OIDCApp, resourceOwner string) (_ *domain.OIDCApp, err error) {
	project, err := r.getProjectByID(ctx, application.AggregateID, resourceOwner)
	if err != nil {
		return nil, err
	}
	addedApplication := NewApplicationOIDCConfigWriteModel(application.AggregateID, resourceOwner)
	projectAgg := ProjectAggregateFromWriteModel(&addedApplication.WriteModel)
	err = r.addOIDCApplication(ctx, projectAgg, project, application, resourceOwner)
	if err != nil {
		return nil, err
	}
	err = r.eventstore.PushAggregate(ctx, addedApplication, projectAgg)
	if err != nil {
		return nil, err
	}

	return oidcWriteModelToOIDCConfig(addedApplication), nil
}

func (r *CommandSide) addOIDCApplication(ctx context.Context, projectAgg *project.Aggregate, proj *domain.Project, oidcApp *domain.OIDCApp, resourceOwner string) (err error) {
	if !oidcApp.IsValid() {
		return caos_errs.ThrowPreconditionFailed(nil, "PROJECT-Bff2g", "Errors.Application.Invalid")
	}
	oidcApp.AppID, err = r.idGenerator.Next()
	if err != nil {
		return err
	}

	projectAgg.PushEvents(project.NewApplicationAddedEvent(ctx, oidcApp.AppID, oidcApp.AppName, resourceOwner, domain.AppTypeOIDC))

	var stringPw string
	err = oidcApp.GenerateNewClientID(r.idGenerator, proj)
	if err != nil {
		return err
	}
	stringPw, err = oidcApp.GenerateClientSecretIfNeeded(r.applicationSecretGenerator)
	if err != nil {
		return err
	}
	projectAgg.PushEvents(project.NewOIDCConfigAddedEvent(ctx,
		oidcApp.OIDCVersion,
		oidcApp.AppID,
		oidcApp.ClientID,
		oidcApp.ClientSecret,
		oidcApp.RedirectUris,
		oidcApp.ResponseTypes,
		oidcApp.GrantTypes,
		oidcApp.ApplicationType,
		oidcApp.AuthMethodType,
		oidcApp.PostLogoutRedirectUris,
		oidcApp.DevMode,
		oidcApp.AccessTokenType,
		oidcApp.AccessTokenRoleAssertion,
		oidcApp.IDTokenRoleAssertion,
		oidcApp.IDTokenUserinfoAssertion,
		oidcApp.ClockSkew))

	_ = stringPw

	return nil
}

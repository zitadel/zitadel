package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/project"
)

func (r *CommandSide) AddApplication(ctx context.Context, application *domain.Application) (_ *domain.Application, err error) {
	project, err := r.getProjectByID(ctx, application.AggregateID)
	if err != nil {
		return nil, err
	}
	addedApplication := NewApplicationWriteModel(application.AggregateID)
	projectAgg := ProjectAggregateFromWriteModel(&addedApplication.WriteModel)
	err = r.addApplication(ctx, projectAgg, project, application)
	if err != nil {
		return nil, err
	}
	err = r.eventstore.PushAggregate(ctx, addedApplication, projectAgg)
	if err != nil {
		return nil, err
	}

	return applicationWriteModelToApplication(addedApplication), nil
}

func (r *CommandSide) addApplication(ctx context.Context, projectAgg *project.Aggregate, proj *domain.Project, application *domain.Application) (err error) {
	if !application.IsValid(true) {
		return caos_errs.ThrowPreconditionFailed(nil, "PROJECT-Bff2g", "Errors.Application.Invalid")
	}
	application.AppID, err = r.idGenerator.Next()
	if err != nil {
		return err
	}

	projectAgg.PushEvents(project.NewApplicationAddedEvent(ctx, application.AppID, application.Name, application.Type))

	var stringPw string
	if application.OIDCConfig != nil {
		application.OIDCConfig.AppID = application.AggregateID
		err = application.OIDCConfig.GenerateNewClientID(r.idGenerator, proj)
		if err != nil {
			return err
		}
		stringPw, err = application.OIDCConfig.GenerateClientSecretIfNeeded(r.applicationSecretGenerator)
		if err != nil {
			return err
		}
		projectAgg.PushEvents(project.NewOIDCConfigAddedEvent(ctx,
			application.OIDCConfig.OIDCVersion,
			application.OIDCConfig.AppID,
			application.OIDCConfig.ClientID,
			application.OIDCConfig.ClientSecret,
			application.OIDCConfig.RedirectUris,
			application.OIDCConfig.ResponseTypes,
			application.OIDCConfig.GrantTypes,
			application.OIDCConfig.ApplicationType,
			application.OIDCConfig.AuthMethodType,
			application.OIDCConfig.PostLogoutRedirectUris,
			application.OIDCConfig.DevMode,
			application.OIDCConfig.AccessTokenType,
			application.OIDCConfig.AccessTokenRoleAssertion,
			application.OIDCConfig.IDTokenRoleAssertion,
			application.OIDCConfig.IDTokenUserinfoAssertion,
			application.OIDCConfig.ClockSkew))
	}
	_ = stringPw

	return nil
}

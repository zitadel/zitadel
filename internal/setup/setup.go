package setup

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/errors"
	caos_errs "github.com/caos/zitadel/internal/errors"
	es_int "github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	es_iam "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	iam_event "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	es_org "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	org_event "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	proj_model "github.com/caos/zitadel/internal/project/model"
	es_proj "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	es_usr "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	usr_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	"github.com/caos/zitadel/internal/v2/business/command"
)

type Setup struct {
	iamID         string
	IamEvents     *iam_event.IAMEventstore
	OrgEvents     *org_event.OrgEventstore
	UserEvents    *usr_event.UserEventstore
	ProjectEvents *proj_event.ProjectEventstore

	Commands *command.CommandSide
}

const (
	OrgOwnerRole                   = "ORG_OWNER"
	SetupUser                      = "SETUP"
	OIDCResponseTypeCode           = "CODE"
	OIDCResponseTypeIDToken        = "ID_TOKEN"
	OIDCResponseTypeToken          = "ID_TOKEN TOKEN"
	OIDCGrantTypeAuthorizationCode = "AUTHORIZATION_CODE"
	OIDCGrantTypeImplicit          = "IMPLICIT"
	OIDCGrantTypeRefreshToken      = "REFRESH_TOKEN"
	OIDCApplicationTypeNative      = "NATIVE"
	OIDCApplicationTypeUserAgent   = "USER_AGENT"
	OIDCApplicationTypeWeb         = "WEB"
	OIDCAuthMethodTypeNone         = "NONE"
	OIDCAuthMethodTypeBasic        = "BASIC"
	OIDCAuthMethodTypePost         = "POST"
)

func StartSetup(esConfig es_int.Config, sd systemdefaults.SystemDefaults) (*Setup, error) {
	setup := &Setup{
		iamID: sd.IamID,
	}
	es, err := es_int.Start(esConfig)
	if err != nil {
		return nil, err
	}

	setup.IamEvents, err = es_iam.StartIAM(es_iam.IAMConfig{
		Eventstore: es,
		Cache:      esConfig.Cache,
	}, sd)
	if err != nil {
		return nil, err
	}

	setup.OrgEvents = es_org.StartOrg(es_org.OrgConfig{Eventstore: es, IAMDomain: sd.Domain}, sd)

	setup.ProjectEvents, err = es_proj.StartProject(es_proj.ProjectConfig{
		Eventstore: es,
		Cache:      esConfig.Cache,
	}, sd)
	if err != nil {
		return nil, err
	}

	setup.UserEvents, err = es_usr.StartUser(es_usr.UserConfig{
		Eventstore: es,
		Cache:      esConfig.Cache,
	}, sd)
	if err != nil {
		return nil, err
	}

	return setup, nil
}

func (s *Setup) Execute(ctx context.Context, setUpConfig IAMSetUp) error {
	logging.Log("SETUP-hwG32").Info("starting setup")

	iam, err := s.IamEvents.IAMByID(ctx, s.iamID)
	if err != nil && !caos_errs.IsNotFound(err) {
		return err
	}
	if iam != nil && (iam.SetUpDone == iam_model.StepCount-1 || iam.SetUpStarted != iam.SetUpDone) {
		logging.Log("SETUP-cWEsn").Info("all steps done")
		return nil
	}

	if iam == nil {
		iam = &iam_model.IAM{ObjectRoot: models.ObjectRoot{AggregateID: s.iamID}}
	}

	steps, err := setUpConfig.steps(iam.SetUpDone)
	if err != nil || len(steps) == 0 {
		return err
	}

	ctx = setSetUpContextData(ctx, s.iamID)

	for _, step := range steps {
		step.init(s)
		if step.step() != iam.SetUpDone+1 {
			logging.LogWithFields("SETUP-rxRM1", "step", step.step(), "previous", iam.SetUpDone).Warn("wrong step order")
			return errors.ThrowPreconditionFailed(nil, "SETUP-wwAqO", "too few steps for this zitadel verison")
		}
		iam, err = s.IamEvents.StartSetup(ctx, s.iamID, step.step())
		if err != nil {
			return err
		}

		iam, err = step.execute(ctx)
		if err != nil {
			return err
		}

		err = s.validateExecutedStep(ctx)
		if err != nil {
			return err
		}
	}

	logging.Log("SETUP-ds31h").Info("setup done")
	return nil
}

func (s *Setup) validateExecutedStep(ctx context.Context) error {
	iam, err := s.IamEvents.IAMByID(ctx, s.iamID)
	if err != nil {
		return err
	}
	if iam.SetUpStarted != iam.SetUpDone {
		return errors.ThrowInternal(nil, "SETUP-QeukK", "started step is not equal to done")
	}
	return nil
}

func getOIDCResponseTypes(responseTypes []string) []proj_model.OIDCResponseType {
	types := make([]proj_model.OIDCResponseType, len(responseTypes))
	for i, t := range responseTypes {
		types[i] = getOIDCResponseType(t)
	}
	return types
}

func getOIDCResponseType(responseType string) proj_model.OIDCResponseType {
	switch responseType {
	case OIDCResponseTypeCode:
		return proj_model.OIDCResponseTypeCode
	case OIDCResponseTypeIDToken:
		return proj_model.OIDCResponseTypeIDToken
	case OIDCResponseTypeToken:
		return proj_model.OIDCResponseTypeIDTokenToken
	}
	return proj_model.OIDCResponseTypeCode
}

func getOIDCGrantTypes(grantTypes []string) []proj_model.OIDCGrantType {
	types := make([]proj_model.OIDCGrantType, len(grantTypes))
	for i, t := range grantTypes {
		types[i] = getOIDCGrantType(t)
	}
	return types
}

func getOIDCGrantType(grantTypes string) proj_model.OIDCGrantType {
	switch grantTypes {
	case OIDCGrantTypeAuthorizationCode:
		return proj_model.OIDCGrantTypeAuthorizationCode
	case OIDCGrantTypeImplicit:
		return proj_model.OIDCGrantTypeImplicit
	case OIDCGrantTypeRefreshToken:
		return proj_model.OIDCGrantTypeRefreshToken
	}
	return proj_model.OIDCGrantTypeAuthorizationCode
}

func getOIDCApplicationType(appType string) proj_model.OIDCApplicationType {
	switch appType {
	case OIDCApplicationTypeNative:
		return proj_model.OIDCApplicationTypeNative
	case OIDCApplicationTypeUserAgent:
		return proj_model.OIDCApplicationTypeUserAgent
	case OIDCApplicationTypeWeb:
		return proj_model.OIDCApplicationTypeWeb
	}
	return proj_model.OIDCApplicationTypeWeb
}

func getOIDCAuthMethod(authMethod string) proj_model.OIDCAuthMethodType {
	switch authMethod {
	case OIDCAuthMethodTypeNone:
		return proj_model.OIDCAuthMethodTypeNone
	case OIDCAuthMethodTypeBasic:
		return proj_model.OIDCAuthMethodTypeBasic
	case OIDCAuthMethodTypePost:
		return proj_model.OIDCAuthMethodTypePost
	}
	return proj_model.OIDCAuthMethodTypeBasic
}

func setSetUpContextData(ctx context.Context, orgID string) context.Context {
	return authz.SetCtxData(ctx, authz.CtxData{UserID: SetupUser, OrgID: orgID})
}

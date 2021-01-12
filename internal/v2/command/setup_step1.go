package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

const (
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

type Step1 struct {
	GlobalOrg          string
	IAMProject         string
	DefaultLoginPolicy LoginPolicy //*iam_model.LoginPolicy
	Orgs               []Org
	Owners             []string

	//setup              *Setup
	//createdUsers       map[string]*usr_model.User
	//createdOrgs        map[string]*org_model.Org
	//createdProjects    map[string]*proj_model.Project
	//pwComplexityPolicy *iam_model.PasswordComplexityPolicyView
}

func (s *Step1) Step() domain.Step {
	return domain.Step1
}

func (s *Step1) execute(ctx context.Context, commandSide *CommandSide) error {
	return commandSide.SetupStep1(ctx, s)
}

type LoginPolicy struct {
	AllowRegister         bool
	AllowUsernamePassword bool
	AllowExternalIdp      bool
}

type User struct {
	FirstName string
	LastName  string
	UserName  string
	Email     string
	Password  string
}

type Org struct {
	Name         string
	Domain       string
	OrgIamPolicy bool
	Owner        User
	Projects     []Project
}

type Project struct {
	Name     string
	Users    []User
	Members  []string
	OIDCApps []OIDCApp
}

type OIDCApp struct {
	Name                   string
	RedirectUris           []string
	ResponseTypes          []string
	GrantTypes             []string
	ApplicationType        string
	AuthMethodType         string
	PostLogoutRedirectUris []string
	DevMode                bool
}

func (r *CommandSide) SetupStep1(ctx context.Context, step1 *Step1) error {
	iamAgg := iam_repo.NewAggregate(domain.IAMID, domain.IAMID, 0)
	//create default login policy
	err := r.addDefaultLoginPolicy(ctx, iamAgg, NewIAMLoginPolicyWriteModel(),
		&domain.LoginPolicy{
			AllowUsernamePassword: step1.DefaultLoginPolicy.AllowUsernamePassword,
			AllowRegister:         step1.DefaultLoginPolicy.AllowRegister,
			AllowExternalIdp:      step1.DefaultLoginPolicy.AllowExternalIdp,
		})
	if err != nil {
		return err
	}
	//create orgs
	aggregates := make([]eventstore.Aggregater, 0)
	for _, organisation := range step1.Orgs {
		orgAgg, userAgg, orgMemberAgg, err := r.setUpOrg(ctx,
			&domain.Org{
				Name:    organisation.Name,
				Domains: []*domain.OrgDomain{{Domain: organisation.Domain}},
			},
			&domain.User{
				UserName: organisation.Owner.UserName,
				Human: &domain.Human{
					Profile: &domain.Profile{
						FirstName: organisation.Owner.FirstName,
						LastName:  organisation.Owner.LastName,
					},
					Password: &domain.Password{
						SecretString: organisation.Owner.Password,
					},
					Email: &domain.Email{
						EmailAddress:    organisation.Owner.Email,
						IsEmailVerified: true,
					},
				},
			})
		if err != nil {
			return err
		}
		if organisation.OrgIamPolicy {
			err = r.addOrgIAMPolicy(ctx, orgAgg, NewORGOrgIAMPolicyWriteModel(orgAgg.ID()), &domain.OrgIAMPolicy{UserLoginMustBeDomain: false})
			if err != nil {
				return err
			}
		}
		aggregates = append(aggregates, orgAgg, userAgg, orgMemberAgg)
		//projects
		for _, proj := range organisation.Projects {
			project := &domain.Project{Name: proj.Name}
			projectAgg, _, err := r.addProject(ctx, project, orgAgg.ID(), userAgg.ID())
			if err != nil {
				return err
			}
			//create applications
			for _, app := range proj.OIDCApps {
				err = r.addApplication(ctx, projectAgg, project, &domain.Application{
					ObjectRoot: models.ObjectRoot{
						AggregateID: projectAgg.ID(),
					},
					Name: app.Name,
					Type: domain.AppTypeOIDC,
					OIDCConfig: &domain.OIDCConfig{
						RedirectUris:    app.RedirectUris,
						ResponseTypes:   getOIDCResponseTypes(app.ResponseTypes),
						GrantTypes:      getOIDCGrantTypes(app.GrantTypes),
						ApplicationType: getOIDCApplicationType(app.ApplicationType),
						AuthMethodType:  getOIDCAuthMethod(app.AuthMethodType),
						DevMode:         app.DevMode,
					},
				})
				if err != nil {
					return err
				}
			}
			aggregates = append(aggregates, projectAgg)
		}
	}

	//set iam owners
	//set global org
	//set iam project id

	/*aggregates:
	  iam:
	  	default login policy
	  	iam owner
	  org:
	  	default
	  	caos
	  		zitadel

	*/
	iamAgg.PushEvents(iam_repo.NewSetupStepDoneEvent(ctx, domain.Step1))

	_, err = r.eventstore.PushAggregates(ctx, append(aggregates, iamAgg)...)
	if err != nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-Gr2hh", "Setup Step1 failed")
	}
	return nil
}

func getOIDCResponseTypes(responseTypes []string) []domain.OIDCResponseType {
	types := make([]domain.OIDCResponseType, len(responseTypes))
	for i, t := range responseTypes {
		types[i] = getOIDCResponseType(t)
	}
	return types
}

func getOIDCResponseType(responseType string) domain.OIDCResponseType {
	switch responseType {
	case OIDCResponseTypeCode:
		return domain.OIDCResponseTypeCode
	case OIDCResponseTypeIDToken:
		return domain.OIDCResponseTypeIDToken
	case OIDCResponseTypeToken:
		return domain.OIDCResponseTypeIDTokenToken
	}
	return domain.OIDCResponseTypeCode
}

func getOIDCGrantTypes(grantTypes []string) []domain.OIDCGrantType {
	types := make([]domain.OIDCGrantType, len(grantTypes))
	for i, t := range grantTypes {
		types[i] = getOIDCGrantType(t)
	}
	return types
}

func getOIDCGrantType(grantTypes string) domain.OIDCGrantType {
	switch grantTypes {
	case OIDCGrantTypeAuthorizationCode:
		return domain.OIDCGrantTypeAuthorizationCode
	case OIDCGrantTypeImplicit:
		return domain.OIDCGrantTypeImplicit
	case OIDCGrantTypeRefreshToken:
		return domain.OIDCGrantTypeRefreshToken
	}
	return domain.OIDCGrantTypeAuthorizationCode
}

func getOIDCApplicationType(appType string) domain.OIDCApplicationType {
	switch appType {
	case OIDCApplicationTypeNative:
		return domain.OIDCApplicationTypeNative
	case OIDCApplicationTypeUserAgent:
		return domain.OIDCApplicationTypeUserAgent
	case OIDCApplicationTypeWeb:
		return domain.OIDCApplicationTypeWeb
	}
	return domain.OIDCApplicationTypeWeb
}

func getOIDCAuthMethod(authMethod string) domain.OIDCAuthMethodType {
	switch authMethod {
	case OIDCAuthMethodTypeNone:
		return domain.OIDCAuthMethodTypeNone
	case OIDCAuthMethodTypeBasic:
		return domain.OIDCAuthMethodTypeBasic
	case OIDCAuthMethodTypePost:
		return domain.OIDCAuthMethodTypePost
	}
	return domain.OIDCAuthMethodTypeBasic
}

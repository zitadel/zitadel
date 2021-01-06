package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/v2/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
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
	return commandSide.SetupStep1(ctx, commandSide.iamID, s)
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
	Users        []User
	Owners       []string
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

func (r *CommandSide) SetupStep1(ctx context.Context, iamID string, step1 *Step1) error {
	iamAgg := iam_repo.NewAggregate(r.iamID, "", 0)
	//create default login policy
	err := r.addDefaultLoginPolicy(ctx, iamAgg, NewIAMLoginPolicyWriteModel(iamAgg.ID()),
		&domain.LoginPolicy{
			AllowUsernamePassword: step1.DefaultLoginPolicy.AllowUsernamePassword,
			AllowRegister:         step1.DefaultLoginPolicy.AllowRegister,
			AllowExternalIdp:      step1.DefaultLoginPolicy.AllowExternalIdp,
		})
	if err != nil {
		return err
	}
	//create orgs
	//create projects
	//create applications
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

	_, err = r.eventstore.PushAggregates(ctx, iamAgg)
	if err != nil {
		return caos_errs.ThrowPreconditionFailed(nil, "EVENT-Gr2hh", "Setup Step1 failed")
	}
	return nil
}

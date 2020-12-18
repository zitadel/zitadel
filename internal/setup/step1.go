package setup

import (
	"context"

	"github.com/caos/logging"

	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/business/command"
)

type Step1 struct {
	//GlobalOrg          string
	//IAMProject         string
	//DefaultLoginPolicy LoginPolicy
	//Orgs               []Org
	//Owners             []string
	command.Step1

	//setup              *Setup
	//createdUsers       map[string]*usr_model.User
	//createdOrgs        map[string]*org_model.Org
	//createdProjects    map[string]*proj_model.Project
	//pwComplexityPolicy *iam_model.PasswordComplexityPolicyView
}

func (s *Step1) isNil() bool {
	return s == nil
}

func (s *Step1) step() iam_model.Step {
	return iam_model.Step1
}

//func (s *Step1) init(setup *Setup) {
//	s.setup = setup
//	s.createdUsers = make(map[string]*usr_model.User)
//	s.createdOrgs = make(map[string]*org_model.Org)
//	s.createdProjects = make(map[string]*proj_model.Project)
//}

func (s *Step1) execute(ctx context.Context, commands command.CommandSide) error {
	step2 := command.Step1{
		GlobalOrg:          "",
		IAMProject:         "",
		DefaultLoginPolicy: nil,
		Orgs:               nil,
		Owners:             nil,
	}
	err := commands.SetupStep1(ctx, "iamID", step2)
	if err != nil {
		logging.Log("SETUP-de342").WithField("step", s.step()).WithError(err).Error("unable to finish setup")
		return err
	}
	return nil
}

//
//func (step *Step1) loginPolicy(ctx context.Context, policy LoginPolicy) error {
//	logging.Log("SETUP-4djul").Info("setting up login policy")
//	loginPolicy := &iam_model.LoginPolicy{
//		ObjectRoot: models.ObjectRoot{
//			AggregateID: step.setup.iamID,
//		},
//		AllowRegister:         policy.AllowRegister,
//		AllowUsernamePassword: policy.AllowUsernamePassword,
//		AllowExternalIdp:      policy.AllowExternalIdp,
//	}
//	_, err := step.setup.Commands.AddDefaultLoginPolicy(ctx, loginPolicy)
//	return err
//}
//
//func (step *Step1) orgs(ctx context.Context, orgs []Org) error {
//	logging.Log("SETUP-dsTh3").Info("setting up orgs")
//	for _, iamOrg := range orgs {
//		org, err := step.org(ctx, iamOrg)
//		if err != nil {
//			logging.LogWithFields("SETUP-IlLif", "Org", iamOrg.Name).WithError(err).Error("unable to create org")
//			return err
//		}
//		step.createdOrgs[iamOrg.Name] = org
//		logging.LogWithFields("SETUP-HR2gh", "name", org.Name, "ID", org.AggregateID).Info("created organisation")
//
//		var policy *iam_model.OrgIAMPolicyView
//		if iamOrg.OrgIamPolicy {
//			policy, err = step.iamorgpolicy(ctx, org)
//			if err != nil {
//				logging.LogWithFields("SETUP-IlLif", "Org IAM Policy", iamOrg.Name).WithError(err).Error("unable to create iam org policy")
//				return err
//			}
//		} else {
//			policy = &iam_model.OrgIAMPolicyView{
//				UserLoginMustBeDomain: true,
//			}
//		}
//
//		ctx = setSetUpContextData(ctx, org.AggregateID)
//		err = step.users(ctx, iamOrg.Users, policy)
//		if err != nil {
//			logging.LogWithFields("SETUP-8zfwz", "Org", iamOrg.Name).WithError(err).Error("unable to set up org users")
//			return err
//		}
//
//		err = step.orgOwners(ctx, org, iamOrg.Owners)
//		if err != nil {
//			logging.LogWithFields("SETUP-0874m", "Org", iamOrg.Name).WithError(err).Error("unable to set up org owners")
//			return err
//		}
//
//		err = step.projects(ctx, iamOrg.Projects, step.createdUsers[iamOrg.Owners[0]].AggregateID)
//		if err != nil {
//			logging.LogWithFields("SETUP-wUzqY", "Org", iamOrg.Name).WithError(err).Error("unable to set up org projects")
//			return err
//		}
//	}
//	logging.Log("SETUP-dgjT4").Info("orgs set up")
//	return nil
//}
//
//func (step *Step1) org(ctx context.Context, org Org) (*org_model.Org, error) {
//	ctx = setSetUpContextData(ctx, "")
//	createOrg := &org_model.Org{
//		Name:    org.Name,
//		Domains: []*org_model.OrgDomain{{Domain: org.Domain}},
//	}
//	return step.setup.OrgEvents.CreateOrg(ctx, createOrg, nil)
//}
//
//func (step *Step1) iamorgpolicy(ctx context.Context, org *org_model.Org) (*iam_model.OrgIAMPolicyView, error) {
//	ctx = setSetUpContextData(ctx, org.AggregateID)
//	policy := &iam_model.OrgIAMPolicy{
//		ObjectRoot:            models.ObjectRoot{AggregateID: org.AggregateID},
//		UserLoginMustBeDomain: false,
//	}
//	createdpolicy, err := step.setup.OrgEvents.AddOrgIAMPolicy(ctx, policy)
//	if err != nil {
//		return nil, err
//	}
//	return &iam_model.OrgIAMPolicyView{
//		AggregateID:           org.AggregateID,
//		UserLoginMustBeDomain: createdpolicy.UserLoginMustBeDomain,
//	}, nil
//}
//
//func (step *Step1) iamOwners(ctx context.Context, owners []string) error {
//	logging.Log("SETUP-dtxfj").Info("setting iam owners")
//	for _, iamOwner := range owners {
//		user, ok := step.createdUsers[iamOwner]
//		if !ok {
//			logging.LogWithFields("SETUP-8siew", "Owner", iamOwner).Error("unable to add user to iam members")
//			return caos_errs.ThrowPreconditionFailedf(nil, "SETUP-su6L3", "unable to add user to iam members")
//		}
//		_, err := step.setup.Commands.AddIAMMember(ctx, &iam_model.IAMMember{ObjectRoot: models.ObjectRoot{AggregateID: step.setup.iamID}, UserID: user.AggregateID, Roles: []string{"IAM_OWNER"}})
//		if err != nil {
//			logging.Log("SETUP-LM7rI").WithError(err).Error("unable to add iam administrator to iam members as owner")
//			return err
//		}
//	}
//	logging.Log("SETUP-fg5aq").Info("iam owners set")
//	return nil
//}
//
//func (step *Step1) setGlobalOrg(ctx context.Context, globalOrgName string) error {
//	logging.Log("SETUP-dsj75").Info("setting global org")
//	globalOrg, ok := step.createdOrgs[globalOrgName]
//	if !ok {
//		logging.LogWithFields("SETUP-FBhs9", "GlobalOrg", globalOrgName).Error("global org not created")
//		return caos_errs.ThrowPreconditionFailedf(nil, "SETUP-4GwU7", "global org not created: %v", globalOrgName)
//	}
//
//	if _, err := step.setup.IamEvents.SetGlobalOrg(ctx, step.setup.iamID, globalOrg.AggregateID); err != nil {
//		logging.Log("SETUP-uGMA3").WithError(err).Error("unable to set global org on iam")
//		return err
//	}
//	logging.Log("SETUP-d32h1").Info("global org set")
//	return nil
//}
//
//func (step *Step1) setIamProject(ctx context.Context, iamProjectName string) error {
//	logging.Log("SETUP-HE3qa").Info("setting iam project")
//	iamProject, ok := step.createdProjects[iamProjectName]
//	if !ok {
//		logging.LogWithFields("SETUP-SJFWP", "IAM Project", iamProjectName).Error("iam project created")
//		return caos_errs.ThrowPreconditionFailedf(nil, "SETUP-sGmQt", "iam project not created: %v", iamProjectName)
//	}
//
//	if _, err := step.setup.IamEvents.SetIAMProject(ctx, step.setup.iamID, iamProject.AggregateID); err != nil {
//		logging.Log("SETUP-i1pNh").WithError(err).Error("unable to set iam project on iam")
//		return err
//	}
//	logging.Log("SETUP-d7WEU").Info("iam project set")
//	return nil
//}
//
//func (step *Step1) users(ctx context.Context, users []User, orgPolicy *iam_model.OrgIAMPolicyView) error {
//	for _, user := range users {
//		created, err := step.user(ctx, user, orgPolicy)
//		if err != nil {
//			logging.LogWithFields("SETUP-9soer", "Email", user.Email).WithError(err).Error("unable to create iam user")
//			return err
//		}
//		step.createdUsers[user.Email] = created
//	}
//	return nil
//}
//
//func (step *Step1) user(ctx context.Context, user User, orgPolicy *iam_model.OrgIAMPolicyView) (*usr_model.User, error) {
//	createUser := &usr_model.User{
//		UserName: user.UserName,
//		Human: &usr_model.Human{
//			Profile: &usr_model.Profile{
//				FirstName: user.FirstName,
//				LastName:  user.LastName,
//			},
//			Email: &usr_model.Email{
//				EmailAddress:    user.Email,
//				IsEmailVerified: true,
//			},
//			Password: &usr_model.Password{
//				SecretString: user.Password,
//			},
//		},
//	}
//	return step.setup.UserEvents.CreateUser(ctx, createUser, step.pwComplexityPolicy, orgPolicy)
//}
//
//func (step *Step1) orgOwners(ctx context.Context, org *org_model.Org, owners []string) error {
//	for _, orgOwner := range owners {
//		user, ok := step.createdUsers[orgOwner]
//		if !ok {
//			logging.LogWithFields("SETUP-s9ilr", "Owner", orgOwner).Error("unable to add user to org members")
//			return caos_errs.ThrowPreconditionFailedf(nil, "SETUP-s0prs", "unable to add user to org members: %v", orgOwner)
//		}
//		err := step.orgOwner(ctx, org, user)
//		if err != nil {
//			logging.Log("SETUP-s90oe").WithError(err).Error("unable to add global org admin to members of global org")
//			return err
//		}
//	}
//	return nil
//}
//
//func (step *Step1) orgOwner(ctx context.Context, org *org_model.Org, user *usr_model.User) error {
//	addMember := &org_model.OrgMember{
//		ObjectRoot: models.ObjectRoot{AggregateID: org.AggregateID},
//		UserID:     user.AggregateID,
//		Roles:      []string{OrgOwnerRole},
//	}
//	_, err := step.setup.OrgEvents.AddOrgMember(ctx, addMember)
//	return err
//}
//
//func (step *Step1) projects(ctx context.Context, projects []Project, ownerID string) error {
//	ctxData := authz.GetCtxData(ctx)
//	ctxData.UserID = ownerID
//	projectCtx := authz.SetCtxData(ctx, ctxData)
//
//	for _, project := range projects {
//		createdProject, err := step.project(projectCtx, project)
//		if err != nil {
//			return err
//		}
//		step.createdProjects[createdProject.Name] = createdProject
//		for _, oidc := range project.OIDCApps {
//			app, err := step.oidcApp(ctx, createdProject, oidc)
//			if err != nil {
//				return err
//			}
//			logging.LogWithFields("SETUP-asd32f", "name", app.Name, "clientID", app.OIDCConfig.ClientID).Info("created OIDC application")
//		}
//	}
//	return nil
//}
//
//func (step *Step1) project(ctx context.Context, project Project) (*proj_model.Project, error) {
//	addProject := &proj_model.Project{
//		Name: project.Name,
//	}
//	return step.setup.ProjectEvents.CreateProject(ctx, addProject, false)
//}
//
//func (step *Step1) oidcApp(ctx context.Context, project *proj_model.Project, oidc OIDCApp) (*proj_model.Application, error) {
//	addOIDCApp := &proj_model.Application{
//		ObjectRoot: models.ObjectRoot{AggregateID: project.AggregateID},
//		Name:       oidc.Name,
//		OIDCConfig: &proj_model.OIDCConfig{
//			RedirectUris:           oidc.RedirectUris,
//			ResponseTypes:          getOIDCResponseTypes(oidc.ResponseTypes),
//			GrantTypes:             getOIDCGrantTypes(oidc.GrantTypes),
//			ApplicationType:        getOIDCApplicationType(oidc.ApplicationType),
//			AuthMethodType:         getOIDCAuthMethod(oidc.AuthMethodType),
//			PostLogoutRedirectUris: oidc.PostLogoutRedirectUris,
//			DevMode:                oidc.DevMode,
//		},
//	}
//	return step.setup.ProjectEvents.AddApplication(ctx, addOIDCApp)
//}

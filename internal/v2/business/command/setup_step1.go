package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/v2/business/domain"
	iam_repo "github.com/caos/zitadel/internal/v2/repository/iam"
)

type Step1 struct {
	GlobalOrg          string
	IAMProject         string
	DefaultLoginPolicy *iam_model.LoginPolicy
	Orgs               []org_model.Org
	Owners             []string

	//setup              *Setup
	//createdUsers       map[string]*usr_model.User
	//createdOrgs        map[string]*org_model.Org
	//createdProjects    map[string]*proj_model.Project
	//pwComplexityPolicy *iam_model.PasswordComplexityPolicyView
}

func (r *CommandSide) SetupStep1(ctx context.Context, iamID string, step1 Step1) error {
	iam, err := r.iamByID(ctx, iamID)
	if err != nil && !caos_errs.IsNotFound(err) {
		return err
	}
	//create default login policy
	_, err = r.addDefaultLoginPolicy(ctx, iam, step1.DefaultLoginPolicy)
	if err != nil {
		return err
	}
	//create orgs
	//create porjects
	//create applications
	//set iam owners
	//set global org
	//set iam project id
	_, err = r.setup(ctx, iam, domain.Step1, iam_repo.NewSetupStepDoneEvent(ctx, domain.Step1))
	return err
}

func (r *CommandSide) addDefaultLoginPolicy(ctx context.Context, iam *IAMWriteModel, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error) {
	if !policy.IsValid() {
		return nil, caos_errs.ThrowPreconditionFailed(nil, "IAM-5Mv0s", "Errors.IAM.LoginPolicyInvalid")
	}

	addedPolicy := NewWriteModel(policy.AggregateID)
	err := r.eventstore.FilterToQueryReducer(ctx, addedPolicy)
	if err != nil {
		return nil, err
	}
	if addedPolicy.IsActive {
		return nil, caos_errs.ThrowAlreadyExists(nil, "IAM-2B0ps", "Errors.IAM.LoginPolicy.AlreadyExists")
	}

	iamAgg.PushEvents(iam_repo.NewLoginPolicyAddedEvent(ctx, policy.AllowUsernamePassword, policy.AllowRegister, policy.AllowExternalIdp, policy.ForceMFA, domain.PasswordlessType(policy.PasswordlessType)))
}

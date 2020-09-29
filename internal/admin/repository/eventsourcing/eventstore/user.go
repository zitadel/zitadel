package eventstore

import (
	"context"
	admin_view "github.com/caos/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_view "github.com/caos/zitadel/internal/iam/repository/view/model"

	"github.com/caos/zitadel/internal/api/authz"
	org_event "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	policy_event "github.com/caos/zitadel/internal/policy/repository/eventsourcing"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type UserRepo struct {
	UserEvents     *usr_event.UserEventstore
	PolicyEvents   *policy_event.PolicyEventstore
	OrgEvents      *org_event.OrgEventstore
	View           *admin_view.View
	SystemDefaults systemdefaults.SystemDefaults
}

func (repo *UserRepo) UserByID(ctx context.Context, id string) (project *usr_model.User, err error) {
	return repo.UserEvents.UserByID(ctx, id)
}

func (repo *UserRepo) CreateUser(ctx context.Context, user *usr_model.User) (*usr_model.User, error) {
	pwPolicy, err := repo.View.PasswordComplexityPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if err != nil && caos_errs.IsNotFound(err) {
		pwPolicy, err = repo.View.PasswordComplexityPolicyByAggregateID(repo.SystemDefaults.IamID)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	pwPolicyView := iam_view.PasswordComplexityViewToModel(pwPolicy)
	orgPolicy, err := repo.OrgEvents.GetOrgIAMPolicy(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return repo.UserEvents.CreateUser(ctx, user, pwPolicyView, orgPolicy)
}

func (repo *UserRepo) RegisterUser(ctx context.Context, user *usr_model.User, resourceOwner string) (*usr_model.User, error) {
	policyResourceOwner := authz.GetCtxData(ctx).OrgID
	if resourceOwner != "" {
		policyResourceOwner = resourceOwner
	}
	pwPolicy, err := repo.View.PasswordComplexityPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if err != nil && caos_errs.IsNotFound(err) {
		pwPolicy, err = repo.View.PasswordComplexityPolicyByAggregateID(repo.SystemDefaults.IamID)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	pwPolicyView := iam_view.PasswordComplexityViewToModel(pwPolicy)

	orgPolicy, err := repo.OrgEvents.GetOrgIAMPolicy(ctx, policyResourceOwner)
	if err != nil {
		return nil, err
	}
	return repo.UserEvents.RegisterUser(ctx, user, pwPolicyView, orgPolicy, resourceOwner)
}

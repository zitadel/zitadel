package eventstore

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	org_event "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	policy_event "github.com/caos/zitadel/internal/policy/repository/eventsourcing"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type UserRepo struct {
	UserEvents   *usr_event.UserEventstore
	PolicyEvents *policy_event.PolicyEventstore
	OrgEvents    *org_event.OrgEventstore
}

func (repo *UserRepo) UserByID(ctx context.Context, id string) (project *usr_model.User, err error) {
	return repo.UserEvents.UserByID(ctx, id)
}

func (repo *UserRepo) CreateUser(ctx context.Context, user *usr_model.User) (*usr_model.User, error) {
	pwPolicy, err := repo.PolicyEvents.GetPasswordComplexityPolicy(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	orgPolicy, err := repo.OrgEvents.GetOrgIamPolicy(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return repo.UserEvents.CreateUser(ctx, user, pwPolicy, orgPolicy)
}

func (repo *UserRepo) RegisterUser(ctx context.Context, user *usr_model.User, resourceOwner string) (*usr_model.User, error) {
	policyResourceOwner := authz.GetCtxData(ctx).OrgID
	if resourceOwner != "" {
		policyResourceOwner = resourceOwner
	}
	pwPolicy, err := repo.PolicyEvents.GetPasswordComplexityPolicy(ctx, policyResourceOwner)
	if err != nil {
		return nil, err
	}
	orgPolicy, err := repo.OrgEvents.GetOrgIamPolicy(ctx, policyResourceOwner)
	if err != nil {
		return nil, err
	}
	return repo.UserEvents.RegisterUser(ctx, user, pwPolicy, orgPolicy, resourceOwner)
}

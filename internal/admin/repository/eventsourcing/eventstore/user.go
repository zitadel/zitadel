package eventstore

import (
	"context"
	admin_view "github.com/caos/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	caos_errs "github.com/caos/zitadel/internal/errors"
	iam_view "github.com/caos/zitadel/internal/iam/repository/view/model"

	"github.com/caos/zitadel/internal/api/authz"
	org_event "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_event "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type UserRepo struct {
	UserEvents     *usr_event.UserEventstore
	OrgEvents      *org_event.OrgEventstore
	View           *admin_view.View
	SystemDefaults systemdefaults.SystemDefaults
}

func (repo *UserRepo) UserByID(ctx context.Context, id string) (project *usr_model.User, err error) {
	return repo.UserEvents.UserByID(ctx, id)
}

func (repo *UserRepo) CreateUser(ctx context.Context, user *usr_model.User) (*usr_model.User, error) {
	pwPolicy, err := repo.View.PasswordComplexityPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if caos_errs.IsNotFound(err) {
		pwPolicy, err = repo.View.PasswordComplexityPolicyByAggregateID(repo.SystemDefaults.IamID)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	pwPolicyView := iam_view.PasswordComplexityViewToModel(pwPolicy)
	orgPolicy, err := repo.View.OrgIAMPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if err != nil && caos_errs.IsNotFound(err) {
		orgPolicy, err = repo.View.OrgIAMPolicyByAggregateID(repo.SystemDefaults.IamID)
		if err != nil {
			return nil, err
		}
	}
	orgPolicyView := iam_view.OrgIAMViewToModel(orgPolicy)
	return repo.UserEvents.CreateUser(ctx, user, pwPolicyView, orgPolicyView)
}

func (repo *UserRepo) RegisterUser(ctx context.Context, user *usr_model.User, resourceOwner string) (*usr_model.User, error) {
	policyResourceOwner := authz.GetCtxData(ctx).OrgID
	if resourceOwner != "" {
		policyResourceOwner = resourceOwner
	}
	pwPolicy, err := repo.View.PasswordComplexityPolicyByAggregateID(policyResourceOwner)
	if caos_errs.IsNotFound(err) {
		pwPolicy, err = repo.View.PasswordComplexityPolicyByAggregateID(repo.SystemDefaults.IamID)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	pwPolicyView := iam_view.PasswordComplexityViewToModel(pwPolicy)

	orgPolicy, err := repo.View.OrgIAMPolicyByAggregateID(policyResourceOwner)
	if caos_errs.IsNotFound(err) {
		orgPolicy, err = repo.View.OrgIAMPolicyByAggregateID(repo.SystemDefaults.IamID)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	orgPolicyView := iam_view.OrgIAMViewToModel(orgPolicy)
	return repo.UserEvents.RegisterUser(ctx, user, pwPolicyView, orgPolicyView, resourceOwner)
}

package eventstore

import (
	"context"

	admin_view "github.com/caos/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/view/model"
)

type OrgRepo struct {
	View *admin_view.View

	SystemDefaults systemdefaults.SystemDefaults
}

func (repo *OrgRepo) GetOrgIAMPolicyByID(ctx context.Context, id string) (*iam_model.OrgIAMPolicyView, error) {
	policy, err := repo.View.OrgIAMPolicyByAggregateID(id)
	if errors.IsNotFound(err) {
		return repo.GetDefaultOrgIAMPolicy(ctx)
	}
	if err != nil {
		return nil, err
	}
	return iam_es_model.OrgIAMViewToModel(policy), err
}

func (repo *OrgRepo) GetDefaultOrgIAMPolicy(ctx context.Context) (*iam_model.OrgIAMPolicyView, error) {
	policy, err := repo.View.OrgIAMPolicyByAggregateID(repo.SystemDefaults.IamID)
	if err != nil {
		return nil, err
	}
	policy.Default = true
	return iam_es_model.OrgIAMViewToModel(policy), err
}

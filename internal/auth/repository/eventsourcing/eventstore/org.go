package eventstore

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_view_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/telemetry/tracing"

	auth_view "github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	org_model "github.com/caos/zitadel/internal/org/model"
	"github.com/caos/zitadel/internal/org/repository/view/model"
)

const (
	orgOwnerRole = "ORG_OWNER"
)

type OrgRepository struct {
	SearchLimit uint64

	View           *auth_view.View
	SystemDefaults systemdefaults.SystemDefaults
}

func (repo *OrgRepository) SearchOrgs(ctx context.Context, request *org_model.OrgSearchRequest) (*org_model.OrgSearchResult, error) {
	err := request.EnsureLimit(repo.SearchLimit)
	if err != nil {
		return nil, err
	}
	sequence, err := repo.View.GetLatestOrgSequence()
	logging.Log("EVENT-7Udhz").OnError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Warn("could not read latest org sequence")
	members, count, err := repo.View.SearchOrgs(request)
	if err != nil {
		return nil, err
	}
	result := &org_model.OrgSearchResult{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      model.OrgsToModel(members),
	}
	if err == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *OrgRepository) OrgByPrimaryDomain(primaryDomain string) (*org_model.OrgView, error) {
	org, err := repo.View.OrgByPrimaryDomain(primaryDomain)
	if err != nil {
		return nil, err
	}
	return model.OrgToModel(org), nil
}

func (repo *OrgRepository) GetDefaultOrgIAMPolicy(ctx context.Context) (*iam_model.OrgIAMPolicyView, error) {
	orgPolicy, err := repo.View.OrgIAMPolicyByAggregateID(repo.SystemDefaults.IamID)
	if err != nil {
		return nil, err
	}
	policy := iam_view_model.OrgIAMViewToModel(orgPolicy)
	policy.IAMDomain = repo.SystemDefaults.Domain
	return policy, err
}

func (repo *OrgRepository) GetOrgIAMPolicy(ctx context.Context, orgID string) (*iam_model.OrgIAMPolicyView, error) {
	orgPolicy, err := repo.View.OrgIAMPolicyByAggregateID(orgID)
	if errors.IsNotFound(err) {
		orgPolicy, err = repo.View.OrgIAMPolicyByAggregateID(repo.SystemDefaults.IamID)
	}
	if err != nil {
		return nil, err
	}
	return iam_view_model.OrgIAMViewToModel(orgPolicy), nil
}

func (repo *OrgRepository) GetIDPConfigByID(ctx context.Context, idpConfigID string) (*iam_model.IDPConfigView, error) {
	idpConfig, err := repo.View.IDPConfigByID(idpConfigID)
	if err != nil {
		return nil, err
	}
	return iam_view_model.IDPConfigViewToModel(idpConfig), nil
}

func (repo *OrgRepository) GetMyPasswordComplexityPolicy(ctx context.Context) (*iam_model.PasswordComplexityPolicyView, error) {
	policy, err := repo.View.PasswordComplexityPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if errors.IsNotFound(err) {
		policy, err = repo.View.PasswordComplexityPolicyByAggregateID(repo.SystemDefaults.IamID)
		if err != nil {
			return nil, err
		}
		policy.Default = true
	}
	if err != nil {
		return nil, err
	}
	return iam_view_model.PasswordComplexityViewToModel(policy), err
}

func (repo *OrgRepository) GetLabelPolicy(ctx context.Context, orgID string) (*iam_model.LabelPolicyView, error) {
	orgPolicy, err := repo.View.LabelPolicyByAggregateIDAndState(orgID, int32(domain.LabelPolicyStateActive))
	if errors.IsNotFound(err) {
		orgPolicy, err = repo.View.LabelPolicyByAggregateIDAndState(repo.SystemDefaults.IamID, int32(domain.LabelPolicyStateActive))
	}
	if err != nil {
		return nil, err
	}
	return iam_view_model.LabelPolicyViewToModel(orgPolicy), nil
}

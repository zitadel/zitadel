package eventstore

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/api/authz"
	auth_view "github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	eventstore "github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_view_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/repository/iam"
)

const (
	orgOwnerRole = "ORG_OWNER"
)

type OrgRepository struct {
	SearchLimit uint64

	Eventstore     eventstore.Eventstore
	View           *auth_view.View
	SystemDefaults systemdefaults.SystemDefaults
	Query          *query.Queries
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
	policy, err := repo.Query.PasswordComplexityPolicyByOrg(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	return iam_view_model.PasswordComplexityViewToModel(policy), err
}

func (repo *OrgRepository) GetLabelPolicy(ctx context.Context, orgID string) (*domain.LabelPolicy, error) {
	orgPolicy, err := repo.View.LabelPolicyByAggregateIDAndState(orgID, int32(domain.LabelPolicyStateActive))
	if errors.IsNotFound(err) {
		orgPolicy, err = repo.View.LabelPolicyByAggregateIDAndState(repo.SystemDefaults.IamID, int32(domain.LabelPolicyStateActive))
	}
	if err != nil {
		return nil, err
	}
	return orgPolicy.ToDomain(), nil
}

func (repo *OrgRepository) GetLoginText(ctx context.Context, orgID string) ([]*domain.CustomText, error) {
	loginTexts, err := repo.View.CustomTextsByAggregateIDAndTemplate(domain.IAMID, domain.LoginCustomText)
	if err != nil {
		return nil, err
	}
	orgLoginTexts, err := repo.View.CustomTextsByAggregateIDAndTemplate(orgID, domain.LoginCustomText)
	if err != nil {
		return nil, err
	}
	return append(iam_view_model.CustomTextViewsToDomain(loginTexts), iam_view_model.CustomTextViewsToDomain(orgLoginTexts)...), nil
}

func (repo *OrgRepository) GetDefaultPrivacyPolicy(ctx context.Context) (*iam_model.PrivacyPolicyView, error) {
	policy, viewErr := repo.View.PrivacyPolicyByAggregateID(repo.SystemDefaults.IamID)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		policy = new(iam_view_model.PrivacyPolicyView)
	}
	events, esErr := repo.getIAMEvents(ctx, policy.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-LPJMp", "Errors.IAM.PrivacyPolicy.NotFound")
	}
	if esErr != nil {
		logging.Log("EVENT-1l7bf").WithError(esErr).Debug("error retrieving new events")
		return iam_view_model.PrivacyViewToModel(policy), nil
	}
	policyCopy := *policy
	for _, event := range events {
		if err := policyCopy.AppendEvent(event); err != nil {
			return iam_view_model.PrivacyViewToModel(policy), nil
		}
	}
	result := iam_view_model.PrivacyViewToModel(policy)
	result.Default = true
	return result, nil
}

func (p *OrgRepository) getIAMEvents(ctx context.Context, sequence uint64) ([]*models.Event, error) {
	return p.Eventstore.FilterEvents(ctx, models.NewSearchQuery().AggregateIDFilter(p.SystemDefaults.IamID).AggregateTypeFilter(iam.AggregateType))
}

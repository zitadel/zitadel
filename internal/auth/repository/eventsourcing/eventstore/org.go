package eventstore

import (
	"context"

	"github.com/zitadel/zitadel/internal/api/authz"
	auth_view "github.com/zitadel/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/domain"
	eventstore "github.com/zitadel/zitadel/internal/eventstore/v1"
	iam_model "github.com/zitadel/zitadel/internal/iam/model"
	iam_view_model "github.com/zitadel/zitadel/internal/iam/repository/view/model"
	"github.com/zitadel/zitadel/internal/query"
)

type OrgRepository struct {
	SearchLimit uint64

	Eventstore     eventstore.Eventstore
	View           *auth_view.View
	SystemDefaults systemdefaults.SystemDefaults
	Query          *query.Queries
}

func (repo *OrgRepository) GetMyPasswordComplexityPolicy(ctx context.Context) (*iam_model.PasswordComplexityPolicyView, error) {
	policy, err := repo.Query.PasswordComplexityPolicyByOrg(ctx, true, authz.GetCtxData(ctx).OrgID, false)
	if err != nil {
		return nil, err
	}
	return iam_view_model.PasswordComplexityViewToModel(policy), err
}

func (repo *OrgRepository) GetLoginText(ctx context.Context, orgID string) ([]*domain.CustomText, error) {
	loginTexts, err := repo.Query.CustomTextListByTemplate(ctx, authz.GetInstance(ctx).InstanceID(), domain.LoginCustomText, false)
	if err != nil {
		return nil, err
	}
	orgLoginTexts, err := repo.Query.CustomTextListByTemplate(ctx, orgID, domain.LoginCustomText, false)
	if err != nil {
		return nil, err
	}
	return append(query.CustomTextsToDomain(loginTexts), query.CustomTextsToDomain(orgLoginTexts)...), nil
}

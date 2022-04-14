package eventstore

import (
	"context"

	"github.com/caos/zitadel/internal/api/authz"
	auth_view "github.com/caos/zitadel/internal/auth/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/domain"
	eventstore "github.com/caos/zitadel/internal/eventstore/v1"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_view_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/query"
)

type OrgRepository struct {
	SearchLimit uint64

	Eventstore     eventstore.Eventstore
	View           *auth_view.View
	SystemDefaults systemdefaults.SystemDefaults
	Query          *query.Queries
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

func (repo *OrgRepository) GetLoginText(ctx context.Context, orgID string) ([]*domain.CustomText, error) {
	loginTexts, err := repo.Query.CustomTextListByTemplate(ctx, authz.GetInstance(ctx).InstanceID(), domain.LoginCustomText)
	if err != nil {
		return nil, err
	}
	orgLoginTexts, err := repo.Query.CustomTextListByTemplate(ctx, orgID, domain.LoginCustomText)
	if err != nil {
		return nil, err
	}
	return append(query.CustomTextsToDomain(loginTexts), query.CustomTextsToDomain(orgLoginTexts)...), nil
}

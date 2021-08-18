package repository

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	org_model "github.com/caos/zitadel/internal/org/model"
)

type OrgRepository interface {
	OrgByPrimaryDomain(primaryDomain string) (*org_model.OrgView, error)
	GetOrgIAMPolicy(ctx context.Context, orgID string) (*iam_model.OrgIAMPolicyView, error)
	GetDefaultOrgIAMPolicy(ctx context.Context) (*iam_model.OrgIAMPolicyView, error)
	GetIDPConfigByID(ctx context.Context, idpConfigID string) (*iam_model.IDPConfigView, error)
	GetMyPasswordComplexityPolicy(ctx context.Context) (*iam_model.PasswordComplexityPolicyView, error)
	GetLabelPolicy(ctx context.Context, orgID string) (*domain.LabelPolicy, error)
	GetLoginText(ctx context.Context, orgID string) ([]*domain.CustomText, error)
	GetDefaultPrivacyPolicy(ctx context.Context) (*iam_model.PrivacyPolicyView, error)
}

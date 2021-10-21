package repository

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type OrgRepository interface {
	GetIDPConfigByID(ctx context.Context, idpConfigID string) (*iam_model.IDPConfigView, error)
	GetMyPasswordComplexityPolicy(ctx context.Context) (*iam_model.PasswordComplexityPolicyView, error)
	GetLabelPolicy(ctx context.Context, orgID string) (*domain.LabelPolicy, error)
	GetLoginText(ctx context.Context, orgID string) ([]*domain.CustomText, error)
}

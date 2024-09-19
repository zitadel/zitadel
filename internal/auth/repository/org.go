package repository

import (
	"context"

	"github.com/zitadel/zitadel/v2/internal/domain"
	iam_model "github.com/zitadel/zitadel/v2/internal/iam/model"
)

type OrgRepository interface {
	GetMyPasswordComplexityPolicy(ctx context.Context) (*iam_model.PasswordComplexityPolicyView, error)
	GetLoginText(ctx context.Context, orgID string) ([]*domain.CustomText, error)
}

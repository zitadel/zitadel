package repository

import (
	"context"

	iam_model "github.com/caos/zitadel/internal/iam/model"
)

type OrgRepository interface {
	GetOrgIAMPolicyByID(ctx context.Context, id string) (*iam_model.OrgIAMPolicyView, error)
}

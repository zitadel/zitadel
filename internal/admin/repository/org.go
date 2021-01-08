package repository

import (
	"context"

	iam_model "github.com/caos/zitadel/internal/iam/model"

	org_model "github.com/caos/zitadel/internal/org/model"
)

type OrgRepository interface {
	IsOrgUnique(ctx context.Context, name, domain string) (bool, error)
	OrgByID(ctx context.Context, id string) (*org_model.Org, error)
	SearchOrgs(ctx context.Context, query *org_model.OrgSearchRequest) (*org_model.OrgSearchResult, error)

	GetOrgIAMPolicyByID(ctx context.Context, id string) (*iam_model.OrgIAMPolicyView, error)
	RemoveOrgIAMPolicy(ctx context.Context, id string) error
}

package command

import (
	"context"
)

type commandProvider interface {
	domainPolicyWriteModel(ctx context.Context, orgID string) (*PolicyDomainWriteModel, error)
}

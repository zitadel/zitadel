package eventstore

import (
	"context"

	"github.com/caos/zitadel/internal/api/auth"
	pol_model "github.com/caos/zitadel/internal/policy/model"
	pol_event "github.com/caos/zitadel/internal/policy/repository/eventsourcing"
)

type PolicyRepo struct {
	PolicyEvents *pol_event.PolicyEventstore
}

func (repo *PolicyRepo) GetMyPasswordComplexityPolicy(ctx context.Context) (*pol_model.PasswordComplexityPolicy, error) {
	ctxData := auth.GetCtxData(ctx)
	return repo.PolicyEvents.GetPasswordComplexityPolicy(ctx, ctxData.OrgID)
}

package command

import (
	"context"

	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/v2/domain"
)

func (r *CommandSide) GetOrg(ctx context.Context, aggregateID string) (*domain.Org, error) {
	orgWriteModel := NewOrgWriteModel(aggregateID)
	err := r.eventstore.FilterToQueryReducer(ctx, orgWriteModel)
	if err != nil {
		return nil, err
	}
	return orgWriteModelToOrg(orgWriteModel), nil
}

func (r *CommandSide) SetUpOrg(ctx context.Context, user *domain.User) (*domain.Org, error) {
	pwPolicy, err := r.GetDefaultPasswordComplexityPolicy(ctx)
	if err != nil {
		return nil, err
	}
	orgPolicy, err := r.GetDefaultOrgIAMPolicy(ctx)
	if err != nil {
		return nil, err
	}
	//TODO: users with verified domain -> domain claimed

	r.addOrg()

	r.addUser()
}

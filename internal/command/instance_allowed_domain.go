package command

import (
	"context"
	"slices"
	"strings"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddAllowedDomain(ctx context.Context, allowedDomain string) (*domain.ObjectDetails, error) {
	allowedDomain = strings.TrimSpace(allowedDomain)
	if allowedDomain == "" {
		return nil, zerrors.ThrowNotFound(nil, "COMMA-Stk21", "Errors.Invalid.Argument")
	}
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	model := NewInstanceAllowedDomainsWriteModel(ctx)
	err := c.eventstore.FilterToQueryReducer(ctx, model)
	if err != nil {
		return nil, err
	}
	if slices.Contains(model.Domains, allowedDomain) {
		return nil, zerrors.ThrowNotFound(nil, "COMMA-hg42a", "Errors.Instance.Domain.AlreadyExists")
	}
	err = c.pushAppendAndReduce(ctx, model, instance.NewAllowedDomainAddedEvent(ctx, &instanceAgg.Aggregate, allowedDomain))
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&model.WriteModel), nil
}

func (c *Commands) RemoveAllowedDomain(ctx context.Context, allowedDomain string) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	model := NewInstanceAllowedDomainsWriteModel(ctx)
	err := c.eventstore.FilterToQueryReducer(ctx, model)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(model.Domains, allowedDomain) {
		return nil, zerrors.ThrowNotFound(nil, "COMMA-de3z9", "Errors.Instance.Domain.NotFound")
	}
	err = c.pushAppendAndReduce(ctx, model, instance.NewAllowedDomainRemovedEvent(ctx, &instanceAgg.Aggregate, allowedDomain))
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&model.WriteModel), nil
}

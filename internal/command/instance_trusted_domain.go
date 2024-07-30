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

func (c *Commands) AddTrustedDomain(ctx context.Context, trustedDomain string) (*domain.ObjectDetails, error) {
	trustedDomain = strings.TrimSpace(trustedDomain)
	if trustedDomain == "" {
		return nil, zerrors.ThrowNotFound(nil, "COMMA-Stk21", "Errors.Invalid.Argument")
	}
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	model := NewInstanceTrustedDomainsWriteModel(ctx)
	err := c.eventstore.FilterToQueryReducer(ctx, model)
	if err != nil {
		return nil, err
	}
	if slices.Contains(model.Domains, trustedDomain) {
		return nil, zerrors.ThrowNotFound(nil, "COMMA-hg42a", "Errors.Instance.Domain.AlreadyExists")
	}
	err = c.pushAppendAndReduce(ctx, model, instance.NewTrustedDomainAddedEvent(ctx, &instanceAgg.Aggregate, trustedDomain))
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&model.WriteModel), nil
}

func (c *Commands) RemoveTrustedDomain(ctx context.Context, trustedDomain string) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	model := NewInstanceTrustedDomainsWriteModel(ctx)
	err := c.eventstore.FilterToQueryReducer(ctx, model)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(model.Domains, trustedDomain) {
		return nil, zerrors.ThrowNotFound(nil, "COMMA-de3z9", "Errors.Instance.Domain.NotFound")
	}
	err = c.pushAppendAndReduce(ctx, model, instance.NewTrustedDomainRemovedEvent(ctx, &instanceAgg.Aggregate, trustedDomain))
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&model.WriteModel), nil
}

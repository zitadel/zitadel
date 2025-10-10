package command

import (
	"context"
	"slices"
	"strings"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddTrustedDomain(ctx context.Context, trustedDomain string) (*domain.ObjectDetails, error) {
	trustedDomain = strings.TrimSpace(trustedDomain)
	if trustedDomain == "" || len(trustedDomain) > 253 {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMA-Stk21", "Errors.Invalid.Argument")
	}
	if !allowDomainRunes.MatchString(trustedDomain) {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMA-S3v3w", "Errors.Instance.Domain.InvalidCharacter")
	}
	model := NewInstanceTrustedDomainsWriteModel(ctx)
	err := c.eventstore.FilterToQueryReducer(ctx, model)
	if err != nil {
		return nil, err
	}
	if slices.Contains(model.Domains, trustedDomain) {
		return nil, zerrors.ThrowPreconditionFailed(nil, "COMMA-hg42a", "Errors.Instance.Domain.AlreadyExists")
	}
	err = c.pushAppendAndReduce(ctx, model, instance.NewTrustedDomainAddedEvent(ctx, InstanceAggregateFromWriteModel(&model.WriteModel), trustedDomain))
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&model.WriteModel), nil
}

func (c *Commands) RemoveTrustedDomain(ctx context.Context, trustedDomain string) (*domain.ObjectDetails, error) {
	trustedDomain = strings.TrimSpace(trustedDomain)
	if trustedDomain == "" || len(trustedDomain) > 253 {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMA-ajAzwu", "Errors.Invalid.Argument")
	}
	if !allowDomainRunes.MatchString(trustedDomain) {
		return nil, zerrors.ThrowInvalidArgument(nil, "COMMA-lfs3Te", "Errors.Instance.Domain.InvalidCharacter")
	}

	model := NewInstanceTrustedDomainsWriteModel(ctx)
	err := c.eventstore.FilterToQueryReducer(ctx, model)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(model.Domains, trustedDomain) {
		return nil, zerrors.ThrowNotFound(nil, "COMMA-de3z9", "Errors.Instance.Domain.NotFound")
	}
	err = c.pushAppendAndReduce(ctx, model, instance.NewTrustedDomainRemovedEvent(ctx, InstanceAggregateFromWriteModel(&model.WriteModel), trustedDomain))
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&model.WriteModel), nil
}

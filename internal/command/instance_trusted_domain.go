package command

import (
	"context"
	"slices"
	"strings"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) AddTrustedDomain(ctx context.Context, trustedDomain string) (details *domain.ObjectDetails, err error) {
	trustedDomain, err = validateTrustedDomain(trustedDomain)
	if err != nil {
		return nil, err
	}
	model := NewInstanceTrustedDomainsWriteModel(ctx)
	err = c.eventstore.FilterToQueryReducer(ctx, model)
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

// setupTrustedDomain is only used for the instance setup process where all events are prepared and pushed in a single transaction
func (c *Commands) setupTrustedDomain(a *instance.Aggregate, trustedDomain string) preparation.Validation {
	return func() (_ preparation.CreateCommands, err error) {
		trustedDomain, err = validateTrustedDomain(trustedDomain)
		if err != nil {
			return nil, err
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			model := NewInstanceTrustedDomainsWriteModel(ctx)
			events, err := filter(ctx, model.Query())
			if err != nil {
				return nil, err
			}
			err = AppendAndReduce(model, events...)
			if err != nil {
				return nil, err
			}
			if slices.Contains(model.Domains, trustedDomain) {
				return nil, zerrors.ThrowPreconditionFailed(nil, "COMMA-hg42a", "Errors.Instance.Domain.AlreadyExists")
			}
			return []eventstore.Command{instance.NewTrustedDomainAddedEvent(ctx, &a.Aggregate, trustedDomain)}, nil
		}, nil
	}
}

func validateTrustedDomain(trustedDomain string) (string, error) {
	trustedDomain = strings.TrimSpace(trustedDomain)
	if trustedDomain == "" || len(trustedDomain) > 253 {
		return "", zerrors.ThrowInvalidArgument(nil, "COMMA-Stk21", "Errors.Invalid.Argument")
	}
	if !allowDomainRunes.MatchString(trustedDomain) {
		return "", zerrors.ThrowInvalidArgument(nil, "COMMA-S3v3w", "Errors.Instance.Domain.InvalidCharacter")
	}
	return trustedDomain, nil
}

func (c *Commands) RemoveTrustedDomain(ctx context.Context, trustedDomain string, errorIfNotFound bool) (*domain.ObjectDetails, error) {
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
		if errorIfNotFound {
			return nil, zerrors.ThrowNotFound(nil, "COMMA-de3z9", "Errors.Instance.Domain.NotFound")
		}
		return nil, nil
	}

	err = c.pushAppendAndReduce(ctx, model, instance.NewTrustedDomainRemovedEvent(ctx, InstanceAggregateFromWriteModel(&model.WriteModel), trustedDomain))
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&model.WriteModel), nil
}

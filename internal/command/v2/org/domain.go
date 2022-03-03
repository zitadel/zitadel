package org

import (
	"context"
	"strings"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
)

func AddDomain(a *org.Aggregate, domain string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if domain = strings.TrimSpace(domain); domain == "" {
			return nil, errors.ThrowInvalidArgument(nil, "ORG-r3h4J", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			return []eventstore.Command{org.NewDomainAddedEvent(ctx, &a.Aggregate, domain)}, nil
		}, nil
	}
}

func VerifyDomain(a *org.Aggregate, domain string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if domain = strings.TrimSpace(domain); domain == "" {
			return nil, errors.ThrowInvalidArgument(nil, "ORG-yqlVQ", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			return []eventstore.Command{org.NewDomainVerifiedEvent(ctx, &a.Aggregate, domain)}, nil
		}, nil
	}
}

func SetDomainPrimary(a *org.Aggregate, domain string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if domain = strings.TrimSpace(domain); domain == "" {
			return nil, errors.ThrowInvalidArgument(nil, "ORG-gmNqY", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			return []eventstore.Command{org.NewDomainPrimarySetEvent(ctx, &a.Aggregate, domain)}, nil
		}, nil
	}
}

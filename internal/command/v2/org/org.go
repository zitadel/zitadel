package org

import (
	"context"
	"strings"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
)

func AddOrg(a *org.Aggregate, name, iamDomain string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if name = strings.TrimSpace(name); name == "" {
			return nil, errors.ThrowInvalidArgument(nil, "ORG-mruNY", "Errors.Invalid.Argument")
		}
		defaultDomain := domain.NewIAMDomainName(name, iamDomain)
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			return []eventstore.Command{
				org.NewOrgAddedEvent(ctx, &a.Aggregate, name),
				org.NewDomainAddedEvent(ctx, &a.Aggregate, defaultDomain),
				org.NewDomainVerifiedEvent(ctx, &a.Aggregate, defaultDomain),
				org.NewDomainPrimarySetEvent(ctx, &a.Aggregate, defaultDomain),
			}, nil
		}, nil
	}
}

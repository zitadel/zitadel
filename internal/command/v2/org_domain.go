package command

import (
	"context"
	"strings"

	"github.com/caos/zitadel/internal/command"
	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
)

func AddOrgDomain(a *org.Aggregate, domain string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if domain = strings.TrimSpace(domain); domain == "" {
			return nil, errors.ThrowInvalidArgument(nil, "ORG-r3h4J", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			existing, err := orgDomain(ctx, filter, a.ID, domain)
			if err != nil {
				return nil, err
			}
			if existing.Verified {
				return nil, errors.ThrowAlreadyExists(nil, "V2-e1wse", "Errors.Already.Exists")
			}
			return []eventstore.Command{org.NewDomainAddedEvent(ctx, &a.Aggregate, domain)}, nil
		}, nil
	}
}

func VerifyOrgDomain(a *org.Aggregate, domain string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if domain = strings.TrimSpace(domain); domain == "" {
			return nil, errors.ThrowInvalidArgument(nil, "ORG-yqlVQ", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			// no checks required because unique constraints handle it
			return []eventstore.Command{org.NewDomainVerifiedEvent(ctx, &a.Aggregate, domain)}, nil
		}, nil
	}
}

func SetPrimaryOrgDomain(a *org.Aggregate, domain string) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if domain = strings.TrimSpace(domain); domain == "" {
			return nil, errors.ThrowInvalidArgument(nil, "ORG-gmNqY", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			existing, err := orgDomain(ctx, filter, a.ID, domain)
			if err != nil || existing.Primary {
				return nil, errors.ThrowAlreadyExists(err, "V2-d0Gyw", "Errors.Already.Exists")
			}
			return []eventstore.Command{org.NewDomainPrimarySetEvent(ctx, &a.Aggregate, domain)}, nil
		}, nil
	}
}

func orgDomain(ctx context.Context, filter preparation.FilterToQueryReducer, orgID, domain string) (*command.OrgDomainWriteModel, error) {
	wm := command.NewOrgDomainWriteModel(orgID, domain)
	events, err := filter(ctx, wm.Query())
	if err != nil {
		return nil, err
	}
	wm.AppendEvents(events...)
	if err = wm.Reduce(); err != nil {
		return nil, err
	}

	return wm, nil
}

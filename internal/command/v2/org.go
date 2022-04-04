package command

import (
	"context"
	"strings"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/id"
	"github.com/caos/zitadel/internal/repository/org"
	user_repo "github.com/caos/zitadel/internal/repository/user"
)

type OrgSetup struct {
	Name  string
	Human AddHuman
}

func (command *Command) SetUpOrg(ctx context.Context, instanceID string, o *OrgSetup) (*domain.ObjectDetails, error) {
	orgID, err := id.SonyFlakeGenerator.Next()
	if err != nil {
		return nil, err
	}

	userID, err := id.SonyFlakeGenerator.Next()
	if err != nil {
		return nil, err
	}

	orgAgg := org.NewAggregate(orgID, orgID)
	userAgg := user_repo.NewAggregate(userID, orgID)

	cmds, err := preparation.PrepareCommands(ctx, command.es.Filter,
		AddOrg(orgAgg, o.Name, command.iamDomain),
		AddHumanCommand(userAgg, &o.Human, command.userPasswordAlg),
		AddOrgMember(orgAgg, userID, domain.RoleOrgOwner),
	)
	if err != nil {
		return nil, err
	}

	events, err := command.es.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	return &domain.ObjectDetails{
		Sequence:      events[len(events)-1].Sequence(),
		EventDate:     events[len(events)-1].CreationDate(),
		ResourceOwner: orgID,
	}, nil
}

//AddOrg defines the commands to create a new org,
// this includes the verified default domain
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

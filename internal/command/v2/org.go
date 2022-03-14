package command

import (
	"context"

	"github.com/caos/zitadel/internal/command/v2/org"
	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/command/v2/user"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/id"
	org_repo "github.com/caos/zitadel/internal/repository/org"
	user_repo "github.com/caos/zitadel/internal/repository/user"
)

type OrgSetup struct {
	Name   string
	Domain string
	Human  user.AddHuman
}

func (command *Command) SetUpOrg(ctx context.Context, o *OrgSetup) (*domain.ObjectDetails, error) {
	orgID, err := id.SonyFlakeGenerator.Next()
	if err != nil {
		return nil, err
	}

	userID, err := id.SonyFlakeGenerator.Next()
	if err != nil {
		return nil, err
	}

	orgAgg := org_repo.NewAggregate(orgID, orgID)
	userAgg := user_repo.NewAggregate(userID, orgID)

	cmds, err := preparation.PrepareCommands(ctx, command.es.Filter,
		org.AddOrg(orgAgg, o.Name, command.iamDomain),
		org.AddDomain(orgAgg, o.Domain),
		user.AddHumanCommand(userAgg, &o.Human),
		org.AddMemberCommand(orgAgg, userID, domain.RoleOrgOwner),
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

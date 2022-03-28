package command

import (
	"context"
	"strings"

	"github.com/caos/zitadel/internal/command/v2/preparation"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/project"
)

func AddProject(
	a *project.Aggregate,
	name string,
	owner string,
	projectRoleAssertion bool,
	projectRoleCheck bool,
	hasProjectCheck bool,
	privateLabelingSetting domain.PrivateLabelingSetting,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if name = strings.TrimSpace(name); name == "" {
			return nil, errors.ThrowInvalidArgument(nil, "PROJE-C01yo", "Errors.Invalid.Argument")
		}
		if !privateLabelingSetting.Valid() {
			return nil, errors.ThrowInvalidArgument(nil, "PROJE-AO52V", "Errors.Invalid.Argument")
		}
		if owner == "" {
			return nil, errors.ThrowPreconditionFailed(nil, "PROJE-hzxwo", "Errors.Invalid.Argument")
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) ([]eventstore.Command, error) {
			return []eventstore.Command{
				project.NewProjectAddedEvent(ctx, &a.Aggregate,
					name,
					projectRoleAssertion,
					projectRoleCheck,
					hasProjectCheck,
					privateLabelingSetting,
				),
				project.NewProjectMemberAddedEvent(ctx, &a.Aggregate,
					owner,
					domain.RoleProjectOwner),
			}, nil
		}, nil
	}
}

func ExistsProject(ctx context.Context, filter preparation.FilterToQueryReducer, projectID, resourceOwner string) (exists bool, err error) {
	events, err := filter(ctx, eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(resourceOwner).
		OrderAsc().
		AddQuery().
		AggregateTypes(project.AggregateType).
		AggregateIDs(projectID).
		EventTypes(
			project.ProjectAddedType,
			project.ProjectRemovedType,
		).Builder())
	if err != nil {
		return false, err
	}

	for _, event := range events {
		switch event.(type) {
		case *project.ProjectAddedEvent:
			exists = true
		case *project.ProjectRemovedEvent:
			exists = false
		}
	}

	return exists, nil
}

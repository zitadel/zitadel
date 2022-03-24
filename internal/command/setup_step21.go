package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/user"
)

type Step21 struct{}

func (s *Step21) Step() domain.Step {
	return domain.Step21
}

func (s *Step21) execute(ctx context.Context, commandSide *Commands) error {
	return commandSide.SetupStep21(ctx, s)
}

func (c *Commands) SetupStep21(ctx context.Context, step *Step21) error {
	fn := func(iam *InstanceWriteModel) ([]eventstore.Command, error) {
		events := make([]eventstore.Command, 0)
		globalMembers := newGlobalOrgMemberWriteModel(iam.GlobalOrgID, domain.RoleOrgProjectCreator)
		orgAgg := OrgAggregateFromWriteModel(&globalMembers.WriteModel)
		if err := c.eventstore.FilterToQueryReducer(ctx, globalMembers); err != nil {
			return nil, err
		}
		for userID, roles := range globalMembers.members {
			for i, role := range roles {
				if role == domain.RoleOrgProjectCreator {
					roles[i] = domain.RoleSelfManagementGlobal
				}
			}
			events = append(events, org.NewMemberChangedEvent(ctx, orgAgg, userID, roles...))
		}
		return events, nil
	}
	return c.setup(ctx, step, fn)
}

type globalOrgMembersWriteModel struct {
	eventstore.WriteModel

	role    string
	members map[string][]string
}

func newGlobalOrgMemberWriteModel(orgID, role string) *globalOrgMembersWriteModel {
	return &globalOrgMembersWriteModel{
		WriteModel: eventstore.WriteModel{
			ResourceOwner: orgID,
			AggregateID:   orgID,
		},
		role:    role,
		members: make(map[string][]string),
	}
}

func (wm *globalOrgMembersWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *org.MemberAddedEvent:
			for _, role := range e.Roles {
				if wm.role == role {
					wm.members[e.UserID] = e.Roles
					break
				}
			}
		case *org.MemberChangedEvent:
			delete(wm.members, e.UserID)
			for _, role := range e.Roles {
				if wm.role == role {
					wm.members[e.UserID] = e.Roles
					break
				}
			}
		case *org.MemberRemovedEvent:
			delete(wm.members, e.UserID)
		case *org.MemberCascadeRemovedEvent:
			delete(wm.members, e.UserID)
		case *user.UserRemovedEvent:
			delete(wm.members, e.Aggregate().ID)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *globalOrgMembersWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			org.MemberAddedEventType,
			org.MemberChangedEventType,
			org.MemberRemovedEventType,
			org.MemberCascadeRemovedEventType,
		).
		Or().
		AggregateTypes(user.AggregateType).
		EventTypes(user.UserRemovedType).
		Builder()
}

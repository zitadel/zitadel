package command

import (
	"context"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/project"
)

func (c *Commands) removeUserMemberships(ctx context.Context, memberships []*domain.UserMembership, cascade bool) (_ []eventstore.EventPusher, err error) {
	events := make([]eventstore.EventPusher, 0)
	for _, membership := range memberships {
		switch membership.MemberType {
		case domain.MemberTypeIam:
			iamAgg := iam.NewAggregate()
			removeEvent := c.removeIAMMember(ctx, &iamAgg.Aggregate, membership.UserID, false)
			events = append(events, removeEvent)
		case domain.MemberTypeOrganisation:
			iamAgg := org.NewAggregate(membership.AggregateID, membership.ResourceOwner)
			removeEvent := c.removeOrgMember(ctx, &iamAgg.Aggregate, membership.UserID, false)
			events = append(events, removeEvent)
		case domain.MemberTypeProject:
			projectAgg := project.NewAggregate(membership.AggregateID, membership.ResourceOwner)
			removeEvent := c.removeProjectMember(ctx, &projectAgg.Aggregate, membership.UserID, false)
			events = append(events, removeEvent)
		case domain.MemberTypeProjectGrant:
			projectAgg := project.NewAggregate(membership.AggregateID, membership.ResourceOwner)
			removeEvent := c.removeProjectGrantMember(ctx, &projectAgg.Aggregate, membership.UserID, membership.ObjectID, false)
			events = append(events, removeEvent)
		}
	}
	return events, nil
}

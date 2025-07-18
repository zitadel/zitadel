package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	settings "github.com/zitadel/zitadel/internal/repository/organization_settings"
	"github.com/zitadel/zitadel/internal/repository/user"
)

type OrganizationSettingsWriteModel struct {
	eventstore.WriteModel

	OrganizationScopedUsernames bool

	OrganizationState domain.OrgState

	State           domain.OrganizationSettingsState
	checkPermission domain.PermissionCheck
}

func (wm *OrganizationSettingsWriteModel) GetWriteModel() *eventstore.WriteModel {
	return &wm.WriteModel
}

func (wm *OrganizationSettingsWriteModel) checkPermissionWrite(
	ctx context.Context,
	resourceOwner string,
	aggregateID string,
) error {
	return wm.checkPermission(ctx, domain.PermissionIAMPolicyWrite, resourceOwner, aggregateID)
}

func (wm *OrganizationSettingsWriteModel) checkPermissionDelete(
	ctx context.Context,
	resourceOwner string,
	aggregateID string,
) error {
	return wm.checkPermission(ctx, domain.PermissionIAMPolicyDelete, resourceOwner, aggregateID)
}

func NewOrganizationSettingsWriteModel(id string, checkPermission domain.PermissionCheck) *OrganizationSettingsWriteModel {
	return &OrganizationSettingsWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   id,
			ResourceOwner: id,
		},
		checkPermission: checkPermission,
	}
}

func (wm *OrganizationSettingsWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *settings.OrganizationSettingsSetEvent:
			wm.OrganizationScopedUsernames = e.OrganizationScopedUsernames
			wm.State = domain.OrganizationSettingsStateActive
		case *settings.OrganizationSettingsRemovedEvent:
			wm.OrganizationScopedUsernames = false
			wm.State = domain.OrganizationSettingsStateRemoved
		case *org.OrgAddedEvent:
			wm.OrganizationState = domain.OrgStateActive
			wm.OrganizationScopedUsernames = false
		case *org.OrgRemovedEvent:
			wm.OrganizationState = domain.OrgStateRemoved
			wm.OrganizationScopedUsernames = false
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *OrganizationSettingsWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(settings.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(settings.OrganizationSettingsSetEventType,
			settings.OrganizationSettingsRemovedEventType).
		Or().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(org.OrgAddedEventType,
			org.OrgRemovedEventType).
		Builder()
}

func (wm *OrganizationSettingsWriteModel) NewSet(
	ctx context.Context,
	organizationScopedUsernames *bool,
	userLoginMustBeDomain bool,
	usernamesF func(ctx context.Context, orgID string) ([]string, error),
) (_ []eventstore.Command, err error) {
	if err := wm.checkPermissionWrite(ctx, wm.ResourceOwner, wm.AggregateID); err != nil {
		return nil, err
	}
	// no changes
	if organizationScopedUsernames == nil || *organizationScopedUsernames == wm.OrganizationScopedUsernames {
		return nil, nil
	}

	var usernames []string
	if (wm.OrganizationScopedUsernames || userLoginMustBeDomain) != (*organizationScopedUsernames || userLoginMustBeDomain) {
		usernames, err = usernamesF(ctx, wm.AggregateID)
		if err != nil {
			return nil, err
		}
	}
	events := []eventstore.Command{
		settings.NewOrganizationSettingsAddedEvent(ctx,
			SettingsAggregateFromWriteModel(&wm.WriteModel),
			usernames,
			*organizationScopedUsernames || userLoginMustBeDomain,
			wm.OrganizationScopedUsernames || userLoginMustBeDomain,
		),
	}
	return events, nil
}

func (wm *OrganizationSettingsWriteModel) NewRemoved(
	ctx context.Context,
	userLoginMustBeDomain bool,
	usernamesF func(ctx context.Context, orgID string) ([]string, error),
) (_ []eventstore.Command, err error) {
	if err := wm.checkPermissionDelete(ctx, wm.ResourceOwner, wm.AggregateID); err != nil {
		return nil, err
	}

	var usernames []string
	if userLoginMustBeDomain != wm.OrganizationScopedUsernames {
		usernames, err = usernamesF(ctx, wm.AggregateID)
		if err != nil {
			return nil, err
		}
	}
	events := []eventstore.Command{
		settings.NewOrganizationSettingsRemovedEvent(ctx,
			SettingsAggregateFromWriteModel(&wm.WriteModel),
			usernames,
			userLoginMustBeDomain,
			wm.OrganizationScopedUsernames || userLoginMustBeDomain,
		),
	}
	return events, nil
}

func SettingsAggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return &eventstore.Aggregate{
		ID:            wm.AggregateID,
		Type:          settings.AggregateType,
		ResourceOwner: wm.ResourceOwner,
		InstanceID:    wm.InstanceID,
		Version:       settings.AggregateVersion,
	}
}

type OrganizationScopedUsernamesWriteModel struct {
	eventstore.WriteModel

	Users []*organizationScopedUser
}

type organizationScopedUser struct {
	id       string
	username string
}

func NewOrganizationScopedUsernamesWriteModel(orgID string) *OrganizationScopedUsernamesWriteModel {
	return &OrganizationScopedUsernamesWriteModel{
		WriteModel: eventstore.WriteModel{
			ResourceOwner: orgID,
		},
		Users: make([]*organizationScopedUser, 0),
	}
}

func (wm *OrganizationScopedUsernamesWriteModel) AppendEvents(events ...eventstore.Event) {
	wm.WriteModel.AppendEvents(events...)
}

func (wm *OrganizationScopedUsernamesWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *user.HumanAddedEvent:
			wm.Users = append(wm.Users, &organizationScopedUser{id: e.Aggregate().ID, username: e.UserName})
		case *user.HumanRegisteredEvent:
			wm.Users = append(wm.Users, &organizationScopedUser{id: e.Aggregate().ID, username: e.UserName})
		case *user.MachineAddedEvent:
			wm.Users = append(wm.Users, &organizationScopedUser{id: e.Aggregate().ID, username: e.UserName})
		case *user.UsernameChangedEvent:
			for _, user := range wm.Users {
				if user.id == e.Aggregate().ID {
					user.username = e.UserName
					break
				}
			}
		case *user.DomainClaimedEvent:
			for _, user := range wm.Users {
				if user.id == e.Aggregate().ID {
					user.username = e.UserName
					break
				}
			}
		case *user.UserRemovedEvent:
			wm.removeUser(e.Aggregate().ID)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *OrganizationScopedUsernamesWriteModel) removeUser(userID string) {
	for i, user := range wm.Users {
		if user.id == userID {
			wm.Users[i] = wm.Users[len(wm.Users)-1]
			wm.Users[len(wm.Users)-1] = nil
			wm.Users = wm.Users[:len(wm.Users)-1]
			return
		}
	}
}

func (wm *OrganizationScopedUsernamesWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(user.AggregateType).
		EventTypes(
			user.HumanAddedType,
			user.HumanRegisteredType,
			user.MachineAddedEventType,
			user.UserUserNameChangedType,
			user.UserDomainClaimedType,
			user.UserRemovedType,
		).
		Builder()
}

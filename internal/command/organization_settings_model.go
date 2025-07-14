package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/settings"
)

type OrganizationSettingsWriteModel struct {
	eventstore.WriteModel

	UserUniqueness bool

	OrganizationState domain.OrgState

	State                domain.OrganizationSettingsState
	writePermissionCheck bool
	checkPermission      domain.PermissionCheck
}

func (wm *OrganizationSettingsWriteModel) GetWriteModel() *eventstore.WriteModel {
	return &wm.WriteModel
}

func (wm *OrganizationSettingsWriteModel) checkPermissionWrite(
	ctx context.Context,
	resourceOwner string,
	aggregateID string,
) error {
	if wm.writePermissionCheck {
		return nil
	}
	if err := wm.checkPermission(ctx, domain.PermissionIAMPolicyWrite, resourceOwner, aggregateID); err != nil {
		return err
	}
	wm.writePermissionCheck = true
	return nil
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
			wm.UserUniqueness = e.UserUniqueness
			wm.State = domain.OrganizationSettingsStateActive
		case *settings.OrganizationSettingsRemovedEvent:
			wm.UserUniqueness = false
			wm.State = domain.OrganizationSettingsStateRemoved
		case *org.OrgAddedEvent:
			wm.OrganizationState = domain.OrgStateActive
		case *org.OrgRemovedEvent:
			wm.OrganizationState = domain.OrgStateRemoved
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
	userUniqueness *bool,
) (_ []eventstore.Command, err error) {
	if err := wm.checkPermissionWrite(ctx, wm.ResourceOwner, wm.AggregateID); err != nil {
		return nil, err
	}
	// no changes
	if userUniqueness == nil || *userUniqueness == wm.UserUniqueness {
		return nil, nil
	}
	events := []eventstore.Command{
		settings.NewOrganizationSettingsAddedEvent(ctx,
			SettingsAggregateFromWriteModel(&wm.WriteModel),
			*userUniqueness,
		),
	}
	return events, nil
}

func (wm *OrganizationSettingsWriteModel) NewRemoved(
	ctx context.Context,
) (_ []eventstore.Command, err error) {
	if err := wm.checkPermissionDelete(ctx, wm.ResourceOwner, wm.AggregateID); err != nil {
		return nil, err
	}
	events := []eventstore.Command{
		settings.NewOrganizationSettingsRemovedEvent(ctx,
			SettingsAggregateFromWriteModel(&wm.WriteModel),
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

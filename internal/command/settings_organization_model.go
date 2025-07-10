package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/settings"
)

type SettingsOrganizationWriteModel struct {
	eventstore.WriteModel

	UserUniqueness bool

	OrganizationState domain.OrgState

	State                domain.OrganizationSettingsState
	writePermissionCheck bool
	checkPermission      domain.PermissionCheck
}

func (wm *SettingsOrganizationWriteModel) GetWriteModel() *eventstore.WriteModel {
	return &wm.WriteModel
}

func (wm *SettingsOrganizationWriteModel) checkPermissionWrite(
	ctx context.Context,
	resourceOwner string,
	aggregateID string,
) error {
	if wm.writePermissionCheck {
		return nil
	}
	if err := wm.checkPermission(ctx, domain.PermissionUserWrite, resourceOwner, aggregateID); err != nil {
		return err
	}
	wm.writePermissionCheck = true
	return nil
}

func (wm *SettingsOrganizationWriteModel) checkPermissionDelete(
	ctx context.Context,
	resourceOwner string,
	aggregateID string,
) error {
	return wm.checkPermission(ctx, domain.PermissionUserDelete, resourceOwner, aggregateID)
}

func NewSettingsOrganizationWriteModel(id string, checkPermission domain.PermissionCheck) *SettingsOrganizationWriteModel {
	return &SettingsOrganizationWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   id,
			ResourceOwner: id,
		},
		checkPermission: checkPermission,
	}
}

func (wm *SettingsOrganizationWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *settings.SettingOrganizationSetEvent:
			wm.UserUniqueness = e.UserUniqueness
			wm.State = domain.OrganizationSettingsStateActive
		case *settings.SettingOrganizationRemovedEvent:
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

func (wm *SettingsOrganizationWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(settings.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(settings.SettingOrganizationSetEventType,
			settings.SettingOrganizationRemovedEventType).
		Or().
		AggregateTypes(org.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(org.OrgAddedEventType,
			org.OrgRemovedEventType).
		Builder()
}

func (wm *SettingsOrganizationWriteModel) NewSet(
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
		settings.NewSettingOrganizationAddedEvent(ctx,
			SettingsAggregateFromWriteModel(&wm.WriteModel),
			*userUniqueness,
		),
	}
	return events, nil
}

func (wm *SettingsOrganizationWriteModel) NewRemoved(
	ctx context.Context,
) (_ []eventstore.Command, err error) {
	if err := wm.checkPermissionDelete(ctx, wm.ResourceOwner, wm.AggregateID); err != nil {
		return nil, err
	}
	events := []eventstore.Command{
		settings.NewSettingOrganizationRemovedEvent(ctx,
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

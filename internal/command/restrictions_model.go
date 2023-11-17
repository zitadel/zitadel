package command

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/restrictions"
)

type restrictionsWriteModel struct {
	eventstore.WriteModel
	disallowPublicOrgRegistrations bool
}

// newRestrictionsWriteModel aggregateId is filled by reducing unit matching events
func newRestrictionsWriteModel(instanceId, resourceOwner string) *restrictionsWriteModel {
	return &restrictionsWriteModel{
		WriteModel: eventstore.WriteModel{
			InstanceID:    instanceId,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *restrictionsWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		InstanceID(wm.InstanceID).
		AddQuery().
		AggregateTypes(restrictions.AggregateType).
		EventTypes(restrictions.SetEventType)

	return query.Builder()
}

func (wm *restrictionsWriteModel) Reduce() error {
	for _, event := range wm.Events {
		wm.ChangeDate = event.CreatedAt()
		if e, ok := event.(*restrictions.SetEvent); ok && e.DisallowPublicOrgRegistrations != nil {
			wm.disallowPublicOrgRegistrations = *e.DisallowPublicOrgRegistrations
		}
	}
	return wm.WriteModel.Reduce()
}

// NewChanges returns all changes that need to be applied to the aggregate.
// nil properties in setRestrictions are ignored
func (wm *restrictionsWriteModel) NewChanges(setRestrictions *SetRestrictions) (changes []restrictions.RestrictionsChange) {
	if setRestrictions == nil {
		return nil
	}
	changes = make([]restrictions.RestrictionsChange, 0, 1)
	if setRestrictions.DisallowPublicOrgRegistration != nil && (wm.disallowPublicOrgRegistrations != *setRestrictions.DisallowPublicOrgRegistration) {
		changes = append(changes, restrictions.ChangePublicOrgRegistrations(*setRestrictions.DisallowPublicOrgRegistration))
	}
	return changes
}

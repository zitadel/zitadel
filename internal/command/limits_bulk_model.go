package command

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/limits"
)

type limitsBulkWriteModel struct {
	eventstore.WriteModel
	writeModels                     map[string]map[string]*limitsWriteModel
	filterInstanceIDs, filterOwners []string
}

// newLimitsBulkWriteModel should be followed by limitsBulkWriteModel.addWriteModel before querying and reducing it.
func newLimitsBulkWriteModel() *limitsBulkWriteModel {
	return &limitsBulkWriteModel{
		writeModels:       make(map[string]map[string]*limitsWriteModel),
		filterInstanceIDs: make([]string, 0),
		filterOwners:      make([]string, 0),
	}
}

func (wm *limitsBulkWriteModel) addWriteModel(instanceID, resourceOwner string) {
	if _, ok := wm.writeModels[instanceID]; !ok {
		wm.writeModels[instanceID] = make(map[string]*limitsWriteModel)
	}
	if _, ok := wm.writeModels[instanceID][resourceOwner]; !ok {
		wm.writeModels[instanceID][resourceOwner] = newLimitsWriteModel(instanceID, resourceOwner)
	}
	wm.filterInstanceIDs = append(wm.filterInstanceIDs, instanceID)
	wm.filterOwners = append(wm.filterOwners, resourceOwner)
}

func (wm *limitsBulkWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwners(wm.filterOwners).
		InstanceIDs(wm.filterInstanceIDs).
		AddQuery().
		AggregateTypes(limits.AggregateType).
		EventTypes(
			limits.SetEventType,
			limits.ResetEventType,
		)

	return query.Builder()
}

func (wm *limitsBulkWriteModel) Reduce() error {
	for _, event := range wm.Events {
		instanceID, resourceOwner := event.Aggregate().InstanceID, event.Aggregate().ResourceOwner
		if _, ok := wm.writeModels[instanceID]; !ok {
			continue
		}
		limitsWm, ok := wm.writeModels[instanceID][resourceOwner]
		if !ok {
			continue
		}
		limitsWm.AppendEvents(event)
		if err := limitsWm.Reduce(); err != nil {
			return err
		}
	}
	return nil
}

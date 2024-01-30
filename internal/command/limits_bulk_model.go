package command

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/limits"
)

type limitsBulkWriteModel struct {
	eventstore.WriteModel
	writeModels       map[string]*limitsWriteModel
	filterInstanceIDs []string
}

// newLimitsBulkWriteModel should be followed by limitsBulkWriteModel.addWriteModel before querying and reducing it.
func newLimitsBulkWriteModel() *limitsBulkWriteModel {
	return &limitsBulkWriteModel{
		writeModels:       make(map[string]*limitsWriteModel),
		filterInstanceIDs: make([]string, 0),
	}
}

func (wm *limitsBulkWriteModel) addWriteModel(instanceID string) {
	if _, ok := wm.writeModels[instanceID]; !ok {
		wm.writeModels[instanceID] = newLimitsWriteModel(instanceID)
	}
	wm.filterInstanceIDs = append(wm.filterInstanceIDs, instanceID)
}

func (wm *limitsBulkWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
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
		instanceID := event.Aggregate().InstanceID
		limitsWm, ok := wm.writeModels[instanceID]
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

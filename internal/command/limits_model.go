package command

import (
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/limits"
)

type limitsWriteModel struct {
	eventstore.WriteModel
	rollingAggregateID string
	auditLogRetention  time.Duration
}

// newLimitsWriteModel aggregateId is filled by reducing unit matching events
func newLimitsWriteModel(instanceId, resourceOwner string) *limitsWriteModel {
	return &limitsWriteModel{
		WriteModel: eventstore.WriteModel{
			InstanceID:    instanceId,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *limitsWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		InstanceID(wm.InstanceID).
		AddQuery().
		AggregateTypes(limits.AggregateType).
		EventTypes(
			limits.SetEventType,
			limits.ResetEventType,
		)

	return query.Builder()
}

func (wm *limitsWriteModel) Reduce() error {
	for _, event := range wm.Events {
		wm.ChangeDate = event.CreatedAt()
		switch e := event.(type) {
		case *limits.SetEvent:
			wm.rollingAggregateID = e.Aggregate().ID
			if e.AuditLogRetention != nil {
				wm.auditLogRetention = *e.AuditLogRetention
			}
		case *limits.ResetEvent:
			wm.rollingAggregateID = ""
		}
	}
	if err := wm.WriteModel.Reduce(); err != nil {
		return err
	}
	// wm.WriteModel.Reduce() sets the aggregateID to the first event's aggregateID, but we need the last one
	wm.AggregateID = wm.rollingAggregateID
	return nil
}

// NewChanges returns all changes that need to be applied to the aggregate.
// If createNew is true, all possible changes are returned.
func (wm *limitsWriteModel) NewChanges(createNew bool, setLimits *SetLimits) (changes []limits.LimitsChange) {
	if setLimits == nil {
		return nil
	}
	if createNew {
		return []limits.LimitsChange{
			limits.ChangeAuditLogRetention(setLimits.AuditLogRetention),
		}
	}
	changes = make([]limits.LimitsChange, 0, 1)
	if wm.auditLogRetention != setLimits.AuditLogRetention {
		changes = append(changes, limits.ChangeAuditLogRetention(setLimits.AuditLogRetention))
	}
	return changes
}

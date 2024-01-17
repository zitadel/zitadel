package command

import (
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/limits"
)

type limitsWriteModel struct {
	eventstore.WriteModel
	rollingAggregateID string
	auditLogRetention  *time.Duration
	block              *bool
}

// newLimitsWriteModel aggregateId is filled by reducing unit matching events
func newLimitsWriteModel(instanceId string) *limitsWriteModel {
	return &limitsWriteModel{
		WriteModel: eventstore.WriteModel{
			InstanceID:    instanceId,
			ResourceOwner: instanceId,
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
				wm.auditLogRetention = e.AuditLogRetention
			}
			if e.Block != nil {
				wm.block = e.Block
			}
		case *limits.ResetEvent:
			wm.rollingAggregateID = ""
			wm.auditLogRetention = nil
			wm.block = nil
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
// nil properties in setLimits are ignored
func (wm *limitsWriteModel) NewChanges(setLimits *SetLimits) (changes []limits.LimitsChange) {
	if setLimits == nil {
		return nil
	}
	changes = make([]limits.LimitsChange, 0, 1)
	if setLimits.AuditLogRetention != nil && (wm.auditLogRetention == nil || *wm.auditLogRetention != *setLimits.AuditLogRetention) {
		changes = append(changes, limits.ChangeAuditLogRetention(setLimits.AuditLogRetention))
	}
	if setLimits.Block != nil && (wm.block == nil || *wm.block != *setLimits.Block) {
		changes = append(changes, limits.ChangeBlock(setLimits.Block))
	}
	return changes
}

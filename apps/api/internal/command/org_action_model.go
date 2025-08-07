package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/action"
)

type ActionWriteModel struct {
	eventstore.WriteModel

	Name          string
	Script        string
	Timeout       time.Duration
	AllowedToFail bool
	State         domain.ActionState
}

func NewActionWriteModel(actionID string, resourceOwner string) *ActionWriteModel {
	return &ActionWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   actionID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *ActionWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *action.AddedEvent:
			wm.Name = e.Name
			wm.Script = e.Script
			wm.Timeout = e.Timeout
			wm.AllowedToFail = e.AllowedToFail
			wm.State = domain.ActionStateActive
		case *action.ChangedEvent:
			if e.Name != nil {
				wm.Name = *e.Name
			}
			if e.Script != nil {
				wm.Script = *e.Script
			}
			if e.Timeout != nil {
				wm.Timeout = *e.Timeout
			}
			if e.AllowedToFail != nil {
				wm.AllowedToFail = *e.AllowedToFail
			}
		case *action.DeactivatedEvent:
			wm.State = domain.ActionStateInactive
		case *action.ReactivatedEvent:
			wm.State = domain.ActionStateActive
		case *action.RemovedEvent:
			wm.State = domain.ActionStateRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *ActionWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(action.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(action.AddedEventType,
			action.ChangedEventType,
			action.DeactivatedEventType,
			action.ReactivatedEventType,
			action.RemovedEventType).
		Builder()
}

func (wm *ActionWriteModel) NewChangedEvent(
	ctx context.Context,
	agg *eventstore.Aggregate,
	name string,
	script string,
	timeout time.Duration,
	allowedToFail bool,
) (*action.ChangedEvent, error) {
	changes := make([]action.ActionChanges, 0)
	if wm.Name != name {
		changes = append(changes, action.ChangeName(name, wm.Name))
	}
	if wm.Script != script {
		changes = append(changes, action.ChangeScript(script))
	}
	if wm.Timeout != timeout {
		changes = append(changes, action.ChangeTimeout(timeout))
	}
	if wm.AllowedToFail != allowedToFail {
		changes = append(changes, action.ChangeAllowedToFail(allowedToFail))
	}
	return action.NewChangedEvent(ctx, agg, changes)
}

func ActionAggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return eventstore.AggregateFromWriteModel(wm, action.AggregateType, action.AggregateVersion)
}

func NewActionAggregate(id, resourceOwner string) *eventstore.Aggregate {
	return ActionAggregateFromWriteModel(&eventstore.WriteModel{
		AggregateID:   id,
		ResourceOwner: resourceOwner,
	})
}

type ActionExistsModel struct {
	eventstore.WriteModel

	actionIDs  []string
	checkedIDs []string
}

func NewActionsExistModel(actionIDs []string, resourceOwner string) *ActionExistsModel {
	return &ActionExistsModel{
		WriteModel: eventstore.WriteModel{
			ResourceOwner: resourceOwner,
		},
		actionIDs: actionIDs,
	}
}

func (wm *ActionExistsModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *action.AddedEvent:
			wm.checkedIDs = append(wm.checkedIDs, e.Aggregate().ID)
		case *action.RemovedEvent:
			for i := len(wm.checkedIDs) - 1; i >= 0; i-- {
				if wm.checkedIDs[i] == e.Aggregate().ID {
					wm.checkedIDs[i] = wm.checkedIDs[len(wm.checkedIDs)-1]
					wm.checkedIDs[len(wm.checkedIDs)-1] = ""
					wm.checkedIDs = wm.checkedIDs[:len(wm.checkedIDs)-1]
					break
				}
			}
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *ActionExistsModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(action.AggregateType).
		AggregateIDs(wm.actionIDs...).
		EventTypes(action.AddedEventType,
			action.RemovedEventType).
		Builder()
}

type ActionsListByOrgModel struct {
	eventstore.WriteModel

	Actions map[string]*ActionWriteModel
}

func NewActionsListByOrgModel(resourceOwner string) *ActionsListByOrgModel {
	return &ActionsListByOrgModel{
		WriteModel: eventstore.WriteModel{
			ResourceOwner: resourceOwner,
		},
		Actions: make(map[string]*ActionWriteModel),
	}
}

func (wm *ActionsListByOrgModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *action.AddedEvent:
			wm.Actions[e.Aggregate().ID] = &ActionWriteModel{
				WriteModel: eventstore.WriteModel{
					AggregateID: e.Aggregate().ID,
					ChangeDate:  e.CreationDate(),
				},
				Name:  e.Name,
				State: domain.ActionStateActive,
			}
		case *action.DeactivatedEvent:
			wm.Actions[e.Aggregate().ID].State = domain.ActionStateInactive
		case *action.ReactivatedEvent:
			wm.Actions[e.Aggregate().ID].State = domain.ActionStateActive
		case *action.RemovedEvent:
			delete(wm.Actions, e.Aggregate().ID)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *ActionsListByOrgModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(action.AggregateType).
		EventTypes(action.AddedEventType,
			action.DeactivatedEventType,
			action.ReactivatedEventType,
			action.RemovedEventType).
		Builder()
}

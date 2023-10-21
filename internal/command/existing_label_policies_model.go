package command

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
)

type ExistingLabelPoliciesReadModel struct {
	eventstore.WriteModel

	aggregateIDs []string
}

func NewExistingLabelPoliciesReadModel(ctx context.Context) *ExistingLabelPoliciesReadModel {
	return &ExistingLabelPoliciesReadModel{}
}

func (rm *ExistingLabelPoliciesReadModel) AppendEvents(events ...eventstore.Event) {
	rm.WriteModel.AppendEvents(events...)
}

func (rm *ExistingLabelPoliciesReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *instance.LabelPolicyAddedEvent,
			*org.LabelPolicyAddedEvent:
			rm.aggregateIDs = append(rm.aggregateIDs, e.Aggregate().ID)
		case *org.LabelPolicyRemovedEvent:
			for i := len(rm.aggregateIDs) - 1; i >= 0; i-- {
				if rm.aggregateIDs[i] == e.Aggregate().ID {
					copy(rm.aggregateIDs[i:], rm.aggregateIDs[i+1:])
					rm.aggregateIDs[len(rm.aggregateIDs)-1] = ""
					rm.aggregateIDs = rm.aggregateIDs[:len(rm.aggregateIDs)-1]
				}
			}
		}
	}
	return rm.WriteModel.Reduce()
}

func (rm *ExistingLabelPoliciesReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		EventTypes(instance.LabelPolicyAddedEventType).
		Or().
		AggregateTypes(org.AggregateType).
		EventTypes(
			org.LabelPolicyAddedEventType,
			org.LabelPolicyRemovedEventType).
		Builder()
}

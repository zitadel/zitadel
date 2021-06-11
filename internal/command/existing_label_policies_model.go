package command

import (
	"context"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/org"
)

type ExistingLabelPoliciesReadModel struct {
	eventstore.WriteModel

	aggregateIDs []string
}

func NewExistingLabelPoliciesReadModel(ctx context.Context) *ExistingLabelPoliciesReadModel {
	return &ExistingLabelPoliciesReadModel{}
}

func (rm *ExistingLabelPoliciesReadModel) AppendEvents(events ...eventstore.EventReader) {
	rm.WriteModel.AppendEvents(events...)
}

func (rm *ExistingLabelPoliciesReadModel) Reduce() error {
	for _, event := range rm.Events {
		switch e := event.(type) {
		case *iam.LabelPolicyAddedEvent,
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
	return nil
}

func (rm *ExistingLabelPoliciesReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(
		eventstore.ColumnsEvent,
		iam.AggregateType,
		org.AggregateType).
		EventTypes(
			iam.LabelPolicyAddedEventType,
			org.LabelPolicyAddedEventType,
			org.LabelPolicyRemovedEventType,
		)
}

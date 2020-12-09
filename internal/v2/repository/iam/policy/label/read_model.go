package label

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/label"
)

type LabelPolicyReadModel struct{ label.LabelPolicyReadModel }

func (rm *LabelPolicyReadModel) AppendEvents(events ...eventstore.EventReader) {
	for _, event := range events {
		switch e := event.(type) {
		case *LabelPolicyAddedEvent:
			rm.ReadModel.AppendEvents(&e.LabelPolicyAddedEvent)
		case *LabelPolicyChangedEvent:
			rm.ReadModel.AppendEvents(&e.LabelPolicyChangedEvent)
		case *label.LabelPolicyAddedEvent, *label.LabelPolicyChangedEvent:
			rm.ReadModel.AppendEvents(e)
		}
	}
}

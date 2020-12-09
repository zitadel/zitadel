package label

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/repository/policy/label"
)

const (
	AggregateType = "iam"
)

type LabelPolicyWriteModel struct {
	eventstore.WriteModel
	Policy label.LabelPolicyWriteModel

	iamID string
}

func NewLabelPolicyWriteModel(iamID string) *LabelPolicyWriteModel {
	return &LabelPolicyWriteModel{
		iamID: iamID,
	}
}

func (wm *LabelPolicyWriteModel) AppendEvents(events ...eventstore.EventReader) {
	wm.WriteModel.AppendEvents(events...)
	for _, event := range events {
		switch e := event.(type) {
		case *LabelPolicyAddedEvent:
			wm.Policy.AppendEvents(&e.LabelPolicyAddedEvent)
		case *LabelPolicyChangedEvent:
			wm.Policy.AppendEvents(&e.LabelPolicyChangedEvent)
		}
	}
}

func (wm *LabelPolicyWriteModel) Reduce() error {
	if err := wm.Policy.Reduce(); err != nil {
		return err
	}
	return wm.WriteModel.Reduce()
}

func (wm *LabelPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent, AggregateType).
		AggregateIDs(wm.iamID)
}

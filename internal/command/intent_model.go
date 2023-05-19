package command

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/intent"
)

type IntentWriteModel struct {
	eventstore.WriteModel

	IDPID      string
	SuccessURL string
	FailureURL string
	Token      string
	UserID     string
	UserInfo   []byte
}

func NewIntentWriteModel(id, resourceOwner string) *IntentWriteModel {
	return &IntentWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   id,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *IntentWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *intent.IntentAddedEvent:
			wm.reduceIntentAddedType(e)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *IntentWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(intent.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			intent.IntentAddedType,
		).
		Builder()
}

func (wm *IntentWriteModel) reduceIntentAddedType(e *intent.IntentAddedEvent) {
	wm.IDPID = e.IDPID
	wm.SuccessURL = e.SuccessURL
	wm.FailureURL = e.FailureURL
}

package command

import (
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/instance"
)

type InstanceCustomTextWriteModel struct {
	CustomTextWriteModel
}

func NewInstanceCustomTextWriteModel(instanceID, key string, language language.Tag) *InstanceCustomTextWriteModel {
	return &InstanceCustomTextWriteModel{
		CustomTextWriteModel{
			WriteModel: eventstore.WriteModel{
				AggregateID:   instanceID,
				ResourceOwner: instanceID,
			},
			Key:      key,
			Language: language,
		},
	}
}

func (wm *InstanceCustomTextWriteModel) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch e := event.(type) {
		case *instance.CustomTextSetEvent:
			wm.CustomTextWriteModel.AppendEvents(&e.CustomTextSetEvent)
		}
	}
}

func (wm *InstanceCustomTextWriteModel) Reduce() error {
	return wm.CustomTextWriteModel.Reduce()
}

func (wm *InstanceCustomTextWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateIDs(wm.CustomTextWriteModel.AggregateID).
		AggregateTypes(instance.AggregateType).
		EventTypes(
			instance.CustomTextSetEventType).
		Builder()
}

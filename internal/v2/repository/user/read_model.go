package user

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
)

type ReadModel struct {
	eventstore.ReadModel
}

func NewReadModel(id string) *ReadModel {
	return &ReadModel{
		ReadModel: eventstore.ReadModel{
			AggregateID: id,
		},
	}
}

func (rm *ReadModel) AppendEvents(events ...eventstore.EventReader) {
	rm.ReadModel.AppendEvents(events...)
	for _, event := range events {
		switch event.(type) {

		}
	}
}

func (rm *ReadModel) Reduce() (err error) {
	for _, event := range rm.Events {
		switch e := event.(type) {
		}
	}
	for _, reduce := range []func() error{
		rm.ReadModel.Reduce,
	} {
		if err = reduce(); err != nil {
			return err
		}
	}

	return nil
}

func (rm *ReadModel) AppendAndReduce(events ...eventstore.EventReader) error {
	rm.AppendEvents(events...)
	return rm.Reduce()
}

func (rm *ReadModel) Query() *eventstore.SearchQueryFactory {
	return eventstore.NewSearchQueryFactory(eventstore.ColumnsEvent, AggregateType).AggregateIDs(rm.AggregateID)
}

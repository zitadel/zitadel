package query

import (
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/user"
)

type UserReadModel struct {
	eventstore.ReadModel
}

func NewUserReadModel(id string) *UserReadModel {
	return &UserReadModel{
		ReadModel: eventstore.ReadModel{
			AggregateID: id,
		},
	}
}

func (rm *UserReadModel) AppendEvents(events ...eventstore.EventReader) {
	rm.ReadModel.AppendEvents(events...)
	for _, event := range events {
		switch event.(type) {
		// TODO: implement append events
		}
	}
}

func (rm *UserReadModel) Reduce() (err error) {
	for _, event := range rm.Events {
		switch event.(type) {
		//TODO: implement reduce
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

func (rm *UserReadModel) AppendAndReduce(events ...eventstore.EventReader) error {
	rm.AppendEvents(events...)
	return rm.Reduce()
}

func (rm *UserReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(rm.AggregateID).
		SearchQueryBuilder()
}

func NewUserEventSearchQuery(userID, orgID string, sequence uint64) *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(orgID).
		AddQuery().
		AggregateTypes(user.AggregateType).
		AggregateIDs(userID).
		SequenceGreater(sequence).
		SearchQueryBuilder()
}

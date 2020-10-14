package eventstore_test

import (
	"context"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
)

type singleAggregateRepo struct {
	events []*repository.Event
}

//Health checks if the connection to the storage is available
func (r *singleAggregateRepo) Health(ctx context.Context) error {
	return nil
}

// PushEvents adds all events of the given aggregates to the eventstreams of the aggregates.
// This call is transaction save. The transaction will be rolled back if one event fails
func (r *singleAggregateRepo) Push(ctx context.Context, events ...*repository.Event) error {
	for _, event := range events {
		if event.AggregateType != "test.agg" || event.AggregateID != "test" {
			return errors.ThrowPreconditionFailed(nil, "V2-ZVDcA", "wrong aggregate")
		}
	}

	r.events = append(r.events, events...)

	return nil
}

// Filter returns all events matching the given search query
func (r *singleAggregateRepo) Filter(ctx context.Context, searchQuery *repository.SearchQuery) (events []*repository.Event, err error) {
	return r.events, nil
}

//LatestSequence returns the latests sequence found by the the search query
func (r *singleAggregateRepo) LatestSequence(ctx context.Context, queryFactory *repository.SearchQuery) (uint64, error) {
	if len(r.events) == 0 {
		return 0, nil
	}
	return r.events[len(r.events)-1].Sequence, nil
}

type UserAggregate struct {
	FirstName string
}

func (a *UserAggregate) ID() string {
	return "test"
}
func (a *UserAggregate) Type() eventstore.AggregateType {
	return "test.agg"
}
func (a *UserAggregate) Events() []eventstore.Event {
	return nil
}
func (a *UserAggregate) ResourceOwner() string {
	return "caos"
}
func (a *UserAggregate) Version() eventstore.Version {
	return "v1"
}
func (a *UserAggregate) PreviousSequence() uint64 {
	return 0
}

type UserAddedEvent struct {
	FirstName string
}

func (e *UserAddedEvent) CheckPrevious() bool {
	return false
}

func (e *UserAddedEvent) EditorService() string {
	return "test.suite"
}

func (e *UserAddedEvent) EditorUser() string {
	return "adlerhurst"
}

func (e *UserAddedEvent) Type() eventstore.EventType {
	return "user.added"
}
func (e *UserAddedEvent) Data() interface{} {
	return e
}

type UserFirstNameChangedEvent struct {
	FirstName string
}

func (e *UserFirstNameChangedEvent) CheckPrevious() bool {
	return false
}

func (e *UserFirstNameChangedEvent) EditorService() string {
	return "test.suite"
}

func (e *UserFirstNameChangedEvent) EditorUser() string {
	return "adlerhurst"
}

func (e *UserFirstNameChangedEvent) Type() eventstore.EventType {
	return "user.changed"
}
func (e *UserFirstNameChangedEvent) Data() interface{} {
	return e
}

type UserReadModel struct {
	eventstore.ReadModel
	FirstName string
}

func (rm *UserReadModel) AppendEvents(events ...eventstore.Event) error {
	rm.ReadModel.Append(events...)
	return nil
}

func (rm *UserReadModel) Reduce() error {
	for _, event := range rm.ReadModel.Events {
		switch e := event.(type) {
		case *UserAddedEvent:
			rm.FirstName = e.FirstName
		case *UserFirstNameChangedEvent:
			rm.FirstName = e.FirstName
		}
	}
	return nil
}

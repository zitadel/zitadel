package eventstore_test

import (
	"context"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	query_repo "github.com/zitadel/zitadel/internal/eventstore/repository/sql"
	v3 "github.com/zitadel/zitadel/internal/eventstore/v3"
)

// ------------------------------------------------------------
// User aggregate start
// ------------------------------------------------------------
func NewUserAggregate(id string) *eventstore.Aggregate {
	return eventstore.NewAggregate(
		authz.NewMockContext("zitadel", "caos", "adlerhurst"),
		id,
		"test.user",
		"v1",
	)
}

// ------------------------------------------------------------
// User added event start
// ------------------------------------------------------------

type UserAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	FirstName string `json:"firstName"`
}

func NewUserAddedEvent(id string, firstName string) *UserAddedEvent {
	return &UserAddedEvent{
		FirstName: firstName,
		BaseEvent: *eventstore.NewBaseEventForPush(
			context.Background(),
			NewUserAggregate(id),
			"user.added"),
	}
}

func UserAddedEventMapper() (eventstore.AggregateType, eventstore.EventType, func(eventstore.Event) (eventstore.Event, error)) {
	return "user", "user.added", func(event eventstore.Event) (eventstore.Event, error) {
		e := &UserAddedEvent{
			BaseEvent: *eventstore.BaseEventFromRepo(event),
		}

		err := event.Unmarshal(e)
		if err != nil {
			return nil, err
		}
		return e, nil
	}
}

func (e *UserAddedEvent) Payload() interface{} {
	return e
}

func (e *UserAddedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *UserAddedEvent) Assets() []*eventstore.Asset {
	return nil
}

// ------------------------------------------------------------
// User first name changed event start
// ------------------------------------------------------------

type UserFirstNameChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	FirstName string `json:"firstName"`
}

func NewUserFirstNameChangedEvent(id, firstName string) *UserFirstNameChangedEvent {
	return &UserFirstNameChangedEvent{
		FirstName: firstName,
		BaseEvent: *eventstore.NewBaseEventForPush(
			context.Background(),
			NewUserAggregate(id),
			"user.firstname.changed"),
	}
}

func UserFirstNameChangedMapper() (eventstore.AggregateType, eventstore.EventType, func(eventstore.Event) (eventstore.Event, error)) {
	return "user", "user.firstName.changed", func(event eventstore.Event) (eventstore.Event, error) {
		e := &UserFirstNameChangedEvent{
			BaseEvent: *eventstore.BaseEventFromRepo(event),
		}
		err := event.Unmarshal(e)
		if err != nil {
			return nil, err
		}
		return e, nil
	}
}

func (e *UserFirstNameChangedEvent) Payload() interface{} {
	return e
}

func (e *UserFirstNameChangedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *UserFirstNameChangedEvent) Assets() []*eventstore.Asset {
	return nil
}

// ------------------------------------------------------------
// User password checked event start
// ------------------------------------------------------------

type UserPasswordCheckedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func NewUserPasswordCheckedEvent(id string) *UserPasswordCheckedEvent {
	return &UserPasswordCheckedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			context.Background(),
			NewUserAggregate(id),
			"user.password.checked"),
	}
}

func UserPasswordCheckedMapper() (eventstore.AggregateType, eventstore.EventType, func(eventstore.Event) (eventstore.Event, error)) {
	return "user", "user.password.checked", func(event eventstore.Event) (eventstore.Event, error) {
		return &UserPasswordCheckedEvent{
			BaseEvent: *eventstore.BaseEventFromRepo(event),
		}, nil
	}
}

func (e *UserPasswordCheckedEvent) Payload() interface{} {
	return nil
}

func (e *UserPasswordCheckedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *UserPasswordCheckedEvent) Assets() []*eventstore.Asset {
	return nil
}

// ------------------------------------------------------------
// User deleted event
// ------------------------------------------------------------

type UserDeletedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func NewUserDeletedEvent(id string) *UserDeletedEvent {
	return &UserDeletedEvent{
		BaseEvent: *eventstore.NewBaseEventForPush(
			context.Background(),
			NewUserAggregate(id),
			"user.deleted"),
	}
}

func UserDeletedMapper() (eventstore.AggregateType, eventstore.EventType, func(eventstore.Event) (eventstore.Event, error)) {
	return "user", "user.deleted", func(event eventstore.Event) (eventstore.Event, error) {
		return &UserDeletedEvent{
			BaseEvent: *eventstore.BaseEventFromRepo(event),
		}, nil
	}
}

func (e *UserDeletedEvent) Payload() interface{} {
	return nil
}

func (e *UserDeletedEvent) UniqueConstraints() []*eventstore.UniqueConstraint {
	return nil
}

func (e *UserDeletedEvent) Assets() []*eventstore.Asset {
	return nil
}

// ------------------------------------------------------------
// Users read model start
// ------------------------------------------------------------

type UsersReadModel struct {
	eventstore.ReadModel

	Users []*UserReadModel
}

func (rm *UsersReadModel) AppendEvents(events ...eventstore.Event) {
	rm.ReadModel.AppendEvents(events...)
	for _, event := range events {
		switch e := event.(type) {
		case *UserAddedEvent:
			//insert
			user := NewUserReadModel(e.Aggregate().ID)
			rm.Users = append(rm.Users, user)
			user.AppendEvents(e)
		case *UserFirstNameChangedEvent, *UserPasswordCheckedEvent:
			//update
			_, user := rm.userByID(e.Aggregate().ID)
			if user == nil {
				return
			}
			user.AppendEvents(e)
		case *UserDeletedEvent:
			idx, _ := rm.userByID(e.Aggregate().ID)
			if idx < 0 {
				return
			}
			copy(rm.Users[idx:], rm.Users[idx+1:])
			rm.Users[len(rm.Users)-1] = nil // or the zero value of T
			rm.Users = rm.Users[:len(rm.Users)-1]
		}
	}
}

func (rm *UsersReadModel) Reduce() error {
	for _, user := range rm.Users {
		err := user.Reduce()
		if err != nil {
			return err
		}
	}
	rm.ReadModel.Reduce()
	return nil
}

func (rm *UsersReadModel) userByID(id string) (idx int, user *UserReadModel) {
	for idx, user = range rm.Users {
		if user.ID == id {
			return idx, user
		}
	}

	return -1, nil
}

// ------------------------------------------------------------
// User read model start
// ------------------------------------------------------------

type UserReadModel struct {
	eventstore.ReadModel

	ID                string
	FirstName         string
	pwCheckCount      int
	lastPasswordCheck time.Time
}

func NewUserReadModel(id string) *UserReadModel {
	return &UserReadModel{
		ID: id,
	}
}

func (rm *UserReadModel) Reduce() error {
	for _, event := range rm.ReadModel.Events {
		switch e := event.(type) {
		case *UserAddedEvent:
			rm.FirstName = e.FirstName
		case *UserFirstNameChangedEvent:
			rm.FirstName = e.FirstName
		case *UserPasswordCheckedEvent:
			rm.pwCheckCount++
			rm.lastPasswordCheck = e.CreationDate()
		}
	}
	rm.ReadModel.Reduce()
	return nil
}

// ------------------------------------------------------------
// Tests
// ------------------------------------------------------------

func TestUserReadModel(t *testing.T) {
	es := eventstore.NewEventstore(
		&eventstore.Config{
			Querier: query_repo.NewPostgres(testClient),
			Pusher:  v3.NewEventstore(testClient),
		},
	)

	eventstore.RegisterFilterEventMapper(UserAddedEventMapper())
	eventstore.RegisterFilterEventMapper(UserFirstNameChangedMapper())
	eventstore.RegisterFilterEventMapper(UserPasswordCheckedMapper())
	eventstore.RegisterFilterEventMapper(UserDeletedMapper())

	events, err := es.Push(context.Background(),
		NewUserAddedEvent("1", "hodor"),
		NewUserAddedEvent("2", "hodor"),
		NewUserPasswordCheckedEvent("2"),
		NewUserPasswordCheckedEvent("2"),
		NewUserFirstNameChangedEvent("2", "ueli"),
		NewUserDeletedEvent("2"))

	if err != nil {
		t.Errorf("unexpected error on push aggregates: %v", err)
	}

	events = append(events, nil)

	t.Logf("%+v\n", events)

	users := UsersReadModel{}

	err = es.FilterToReducer(context.Background(), eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).AddQuery().AggregateTypes("test.user").Builder(), &users)
	if err != nil {
		t.Errorf("unexpected error on filter to reducer: %v", err)
	}
	t.Logf("%+v", users)
}

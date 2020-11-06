package eventstore_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/eventstore/v2/repository"
	"github.com/caos/zitadel/internal/eventstore/v2/repository/sql"
)

// ------------------------------------------------------------
// User aggregate start
// ------------------------------------------------------------

type UserAggregate struct {
	eventstore.BaseEvent

	eventstore.Aggregate
	FirstName string
}

func (a *UserAggregate) ID() string {
	return a.Aggregate.ID
}
func (a *UserAggregate) Type() eventstore.AggregateType {
	return "test.user"
}
func (a *UserAggregate) Events() []eventstore.EventPusher {
	events := make([]eventstore.EventPusher, len(a.Aggregate.Events))
	for i, event := range a.Aggregate.Events {
		events[i] = event
	}

	return events
}
func (a *UserAggregate) ResourceOwner() string {
	return "caos"
}
func (a *UserAggregate) Version() eventstore.Version {
	return "v1"
}
func (a *UserAggregate) PreviousSequence() uint64 {
	return a.Aggregate.PreviousSequence
}

func NewUserAggregate(id string) *UserAggregate {
	return &UserAggregate{
		Aggregate: *eventstore.NewAggregate(id),
	}
}

func (rm *UserAggregate) AppendEvents(events ...eventstore.Event) *UserAggregate {
	rm.Aggregate.AppendEvents(events...)
	return rm
}

func (rm *UserAggregate) Reduce() error {
	for _, event := range rm.Aggregate.Events {
		switch e := event.(type) {
		case *UserAddedEvent:
			rm.FirstName = e.FirstName
		case *UserFirstNameChangedEvent:
			rm.FirstName = e.FirstName
		}
	}
	return rm.Aggregate.Reduce()
}

// ------------------------------------------------------------
// User added event start
// ------------------------------------------------------------

type UserAddedEvent struct {
	eventstore.BaseEvent `json:"-"`

	FirstName string `json:"firstName"`
}

func NewUserAddedEvent(firstName string) *UserAddedEvent {
	return &UserAddedEvent{
		FirstName: firstName,
		BaseEvent: eventstore.BaseEvent{
			Service:   "test.suite",
			User:      "adlerhurst",
			EventType: "user.added",
		},
	}
}

func UserAddedEventMapper() (eventstore.EventType, func(*repository.Event) (eventstore.EventReader, error)) {
	return "user.added", func(event *repository.Event) (eventstore.EventReader, error) {
		e := &UserAddedEvent{
			BaseEvent: *eventstore.BaseEventFromRepo(event),
		}
		err := json.Unmarshal(event.Data, e)
		if err != nil {
			return nil, err
		}
		return e, nil
	}
}

func (e *UserAddedEvent) CheckPrevious() bool {
	return true
}

func (e *UserAddedEvent) Data() interface{} {
	return e
}

// ------------------------------------------------------------
// User first name changed event start
// ------------------------------------------------------------

type UserFirstNameChangedEvent struct {
	eventstore.BaseEvent `json:"-"`

	FirstName string `json:"firstName"`
}

func NewUserFirstNameChangedEvent(firstName string) *UserFirstNameChangedEvent {
	return &UserFirstNameChangedEvent{
		FirstName: firstName,
		BaseEvent: eventstore.BaseEvent{
			Service:   "test.suite",
			User:      "adlerhurst",
			EventType: "user.firstName.changed",
		},
	}
}

func UserFirstNameChangedMapper() (eventstore.EventType, func(*repository.Event) (eventstore.EventReader, error)) {
	return "user.firstName.changed", func(event *repository.Event) (eventstore.EventReader, error) {
		e := &UserFirstNameChangedEvent{
			BaseEvent: *eventstore.BaseEventFromRepo(event),
		}
		err := json.Unmarshal(event.Data, e)
		if err != nil {
			return nil, err
		}
		return e, nil
	}
}

func (e *UserFirstNameChangedEvent) CheckPrevious() bool {
	return true
}

func (e *UserFirstNameChangedEvent) Data() interface{} {
	return e
}

// ------------------------------------------------------------
// User password checked event start
// ------------------------------------------------------------

type UserPasswordCheckedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func NewUserPasswordCheckedEvent() *UserPasswordCheckedEvent {
	return &UserPasswordCheckedEvent{
		BaseEvent: eventstore.BaseEvent{
			Service:   "test.suite",
			User:      "adlerhurst",
			EventType: "user.password.checked",
		},
	}
}

func UserPasswordCheckedMapper() (eventstore.EventType, func(*repository.Event) (eventstore.EventReader, error)) {
	return "user.password.checked", func(event *repository.Event) (eventstore.EventReader, error) {
		return &UserPasswordCheckedEvent{
			BaseEvent: *eventstore.BaseEventFromRepo(event),
		}, nil
	}
}

func (e *UserPasswordCheckedEvent) CheckPrevious() bool {
	return false
}

func (e *UserPasswordCheckedEvent) Data() interface{} {
	return nil
}

// ------------------------------------------------------------
// User deleted event
// ------------------------------------------------------------

type UserDeletedEvent struct {
	eventstore.BaseEvent `json:"-"`
}

func NewUserDeletedEvent() *UserDeletedEvent {
	return &UserDeletedEvent{
		BaseEvent: eventstore.BaseEvent{
			Service:   "test.suite",
			User:      "adlerhurst",
			EventType: "user.deleted",
		},
	}
}

func UserDeletedMapper() (eventstore.EventType, func(*repository.Event) (eventstore.EventReader, error)) {
	return "user.deleted", func(event *repository.Event) (eventstore.EventReader, error) {
		return &UserDeletedEvent{
			BaseEvent: *eventstore.BaseEventFromRepo(event),
		}, nil
	}
}

func (e *UserDeletedEvent) CheckPrevious() bool {
	return false
}

func (e *UserDeletedEvent) Data() interface{} {
	return nil
}

// ------------------------------------------------------------
// Users read model start
// ------------------------------------------------------------

type UsersReadModel struct {
	eventstore.ReadModel
	Users []*UserReadModel
}

func NewUsersReadModel() *UsersReadModel {
	return &UsersReadModel{
		ReadModel: *eventstore.NewReadModel(""),
		Users:     []*UserReadModel{},
	}
}

func (rm *UsersReadModel) AppendEvents(events ...eventstore.EventReader) (err error) {
	rm.ReadModel.AppendEvents(events...)
	for _, event := range events {
		switch e := event.(type) {
		case *UserAddedEvent:
			//insert
			user := NewUserReadModel(e.AggregateID())
			rm.Users = append(rm.Users, user)
			err = user.AppendEvents(e)
		case *UserFirstNameChangedEvent, *UserPasswordCheckedEvent:
			//update
			_, user := rm.userByID(e.AggregateID())
			if user == nil {
				return errors.New("user not found")
			}
			err = user.AppendEvents(e)
		case *UserDeletedEvent:
			idx, _ := rm.userByID(e.AggregateID())
			if idx < 0 {
				return nil
			}
			copy(rm.Users[idx:], rm.Users[idx+1:])
			rm.Users[len(rm.Users)-1] = nil // or the zero value of T
			rm.Users = rm.Users[:len(rm.Users)-1]
		}
		if err != nil {
			return err
		}
	}

	return nil
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
		if user.ReadModel.ID == id {
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
	FirstName         string
	pwCheckCount      int
	lastPasswordCheck time.Time
}

func NewUserReadModel(id string) *UserReadModel {
	return &UserReadModel{
		ReadModel: *eventstore.NewReadModel(id),
	}
}

func (rm *UserReadModel) AppendEvents(events ...eventstore.EventReader) error {
	rm.ReadModel.AppendEvents(events...)
	return nil
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
	es := eventstore.NewEventstore(sql.NewCRDB(testCRDBClient))
	es.RegisterFilterEventMapper(UserAddedEventMapper()).
		RegisterFilterEventMapper(UserFirstNameChangedMapper()).
		RegisterFilterEventMapper(UserPasswordCheckedMapper()).
		RegisterFilterEventMapper(UserDeletedMapper())

	events, err := es.PushAggregates(context.Background(),
		NewUserAggregate("1").AppendEvents(NewUserAddedEvent("hodor")),
		NewUserAggregate("2").AppendEvents(NewUserAddedEvent("hodor"), NewUserPasswordCheckedEvent(), NewUserPasswordCheckedEvent(), NewUserFirstNameChangedEvent("ueli")),
		NewUserAggregate("2").AppendEvents(NewUserDeletedEvent()),
	)
	if err != nil {
		t.Errorf("unexpected error on push aggregates: %v", err)
	}

	events = append(events, nil)

	fmt.Printf("%+v\n", events)

	users := NewUsersReadModel()
	err = es.FilterToReducer(context.Background(), eventstore.NewSearchQueryFactory(eventstore.ColumnsEvent, "test.user"), users)
	if err != nil {
		t.Errorf("unexpected error on filter to reducer: %v", err)
	}
	fmt.Printf("%+v", users)
}

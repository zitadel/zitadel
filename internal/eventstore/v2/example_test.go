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
	eventstore.Aggregate
	FirstName string
}

func (a *UserAggregate) ID() string {
	return a.Aggregate.ID
}
func (a *UserAggregate) Type() eventstore.AggregateType {
	return "test.user"
}
func (a *UserAggregate) Events() []eventstore.Event {
	return a.Aggregate.Events
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
	FirstName string `json:"firstName"`
	metaData  *eventstore.EventMetaData
}

func UserAddedEventMapper() (eventstore.EventType, func(*repository.Event) (eventstore.Event, error)) {
	return "user.added", func(event *repository.Event) (eventstore.Event, error) {
		e := &UserAddedEvent{
			metaData: eventstore.MetaDataFromRepo(event),
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

func (e *UserAddedEvent) MetaData() *eventstore.EventMetaData {
	return e.metaData
}

// ------------------------------------------------------------
// User first name changed event start
// ------------------------------------------------------------

type UserFirstNameChangedEvent struct {
	FirstName string                    `json:"firstName"`
	metaData  *eventstore.EventMetaData `json:"-"`
}

func UserFirstNameChangedMapper() (eventstore.EventType, func(*repository.Event) (eventstore.Event, error)) {
	return "user.firstName.changed", func(event *repository.Event) (eventstore.Event, error) {
		e := &UserFirstNameChangedEvent{
			metaData: eventstore.MetaDataFromRepo(event),
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

func (e *UserFirstNameChangedEvent) EditorService() string {
	return "test.suite"
}

func (e *UserFirstNameChangedEvent) EditorUser() string {
	return "adlerhurst"
}

func (e *UserFirstNameChangedEvent) Type() eventstore.EventType {
	return "user.firstName.changed"
}

func (e *UserFirstNameChangedEvent) Data() interface{} {
	return e
}

func (e *UserFirstNameChangedEvent) MetaData() *eventstore.EventMetaData {
	return e.metaData
}

// ------------------------------------------------------------
// User password checked event start
// ------------------------------------------------------------

type UserPasswordCheckedEvent struct {
	metaData *eventstore.EventMetaData `json:"-"`
}

func UserPasswordCheckedMapper() (eventstore.EventType, func(*repository.Event) (eventstore.Event, error)) {
	return "user.password.checked", func(event *repository.Event) (eventstore.Event, error) {
		return &UserPasswordCheckedEvent{
			metaData: eventstore.MetaDataFromRepo(event),
		}, nil
	}
}

func (e *UserPasswordCheckedEvent) CheckPrevious() bool {
	return false
}

func (e *UserPasswordCheckedEvent) EditorService() string {
	return "test.suite"
}

func (e *UserPasswordCheckedEvent) EditorUser() string {
	return "adlerhurst"
}

func (e *UserPasswordCheckedEvent) Type() eventstore.EventType {
	return "user.password.checked"
}

func (e *UserPasswordCheckedEvent) Data() interface{} {
	return nil
}

func (e *UserPasswordCheckedEvent) MetaData() *eventstore.EventMetaData {
	return e.metaData
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

func (rm *UsersReadModel) AppendEvents(events ...eventstore.Event) (err error) {
	rm.ReadModel.AppendEvents(events...)
	for _, event := range events {
		switch e := event.(type) {
		case *UserAddedEvent:
			user := NewUserReadModel(e.MetaData().AggregateID)
			rm.Users = append(rm.Users, user)
			err = user.AppendEvents(e)
		case *UserFirstNameChangedEvent, *UserPasswordCheckedEvent:
			_, user := rm.userByID(e.MetaData().AggregateID)
			if user == nil {
				return errors.New("user not found")
			}
			err = user.AppendEvents(e)
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

func (rm *UserReadModel) AppendEvents(events ...eventstore.Event) error {
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
			rm.lastPasswordCheck = e.metaData.CreationDate
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
		RegisterFilterEventMapper(UserPasswordCheckedMapper())

	events, err := es.PushAggregates(context.Background(),
		NewUserAggregate("1").AppendEvents(&UserAddedEvent{FirstName: "hodor"}),
		NewUserAggregate("2").AppendEvents(&UserAddedEvent{FirstName: "hodor"}, &UserPasswordCheckedEvent{}, &UserPasswordCheckedEvent{}, &UserFirstNameChangedEvent{FirstName: "ueli"}),
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

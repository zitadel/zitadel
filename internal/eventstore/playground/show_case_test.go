package eventstore_test

import (
	"encoding/json"
	"fmt"
	"log"
)

//MemberReadModel is the minimum representation of a View model.
// it might be saved in a database or in memory
type ReadModel struct {
	ProcessedSequence uint64
	ID                string
	events            []Event
}

//Append adds all the events to the aggregate.
// The function doesn't compute the new state of the read model
func (a *ReadModel) Append(events ...Event) {
	a.events = append(a.events, events...)
}

type ProjectReadModel struct {
	ReadModel
	Apps []*AppReadModel
	Name string
}

func (p *ProjectReadModel) Append(events ...Event) {
	for _, event := range events {
		switch event.(type) {
		case *AddAppEvent:
			app := new(AppReadModel)
			app.Append(event)
			p.Apps = append(p.Apps, app)
		case *UpdateAppEvent:
			for _, app := range p.Apps {
				app.Append(event)
			}
		}
	}
	p.events = append(p.events, events...)
}

type AppReadModel struct {
	ReadModel
	Name string
}

//Reduce calculates the new state of the read model
func (a *AppReadModel) Reduce() error {
	for _, event := range a.events {
		switch e := event.(type) {
		case *AddAppEvent:
			a.Name = e.Name
			a.ID = e.GetID()
		case *UpdateAppEvent:
			a.Name = e.Name
		}
		a.ProcessedSequence = event.GetSequence()
	}
	return nil
}

//Reduce calculates the new state of the read model
func (p *ProjectReadModel) Reduce() error {
	for i := range p.Apps {
		if err := p.Apps[i].Reduce(); err != nil {
			return err
		}
	}
	for _, event := range p.events {
		switch e := event.(type) {
		case *CreateProjectEvent:
			p.ID = e.ID
			p.Name = e.Name
		case *RemoveAppEvent:
			for i := len(p.Apps) - 1; i >= 0; i-- {
				app := p.Apps[i]
				if app.ID == e.GetID() {
					p.Apps[i] = p.Apps[len(p.Apps)-1]
					p.Apps[len(p.Apps)-1] = nil
					p.Apps = p.Apps[:len(p.Apps)-1]
				}
			}
		}
		p.ProcessedSequence = event.GetSequence()
	}
	return nil
}

//Event is the minimal representation of a event
// which can be processed by the read models
type Event interface {
	//GetSequence returns the event sequence
	GetSequence() uint64
	//GetID returns the id of the aggregate. It's not the id of the event
	GetID() string
}

//DefaultEvent is the implementation of Event
type DefaultEvent struct {
	Sequence uint64 `json:"-"`
	ID       string `json:"-"`
}

func (e *DefaultEvent) GetID() string {
	return e.ID
}

func (e *DefaultEvent) GetSequence() uint64 {
	return e.Sequence
}

type CreateProjectEvent struct {
	DefaultEvent
	Name string `json:"name,omitempty"`
}

//CreateProjectEventFromEventstore returns the specific type
// of the general EventstoreEvent
func CreateProjectEventFromEventstore(event *EventstoreEvent) (Event, error) {
	e := &CreateProjectEvent{
		DefaultEvent: DefaultEvent{Sequence: event.Sequence, ID: event.AggregateID},
	}
	err := json.Unmarshal(event.Data, e)

	return e, err
}

type AddAppEvent struct {
	ProjectID string `json:"-"`
	AppID     string `json:"id"`
	Sequence  uint64 `json:"-"`
	Name      string `json:"name,omitempty"`
}

func (e *AddAppEvent) GetID() string {
	return e.AppID
}

func (e *AddAppEvent) GetSequence() uint64 {
	return e.Sequence
}

func AppAddedEventFromEventstore(event *EventstoreEvent) (Event, error) {
	e := &AddAppEvent{
		Sequence:  event.Sequence,
		ProjectID: event.AggregateID,
	}
	err := json.Unmarshal(event.Data, e)

	return e, err
}

type UpdateAppEvent struct {
	ProjectID string `json:"-"`
	AppID     string `json:"id"`
	Sequence  uint64 `json:"-"`
	Name      string `json:"name,omitempty"`
}

func (e *UpdateAppEvent) GetID() string {
	return e.AppID
}

func (e *UpdateAppEvent) GetSequence() uint64 {
	return e.Sequence
}

func AppUpdatedEventFromEventstore(event *EventstoreEvent) (Event, error) {
	e := &UpdateAppEvent{
		Sequence:  event.Sequence,
		ProjectID: event.AggregateID,
	}
	err := json.Unmarshal(event.Data, e)

	return e, err
}

type RemoveAppEvent struct {
	ProjectID string `json:"-"`
	AppID     string `json:"id"`
	Sequence  uint64 `json:"-"`
}

func (e *RemoveAppEvent) GetID() string {
	return e.AppID
}

func (e *RemoveAppEvent) GetSequence() uint64 {
	return e.Sequence
}

func AppRemovedEventFromEventstore(event *EventstoreEvent) (Event, error) {
	e := &RemoveAppEvent{
		Sequence:  event.Sequence,
		ProjectID: event.AggregateID,
	}
	err := json.Unmarshal(event.Data, e)

	return e, err
}

func main() {
	eventstore := &Eventstore{
		eventMapper: map[string]func(*EventstoreEvent) (Event, error){
			"project.added": CreateProjectEventFromEventstore,
			"app.added":     AppAddedEventFromEventstore,
			"app.updated":   AppUpdatedEventFromEventstore,
			"app.removed":   AppRemovedEventFromEventstore,
		},
		events: []*EventstoreEvent{
			{
				AggregateID: "p1",
				EventType:   "project.added",
				Sequence:    1,
				Data:        []byte(`{"name":"hodor"}`),
			},
			{
				AggregateID: "123",
				EventType:   "app.added",
				Sequence:    2,
				Data:        []byte(`{"id":"a1", "name": "ap 1"}`),
			},
			{
				AggregateID: "123",
				EventType:   "app.updated",
				Sequence:    3,
				Data:        []byte(`{"id":"a1", "name":"app 1"}`),
			},
			{
				AggregateID: "123",
				EventType:   "app.added",
				Sequence:    4,
				Data:        []byte(`{"id":"a2", "name": "app 2"}`),
			},
			{
				AggregateID: "123",
				EventType:   "app.removed",
				Sequence:    5,
				Data:        []byte(`{"id":"a1"}`),
			},
		},
	}
	events, err := eventstore.GetEvents()
	if err != nil {
		log.Panic(err)
	}

	p := &ProjectReadModel{Apps: []*AppReadModel{}}
	p.Append(events...)
	p.Reduce()

	fmt.Printf("%+v\n", p)
	for _, app := range p.Apps {
		fmt.Printf("%+v\n", app)
	}
}

//Eventstore is a simple abstraction of the eventstore framework
type Eventstore struct {
	eventMapper map[string]func(*EventstoreEvent) (Event, error)
	events      []*EventstoreEvent
}

func (es *Eventstore) GetEvents() (events []Event, err error) {
	events = make([]Event, len(es.events))
	for i, event := range es.events {
		events[i], err = es.eventMapper[event.EventType](event)
		if err != nil {
			return nil, err
		}
	}
	return events, nil
}

type EventstoreEvent struct {
	AggregateID string
	Sequence    uint64
	EventType   string
	Data        []byte
}

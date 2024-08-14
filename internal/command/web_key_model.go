package command

import (
	"github.com/go-jose/go-jose/v4"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/webkey"
)

type WebKeyWriteModel struct {
	eventstore.WriteModel
	State      domain.WebKeyState
	PrivateKey *crypto.CryptoValue
	PublicKey  *jose.JSONWebKey
}

func NewWebKeyWriteModel(keyID, resourceOwner string) *WebKeyWriteModel {
	return &WebKeyWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   keyID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *WebKeyWriteModel) AppendEvents(events ...eventstore.Event) {
	wm.WriteModel.AppendEvents(events...)
}

func (wm *WebKeyWriteModel) Reduce() error {
	for _, event := range wm.Events {
		if event.Aggregate().ID != wm.AggregateID {
			continue
		}
		switch e := event.(type) {
		case *webkey.AddedEvent:
			wm.State = domain.WebKeyStateInitial
			wm.PrivateKey = e.PrivateKey
			wm.PublicKey = e.PublicKey
		case *webkey.ActivatedEvent:
			wm.State = domain.WebKeyStateActive
		case *webkey.DeactivatedEvent:
			wm.State = domain.WebKeyStateInactive
		case *webkey.RemovedEvent:
			wm.State = domain.WebKeyStateRemoved
			wm.PrivateKey = nil
			wm.PublicKey = nil
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *WebKeyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(webkey.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			webkey.AddedEventType,
			webkey.ActivatedEventType,
			webkey.DeactivatedEventType,
			webkey.RemovedEventType,
		).
		Builder()
}

type webKeyWriteModels struct {
	resourceOwner string
	events        []eventstore.Event
	keys          map[string]*WebKeyWriteModel
	activeID      string
}

func newWebKeyWriteModels(resourceOwner string) *webKeyWriteModels {
	return &webKeyWriteModels{
		resourceOwner: resourceOwner,
		keys:          make(map[string]*WebKeyWriteModel),
	}
}

func (models *webKeyWriteModels) AppendEvents(events ...eventstore.Event) {
	models.events = append(models.events, events...)
}

func (models *webKeyWriteModels) Reduce() error {
	for _, event := range models.events {
		aggregate := event.Aggregate()
		if models.keys[aggregate.ID] == nil {
			models.keys[aggregate.ID] = NewWebKeyWriteModel(aggregate.ID, aggregate.ResourceOwner)
		}

		switch event.(type) {
		case *webkey.AddedEvent:
			break
		case *webkey.ActivatedEvent:
			models.activeID = aggregate.ID
		case *webkey.DeactivatedEvent:
			if models.activeID == aggregate.ID {
				models.activeID = ""
			}
		case *webkey.RemovedEvent:
			delete(models.keys, aggregate.ID)
			continue
		}

		model := models.keys[aggregate.ID]
		model.AppendEvents(event)
		if err := model.Reduce(); err != nil {
			return err
		}
	}
	models.events = models.events[0:0]
	return nil
}

func (models *webKeyWriteModels) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(models.resourceOwner).
		AddQuery().
		AggregateTypes(webkey.AggregateType).
		EventTypes(
			webkey.AddedEventType,
			webkey.ActivatedEventType,
			webkey.DeactivatedEventType,
			webkey.RemovedEventType,
		).
		Builder()
}

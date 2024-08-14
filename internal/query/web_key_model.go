package query

import (
	"github.com/go-jose/go-jose/v4"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/webkey"
)

type WebKeyReadModel struct {
	eventstore.ReadModel
	State      domain.WebKeyState
	PrivateKey *crypto.CryptoValue
	PublicKey  *jose.JSONWebKey
	Config     crypto.WebKeyConfig
}

func NewWebKeyReadModel(keyID, resourceOwner string) *WebKeyReadModel {
	return &WebKeyReadModel{
		ReadModel: eventstore.ReadModel{
			AggregateID:   keyID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *WebKeyReadModel) AppendEvents(events ...eventstore.Event) {
	wm.ReadModel.AppendEvents(events...)
}

func (wm *WebKeyReadModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *webkey.AddedEvent:
			if err := wm.reduceAdded(e); err != nil {
				return err
			}
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
	return wm.ReadModel.Reduce()
}

func (wm *WebKeyReadModel) reduceAdded(e *webkey.AddedEvent) (err error) {
	wm.State = domain.WebKeyStateInitial
	wm.PrivateKey = e.PrivateKey
	wm.PublicKey = e.PublicKey
	wm.Config, err = crypto.UnmarshalWebKeyConfig(e.Config, e.ConfigType)
	return err
}

func (wm *WebKeyReadModel) Query() *eventstore.SearchQueryBuilder {
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

package command

import (
	"net/url"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/idpintent"
)

type IDPIntentWriteModel struct {
	eventstore.WriteModel

	SuccessURL *url.URL
	FailureURL *url.URL
	IDPID      string
	Token      *crypto.CryptoValue

	UserID string
	State  domain.IDPIntentState
}

func NewIDPIntentWriteModel(id, resourceOwner string) *IDPIntentWriteModel {
	return &IDPIntentWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   id,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *IDPIntentWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idpintent.StartedEvent:
			wm.reduceStartedEvent(e)
		case *idpintent.SucceededEvent:
			wm.reduceSucceededEvent(e)
		case *idpintent.FailedEvent:
			wm.reduceFailedEvent(e)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *IDPIntentWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(idpintent.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			idpintent.StartedEventType,
			idpintent.SucceededEventType,
			idpintent.FailedEventType,
		).
		Builder()
}

func (wm *IDPIntentWriteModel) reduceStartedEvent(e *idpintent.StartedEvent) {
	wm.SuccessURL = e.SuccessURL
	wm.FailureURL = e.FailureURL
	wm.IDPID = e.IDPID
}

func (wm *IDPIntentWriteModel) reduceSucceededEvent(e *idpintent.SucceededEvent) {
	wm.Token = e.Token
	wm.UserID = e.UserID
}

func (wm *IDPIntentWriteModel) reduceFailedEvent(e *idpintent.FailedEvent) {

}

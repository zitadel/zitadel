package command

import (
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/samlsession"
)

type SAMLSessionWriteModel struct {
	eventstore.WriteModel

	UserID                 string
	UserResourceOwner      string
	PreferredLanguage      *language.Tag
	SessionID              string
	EntityID               string
	Audience               []string
	AuthMethods            []domain.UserAuthMethodType
	AuthTime               time.Time
	UserAgent              *domain.UserAgent
	State                  domain.SAMLSessionState
	SAMLResponseID         string
	SAMLResponseCreation   time.Time
	SAMLResponseExpiration time.Time

	aggregate *eventstore.Aggregate
}

func NewSAMLSessionWriteModel(id string, resourceOwner string) *SAMLSessionWriteModel {
	return &SAMLSessionWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   id,
			ResourceOwner: resourceOwner,
		},
		aggregate: &samlsession.NewAggregate(id, resourceOwner).Aggregate,
	}
}

func (wm *SAMLSessionWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *samlsession.AddedEvent:
			wm.reduceAdded(e)
		case *samlsession.SAMLResponseAddedEvent:
			wm.reduceSAMLResponseAdded(e)
		case *samlsession.SAMLResponseRevokedEvent:
			wm.reduceSAMLResponseRevoked(e)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *SAMLSessionWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(samlsession.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			samlsession.AddedType,
			samlsession.SAMLResponseAddedType,
			samlsession.SAMLResponseRevokedType,
		).
		Builder()

	if wm.ResourceOwner != "" {
		query.ResourceOwner(wm.ResourceOwner)
	}
	return query
}

func (wm *SAMLSessionWriteModel) reduceAdded(e *samlsession.AddedEvent) {
	wm.UserID = e.UserID
	wm.UserResourceOwner = e.UserResourceOwner
	wm.SessionID = e.SessionID
	wm.EntityID = e.EntityID
	wm.Audience = e.Audience
	wm.AuthMethods = e.AuthMethods
	wm.AuthTime = e.AuthTime
	wm.PreferredLanguage = e.PreferredLanguage
	wm.UserAgent = e.UserAgent
	wm.State = domain.SAMLSessionStateActive
	// the write model might be initialized without resource owner,
	// so update the aggregate
	if wm.ResourceOwner == "" {
		wm.aggregate = &samlsession.NewAggregate(wm.AggregateID, e.Aggregate().ResourceOwner).Aggregate
	}
}

func (wm *SAMLSessionWriteModel) reduceSAMLResponseAdded(e *samlsession.SAMLResponseAddedEvent) {
	wm.SAMLResponseID = e.ID
	wm.SAMLResponseCreation = e.CreationDate()
	wm.SAMLResponseExpiration = e.CreationDate().Add(e.Lifetime)
}

func (wm *SAMLSessionWriteModel) reduceSAMLResponseRevoked(e *samlsession.SAMLResponseRevokedEvent) {
	wm.SAMLResponseID = ""
	wm.SAMLResponseExpiration = e.CreationDate()
}

package command

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/sessionlogout"
)

type SessionLogoutWriteModel struct {
	eventstore.WriteModel

	UserID                string
	OIDCSessionID         string
	ClientID              string
	BackChannelLogoutURI  string
	BackChannelLogoutSent bool

	aggregate *eventstore.Aggregate
}

func NewSessionLogoutWriteModel(id string, instanceID string, sessionID string) *SessionLogoutWriteModel {
	return &SessionLogoutWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   id,
			ResourceOwner: instanceID,
			InstanceID:    instanceID,
		},
		aggregate:     &sessionlogout.NewAggregate(id, instanceID).Aggregate,
		OIDCSessionID: sessionID,
	}
}

func (wm *SessionLogoutWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *sessionlogout.BackChannelLogoutRegisteredEvent:
			wm.reduceRegistered(e)
		case *sessionlogout.BackChannelLogoutSentEvent:
			wm.reduceSent(e)
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *SessionLogoutWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(sessionlogout.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			sessionlogout.BackChannelLogoutRegisteredType,
			sessionlogout.BackChannelLogoutSentType,
		).
		EventData(map[string]interface{}{
			"oidc_session_id": wm.OIDCSessionID,
		}).
		Builder()
	return query
}

func (wm *SessionLogoutWriteModel) reduceRegistered(e *sessionlogout.BackChannelLogoutRegisteredEvent) {
	if wm.OIDCSessionID != e.OIDCSessionID {
		return
	}
	wm.UserID = e.UserID
	wm.ClientID = e.ClientID
	wm.BackChannelLogoutURI = e.BackChannelLogoutURI
}

func (wm *SessionLogoutWriteModel) reduceSent(e *sessionlogout.BackChannelLogoutSentEvent) {
	if wm.OIDCSessionID != e.OIDCSessionID {
		return
	}
	wm.BackChannelLogoutSent = true
}

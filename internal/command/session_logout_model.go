package command

import (
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/sessionlogout"
)

type FederatedLogoutWriteModel struct {
	eventstore.WriteModel

	SessionID             string
	IDPID                 string
	UserID                string
	PostLogoutRedirectURI string

	SAMLRequestID   string
	SAMLBindingType string
	SAMLRedirectURL string
	SAMLPostURL     string
	SAMLRequest     string
	SAMLRelayState  string

	State SessionLogoutState
}

type SessionLogoutState int

const (
	SessionLogoutStateUnspecified SessionLogoutState = iota
	SessionLogoutStateStarted
	SessionLogoutStateSAMLRequestCreated
	SessionLogoutStateSAMLResponseReceived
	SessionLogoutStateCompleted
	SessionLogoutStateFailed
)

func NewFederatedLogoutWriteModel(logoutID, instanceID string) *FederatedLogoutWriteModel {
	return &FederatedLogoutWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   logoutID,
			InstanceID:    instanceID,
			ResourceOwner: instanceID,
		},
	}
}

func (wm *FederatedLogoutWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *sessionlogout.StartedEvent:
			wm.SessionID = e.SessionID
			wm.IDPID = e.IDPID
			wm.UserID = e.UserID
			wm.PostLogoutRedirectURI = e.PostLogoutRedirectURI
			wm.State = SessionLogoutStateStarted
		case *sessionlogout.SAMLRequestCreatedEvent:
			wm.SAMLRequestID = e.RequestID
			wm.SAMLBindingType = e.BindingType
			wm.SAMLRedirectURL = e.RedirectURL
			wm.SAMLPostURL = e.PostURL
			wm.SAMLRequest = e.SAMLRequest
			wm.SAMLRelayState = e.RelayState
			wm.State = SessionLogoutStateSAMLRequestCreated
		case *sessionlogout.SAMLResponseReceivedEvent:
			wm.State = SessionLogoutStateSAMLResponseReceived
		case *sessionlogout.CompletedEvent:
			wm.State = SessionLogoutStateCompleted
		case *sessionlogout.FailedEvent:
			wm.State = SessionLogoutStateFailed
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *FederatedLogoutWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(sessionlogout.AggregateType).
		AggregateIDs(wm.AggregateID).
		Builder()
}

func (wm *FederatedLogoutWriteModel) IsActive() bool {
	return wm.State == SessionLogoutStateStarted || wm.State == SessionLogoutStateSAMLRequestCreated
}

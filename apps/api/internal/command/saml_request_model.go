package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/samlrequest"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type SAMLRequestWriteModel struct {
	eventstore.WriteModel
	aggregate *eventstore.Aggregate

	LoginClient    string
	ApplicationID  string
	ACSURL         string
	RelayState     string
	RequestID      string
	Binding        string
	Issuer         string
	Destination    string
	ResponseIssuer string

	SessionID        string
	UserID           string
	AuthTime         time.Time
	AuthMethods      []domain.UserAuthMethodType
	SAMLRequestState domain.SAMLRequestState
}

func NewSAMLRequestWriteModel(ctx context.Context, id string) *SAMLRequestWriteModel {
	return &SAMLRequestWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID: id,
		},
		aggregate: &samlrequest.NewAggregate(id, authz.GetInstance(ctx).InstanceID()).Aggregate,
	}
}

func (m *SAMLRequestWriteModel) Reduce() error {
	for _, event := range m.Events {
		switch e := event.(type) {
		case *samlrequest.AddedEvent:
			m.LoginClient = e.LoginClient
			m.ApplicationID = e.ApplicationID
			m.ACSURL = e.ACSURL
			m.RelayState = e.RelayState
			m.RequestID = e.RequestID
			m.Binding = e.Binding
			m.Issuer = e.Issuer
			m.Destination = e.Destination
			m.ResponseIssuer = e.ResponseIssuer
			m.SAMLRequestState = domain.SAMLRequestStateAdded
		case *samlrequest.SessionLinkedEvent:
			m.SessionID = e.SessionID
			m.UserID = e.UserID
			m.AuthTime = e.AuthTime
			m.AuthMethods = e.AuthMethods
		case *samlrequest.FailedEvent:
			m.SAMLRequestState = domain.SAMLRequestStateFailed
		case *samlrequest.SucceededEvent:
			m.SAMLRequestState = domain.SAMLRequestStateSucceeded
		}
	}
	return m.WriteModel.Reduce()
}

func (m *SAMLRequestWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(samlrequest.AggregateType).
		AggregateIDs(m.AggregateID).
		Builder()
}

// CheckAuthenticated checks that the auth request exists, a session must have been linked
func (m *SAMLRequestWriteModel) CheckAuthenticated() error {
	if m.SessionID == "" {
		return zerrors.ThrowPreconditionFailed(nil, "AUTHR-3dNRNwSYeC", "Errors.SAMLRequest.NotAuthenticated")
	}
	// check that the requests exists, but has not succeeded yet
	if m.SAMLRequestState == domain.SAMLRequestStateAdded {
		return nil
	}
	return zerrors.ThrowPreconditionFailed(nil, "AUTHR-krQV50AlnJ", "Errors.SAMLRequest.NotAuthenticated")
}

package query

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/deviceauth"
)

type DeviceAuthReadModel struct {
	eventstore.ReadModel
	DeviceAuth
}

func NewDeviceAuthReadModel(deviceCode, resourceOwner string) *DeviceAuthReadModel {
	return &DeviceAuthReadModel{
		ReadModel: eventstore.ReadModel{
			AggregateID:   deviceCode,
			ResourceOwner: resourceOwner,
		},
	}
}

func (m *DeviceAuthReadModel) Reduce() error {
	for _, event := range m.Events {
		switch e := event.(type) {
		case *deviceauth.AddedEvent:
			m.ClientID = e.ClientID
			m.DeviceCode = e.DeviceCode
			m.UserCode = e.UserCode
			m.Expires = e.Expires
			m.Scopes = e.Scopes
			m.Audience = e.Audience
			m.State = e.State
		case *deviceauth.ApprovedEvent:
			m.State = domain.DeviceAuthStateApproved
			m.Subject = e.Subject
			m.UserAuthMethods = e.UserAuthMethods
			m.AuthTime = e.AuthTime
		case *deviceauth.CanceledEvent:
			m.State = e.Reason.State()
		}
	}

	return m.ReadModel.Reduce()
}

func (m *DeviceAuthReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(m.ResourceOwner).
		AddQuery().
		AggregateTypes(deviceauth.AggregateType).
		AggregateIDs(m.AggregateID).
		EventTypes(
			deviceauth.AddedEventType,
			deviceauth.ApprovedEventType,
			deviceauth.CanceledEventType,
		).
		Builder()
}

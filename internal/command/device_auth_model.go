package command

import (
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/deviceauth"
)

type DeviceAuthWriteModel struct {
	eventstore.WriteModel

	ClientID        string
	DeviceCode      string
	UserCode        string
	Expires         time.Time
	Scopes          []string
	State           domain.DeviceAuthState
	Subject         string
	UserAuthMethods []domain.UserAuthMethodType
	AuthTime        time.Time
}

func NewDeviceAuthWriteModel(deviceCode, resourceOwner string) *DeviceAuthWriteModel {
	return &DeviceAuthWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   deviceCode,
			ResourceOwner: resourceOwner,
		},
	}
}

func (m *DeviceAuthWriteModel) Reduce() error {
	for _, event := range m.Events {
		switch e := event.(type) {
		case *deviceauth.AddedEvent:
			m.ClientID = e.ClientID
			m.DeviceCode = e.DeviceCode
			m.UserCode = e.UserCode
			m.Expires = e.Expires
			m.Scopes = e.Scopes
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

	return m.WriteModel.Reduce()
}

func (m *DeviceAuthWriteModel) Query() *eventstore.SearchQueryBuilder {
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

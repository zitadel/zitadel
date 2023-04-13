package command

import (
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/deviceauth"
)

type DeviceAuthWriteModel struct {
	eventstore.WriteModel

	ClientID   string
	DeviceCode string
	UserCode   string
	Expires    time.Time
	Scopes     []string
	Subject    string
	State      domain.DeviceAuthState
}

func NewDeviceAuthWriteModel(aggrID, resourceOwner string) *DeviceAuthWriteModel {
	return &DeviceAuthWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   aggrID,
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
			m.Subject = e.Subject
			m.State = domain.DeviceAuthStateApproved
		case *deviceauth.CanceledEvent:
			m.State = e.Reason.State()
		case *deviceauth.RemovedEvent:
			m.State = domain.DeviceAuthStateRemoved
		}
	}

	return m.WriteModel.Reduce()
}

func (m *DeviceAuthWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(deviceauth.AggregateType).
		AggregateIDs(m.AggregateID).
		Builder()
}

package command

import (
	"time"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/deviceauth"
)

type DeviceAuthWriteModel struct {
	eventstore.WriteModel
	aggregate *eventstore.Aggregate

	ClientID          string
	DeviceCode        string
	UserCode          string
	Expires           time.Time
	Scopes            []string
	Audience          []string
	State             domain.DeviceAuthState
	UserID            string
	UserOrgID         string
	UserAuthMethods   []domain.UserAuthMethodType
	AuthTime          time.Time
	PreferredLanguage *language.Tag
	UserAgent         *domain.UserAgent
	NeedRefreshToken  bool
}

func NewDeviceAuthWriteModel(deviceCode, resourceOwner string) *DeviceAuthWriteModel {
	return &DeviceAuthWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   deviceCode,
			ResourceOwner: resourceOwner,
		},
		aggregate: deviceauth.NewAggregate(deviceCode, resourceOwner),
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
			m.Audience = e.Audience
			m.State = e.State
			m.NeedRefreshToken = e.NeedRefreshToken
		case *deviceauth.ApprovedEvent:
			m.State = domain.DeviceAuthStateApproved
			m.UserID = e.UserID
			m.UserOrgID = e.UserOrgID
			m.UserAuthMethods = e.UserAuthMethods
			m.AuthTime = e.AuthTime
			m.PreferredLanguage = e.PreferredLanguage
			m.UserAgent = e.UserAgent
		case *deviceauth.CanceledEvent:
			m.State = e.Reason.State()
		case *deviceauth.DoneEvent:
			m.State = domain.DeviceAuthStateDone
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

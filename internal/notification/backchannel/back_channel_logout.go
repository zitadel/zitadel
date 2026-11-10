package backchannel

import "github.com/zitadel/zitadel/internal/eventstore"

const (
	QueueName = "back_channel_logout"
)

type LogoutRequest struct {
	Aggregate            *eventstore.Aggregate
	SessionID            string
	TriggeredAtOrigin    string
	TriggeringEventType  eventstore.EventType
	TokenID              string
	UserID               string
	OIDCSessionID        string
	ClientID             string
	BackChannelLogoutURI string
}

func (l *LogoutRequest) Kind() string {
	return "back_channel_logout_request"
}

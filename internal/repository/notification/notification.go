package notification

import (
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
)

const (
	QueueName = "notification"
)

type Request struct {
	Aggregate                     *eventstore.Aggregate         `json:"aggregate"`
	UserID                        string                        `json:"userID"`
	UserResourceOwner             string                        `json:"userResourceOwner"`
	TriggeredAtOrigin             string                        `json:"triggeredAtOrigin"`
	EventType                     eventstore.EventType          `json:"eventType"`
	MessageType                   string                        `json:"messageType"`
	NotificationType              domain.NotificationType       `json:"notificationType"`
	URLTemplate                   string                        `json:"urlTemplate,omitempty"`
	CodeExpiry                    time.Duration                 `json:"codeExpiry,omitempty"`
	Code                          *crypto.CryptoValue           `json:"code,omitempty"`
	UnverifiedNotificationChannel bool                          `json:"unverifiedNotificationChannel,omitempty"`
	IsOTP                         bool                          `json:"isOTP,omitempty"`
	RequiresPreviousDomain        bool                          `json:"requiresPreviousDomain,omitempty"`
	Args                          *domain.NotificationArguments `json:"args,omitempty"`
}

func (e *Request) Kind() string {
	return "notification_request"
}

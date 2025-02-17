package command

import (
	"time"

	"github.com/riverqueue/river"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
)

type NotificationRequest struct {
	Aggregate                     *eventstore.Aggregate
	UserID                        string
	UserResourceOwner             string
	TriggerOrigin                 string
	URLTemplate                   string
	Code                          *crypto.CryptoValue
	CodeExpiry                    time.Duration
	EventType                     eventstore.EventType
	NotificationType              domain.NotificationType
	MessageType                   string
	UnverifiedNotificationChannel bool
	Args                          *domain.NotificationArguments
	IsOTP                         bool
	RequiresPreviousDomain        bool
}

type NotificationRetryRequest struct {
	NotificationRequest
	BackOff    time.Duration
	NotifyUser *query.NotifyUser
}

var _ river.JobArgs = (*NotificationRequest)(nil)

func NewNotificationRequest(
	aggregate *eventstore.Aggregate,
	userID, resourceOwner, triggerOrigin string,
	eventType eventstore.EventType,
	notificationType domain.NotificationType,
	messageType string,
) *NotificationRequest {
	return &NotificationRequest{
		Aggregate:         aggregate,
		UserID:            userID,
		UserResourceOwner: resourceOwner,
		TriggerOrigin:     triggerOrigin,
		EventType:         eventType,
		NotificationType:  notificationType,
		MessageType:       messageType,
	}
}

// Kind implements [river.JobArgs].
func (r *NotificationRequest) Kind() string {
	return "notification_requested"
}

func (r *NotificationRequest) WithCode(code *crypto.CryptoValue, expiry time.Duration) *NotificationRequest {
	r.Code = code
	r.CodeExpiry = expiry
	return r
}

func (r *NotificationRequest) WithURLTemplate(urlTemplate string) *NotificationRequest {
	r.URLTemplate = urlTemplate
	return r
}

func (r *NotificationRequest) WithUnverifiedChannel() *NotificationRequest {
	r.UnverifiedNotificationChannel = true
	return r
}

func (r *NotificationRequest) WithArgs(args *domain.NotificationArguments) *NotificationRequest {
	r.Args = args
	return r
}

func (r *NotificationRequest) WithAggregate(id, resourceOwner string) *NotificationRequest {
	r.Aggregate.ID = id
	r.Aggregate.ResourceOwner = resourceOwner
	return r
}

func (r *NotificationRequest) WithOTP() *NotificationRequest {
	r.IsOTP = true
	return r
}

func (r *NotificationRequest) WithPreviousDomain() *NotificationRequest {
	r.RequiresPreviousDomain = true
	return r
}

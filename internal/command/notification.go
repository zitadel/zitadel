package command

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/notification"
)

type NotificationRequest struct {
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
	Args                          map[string]any
	AggregateID                   string
	AggregateResourceOwner        string
	IsOTP                         bool
}

func NewNotificationRequest(
	userID, resourceOwner, triggerOrigin string,
	eventType eventstore.EventType,
	notificationType domain.NotificationType,
	messageType string,
) *NotificationRequest {
	return &NotificationRequest{
		UserID:            userID,
		UserResourceOwner: resourceOwner,
		TriggerOrigin:     triggerOrigin,
		EventType:         eventType,
		NotificationType:  notificationType,
		MessageType:       messageType,
	}
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

func (r *NotificationRequest) WithArgs(args map[string]interface{}) *NotificationRequest {
	r.Args = args
	return r
}

func (r *NotificationRequest) WithAggregate(id, resourceOwner string) *NotificationRequest {
	r.AggregateID = id
	r.AggregateResourceOwner = resourceOwner
	return r
}

func (r *NotificationRequest) WithOTP() *NotificationRequest {
	r.IsOTP = true
	return r
}

// RequestNotification writes a new notification.RequestEvent with the notification.Aggregate to the eventstore
func (c *Commands) RequestNotification(
	ctx context.Context,
	instanceID string,
	request *NotificationRequest,
) error {
	id, err := c.idGenerator.Next()
	if err != nil {
		return err
	}
	_, err = c.eventstore.Push(ctx, notification.NewRequestedEvent(ctx, &notification.NewAggregate(id, instanceID).Aggregate,
		request.UserID,
		request.UserResourceOwner,
		request.AggregateID,
		request.AggregateResourceOwner,
		request.TriggerOrigin,
		request.URLTemplate,
		request.Code,
		request.CodeExpiry,
		request.EventType,
		request.NotificationType,
		request.MessageType,
		request.UnverifiedNotificationChannel,
		request.IsOTP,
		request.Args))
	return err
}

// NotificationFailed writes a new notification.RequestEvent with the notification.Aggregate to the eventstore
func (c *Commands) NotificationFailed(ctx context.Context, id, instanceID string, err error) error {
	_, err = c.eventstore.Push(ctx, notification.NewFailedEvent(ctx, &notification.NewAggregate(id, instanceID).Aggregate,
		err))
	return err
}

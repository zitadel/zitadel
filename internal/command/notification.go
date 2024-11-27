package command

import (
	"context"
	"database/sql"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query"
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
	Args                          *domain.NotificationArguments
	AggregateID                   string
	AggregateResourceOwner        string
	IsOTP                         bool
	RequiresPreviousDomain        bool
}

type NotificationRetryRequest struct {
	NotificationRequest
	BackOff    time.Duration
	NotifyUser *query.NotifyUser
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

func (r *NotificationRequest) WithArgs(args *domain.NotificationArguments) *NotificationRequest {
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

func (r *NotificationRequest) WithPreviousDomain() *NotificationRequest {
	r.RequiresPreviousDomain = true
	return r
}

// RequestNotification writes a new notification.RequestEvent with the notification.Aggregate to the eventstore
func (c *Commands) RequestNotification(
	ctx context.Context,
	resourceOwner string,
	request *NotificationRequest,
) error {
	id, err := c.idGenerator.Next()
	if err != nil {
		return err
	}
	_, err = c.eventstore.Push(ctx, notification.NewRequestedEvent(ctx, &notification.NewAggregate(id, resourceOwner).Aggregate,
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
		request.RequiresPreviousDomain,
		request.Args))
	return err
}

// NotificationCanceled writes a new notification.CanceledEvent with the notification.Aggregate to the eventstore
func (c *Commands) NotificationCanceled(ctx context.Context, tx *sql.Tx, id, resourceOwner string, requestError error) error {
	var errorMessage string
	if requestError != nil {
		errorMessage = requestError.Error()
	}
	_, err := c.eventstore.PushWithClient(ctx, tx, notification.NewCanceledEvent(ctx, &notification.NewAggregate(id, resourceOwner).Aggregate, errorMessage))
	return err
}

// NotificationSent writes a new notification.SentEvent with the notification.Aggregate to the eventstore
func (c *Commands) NotificationSent(ctx context.Context, tx *sql.Tx, id, resourceOwner string) error {
	_, err := c.eventstore.PushWithClient(ctx, tx, notification.NewSentEvent(ctx, &notification.NewAggregate(id, resourceOwner).Aggregate))
	return err
}

// NotificationRetryRequested writes a new notification.RetryRequestEvent with the notification.Aggregate to the eventstore
func (c *Commands) NotificationRetryRequested(ctx context.Context, tx *sql.Tx, id, resourceOwner string, request *NotificationRetryRequest, requestError error) error {
	var errorMessage string
	if requestError != nil {
		errorMessage = requestError.Error()
	}
	_, err := c.eventstore.PushWithClient(ctx, tx, notification.NewRetryRequestedEvent(ctx, &notification.NewAggregate(id, resourceOwner).Aggregate,
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
		request.Args,
		request.NotifyUser,
		request.BackOff,
		errorMessage))
	return err
}

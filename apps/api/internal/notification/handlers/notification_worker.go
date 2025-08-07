package handlers

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/riverqueue/river"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/notification"
)

const (
	Code = "Code"
	OTP  = "OTP"
)

type NotificationWorker struct {
	river.WorkerDefaults[*notification.Request]

	commands Commands
	queries  *NotificationQueries
	channels types.ChannelChains
	config   WorkerConfig
	now      nowFunc
}

// Timeout implements the Timeout-function of [river.Worker].
// Maximum time a job can run before the context gets cancelled.
func (w *NotificationWorker) Timeout(*river.Job[*notification.Request]) time.Duration {
	return w.config.TransactionDuration
}

// Work implements [river.Worker].
func (w *NotificationWorker) Work(ctx context.Context, job *river.Job[*notification.Request]) error {
	ctx = ContextWithNotifier(ctx, job.Args.Aggregate)

	// if the notification is too old, we can directly cancel
	if job.CreatedAt.Add(w.config.MaxTtl).Before(w.now()) {
		return river.JobCancel(errors.New("notification is too old"))
	}

	// We do not trigger the projection to reduce load on the database. By the time the notification is processed,
	// the user should be projected anyway. If not, it will just wait for the next run.
	// We are aware that the user can change during the time the notification is in the queue.
	notifyUser, err := w.queries.GetNotifyUserByID(ctx, false, job.Args.UserID)
	if err != nil {
		return err
	}

	// The domain claimed event requires the domain as argument, but lacks the user when creating the request event.
	// Since we set it into the request arguments, it will be passed into a potential retry event.
	if job.Args.RequiresPreviousDomain && job.Args.Args != nil && job.Args.Args.Domain == "" {
		index := strings.LastIndex(notifyUser.LastEmail, "@")
		job.Args.Args.Domain = notifyUser.LastEmail[index+1:]
	}

	err = w.sendNotificationQueue(ctx, job.Args, strconv.Itoa(int(job.ID)), notifyUser)
	if err == nil {
		return nil
	}
	// if the error explicitly specifies, we cancel the notification
	if errors.Is(err, &channels.CancelError{}) {
		return river.JobCancel(err)
	}
	return err
}

type WorkerConfig struct {
	LegacyEnabled       bool
	Workers             uint8
	TransactionDuration time.Duration
	MaxTtl              time.Duration
	MaxAttempts         uint8
}

// nowFunc makes [time.Now] mockable
type nowFunc func() time.Time

type Sent func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error

var sentHandlers map[eventstore.EventType]Sent

func RegisterSentHandler(eventType eventstore.EventType, sent Sent) {
	if sentHandlers == nil {
		sentHandlers = make(map[eventstore.EventType]Sent)
	}
	sentHandlers[eventType] = sent
}

func NewNotificationWorker(
	config WorkerConfig,
	commands Commands,
	queries *NotificationQueries,
	channels types.ChannelChains,
) *NotificationWorker {
	return &NotificationWorker{
		config:   config,
		commands: commands,
		queries:  queries,
		channels: channels,
		now:      time.Now,
	}
}

var _ river.Worker[*notification.Request] = (*NotificationWorker)(nil)

func (w *NotificationWorker) Register(workers *river.Workers, queues map[string]river.QueueConfig) {
	river.AddWorker(workers, w)
	queues[notification.QueueName] = river.QueueConfig{
		MaxWorkers: int(w.config.Workers),
	}
}

func (w *NotificationWorker) sendNotificationQueue(ctx context.Context, request *notification.Request, jobID string, notifyUser *query.NotifyUser) error {
	// check early that a "sent" handler exists, otherwise we can cancel early
	sentHandler, ok := sentHandlers[request.EventType]
	if !ok {
		logging.Errorf(`no "sent" handler registered for %s`, request.EventType)
		return channels.NewCancelError(fmt.Errorf("no sent handler registered for %s", request.EventType))
	}

	ctx, err := enrichCtx(ctx, request.TriggeredAtOrigin)
	if err != nil {
		return channels.NewCancelError(err)
	}

	var code string
	if request.Code != nil {
		code, err = crypto.DecryptString(request.Code, w.queries.UserDataCrypto)
		if err != nil {
			return err
		}
	}

	colors, err := w.queries.ActiveLabelPolicyByOrg(ctx, request.UserResourceOwner, false)
	if err != nil {
		return err
	}

	translator, err := w.queries.GetTranslatorWithOrgTexts(ctx, request.UserResourceOwner, request.MessageType)
	if err != nil {
		return err
	}

	generatorInfo := new(senders.CodeGeneratorInfo)
	var notify types.Notify
	switch request.NotificationType {
	case domain.NotificationTypeEmail:
		template, err := w.queries.MailTemplateByOrg(ctx, notifyUser.ResourceOwner, false)
		if err != nil {
			return err
		}
		notify = types.SendEmail(ctx, w.channels, string(template.Template), translator, notifyUser, colors, request.EventType)
	case domain.NotificationTypeSms:
		notify = types.SendSMS(ctx, w.channels, translator, notifyUser, colors, request.EventType, request.Aggregate.InstanceID, jobID, generatorInfo)
	}

	args := request.Args.ToMap()
	args[Code] = code
	// existing notifications use `OTP` as argument for the code
	if request.IsOTP {
		args[OTP] = code
	}

	if err = notify(request.URLTemplate, args, request.MessageType, request.UnverifiedNotificationChannel); err != nil {
		return err
	}

	err = sentHandler(authz.WithInstanceID(ctx, request.Aggregate.InstanceID), w.commands, request.Aggregate.ID, request.Aggregate.ResourceOwner, generatorInfo, args)
	logging.WithFields("instanceID", request.Aggregate.InstanceID, "notification", request.Aggregate.ID).
		OnError(err).Error("could not set notification event on aggregate")
	return nil
}

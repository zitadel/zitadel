package handlers

import (
	"context"
	"errors"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/riverqueue/river"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/queue"
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
	es       *eventstore.Eventstore
	client   *database.DB
	channels types.ChannelChains
	config   WorkerConfig
	now      nowFunc
	backOff  func(current time.Duration) time.Duration
}

// Work implements [river.Worker].
func (w *NotificationWorker) Work(ctx context.Context, job *river.Job[*notification.Request]) error {
	ctx = ContextWithNotifier(ctx, job.Args.Aggregate)

	// if the notification is too old, we can directly cancel
	if job.CreatedAt.Add(w.config.MaxTtl).Before(w.now()) {
		return river.JobCancel(errors.New("notification is too old"))
	}

	// Get the notify user first, so if anything fails afterward we have the current state of the user
	// and can pass that to the retry request.
	// We do not trigger the projection to reduce load on the database. By the time the notification is processed,
	// the user should be projected anyway. If not, it will just wait for the next run.
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

	err = w.sendNotificationQueue(ctx, job.Args, notifyUser)
	if err == nil {
		return nil
	}
	// if retries are disabled or if the error explicitly specifies, we cancel the notification
	// TODO: move max attempts to job config
	if w.config.MaxAttempts <= 1 || errors.Is(err, &channels.CancelError{}) {
		return river.JobCancel(errors.New("notification is too old"))
	}
	return err
	// TODO: handle backoff
}

type WorkerConfig struct {
	LegacyEnabled       bool
	Workers             uint8
	BulkLimit           uint16
	RequeueEvery        time.Duration
	RetryWorkers        uint8
	RetryRequeueEvery   time.Duration
	TransactionDuration time.Duration
	MaxAttempts         uint8
	MaxTtl              time.Duration
	MinRetryDelay       time.Duration
	MaxRetryDelay       time.Duration
	RetryDelayFactor    float32
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
	es *eventstore.Eventstore,
	client *database.DB,
	channels types.ChannelChains,
	queue *queue.Queue,
) *NotificationWorker {
	// make sure the delay does not get less
	if config.RetryDelayFactor < 1 {
		config.RetryDelayFactor = 1
	}
	w := &NotificationWorker{
		config:   config,
		commands: commands,
		queries:  queries,
		es:       es,
		client:   client,
		channels: channels,
		now:      time.Now,
	}
	if !config.LegacyEnabled {
		queue.AddWorkers(w)
	}
	w.backOff = w.exponentialBackOff
	return w
}

var _ river.Worker[*notification.Request] = (*NotificationWorker)(nil)

func (w *NotificationWorker) Register(workers *river.Workers) {
	river.AddWorker(workers, w)
}

func (w *NotificationWorker) Start(ctx context.Context) {
	if w.config.LegacyEnabled {
		return
	}
}

func (w *NotificationWorker) sendNotificationQueue(ctx context.Context, request *notification.Request, notifyUser *query.NotifyUser) error {
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
		notify = types.SendEmail(ctx, w.channels, string(template.Template), translator, notifyUser, colors)
	case domain.NotificationTypeSms:
		notify = types.SendSMS(ctx, w.channels, translator, notifyUser, colors, generatorInfo)
	}

	args := request.Args.ToMap()
	args[Code] = code
	// existing notifications use `OTP` as argument for the code
	if request.IsOTP {
		args[OTP] = code
	}

	return notify(request.URLTemplate, args, request.MessageType, request.UnverifiedNotificationChannel)
}

func (w *NotificationWorker) exponentialBackOff(current time.Duration) time.Duration {
	if current >= w.config.MaxRetryDelay {
		return w.config.MaxRetryDelay
	}
	if current < w.config.MinRetryDelay {
		current = w.config.MinRetryDelay
	}
	t := time.Duration(rand.Int64N(int64(w.config.RetryDelayFactor*float32(current.Nanoseconds()))-current.Nanoseconds()) + current.Nanoseconds())
	if t > w.config.MaxRetryDelay {
		return w.config.MaxRetryDelay
	}
	return t
}

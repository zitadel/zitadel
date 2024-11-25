package handlers

import (
	"context"
	"database/sql"
	"math/rand/v2"
	"slices"
	"strings"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/notification"
)

type NotificationWorker struct {
	commands Commands
	queries  *NotificationQueries
	es       *eventstore.Eventstore
	client   *database.DB
	channels types.ChannelChains
	config   WorkerConfig
	now      nowFunc
}

type WorkerConfig struct {
	Workers               uint8
	BulkLimit             uint16
	RequeueEvery          time.Duration
	HandleActiveInstances time.Duration
	TransactionDuration   time.Duration
	MaxAttempts           uint8
	MaxTtl                time.Duration
	MinRetryDelay         time.Duration
	MaxRetryDelay         time.Duration
	RetryDelayFactor      float32
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
) *NotificationWorker {
	// make sure the delay does not get less
	if config.RetryDelayFactor < 1 {
		config.RetryDelayFactor = 1
	}
	return &NotificationWorker{
		config:   config,
		commands: commands,
		queries:  queries,
		es:       es,
		client:   client,
		channels: channels,
		now:      time.Now,
	}
}

func (w *NotificationWorker) Start(ctx context.Context) {
	for i := 0; i < int(w.config.Workers); i++ {
		go w.schedule(ctx, i)
	}
}

func (w *NotificationWorker) reduceNotificationRequested(ctx context.Context, tx *sql.Tx, event *notification.RequestedEvent) (err error) {
	ctx = HandlerContext(event.Aggregate())

	// if the notification is too old, we can directly cancel
	if event.CreatedAt().Add(w.config.MaxTtl).Before(w.now()) {
		return w.commands.NotificationCanceled(ctx, tx, event.Aggregate().ID, authz.GetInstance(ctx).InstanceID(), err)
	}

	// Get the notify user first, so if anything fails afterward we have the current state of the user
	// and can pass that to the retry request.
	notifyUser, err := w.queries.GetNotifyUserByID(ctx, true, event.UserID)
	if err != nil {
		return err
	}

	// The domain claimed event requires the domain as argument, but lacks the user when creating the request event.
	// Since we set it into the request arguments, it will be passed into a potential retry event.
	if _, ok := event.Request.Args["Domain"]; !ok && event.RequiresPreviousDomain {
		index := strings.LastIndex(notifyUser.LastEmail, "@")
		event.Request.Args["Domain"] = notifyUser.LastEmail[index+1:]
	}

	err = w.sendNotification(ctx, tx, event.Request, notifyUser, event)
	if err != nil {
		return w.commands.NotificationRetryRequested(
			ctx,
			tx,
			event.Aggregate().ID,
			authz.GetInstance(ctx).InstanceID(),
			notificationEventToRequest(event.Request, notifyUser, w.backOff(0)),
			err,
		)
	}
	return nil
}

func (w *NotificationWorker) reduceNotificationRetry(ctx context.Context, tx *sql.Tx, event *notification.RetryRequestedEvent) (err error) {
	ctx = HandlerContext(event.Aggregate())

	// if the notification is too old, we can directly cancel
	if event.CreatedAt().Add(w.config.MaxTtl).Before(w.now()) {
		return w.commands.NotificationCanceled(ctx, tx, event.Aggregate().ID, authz.GetInstance(ctx).InstanceID(), err)
	}

	if event.CreatedAt().Add(event.BackOff).After(w.now()) {
		return nil
	}
	err = w.sendNotification(ctx, tx, event.Request, event.NotifyUser, event)
	if err == nil {
		return nil
	}
	// if the max attempts are reached, we cancel the notification
	if event.Sequence() >= uint64(w.config.MaxAttempts) {
		return w.commands.NotificationCanceled(ctx, tx, event.Aggregate().ID, authz.GetInstance(ctx).InstanceID(), err)
	}
	// otherwise we retry after a backoff delay
	return w.commands.NotificationRetryRequested(ctx, tx, event.Aggregate().ID, authz.GetInstance(ctx).InstanceID(), notificationEventToRequest(
		event.Request,
		event.NotifyUser,
		w.backOff(event.BackOff),
	), err)
}

func (w *NotificationWorker) sendNotification(ctx context.Context, tx *sql.Tx, request notification.Request, notifyUser *query.NotifyUser, e eventstore.Event) error {
	ctx, err := enrichCtx(ctx, request.TriggeredAtOrigin)
	if err != nil {
		err := w.commands.NotificationCanceled(ctx, tx, e.Aggregate().ID, authz.GetInstance(ctx).InstanceID(), err)
		logging.WithFields("instanceID", authz.GetInstance(ctx).InstanceID(), "notification", e.Aggregate().ID).
			OnError(err).Error("could not cancel notification")
		return nil
	}

	// check early that a "sent" handler exists, otherwise we can cancel early
	sender, ok := sentHandlers[request.EventType]
	if !ok {
		err := w.commands.NotificationCanceled(ctx, tx, e.Aggregate().ID, authz.GetInstance(ctx).InstanceID(), err)
		logging.WithFields("instanceID", authz.GetInstance(ctx).InstanceID(), "notification", e.Aggregate().ID).
			OnError(err).Errorf(`no "sent" handler registered for %s`, request.EventType)
		return nil
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
		notify = types.SendEmail(ctx, w.channels, string(template.Template), translator, notifyUser, colors, e)
	case domain.NotificationTypeSms:
		notify = types.SendSMS(ctx, w.channels, translator, notifyUser, colors, e, generatorInfo)
	}

	args := request.Args
	if len(args) == 0 {
		args = make(map[string]any)
	}
	if code != "" {
		args["Code"] = code
		// existing notifications use `OTP` as argument for the code
		if request.IsOTP {
			args["OTP"] = code
		}
	}

	if err := notify(request.URLTemplate, args, request.MessageType, request.UnverifiedNotificationChannel); err != nil {
		return err
	}
	err = w.commands.NotificationSent(ctx, tx, e.Aggregate().ID, authz.GetInstance(ctx).InstanceID())
	if err != nil {
		logging.WithFields("instanceID", authz.GetInstance(ctx).InstanceID(), "notification", e.Aggregate().ID).
			OnError(err).Error("could not set sent notification event")
	}
	err = sender(ctx, w.commands, request.NotificationAggregateID(), request.NotificationAggregateResourceOwner(), generatorInfo, args)
	logging.WithFields("instanceID", authz.GetInstance(ctx).InstanceID(), "notification", e.Aggregate().ID).
		OnError(err).Error("could not set notification event on aggregate")
	return nil
}

func (w *NotificationWorker) backOff(current time.Duration) time.Duration {
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

func notificationEventToRequest(e notification.Request, notifyUser *query.NotifyUser, backoff time.Duration) *command.NotificationRetryRequest {
	return &command.NotificationRetryRequest{
		NotificationRequest: command.NotificationRequest{
			UserID:                        e.UserID,
			UserResourceOwner:             e.UserResourceOwner,
			TriggerOrigin:                 e.TriggeredAtOrigin,
			URLTemplate:                   e.URLTemplate,
			Code:                          e.Code,
			CodeExpiry:                    e.CodeExpiry,
			EventType:                     e.EventType,
			NotificationType:              e.NotificationType,
			MessageType:                   e.MessageType,
			UnverifiedNotificationChannel: e.UnverifiedNotificationChannel,
			Args:                          e.Args,
			AggregateID:                   e.AggregateID,
			AggregateResourceOwner:        e.AggregateResourceOwner,
			IsOTP:                         e.IsOTP,
		},
		BackOff:    backoff,
		NotifyUser: notifyUser,
	}
}

func (w *NotificationWorker) schedule(ctx context.Context, workerID int) {
	t := time.NewTimer(0)

	for {
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
			instances, err := w.queryInstances(ctx)
			w.log(workerID).OnError(err).Debug("unable to query instances")

			w.triggerInstances(call.WithTimestamp(ctx), instances, workerID)
			t.Reset(w.config.RequeueEvery)
		}
	}
}

func (w *NotificationWorker) log(workerID int) *logging.Entry {
	return logging.WithFields("notification worker", workerID)
}

func (w *NotificationWorker) queryInstances(ctx context.Context) ([]string, error) {
	if w.config.HandleActiveInstances == 0 {
		return w.existingInstances(ctx)
	}

	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsInstanceIDs).
		AwaitOpenTransactions().
		AllowTimeTravel().
		CreationDateAfter(w.now().Add(-1 * w.config.HandleActiveInstances))

	return w.es.InstanceIDs(ctx, w.config.RequeueEvery, false, query)
}

func (w *NotificationWorker) existingInstances(ctx context.Context) ([]string, error) {
	ai := existingInstances{}
	if err := w.es.FilterToQueryReducer(ctx, &ai); err != nil {
		return nil, err
	}

	return ai, nil
}

func (w *NotificationWorker) triggerInstances(ctx context.Context, instances []string, workerID int) {
	for _, instance := range instances {
		instanceCtx := authz.WithInstanceID(ctx, instance)

		err := w.trigger(instanceCtx, workerID)
		w.log(workerID).WithField("instance", instance).OnError(err).Info("trigger failed")
	}
}

func (w *NotificationWorker) trigger(ctx context.Context, workerID int) (err error) {
	txCtx := ctx
	if w.config.TransactionDuration > 0 {
		var cancel, cancelTx func()
		txCtx, cancelTx = context.WithCancel(ctx)
		defer cancelTx()
		ctx, cancel = context.WithTimeout(ctx, w.config.TransactionDuration)
		defer cancel()
	}
	tx, err := w.client.BeginTx(txCtx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			w.log(workerID).OnError(rollbackErr).Debug("unable to rollback tx")
			return
		}
	}()

	// query events and lock them for update (with skip locked)
	searchQuery := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		LockRowsDuringTx(tx, eventstore.LockOptionSkipLocked).
		// Messages older than the MaxTTL, we can be ignored.
		// The first attempt of a retry might still be older than the TTL and needs to be filtered out later on.
		CreationDateAfter(w.now().Add(-1*w.config.MaxTtl)).
		Limit(uint64(w.config.BulkLimit)).
		AddQuery().
		AggregateTypes(notification.AggregateType).
		EventTypes(notification.RequestedType, notification.RetryRequestedType, notification.CanceledType, notification.SentType).
		Builder()

	events, err := w.es.Filter(ctx, searchQuery)
	if err != nil {
		return err
	}

	for i := len(events) - 1; i > 0; i-- {
		if len(events)-1 < i {
			continue
		}
		event := events[i]
		if event != nil && (event.Type() != notification.RequestedType) {
			events = slices.DeleteFunc(events, func(e eventstore.Event) bool {
				if e.Aggregate().ID != event.Aggregate().ID {
					return false
				}
				a := e.Sequence() < event.Sequence()
				return a
			})
		}
	}
	// If there aren't any events or no unlocked event terminate early and start a new run.
	if len(events) == 0 {
		return nil
	}

	w.log(workerID).
		WithField("instanceID", authz.GetInstance(ctx).InstanceID()).
		WithField("events", len(events)).
		Info("handling notification events")

	for _, event := range events {
		var err error
		switch e := event.(type) {
		case *notification.RequestedEvent:
			err = w.reduceNotificationRequested(ctx, tx, e)
		case *notification.RetryRequestedEvent:
			err = w.reduceNotificationRetry(ctx, tx, e)
		}
		w.log(workerID).OnError(err).
			WithField("instanceID", authz.GetInstance(ctx).InstanceID()).
			WithField("notificationID", event.Aggregate().ID).
			WithField("sequence", event.Sequence()).
			WithField("type", event.Type()).
			Error("could not push notification event")
	}
	return database.CloseTransaction(tx, nil) //TODO: do we need to rollback on an error?
}

type existingInstances []string

// AppendEvents implements eventstore.QueryReducer.
func (ai *existingInstances) AppendEvents(events ...eventstore.Event) {
	for _, event := range events {
		switch event.Type() {
		case instance.InstanceAddedEventType:
			*ai = append(*ai, event.Aggregate().InstanceID)
		case instance.InstanceRemovedEventType:
			*ai = slices.DeleteFunc(*ai, func(s string) bool {
				return s == event.Aggregate().InstanceID
			})
		}
	}
}

// Query implements eventstore.QueryReducer.
func (*existingInstances) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		EventTypes(
			instance.InstanceAddedEventType,
			instance.InstanceRemovedEventType,
		).
		Builder()
}

// Reduce implements eventstore.QueryReducer.
// reduce is not used as events are reduced during AppendEvents
func (*existingInstances) Reduce() error {
	return nil
}

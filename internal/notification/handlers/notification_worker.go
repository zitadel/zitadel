package handlers

import (
	"context"
	"database/sql"
	"errors"
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
	"github.com/zitadel/zitadel/internal/notification/channels"
	"github.com/zitadel/zitadel/internal/notification/senders"
	"github.com/zitadel/zitadel/internal/notification/types"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/notification"
)

const (
	Code = "Code"
	OTP  = "OTP"
)

type NotificationWorker struct {
	commands Commands
	queries  *NotificationQueries
	es       *eventstore.Eventstore
	client   *database.DB
	channels types.ChannelChains
	config   WorkerConfig
	now      nowFunc
	backOff  func(current time.Duration) time.Duration
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
	w.backOff = w.exponentialBackOff
	return w
}

func (w *NotificationWorker) Start(ctx context.Context) {
	if w.config.LegacyEnabled {
		return
	}
	for i := 0; i < int(w.config.Workers); i++ {
		go w.schedule(ctx, i, false)
	}
	for i := 0; i < int(w.config.RetryWorkers); i++ {
		go w.schedule(ctx, i, true)
	}
}

func (w *NotificationWorker) reduceNotificationRequested(ctx, txCtx context.Context, tx *sql.Tx, event *notification.RequestedEvent) (err error) {
	ctx = ContextWithNotifier(ctx, event.Aggregate())

	// if the notification is too old, we can directly cancel
	if event.CreatedAt().Add(w.config.MaxTtl).Before(w.now()) {
		return w.commands.NotificationCanceled(txCtx, tx, event.Aggregate().ID, event.Aggregate().ResourceOwner, nil)
	}

	// Get the notify user first, so if anything fails afterward we have the current state of the user
	// and can pass that to the retry request.
	// We do not trigger the projection to reduce load on the database. By the time the notification is processed,
	// the user should be projected anyway. If not, it will just wait for the next run.
	notifyUser, err := w.queries.GetNotifyUserByID(ctx, false, event.UserID)
	if err != nil {
		return err
	}

	// The domain claimed event requires the domain as argument, but lacks the user when creating the request event.
	// Since we set it into the request arguments, it will be passed into a potential retry event.
	if event.RequiresPreviousDomain && event.Request.Args != nil && event.Request.Args.Domain == "" {
		index := strings.LastIndex(notifyUser.LastEmail, "@")
		event.Request.Args.Domain = notifyUser.LastEmail[index+1:]
	}

	err = w.sendNotification(ctx, txCtx, tx, event.Request, notifyUser, event)
	if err == nil {
		return nil
	}
	// if retries are disabled or if the error explicitly specifies, we cancel the notification
	if w.config.MaxAttempts <= 1 || errors.Is(err, &channels.CancelError{}) {
		return w.commands.NotificationCanceled(txCtx, tx, event.Aggregate().ID, event.Aggregate().ResourceOwner, err)
	}
	// otherwise we retry after a backoff delay
	return w.commands.NotificationRetryRequested(
		txCtx,
		tx,
		event.Aggregate().ID,
		event.Aggregate().ResourceOwner,
		notificationEventToRequest(event.Request, notifyUser, w.backOff(0)),
		err,
	)
}

func (w *NotificationWorker) reduceNotificationRetry(ctx, txCtx context.Context, tx *sql.Tx, event *notification.RetryRequestedEvent) (err error) {
	ctx = ContextWithNotifier(ctx, event.Aggregate())

	// if the notification is too old, we can directly cancel
	if event.CreatedAt().Add(w.config.MaxTtl).Before(w.now()) {
		return w.commands.NotificationCanceled(txCtx, tx, event.Aggregate().ID, event.Aggregate().ResourceOwner, err)
	}

	if event.CreatedAt().Add(event.BackOff).After(w.now()) {
		return nil
	}
	err = w.sendNotification(ctx, txCtx, tx, event.Request, event.NotifyUser, event)
	if err == nil {
		return nil
	}
	// if the max attempts are reached or if the error explicitly specifies, we cancel the notification
	if event.Sequence() >= uint64(w.config.MaxAttempts) || errors.Is(err, &channels.CancelError{}) {
		return w.commands.NotificationCanceled(txCtx, tx, event.Aggregate().ID, event.Aggregate().ResourceOwner, err)
	}
	// otherwise we retry after a backoff delay
	return w.commands.NotificationRetryRequested(txCtx, tx, event.Aggregate().ID, event.Aggregate().ResourceOwner, notificationEventToRequest(
		event.Request,
		event.NotifyUser,
		w.backOff(event.BackOff),
	), err)
}

func (w *NotificationWorker) sendNotification(ctx, txCtx context.Context, tx *sql.Tx, request notification.Request, notifyUser *query.NotifyUser, e eventstore.Event) error {
	ctx, err := enrichCtx(ctx, request.TriggeredAtOrigin)
	if err != nil {
		return channels.NewCancelError(err)
	}

	// check early that a "sent" handler exists, otherwise we can cancel early
	sentHandler, ok := sentHandlers[request.EventType]
	if !ok {
		logging.Errorf(`no "sent" handler registered for %s`, request.EventType)
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
		notify = types.SendEmail(ctx, w.channels, string(template.Template), translator, notifyUser, colors, e)
	case domain.NotificationTypeSms:
		notify = types.SendSMS(ctx, w.channels, translator, notifyUser, colors, e, generatorInfo)
	}

	args := request.Args.ToMap()
	args[Code] = code
	// existing notifications use `OTP` as argument for the code
	if request.IsOTP {
		args[OTP] = code
	}

	if err := notify(request.URLTemplate, args, request.MessageType, request.UnverifiedNotificationChannel); err != nil {
		return err
	}
	err = w.commands.NotificationSent(txCtx, tx, e.Aggregate().ID, e.Aggregate().ResourceOwner)
	if err != nil {
		// In case the notification event cannot be pushed, we most likely cannot create a retry or cancel event.
		// Therefore, we'll only log the error and also do not need to try to push to the user / session.
		logging.WithFields("instanceID", authz.GetInstance(ctx).InstanceID(), "notification", e.Aggregate().ID).
			OnError(err).Error("could not set sent notification event")
		return nil
	}
	err = sentHandler(txCtx, w.commands, request.NotificationAggregateID(), request.NotificationAggregateResourceOwner(), generatorInfo, args)
	logging.WithFields("instanceID", authz.GetInstance(ctx).InstanceID(), "notification", e.Aggregate().ID).
		OnError(err).Error("could not set notification event on aggregate")
	return nil
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

func (w *NotificationWorker) schedule(ctx context.Context, workerID int, retry bool) {
	t := time.NewTimer(0)

	for {
		select {
		case <-ctx.Done():
			t.Stop()
			w.log(workerID, retry).Info("scheduler stopped")
			return
		case <-t.C:
			instances, err := w.queryInstances(ctx, retry)
			w.log(workerID, retry).OnError(err).Error("unable to query instances")

			w.triggerInstances(call.WithTimestamp(ctx), instances, workerID, retry)
			if retry {
				t.Reset(w.config.RetryRequeueEvery)
				continue
			}
			t.Reset(w.config.RequeueEvery)
		}
	}
}

func (w *NotificationWorker) log(workerID int, retry bool) *logging.Entry {
	return logging.WithFields("notification worker", workerID, "retries", retry)
}

func (w *NotificationWorker) queryInstances(ctx context.Context, retry bool) ([]string, error) {
	return w.queries.ActiveInstances(), nil
}

func (w *NotificationWorker) triggerInstances(ctx context.Context, instances []string, workerID int, retry bool) {
	for _, instance := range instances {
		instanceCtx := authz.WithInstanceID(ctx, instance)

		err := w.trigger(instanceCtx, workerID, retry)
		w.log(workerID, retry).WithField("instance", instance).OnError(err).Info("trigger failed")
	}
}

func (w *NotificationWorker) trigger(ctx context.Context, workerID int, retry bool) (err error) {
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
		err = database.CloseTransaction(tx, err)
	}()

	events, err := w.searchEvents(txCtx, tx, retry)
	if err != nil {
		return err
	}

	// If there aren't any events or no unlocked event terminate early and start a new run.
	if len(events) == 0 {
		return nil
	}

	w.log(workerID, retry).
		WithField("instanceID", authz.GetInstance(ctx).InstanceID()).
		WithField("events", len(events)).
		Info("handling notification events")

	for _, event := range events {
		var err error
		switch e := event.(type) {
		case *notification.RequestedEvent:
			w.createSavepoint(txCtx, tx, event, workerID, retry)
			err = w.reduceNotificationRequested(ctx, txCtx, tx, e)
		case *notification.RetryRequestedEvent:
			w.createSavepoint(txCtx, tx, event, workerID, retry)
			err = w.reduceNotificationRetry(ctx, txCtx, tx, e)
		}
		if err != nil {
			w.log(workerID, retry).OnError(err).
				WithField("instanceID", authz.GetInstance(ctx).InstanceID()).
				WithField("notificationID", event.Aggregate().ID).
				WithField("sequence", event.Sequence()).
				WithField("type", event.Type()).
				Error("could not handle notification event")
			// if we have an error, we rollback to the savepoint and continue with the next event
			// we use the txCtx to make sure we can rollback the transaction in case the ctx is canceled
			w.rollbackToSavepoint(txCtx, tx, event, workerID, retry)
		}
		// if the context is canceled, we stop the processing
		if ctx.Err() != nil {
			return nil
		}
	}
	return nil
}

func (w *NotificationWorker) latestRetries(events []eventstore.Event) []eventstore.Event {
	for i := len(events) - 1; i > 0; i-- {
		// since we delete during the iteration, we need to make sure we don't panic
		if len(events) <= i {
			continue
		}
		// delete all the previous retries of the same notification
		events = slices.DeleteFunc(events, func(e eventstore.Event) bool {
			return e.Aggregate().ID == events[i].Aggregate().ID &&
				e.Sequence() < events[i].Sequence()
		})
	}
	return events
}

func (w *NotificationWorker) createSavepoint(ctx context.Context, tx *sql.Tx, event eventstore.Event, workerID int, retry bool) {
	_, err := tx.ExecContext(ctx, "SAVEPOINT notification_send")
	w.log(workerID, retry).OnError(err).
		WithField("instanceID", authz.GetInstance(ctx).InstanceID()).
		WithField("notificationID", event.Aggregate().ID).
		WithField("sequence", event.Sequence()).
		WithField("type", event.Type()).
		Error("could not create savepoint for notification event")
}

func (w *NotificationWorker) rollbackToSavepoint(ctx context.Context, tx *sql.Tx, event eventstore.Event, workerID int, retry bool) {
	_, err := tx.ExecContext(ctx, "ROLLBACK TO SAVEPOINT notification_send")
	w.log(workerID, retry).OnError(err).
		WithField("instanceID", authz.GetInstance(ctx).InstanceID()).
		WithField("notificationID", event.Aggregate().ID).
		WithField("sequence", event.Sequence()).
		WithField("type", event.Type()).
		Error("could not rollback to savepoint for notification event")
}

func (w *NotificationWorker) searchEvents(ctx context.Context, tx *sql.Tx, retry bool) ([]eventstore.Event, error) {
	if retry {
		return w.searchRetryEvents(ctx, tx)
	}
	// query events and lock them for update (with skip locked)
	searchQuery := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		LockRowsDuringTx(tx, eventstore.LockOptionSkipLocked).
		// Messages older than the MaxTTL, we can be ignored.
		// The first attempt of a retry might still be older than the TTL and needs to be filtered out later on.
		CreationDateAfter(w.now().Add(-1*w.config.MaxTtl)).
		Limit(uint64(w.config.BulkLimit)).
		AddQuery().
		AggregateTypes(notification.AggregateType).
		EventTypes(notification.RequestedType).
		Builder().
		ExcludeAggregateIDs().
		EventTypes(notification.RetryRequestedType, notification.CanceledType, notification.SentType).
		Builder()
	//nolint:staticcheck
	return w.es.Filter(ctx, searchQuery)
}

func (w *NotificationWorker) searchRetryEvents(ctx context.Context, tx *sql.Tx) ([]eventstore.Event, error) {
	// query events and lock them for update (with skip locked)
	searchQuery := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		LockRowsDuringTx(tx, eventstore.LockOptionSkipLocked).
		// Messages older than the MaxTTL, we can be ignored.
		// The first attempt of a retry might still be older than the TTL and needs to be filtered out later on.
		CreationDateAfter(w.now().Add(-1*w.config.MaxTtl)).
		AddQuery().
		AggregateTypes(notification.AggregateType).
		EventTypes(notification.RetryRequestedType).
		Builder().
		ExcludeAggregateIDs().
		EventTypes(notification.CanceledType, notification.SentType).
		Builder()
	//nolint:staticcheck
	events, err := w.es.Filter(ctx, searchQuery)
	if err != nil {
		return nil, err
	}
	return w.latestRetries(events), nil
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

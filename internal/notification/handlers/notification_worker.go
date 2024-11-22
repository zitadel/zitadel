package handlers

import (
	"context"
	"database/sql"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/i18n"
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
	//BulkLimit             uint16
	RequeueEvery time.Duration
	//RetryFailedAfter      time.Duration
	HandleActiveInstances time.Duration
	TransactionDuration   time.Duration
	MaxAttempts           uint8
	MaxTtl                time.Duration
	MinRetryDelay         time.Duration
	MaxRetryDelay         time.Duration
}

var sentHandlers map[eventstore.EventType]Sent

// nowFunc makes [time.Now] mockable
type nowFunc func() time.Time

type Sent func(ctx context.Context, commands Commands, id, orgID string, generatorInfo *senders.CodeGeneratorInfo, args map[string]any) error

func NewNotificationWorker(
	ctx context.Context,
	config WorkerConfig,
	commands Commands,
	queries *NotificationQueries,
	es *eventstore.Eventstore,
	client *database.DB,
	channels types.ChannelChains,
) *NotificationWorker {
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
	go w.schedule(ctx)
}

func RegisterSentHandler(eventType eventstore.EventType, sent Sent) {
	if sentHandlers == nil {
		sentHandlers = make(map[eventstore.EventType]Sent)
	}
	sentHandlers[eventType] = sent
}

func (w *NotificationWorker) reduceNotificationAdded(ctx context.Context, tx *sql.Tx, event *notification.RequestedEvent) (err error) {
	ctx = HandlerContext(event.Aggregate())

	// Get the notify user first, so if anything fails afterward we have the current state of the user
	// and can pass that to the retry request.
	notifyUser, err := w.queries.GetNotifyUserByID(ctx, true, event.UserID)
	if err != nil {
		return err
	}

	var code string
	if event.Code != nil {
		code, err = crypto.DecryptString(event.Code, w.queries.UserDataCrypto)
		if err != nil {
			// TODO: do we need to retry that?
			return w.commands.NotificationCanceled(ctx, event.Aggregate().ID, authz.GetInstance(ctx).InstanceID(), err)
		}
	}

	colors, err := w.queries.ActiveLabelPolicyByOrg(ctx, event.UserResourceOwner, false)
	if err != nil {
		return w.commands.NotificationRetryRequested(ctx, event.Aggregate().ID, authz.GetInstance(ctx).InstanceID(), notificationEventToRequest(event.Request, notifyUser, w.backOff(0)), err)
	}

	translator, err := w.queries.GetTranslatorWithOrgTexts(ctx, event.UserResourceOwner, event.MessageType)
	if err != nil {
		return w.commands.NotificationRetryRequested(ctx, event.Aggregate().ID, authz.GetInstance(ctx).InstanceID(), notificationEventToRequest(event.Request, notifyUser, w.backOff(0)), err)
	}
	err = w.sendNotification(ctx, event.Request, code, notifyUser, colors, translator, event)
	if err != nil {
		return w.commands.NotificationRetryRequested(ctx, event.Aggregate().ID, authz.GetInstance(ctx).InstanceID(), notificationEventToRequest(event.Request, notifyUser, w.backOff(0)), err)
	}
	return nil
}

func (w *NotificationWorker) backOff(current time.Duration) time.Duration {
	if current >= w.config.MaxRetryDelay {
		return w.config.MaxRetryDelay
	}
	if current < w.config.MinRetryDelay {
		current = w.config.MinRetryDelay
	}
	t := time.Duration(rand.Int64N(int64(1.5*float64(current.Nanoseconds()))-current.Nanoseconds()) + current.Nanoseconds())
	if t > w.config.MaxRetryDelay {
		return w.config.MaxRetryDelay
	}
	return t
}

func (w *NotificationWorker) reduceNotificationRetry(ctx context.Context, tx *sql.Tx, event *notification.RetryRequestedEvent) (err error) {
	ctx = HandlerContext(event.Aggregate())

	var code string
	if event.Code != nil {
		code, err = crypto.DecryptString(event.Code, w.queries.UserDataCrypto)
		if err != nil {
			// TODO: do we need to retry that?
			return w.commands.NotificationCanceled(ctx, event.Aggregate().ID, authz.GetInstance(ctx).InstanceID(), err)
		}
	}

	colors, err := w.queries.ActiveLabelPolicyByOrg(ctx, event.UserResourceOwner, false)
	if err != nil {
		return err
	}

	translator, err := w.queries.GetTranslatorWithOrgTexts(ctx, event.UserResourceOwner, event.MessageType)
	if err != nil {
		return err
	}
	err = w.sendNotification(ctx, event.Request, code, event.NotifyUser, colors, translator, event)
	if err != nil {
		if event.Sequence() > uint64(w.config.MaxAttempts) {
			return w.commands.NotificationCanceled(ctx, event.Aggregate().ID, authz.GetInstance(ctx).InstanceID(), err)
		}
		return w.commands.NotificationRetryRequested(ctx, event.Aggregate().ID, authz.GetInstance(ctx).InstanceID(), notificationEventToRequest(event.Request, event.NotifyUser, w.backOff(event.BackOff)), err)
	}
	return nil
}

func (w *NotificationWorker) sendNotification(
	ctx context.Context,
	request notification.Request,
	code string,
	notifyUser *query.NotifyUser,
	colors *query.LabelPolicy,
	translator *i18n.Translator,
	e eventstore.Event,
) error {
	ctx, err := enrichCtx(ctx, request.TriggeredAtOrigin)
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
	sender, ok := sentHandlers[request.EventType]
	if !ok {
		return nil //TODO: !
	}
	if err := w.commands.NotificationSent(ctx, e.Aggregate().ID, authz.GetInstance(ctx).InstanceID()); err != nil {
		return err
	}
	return sender(ctx, w.commands, request.NotificationAggregateID(), request.NotificationAggregateResourceOwner(), generatorInfo, args)
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

func (w *NotificationWorker) schedule(ctx context.Context) {
	t := time.NewTimer(0)

	for {
		select {
		case <-ctx.Done():
			t.Stop()
			return
		case <-t.C:
			instances, err := w.queryInstances(ctx)
			w.log().OnError(err).Debug("unable to query instances")

			w.triggerInstances(call.WithTimestamp(ctx), instances)
			t.Reset(w.config.RequeueEvery)
		}
	}
}

func (w *NotificationWorker) log() *logging.Entry {
	return logging.WithFields("projection", "notification worker")
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

func (w *NotificationWorker) existingInstances(ctx context.Context) ([]string, error) {
	ai := existingInstances{}
	if err := w.es.FilterToQueryReducer(ctx, &ai); err != nil {
		return nil, err
	}

	return ai, nil
}

func (w *NotificationWorker) triggerInstances(ctx context.Context, instances []string /*triggerOpts ...TriggerOpt*/) {
	for _, instance := range instances {
		instanceCtx := authz.WithInstanceID(ctx, instance)

		err := w.Trigger(instanceCtx /*triggerOpts...*/)
		w.log().WithField("instance", instance).OnError(err).Info("trigger failed")
	}
}

func (w *NotificationWorker) Trigger(ctx context.Context /*opts ...TriggerOpt*/) (err error) {
	//config := new(triggerConfig)
	//for _, opt := range opts {
	//	opt(config)
	//}

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
		//if err != nil {
		rollbackErr := tx.Rollback()
		w.log().OnError(rollbackErr).Debug("unable to rollback tx")
		return
		//}
	}()

	// query events with skip locked
	searchQuery := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		SetTx(tx).
		CreationDateAfter(w.now().Add(-1*w.config.MaxTtl)).
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
	if len(events) == 0 {
		return nil
	}

	w.log().Info("1")

	for _, event := range events {
		var err error
		switch e := event.(type) {
		case *notification.RequestedEvent:
			err = w.reduceNotificationAdded(ctx, tx, e)
		case *notification.RetryRequestedEvent:
			if e.CreatedAt().Add(e.BackOff).After(time.Now()) {
				continue
			}
			err = w.reduceNotificationRetry(ctx, tx, e)
		}
		if err != nil {
			tt := w.commands.NotificationCanceled(ctx, event.Aggregate().ID, authz.GetInstance(ctx).InstanceID(), err)
			w.log().OnError(tt).Info("could not cancel notification")
		}
	}
	return nil
}

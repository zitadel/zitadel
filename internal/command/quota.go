package command

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/zitadel/logging"
	caos_errs "github.com/zitadel/zitadel/internal/errors"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/repository/quota"

	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"

	"github.com/zitadel/zitadel/internal/domain"
)

type QuotaUnit string

const (
	QuotaRequestsAllAuthenticated QuotaUnit = "requests.all.authenticated"
	QuotaActionsAllRunsSeconds              = "actions.all.runs.seconds"
)

func (q *QuotaUnit) Enum() quota.Unit {
	switch *q {
	case QuotaRequestsAllAuthenticated:
		return quota.RequestsAllAuthenticated
	case QuotaActionsAllRunsSeconds:
		return quota.ActionsAllRunsSeconds
	default:
		return quota.Unimplemented
	}
}

func (c *Commands) AddQuota(
	ctx context.Context,
	q *AddQuota,
) (*domain.ObjectDetails, error) {
	instanceId := authz.GetInstance(ctx).InstanceID()
	aggregateId, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
	}

	wm, err := c.getQuotaWriteModel(ctx, instanceId, instanceId, q.Unit.Enum())
	if err != nil {
		return nil, err
	}

	if wm.active {
		return nil, caos_errs.ThrowAlreadyExists(nil, "COMMAND-WDfFf", "Errors.Quota.AlreadyExists")
	}

	aggregate := quota.NewAggregate(aggregateId, instanceId, instanceId)
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.AddQuotaCommand(aggregate, q))
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(wm, events...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

func (c *Commands) RemoveQuota(ctx context.Context, unit QuotaUnit) (*domain.ObjectDetails, error) {
	instanceId := authz.GetInstance(ctx).InstanceID()

	q, err := c.getQuotaWriteModel(ctx, instanceId, instanceId, unit.Enum())
	if err != nil {
		return nil, err
	}

	if !q.active {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-WDfFf", "Errors.Quota.NotFound")
	}

	aggregate := quota.NewAggregate(q.AggregateID, instanceId, instanceId)

	events := []eventstore.Command{
		quota.NewRemovedEvent(ctx, &aggregate.Aggregate, unit.Enum()),
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(q, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&q.WriteModel), nil
}

func (c *Commands) getQuotaWriteModel(ctx context.Context, instanceId, resourceOwner string, unit quota.Unit) (*quotaWriteModel, error) {
	wm := newQuotaWriteModel(instanceId, resourceOwner, unit)
	err := c.eventstore.FilterToQueryReducer(ctx, wm)
	if err != nil {
		return nil, err
	}
	return wm, nil
}

type QuotaNotification struct {
	Percent uint64
	Repeat  bool
	CallURL string
}

type QuotaNotifications []*QuotaNotification

func (q *QuotaNotifications) toAddedEventNotifications(genID func() string) []*quota.AddedEventNotification {
	if q == nil {
		return nil
	}

	notifications := make([]*quota.AddedEventNotification, len(*q))
	for idx, notification := range *q {

		notifications[idx] = &quota.AddedEventNotification{
			ID:      genID(),
			Percent: notification.Percent,
			Repeat:  notification.Repeat,
			CallURL: notification.CallURL,
		}
	}

	return notifications
}

type AddQuota struct {
	Unit          QuotaUnit
	From          time.Time
	Interval      time.Duration
	Amount        uint64
	Limit         bool
	Notifications QuotaNotifications
}

func (q *AddQuota) isValid() bool {
	for _, notification := range q.Notifications {
		if err := isUrl(notification.CallURL); err != nil || notification.Percent < 1 {
			return false
		}
	}

	return q.Unit.Enum() != quota.Unimplemented &&
		!q.From.IsZero() &&
		q.Amount > 0 &&
		q.Interval > time.Minute &&
		(q.Limit || len(q.Notifications) > 0)
}

func (c *Commands) AddQuotaCommand(a *quota.Aggregate, q *AddQuota) preparation.Validation {
	return func() (preparation.CreateCommands, error) {

		if !q.isValid() {
			return nil, errors.ThrowInvalidArgument(nil, "QUOTA-pBfjq", "Errors.Invalid.Argument") // TODO: Better error message?
		}

		return func(ctx context.Context, filter preparation.FilterToQueryReducer) (cmd []eventstore.Command, err error) {

				genID := func() string {
					id, genErr := c.idGenerator.Next()
					if genErr != nil {
						err = genErr
					}
					return id
				}

				return []eventstore.Command{quota.NewAddedEvent(
					ctx,
					&a.Aggregate,
					q.Unit.Enum(),
					q.From,
					q.Interval,
					q.Amount,
					q.Limit,
					q.Notifications.toAddedEventNotifications(genID),
				)}, err
			},
			nil
	}
}

func isUrl(str string) error {
	u, err := url.Parse(str)
	if err != nil {
		return err
	}

	if u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("url %s is invalid", str)
	}

	return nil
}

// ReportUsage calls notification hooks if necessary and returns if usage should be limited
func (c *Commands) ReportUsage(ctx context.Context, dueNotifications []*quota.NotifiedEvent) error {

	for _, notification := range dueNotifications {
		alreadyNotified, err := isAlreadNotified(ctx, c.eventstore, notification)
		if err != nil {
			return err
		}

		if alreadyNotified {
			// TODO: Debugf
			logging.Infof(
				"quota notification with ID %s and threshold %d was already notified in this period",
				notification.ID,
				notification.Threshold,
			)
			continue
		}

		if err = notify(ctx, notification); err != nil {
			if err != nil {
				return err
			}
		}

		if _, err = c.eventstore.Push(ctx, notification); err != nil {
			return err
		}
	}

	return nil
}

func isAlreadNotified(ctx context.Context, es *eventstore.Eventstore, notification *quota.NotifiedEvent) (bool, error) {

	events, err := es.Filter(
		ctx,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
			InstanceID(notification.Aggregate().InstanceID).
			AddQuery().
			AggregateTypes(quota.AggregateType).
			AggregateIDs(notification.Aggregate().ID).
			SequenceGreater(notification.Sequence()).
			EventTypes(quota.NotifiedEventType).
			CreationDateAfter(notification.PeriodStart).
			EventData(map[string]interface{}{
				"id":        notification.ID,
				"threshold": notification.Threshold,
			}).
			Builder(),
	)
	return len(events) > 0, err
}

func notify(ctx context.Context, notification *quota.NotifiedEvent) error {

	payload, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, notification.CallURL, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if err = resp.Body.Close(); err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("calling url %s returned %s", notification.CallURL, resp.Status)
	}

	return nil
}

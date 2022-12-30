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
	"github.com/zitadel/zitadel/internal/repository/instance"

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

func (c *Commands) AddInstanceQuota(
	ctx context.Context,
	quota *Quota,
) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.AddInstanceQuotaCommand(instanceAgg, quota))
	if err != nil {
		return nil, err
	}
	events, err := c.eventstore.Push(ctx, cmds...)
	if err != nil {
		return nil, err
	}
	wm := &eventstore.WriteModel{
		AggregateID:   authz.GetInstance(ctx).InstanceID(),
		ResourceOwner: authz.GetInstance(ctx).InstanceID(),
	}
	err = AppendAndReduce(wm, events...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(wm), nil
}

func (c *Commands) RemoveInstanceQuota(ctx context.Context, unit QuotaUnit) (*domain.ObjectDetails, error) {
	instanceAgg := instance.NewAggregate(authz.GetInstance(ctx).InstanceID())
	if instanceAgg.InstanceID == "" {
		return nil, caos_errs.ThrowInvalidArgument(nil, "COMMAND-Gfg3g", "Errors.IDMissing")
	}

	domainQuotaUnit := unit.Enum()

	if domainQuotaUnit == quota.Unimplemented {
		return nil, errors.ThrowInvalidArgument(nil, "INSTA-SDSfs", "Errors.Invalid.Argument") // TODO: Better error message?
	}

	quota, err := c.getQuotaWriteModel(ctx, domainQuotaUnit)
	if err != nil {
		return nil, err
	}

	if !quota.exists {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-WDfFf", "Errors.Quota.NotFound")
	}

	events := []eventstore.Command{
		instance.NewQuotaRemovedEvent(ctx, &instanceAgg.Aggregate, unit.Enum()),
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(quota, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&quota.WriteModel), nil
}

func (c *Commands) getQuotaWriteModel(ctx context.Context, unit quota.Unit) (*quotaWriteModel, error) {
	wm := newQuotaWriteModel(ctx, unit)
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

type Quota struct {
	Unit          QuotaUnit
	From          time.Time
	Interval      time.Duration
	Amount        uint64
	Limit         bool
	Notifications QuotaNotifications
}

func (q *Quota) isValid() bool {
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

func (c *Commands) AddInstanceQuotaCommand(
	a *instance.Aggregate,
	q *Quota,
) preparation.Validation {
	return func() (preparation.CreateCommands, error) {

		if !q.isValid() {
			return nil, errors.ThrowInvalidArgument(nil, "INSTA-pBfjq", "Errors.Invalid.Argument") // TODO: Better error message?
		}

		return func(ctx context.Context, filter preparation.FilterToQueryReducer) (cmd []eventstore.Command, err error) {
				// TODO: Validations with side effects
				genID := func() string {
					id, genErr := c.idGenerator.Next()
					if genErr != nil {
						err = genErr
					}
					return id
				}

				return []eventstore.Command{instance.NewQuotaAddedEvent(
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
func (c *Commands) ReportUsage(ctx context.Context, dueNotifications []*instance.QuotaNotifiedEvent) error {

	for _, notification := range dueNotifications {
		alreadyNotified, err := isAlreadNotified(ctx, c.eventstore, notification)
		if err != nil {
			return err
		}

		if alreadyNotified {
			// TODO: Debugf
			logging.Infof(
				"quota notification with ID %s and threshold %d was already notified in this period",
				notification.NotifiedEvent.ID,
				notification.NotifiedEvent.Threshold,
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

func isAlreadNotified(ctx context.Context, es *eventstore.Eventstore, notification *instance.QuotaNotifiedEvent) (bool, error) {

	events, err := es.Filter(
		ctx,
		eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
			InstanceID(notification.Aggregate().InstanceID).
			AddQuery().
			AggregateTypes(instance.AggregateType).
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

func notify(ctx context.Context, notification *instance.QuotaNotifiedEvent) error {

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

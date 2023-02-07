package command

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

type QuotaUnit string

const (
	QuotaRequestsAllAuthenticated QuotaUnit = "requests.all.authenticated"
	QuotaActionsAllRunsSeconds    QuotaUnit = "actions.all.runs.seconds"
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

	wm, err := c.getQuotaWriteModel(ctx, instanceId, instanceId, q.Unit.Enum())
	if err != nil {
		return nil, err
	}

	if wm.active {
		return nil, caos_errs.ThrowAlreadyExists(nil, "COMMAND-WDfFf", "Errors.Quota.AlreadyExists")
	}

	aggregateId, err := c.idGenerator.Next()
	if err != nil {
		return nil, err
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

	wm, err := c.getQuotaWriteModel(ctx, instanceId, instanceId, unit.Enum())
	if err != nil {
		return nil, err
	}

	if !wm.active {
		return nil, caos_errs.ThrowNotFound(nil, "COMMAND-WDfFf", "Errors.Quota.NotFound")
	}

	aggregate := quota.NewAggregate(wm.AggregateID, instanceId, instanceId)

	events := []eventstore.Command{
		quota.NewRemovedEvent(ctx, &aggregate.Aggregate, unit.Enum()),
	}
	pushedEvents, err := c.eventstore.Push(ctx, events...)
	if err != nil {
		return nil, err
	}
	err = AppendAndReduce(wm, pushedEvents...)
	if err != nil {
		return nil, err
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

func (c *Commands) getQuotaWriteModel(ctx context.Context, instanceId, resourceOwner string, unit quota.Unit) (*quotaWriteModel, error) {
	wm := newQuotaWriteModel(instanceId, resourceOwner, unit)
	return wm, c.eventstore.FilterToQueryReducer(ctx, wm)
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

func (q *AddQuota) validate() error {

	isUrl := func(str string) error {
		u, err := url.Parse(str)
		if err != nil {
			return err
		}

		if u.Scheme == "" || u.Host == "" {
			return fmt.Errorf("url %s is invalid", str)
		}

		return nil
	}

	for _, notification := range q.Notifications {
		if err := isUrl(notification.CallURL); err != nil {
			return errors.ThrowInvalidArgument(err, "QUOTA-HAYmN", "Errors.Quota.Invalid.CallURL")
		}
		if notification.Percent < 1 {
			return errors.ThrowInvalidArgument(nil, "QUOTA-pBfjq", "Errors.Quota.Invalid.Percent")
		}
	}

	if q.Unit.Enum() == quota.Unimplemented {
		return errors.ThrowInvalidArgument(nil, "QUOTA-OTeSh", "Errors.Quota.Invalid.Unimplemented")
	}

	if q.Amount < 1 {
		return errors.ThrowInvalidArgument(nil, "QUOTA-hOKSJ", "Errors.Quota.Invalid.Amount")
	}

	if q.Interval < time.Minute {
		return errors.ThrowInvalidArgument(nil, "QUOTA-R5otd", "Errors.Quota.Invalid.Interval")
	}

	if !q.Limit && len(q.Notifications) == 0 {
		return errors.ThrowInvalidArgument(nil, "QUOTA-4Nv68", "Errors.Quota.Invalid.Noop")
	}

	return nil
}

func (c *Commands) AddQuotaCommand(a *quota.Aggregate, q *AddQuota) preparation.Validation {
	return func() (preparation.CreateCommands, error) {

		if err := q.validate(); err != nil {
			return nil, err
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

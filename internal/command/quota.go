package command

import (
	"context"
	"net/url"
	"time"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/command/preparation"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/repository/quota"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type QuotaUnit string

const (
	QuotaRequestsAllAuthenticated QuotaUnit = "requests.all.authenticated"
	QuotaActionsAllRunsSeconds    QuotaUnit = "actions.all.runs.seconds"
)

func (q QuotaUnit) Enum() quota.Unit {
	switch q {
	case QuotaRequestsAllAuthenticated:
		return quota.RequestsAllAuthenticated
	case QuotaActionsAllRunsSeconds:
		return quota.ActionsAllRunsSeconds
	default:
		return quota.Unimplemented
	}
}

// AddQuota returns and error if the quota already exists.
// AddQuota is deprecated. Use SetQuota instead.
func (c *Commands) AddQuota(
	ctx context.Context,
	q *SetQuota,
) (*domain.ObjectDetails, error) {
	instanceId := authz.GetInstance(ctx).InstanceID()
	wm, err := c.getQuotaWriteModel(ctx, instanceId, instanceId, q.Unit.Enum())
	if err != nil {
		return nil, err
	}
	if wm.AggregateID != "" {
		return nil, zerrors.ThrowAlreadyExists(nil, "COMMAND-WDfFf", "Errors.Quota.AlreadyExists")
	}
	aggregateId, err := id_generator.Next()
	if err != nil {
		return nil, err
	}
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.SetQuotaCommand(quota.NewAggregate(aggregateId, instanceId), wm, true, q))
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

// SetQuota creates a new quota or updates an existing one.
func (c *Commands) SetQuota(
	ctx context.Context,
	q *SetQuota,
) (*domain.ObjectDetails, error) {
	instanceId := authz.GetInstance(ctx).InstanceID()
	wm, err := c.getQuotaWriteModel(ctx, instanceId, instanceId, q.Unit.Enum())
	if err != nil {
		return nil, err
	}
	aggregateId := wm.AggregateID
	createNewQuota := aggregateId == ""
	if aggregateId == "" {
		aggregateId, err = id_generator.Next()
		if err != nil {
			return nil, err
		}
	}
	cmds, err := preparation.PrepareCommands(ctx, c.eventstore.Filter, c.SetQuotaCommand(quota.NewAggregate(aggregateId, instanceId), wm, createNewQuota, q))
	if err != nil {
		return nil, err
	}
	if len(cmds) > 0 {
		events, err := c.eventstore.Push(ctx, cmds...)
		if err != nil {
			return nil, err
		}
		err = AppendAndReduce(wm, events...)
		if err != nil {
			return nil, err
		}
	}
	return writeModelToObjectDetails(&wm.WriteModel), nil
}

func (c *Commands) RemoveQuota(ctx context.Context, unit QuotaUnit) (*domain.ObjectDetails, error) {
	instanceId := authz.GetInstance(ctx).InstanceID()
	wm, err := c.getQuotaWriteModel(ctx, instanceId, instanceId, unit.Enum())
	if err != nil {
		return nil, err
	}
	if wm.AggregateID == "" {
		return nil, zerrors.ThrowNotFound(nil, "COMMAND-WDfFf", "Errors.Quota.NotFound")
	}
	aggregate := quota.NewAggregate(wm.AggregateID, instanceId)
	events := []eventstore.Command{quota.NewRemovedEvent(ctx, &aggregate.Aggregate, unit.Enum())}
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
	Percent uint16
	Repeat  bool
	CallURL string
}

// SetQuota configures a quota and activates it if it isn't active already
type SetQuota struct {
	Unit          QuotaUnit          `json:"unit"`
	From          time.Time          `json:"from"`
	ResetInterval time.Duration      `json:"ResetInterval,omitempty"`
	Amount        uint64             `json:"Amount,omitempty"`
	Limit         bool               `json:"Limit,omitempty"`
	Notifications QuotaNotifications `json:"Notifications,omitempty"`
}

type QuotaNotifications []*QuotaNotification

func (q *QuotaNotification) validate() error {
	u, err := url.Parse(q.CallURL)
	if err != nil {
		return zerrors.ThrowInvalidArgument(err, "QUOTA-bZ0Fj", "Errors.Quota.Invalid.CallURL")
	}
	if !u.IsAbs() || u.Host == "" {
		return zerrors.ThrowInvalidArgument(nil, "QUOTA-HAYmN", "Errors.Quota.Invalid.CallURL")
	}
	if q.Percent < 1 {
		return zerrors.ThrowInvalidArgument(nil, "QUOTA-pBfjq", "Errors.Quota.Invalid.Percent")
	}
	return nil
}

func (q *SetQuota) validate() error {
	for _, notification := range q.Notifications {
		if err := notification.validate(); err != nil {
			return err
		}
	}
	if q.Unit.Enum() == quota.Unimplemented {
		return zerrors.ThrowInvalidArgument(nil, "QUOTA-OTeSh", "Errors.Quota.Invalid.Unimplemented")
	}
	if q.ResetInterval < time.Minute {
		return zerrors.ThrowInvalidArgument(nil, "QUOTA-R5otd", "Errors.Quota.Invalid.ResetInterval")
	}
	return nil
}

func (c *Commands) SetQuotaCommand(a *quota.Aggregate, wm *quotaWriteModel, createNew bool, q *SetQuota) preparation.Validation {
	return func() (preparation.CreateCommands, error) {
		if err := q.validate(); err != nil {
			return nil, err
		}
		return func(ctx context.Context, filter preparation.FilterToQueryReducer) (cmd []eventstore.Command, err error) {
				changes, err := wm.NewChanges(createNew, q.Amount, q.From, q.ResetInterval, q.Limit, q.Notifications...)
				if len(changes) == 0 {
					return nil, err
				}
				return []eventstore.Command{quota.NewSetEvent(
					eventstore.NewBaseEventForPush(
						ctx,
						&a.Aggregate,
						quota.SetEventType,
					),
					q.Unit.Enum(),
					changes...,
				)}, err
			},
			nil
	}
}

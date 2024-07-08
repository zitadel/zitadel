package command

import (
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/id_generator"
	"github.com/zitadel/zitadel/internal/repository/quota"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type quotaWriteModel struct {
	eventstore.WriteModel
	rollingAggregateID string
	unit               quota.Unit
	from               time.Time
	resetInterval      time.Duration
	amount             uint64
	limit              bool
	notifications      []*quota.SetEventNotification
}

// newQuotaWriteModel aggregateId is filled by reducing unit matching events
func newQuotaWriteModel(instanceId, resourceOwner string, unit quota.Unit) *quotaWriteModel {
	return &quotaWriteModel{
		WriteModel: eventstore.WriteModel{
			InstanceID:    instanceId,
			ResourceOwner: resourceOwner,
		},
		unit: unit,
	}
}

func (wm *quotaWriteModel) Query() *eventstore.SearchQueryBuilder {
	query := eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		InstanceID(wm.InstanceID).
		AddQuery().
		AggregateTypes(quota.AggregateType).
		EventTypes(
			quota.AddedEventType,
			quota.SetEventType,
			quota.RemovedEventType,
		).EventData(map[string]interface{}{"unit": wm.unit})

	return query.Builder()
}

func (wm *quotaWriteModel) Reduce() error {
	for _, event := range wm.Events {
		wm.ChangeDate = event.CreatedAt()
		switch e := event.(type) {
		case *quota.SetEvent:
			wm.rollingAggregateID = e.Aggregate().ID
			if e.Amount != nil {
				wm.amount = *e.Amount
			}
			if e.From != nil {
				wm.from = *e.From
			}
			if e.Limit != nil {
				wm.limit = *e.Limit
			}
			if e.ResetInterval != nil {
				wm.resetInterval = *e.ResetInterval
			}
			if e.Notifications != nil {
				wm.notifications = *e.Notifications
			}
		case *quota.RemovedEvent:
			wm.rollingAggregateID = ""
		}
	}
	if err := wm.WriteModel.Reduce(); err != nil {
		return err
	}
	// wm.WriteModel.Reduce() sets the aggregateID to the first event's aggregateID, but we need the last one
	wm.AggregateID = wm.rollingAggregateID
	return wm.WriteModel.Reduce()
}

// NewChanges returns all changes that need to be applied to the aggregate.
// If createNew is true, all quota properties are set.
func (wm *quotaWriteModel) NewChanges(
	createNew bool,
	amount uint64,
	from time.Time,
	resetInterval time.Duration,
	limit bool,
	notifications ...*QuotaNotification,
) (changes []quota.QuotaChange, err error) {
	setEventNotifications, err := QuotaNotifications(notifications).newSetEventNotifications()
	if err != nil {
		return nil, err
	}
	// we sort the input notifications already, so we can return early if they have duplicates
	err = sortSetEventNotifications(setEventNotifications)
	if err != nil {
		return nil, err
	}
	if createNew {
		return []quota.QuotaChange{
			quota.ChangeAmount(amount),
			quota.ChangeFrom(from),
			quota.ChangeResetInterval(resetInterval),
			quota.ChangeLimit(limit),
			quota.ChangeNotifications(setEventNotifications),
		}, nil
	}
	changes = make([]quota.QuotaChange, 0, 5)
	if wm.amount != amount {
		changes = append(changes, quota.ChangeAmount(amount))
	}
	if wm.from != from {
		changes = append(changes, quota.ChangeFrom(from))
	}
	if wm.resetInterval != resetInterval {
		changes = append(changes, quota.ChangeResetInterval(resetInterval))
	}
	if wm.limit != limit {
		changes = append(changes, quota.ChangeLimit(limit))
	}
	// If the number of notifications differs, we renew the notifications and we can return early
	if len(setEventNotifications) != len(wm.notifications) {
		changes = append(changes, quota.ChangeNotifications(setEventNotifications))
		return changes, nil
	}
	// Now we sort the existing notifications too, so comparing the input properties with the existing ones is easier.
	// We ignore the sorting error for the existing notifications, because this is system state, not user input.
	// If sorting fails this time, the notifications are listed in the event payload and the projection cleans them up anyway.
	_ = sortSetEventNotifications(wm.notifications)
	for i, notification := range setEventNotifications {
		if notification.CallURL != wm.notifications[i].CallURL ||
			notification.Percent != wm.notifications[i].Percent ||
			notification.Repeat != wm.notifications[i].Repeat {
			changes = append(changes, quota.ChangeNotifications(setEventNotifications))
			return changes, nil
		}
	}
	return changes, err
}

// newSetEventNotifications returns quota.SetEventNotification elements with generated IDs.
func (q QuotaNotifications) newSetEventNotifications() (setNotifications []*quota.SetEventNotification, err error) {
	if q == nil {
		return make([]*quota.SetEventNotification, 0), nil
	}
	notifications := make([]*quota.SetEventNotification, len(q))
	for idx, notification := range q {
		notifications[idx] = &quota.SetEventNotification{
			Percent: notification.Percent,
			Repeat:  notification.Repeat,
			CallURL: notification.CallURL,
		}
		notifications[idx].ID, err = id_generator.Next()
		if err != nil {
			return nil, err
		}
	}
	return notifications, nil
}

// sortSetEventNotifications reports an error if there are duplicate notifications or if a pointer is nil
func sortSetEventNotifications(notifications []*quota.SetEventNotification) (err error) {
	slices.SortFunc(notifications, func(i, j *quota.SetEventNotification) int {
		if i == nil || j == nil {
			err = zerrors.ThrowInternal(errors.New("sorting slices of *quota.SetEventNotification with nil pointers is not supported"), "QUOTA-8YXPk", "Errors.Internal")
			return 0
		}
		if i.Percent == j.Percent && i.CallURL == j.CallURL && i.Repeat == j.Repeat {
			// TODO: translate
			err = zerrors.ThrowInternal(fmt.Errorf("%+v", i), "QUOTA-Pty2n", "Errors.Quota.Notifications.Duplicate")
			return 0
		}
		if i.Percent < j.Percent ||
			i.Percent == j.Percent && i.CallURL < j.CallURL ||
			i.Percent == j.Percent && i.CallURL == j.CallURL && !i.Repeat && j.Repeat {
			return -1
		}
		return +1
	})
	return err
}

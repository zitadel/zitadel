package command

import (
	"slices"
	"strings"
	"time"

	"github.com/zitadel/zitadel/internal/errors"

	"github.com/zitadel/zitadel/internal/id"

	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/quota"
)

type quotaWriteModel struct {
	eventstore.WriteModel
	active        bool
	unit          quota.Unit
	from          time.Time
	resetInterval time.Duration
	amount        uint64
	limit         bool
	notifications []*quota.SetEventNotification
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
		AddQuery().
		InstanceID(wm.InstanceID).
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
		wm.ChangeDate = event.CreationDate()
		switch e := event.(type) {
		case *quota.SetEvent:
			wm.active = true
			wm.AggregateID = event.Aggregate().ID
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
				wm.notifications = e.Notifications
			}
		case *quota.RemovedEvent:
			wm.active = false
			wm.AggregateID = ""
		}
	}
	return wm.WriteModel.Reduce()
}

// NewChanges returns all changes that need to be applied to the aggregate.
// If wm is nil, all properties are set.
func (wm *quotaWriteModel) NewChanges(
	idGenerator id.Generator,
	amount uint64,
	from time.Time,
	resetInterval time.Duration,
	limit bool,
	notifications QuotaNotifications,
) (changes []quota.QuotaChange, err error) {
	setEventNotifications, err := notifications.newSetEventNotifications(idGenerator)
	if err != nil {
		return nil, err
	}
	if wm == nil || wm.active == false {
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
	if len(setEventNotifications) != len(wm.notifications) {
		changes = append(changes, quota.ChangeNotifications(setEventNotifications))
		return changes, nil
	}
	replaceIDs(wm.notifications, setEventNotifications)
	// All IDs are passed and the number of notifications didn't change.
	// Now we check if the properties changed.
	err = sortSetEventNotifications(setEventNotifications)
	if err != nil {
		return nil, err
	}
	// We ignore the sorting error for the existing notifications, because this is system state, not user input.
	// If the sorting fails, the notifications will be cleaned up and triggered again.
	_ = sortSetEventNotifications(wm.notifications)
	for i, notification := range setEventNotifications {
		if notification.ID != wm.notifications[i].ID ||
			notification.CallURL != wm.notifications[i].CallURL ||
			notification.Percent != wm.notifications[i].Percent ||
			notification.Repeat != wm.notifications[i].Repeat {
			changes = append(changes, quota.ChangeNotifications(setEventNotifications))
			return changes, nil
		}
	}
	return changes, err
}

// newSetEventNotifications returns quota.SetEventNotification elements with generated IDs.
func (q *QuotaNotifications) newSetEventNotifications(idGenerator id.Generator) (setNotifications []*quota.SetEventNotification, err error) {
	if q == nil {
		return nil, nil
	}
	notifications := make([]*quota.SetEventNotification, len(*q))
	for idx, notification := range *q {
		notifications[idx] = &quota.SetEventNotification{
			Percent: notification.Percent,
			Repeat:  notification.Repeat,
			CallURL: notification.CallURL,
		}
		notifications[idx].ID, err = idGenerator.Next()
		if err != nil {
			return nil, err
		}
	}
	return notifications, nil
}

func replaceIDs(srcs, dsts []*quota.SetEventNotification) {
	for _, dst := range dsts {
		for _, src := range srcs {
			if dst.CallURL == src.CallURL &&
				dst.Percent == src.Percent &&
				dst.Repeat == src.Repeat {
				dst.ID = src.ID
				break
			}
		}
	}
}

func sortSetEventNotifications(notifications []*quota.SetEventNotification) (err error) {
	slices.SortFunc(notifications, func(i, j *quota.SetEventNotification) int {
		comp := strings.Compare(i.ID, j.ID)
		if comp == 0 {
			// TODO: translate
			err = errors.ThrowPreconditionFailed(nil, "EVENT-3M9fs", "Errors.Quota.Notifications.Duplicate")
		}
		return comp
	})
	return err
}

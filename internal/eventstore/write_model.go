package eventstore

import (
	"time"
)

type PermissionCheck func(resourceOwner, aggregateID string) error

// WriteModel is the minimum representation of a command side write model.
// It implements a basic reducer
// it's purpose is to reduce events to create new ones
type WriteModel struct {
	AggregateID       string          `json:"-"`
	ProcessedSequence uint64          `json:"-"`
	Events            []Event         `json:"-"`
	ResourceOwner     string          `json:"-"`
	InstanceID        string          `json:"-"`
	ChangeDate        time.Time       `json:"-"`
	PermissionCheck   PermissionCheck `json:"-"`
}

// AppendEvents adds all the events to the read model.
// The function doesn't compute the new state of the read model
func (rm *WriteModel) AppendEvents(events ...Event) {
	rm.Events = append(rm.Events, events...)
}

// Reduce is the basic implementation of reducer
// If this function is extended the extending function should be the last step
func (wm *WriteModel) Reduce() error {
	// We have to call checkPermission before we overwrite wm.ResourceOwner from the events.
	if err := wm.checkPermission(); err != nil {
		return err
	}
	if len(wm.Events) == 0 {
		return nil
	}
	if wm.AggregateID == "" {
		wm.AggregateID = wm.Events[0].Aggregate().ID
	}

	if wm.ResourceOwner == "" {
		// TODO: use the latest events resource owner for cases where the resource is recreated with the same ID and different resource owner?
		wm.ResourceOwner = wm.Events[0].Aggregate().ResourceOwner
	}

	if wm.InstanceID == "" {
		// TODO: use the latest events instance ID for cases where the resource is recreated with the same ID and different instance?
		wm.InstanceID = wm.Events[0].Aggregate().InstanceID
	}

	wm.ProcessedSequence = wm.Events[len(wm.Events)-1].Sequence()
	wm.ChangeDate = wm.Events[len(wm.Events)-1].CreatedAt()

	// all events processed and not needed anymore
	wm.Events = nil
	wm.Events = []Event{}

	return nil
}

// checkPermission succeeds, if no permission check is set
// If there are no events on the write model, it checks the permission on the given resource owner
// If there are events on the write model, it checks the permission on the resource owner of the last event.
// This makes sure that the correct resource owner is checked in cases where an aggregate is recreated with the same ID and different resource owner.
// checkPermission has to be called before the resource owner is set from the events.
func (wm *WriteModel) checkPermission() error {
	if wm.PermissionCheck == nil {
		return nil
	}
	resourceOwner := wm.ResourceOwner
	if len(wm.Events) > 0 {
		resourceOwner = wm.Events[len(wm.Events)-1].Aggregate().ResourceOwner
	}
	return wm.PermissionCheck(resourceOwner, wm.AggregateID)
}

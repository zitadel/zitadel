package command

import (
	"context"
	"slices"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/target"
)

type TargetWriteModel struct {
	eventstore.WriteModel

	Name             string
	TargetType       domain.TargetType
	Endpoint         string
	Timeout          time.Duration
	InterruptOnError bool
	SigningKey       *crypto.CryptoValue

	State domain.TargetState
}

func NewTargetWriteModel(id string, resourceOwner string) *TargetWriteModel {
	return &TargetWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   id,
			ResourceOwner: resourceOwner,
			InstanceID:    resourceOwner,
		},
	}
}

func (wm *TargetWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *target.AddedEvent:
			wm.Name = e.Name
			wm.TargetType = e.TargetType
			wm.Endpoint = e.Endpoint
			wm.Timeout = e.Timeout
			wm.State = domain.TargetActive
			wm.SigningKey = e.SigningKey
		case *target.ChangedEvent:
			if e.Name != nil {
				wm.Name = *e.Name
			}
			if e.TargetType != nil {
				wm.TargetType = *e.TargetType
			}
			if e.Endpoint != nil {
				wm.Endpoint = *e.Endpoint
			}
			if e.Timeout != nil {
				wm.Timeout = *e.Timeout
			}
			if e.InterruptOnError != nil {
				wm.InterruptOnError = *e.InterruptOnError
			}
			if e.SigningKey != nil {
				wm.SigningKey = e.SigningKey
			}
		case *target.RemovedEvent:
			wm.State = domain.TargetRemoved
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *TargetWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(target.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(target.AddedEventType,
			target.ChangedEventType,
			target.RemovedEventType).
		Builder()
}

func (wm *TargetWriteModel) NewChangedEvent(
	ctx context.Context,
	agg *eventstore.Aggregate,
	name *string,
	targetType *domain.TargetType,
	endpoint *string,
	timeout *time.Duration,
	interruptOnError *bool,
	signingKey *crypto.CryptoValue,
) *target.ChangedEvent {
	changes := make([]target.Changes, 0)
	if name != nil && wm.Name != *name {
		changes = append(changes, target.ChangeName(wm.Name, *name))
	}
	if targetType != nil && wm.TargetType != *targetType {
		changes = append(changes, target.ChangeTargetType(*targetType))
	}
	if endpoint != nil && wm.Endpoint != *endpoint {
		changes = append(changes, target.ChangeEndpoint(*endpoint))
	}
	if timeout != nil && wm.Timeout != *timeout {
		changes = append(changes, target.ChangeTimeout(*timeout))
	}
	if interruptOnError != nil && wm.InterruptOnError != *interruptOnError {
		changes = append(changes, target.ChangeInterruptOnError(*interruptOnError))
	}
	// if signingkey is set, update it as it is encrypted
	if signingKey != nil {
		changes = append(changes, target.ChangeSigningKey(signingKey))
	}
	if len(changes) == 0 {
		return nil
	}
	return target.NewChangedEvent(ctx, agg, changes)
}

type TargetsExistsWriteModel struct {
	eventstore.WriteModel
	ids         []string
	existingIDs []string
}

func (e *TargetsExistsWriteModel) AllExists() bool {
	return len(e.ids) == len(e.existingIDs)
}

func NewTargetsExistsWriteModel(ids []string, resourceOwner string) *TargetsExistsWriteModel {
	return &TargetsExistsWriteModel{
		WriteModel: eventstore.WriteModel{
			ResourceOwner: resourceOwner,
			InstanceID:    resourceOwner,
		},
		ids: ids,
	}
}

func (wm *TargetsExistsWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *target.AddedEvent:
			if !slices.Contains(wm.existingIDs, e.Aggregate().ID) {
				wm.existingIDs = append(wm.existingIDs, e.Aggregate().ID)
			}
		case *target.RemovedEvent:
			i := slices.Index(wm.existingIDs, e.Aggregate().ID)
			if i >= 0 {
				wm.existingIDs = slices.Delete(wm.existingIDs, i, i+1)
			}
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *TargetsExistsWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(target.AggregateType).
		AggregateIDs(wm.ids...).
		EventTypes(target.AddedEventType,
			target.RemovedEventType).
		Builder()
}

func TargetAggregateFromWriteModel(wm *eventstore.WriteModel) *eventstore.Aggregate {
	return &eventstore.Aggregate{
		ID:            wm.AggregateID,
		Type:          target.AggregateType,
		ResourceOwner: wm.ResourceOwner,
		InstanceID:    wm.InstanceID,
		Version:       target.AggregateVersion,
	}
}

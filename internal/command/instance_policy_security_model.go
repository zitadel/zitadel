package command

import (
	"context"
	"reflect"

	"github.com/zitadel/zitadel/v2/internal/api/authz"
	"github.com/zitadel/zitadel/v2/internal/eventstore"
	"github.com/zitadel/zitadel/v2/internal/repository/instance"
)

type InstanceSecurityPolicyWriteModel struct {
	eventstore.WriteModel

	Enabled        bool
	AllowedOrigins []string
}

func NewInstanceSecurityPolicyWriteModel(ctx context.Context) *InstanceSecurityPolicyWriteModel {
	return &InstanceSecurityPolicyWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID:   authz.GetInstance(ctx).InstanceID(),
			ResourceOwner: authz.GetInstance(ctx).InstanceID(),
		},
	}
}

func (wm *InstanceSecurityPolicyWriteModel) Reduce() error {
	for _, event := range wm.Events {
		if e, ok := event.(*instance.SecurityPolicySetEvent); ok {
			if e.Enabled != nil {
				wm.Enabled = *e.Enabled
			}
			if e.AllowedOrigins != nil {
				wm.AllowedOrigins = *e.AllowedOrigins
			}
		}
	}
	return wm.WriteModel.Reduce()
}

func (wm *InstanceSecurityPolicyWriteModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(instance.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(
			instance.SecurityPolicySetEventType).
		Builder()
}

func (wm *InstanceSecurityPolicyWriteModel) NewSetEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	enabled bool,
	allowedOrigins []string,
) (*instance.SecurityPolicySetEvent, error) {
	changes := make([]instance.SecurityPolicyChanges, 0, 2)
	var err error

	if wm.Enabled != enabled {
		changes = append(changes, instance.ChangeSecurityPolicyEnabled(enabled))
	}
	if enabled && !reflect.DeepEqual(wm.AllowedOrigins, allowedOrigins) {
		changes = append(changes, instance.ChangeSecurityPolicyAllowedOrigins(allowedOrigins))
	}
	changeEvent, err := instance.NewSecurityPolicySetEvent(ctx, aggregate, changes)
	if err != nil {
		return nil, err
	}
	return changeEvent, nil
}

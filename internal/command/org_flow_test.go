package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/action"
	"github.com/zitadel/zitadel/internal/repository/org"
)

func TestCommands_ClearFlow(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		flowType      domain.FlowType
		resourceOwner string
	}
	type res struct {
		details *domain.ObjectDetails
		err     func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid flow type, error",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				flowType:      domain.FlowTypeUnspecified,
				resourceOwner: "org1",
			},
			res{
				details: nil,
				err:     errors.IsErrorInvalidArgument,
			},
		},
		{
			"already empty, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args{
				ctx:           context.Background(),
				flowType:      domain.FlowTypeExternalAuthentication,
				resourceOwner: "org1",
			},
			res{
				details: nil,
				err:     errors.IsPreconditionFailed,
			},
		},
		{
			"clear ok",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							org.NewTriggerActionsSetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								domain.FlowTypeExternalAuthentication,
								domain.TriggerTypePostAuthentication,
								[]string{"actionID1"},
							),
						),
					),
					expectPush(
						org.NewFlowClearedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							domain.FlowTypeExternalAuthentication,
						),
					),
				),
			},
			args{
				ctx:           context.Background(),
				flowType:      domain.FlowTypeExternalAuthentication,
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
			}
			details, err := c.ClearFlow(tt.args.ctx, tt.args.flowType, tt.args.resourceOwner)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.details, details)
			}
		})
	}
}

func TestCommands_SetTriggerActions(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		flowType      domain.FlowType
		resourceOwner string
		triggerType   domain.TriggerType
		actionIDs     []string
	}
	type res struct {
		details *domain.ObjectDetails
		err     func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"invalid flow type, error",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				flowType:      domain.FlowTypeUnspecified,
				triggerType:   domain.TriggerTypePostAuthentication,
				actionIDs:     []string{"actionID1"},
				resourceOwner: "org1",
			},
			res{
				details: nil,
				err:     errors.IsErrorInvalidArgument,
			},
		},
		//TODO: combination not possible at the moment, add when more flow types available
		//{
		//	"impossible flow / trigger type, error",
		//	fields{
		//		eventstore: eventstoreExpect(t,),
		//	},
		//	args{
		//		ctx:           context.Background(),
		//		flowType:      domain.FlowTypeUnspecified,
		//		triggerType:   domain.TriggerTypePostAuthentication,
		//		actionIDs:     []string{"actionID1"},
		//		resourceOwner: "org1",
		//	},
		//	res{
		//		details: nil,
		//		err:     errors.IsErrorInvalidArgument,
		//	},
		//},
		{
			"no changes, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							org.NewTriggerActionsSetEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								domain.FlowTypeExternalAuthentication,
								domain.TriggerTypePostAuthentication,
								[]string{"actionID1"},
							),
						),
					),
				),
			},
			args{
				ctx:           context.Background(),
				flowType:      domain.FlowTypeExternalAuthentication,
				triggerType:   domain.TriggerTypePostAuthentication,
				actionIDs:     []string{"actionID1"},
				resourceOwner: "org1",
			},
			res{
				details: nil,
				err:     errors.IsPreconditionFailed,
			},
		},
		{
			"actionID not exists, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectFilter(),
				),
			},
			args{
				ctx:           context.Background(),
				flowType:      domain.FlowTypeExternalAuthentication,
				triggerType:   domain.TriggerTypePostAuthentication,
				actionIDs:     []string{"actionID1"},
				resourceOwner: "org1",
			},
			res{
				details: nil,
				err:     errors.IsPreconditionFailed,
			},
		},
		{
			"set ok",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							action.NewAddedEvent(context.Background(),
								&action.NewAggregate("action1", "org1").Aggregate,
								"actionID1",
								"function(ctx, api) action {};",
								0,
								false,
							),
						),
					),
					expectPush(
						org.NewTriggerActionsSetEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							domain.FlowTypeExternalAuthentication,
							domain.TriggerTypePostAuthentication,
							[]string{"actionID1"},
						),
					),
				),
			},
			args{
				ctx:           context.Background(),
				flowType:      domain.FlowTypeExternalAuthentication,
				triggerType:   domain.TriggerTypePostAuthentication,
				actionIDs:     []string{"actionID1"},
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
			}
			details, err := c.SetTriggerActions(tt.args.ctx, tt.args.flowType, tt.args.triggerType, tt.args.actionIDs, tt.args.resourceOwner)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.details, details)
			}
		})
	}
}

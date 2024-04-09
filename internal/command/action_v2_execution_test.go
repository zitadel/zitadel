package command

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/execution"
	"github.com/zitadel/zitadel/internal/repository/target"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func existsMock(exists bool) func(method string) bool {
	return func(method string) bool {
		return exists
	}
}
func TestCommands_SetExecutionRequest(t *testing.T) {
	type fields struct {
		eventstore        func(t *testing.T) *eventstore.Eventstore
		grpcMethodExists  func(method string) bool
		grpcServiceExists func(method string) bool
	}
	type args struct {
		ctx           context.Context
		cond          *ExecutionAPICondition
		set           *SetExecution
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
			"no resourceowner, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				cond:          &ExecutionAPICondition{},
				set:           &SetExecution{},
				resourceOwner: "",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"no cond, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				cond:          &ExecutionAPICondition{},
				set:           &SetExecution{},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"no valid cond, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"notvalid",
					"notvalid",
					false,
				},
				set:           &SetExecution{},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"empty executionType, error",
			fields{
				eventstore:       expectEventstore(),
				grpcMethodExists: existsMock(true),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"notvalid",
					"",
					false,
				},
				set:           &SetExecution{},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"empty target, error",
			fields{
				eventstore:       expectEventstore(),
				grpcMethodExists: existsMock(true),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"notvalid",
					"",
					false,
				},
				set:           &SetExecution{},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"method not found, error",
			fields{
				eventstore:       expectEventstore(),
				grpcMethodExists: existsMock(false),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"method",
					"",
					false,
				},
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeTarget, "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"service not found, error",
			fields{
				eventstore:        expectEventstore(),
				grpcServiceExists: existsMock(false),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"",
					"service",
					false,
				},
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeTarget, "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push ok, method target",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("target", "instance"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								time.Second,
								true,
							),
						),
					),
					expectPush(
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("request.method", "instance"),
							[]*execution.Target{
								{domain.ExecutionTargetTypeTarget, "target"},
							},
						),
					),
				),
				grpcMethodExists: existsMock(true),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"method",
					"",
					false,
				},
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeTarget, "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
		{
			"push ok, service target",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("target", "instance"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								time.Second,
								true,
							),
						),
					),
					expectPush(
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("request.service", "instance"),
							[]*execution.Target{
								{domain.ExecutionTargetTypeTarget, "target"},
							},
						),
					),
				),
				grpcServiceExists: existsMock(true),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"",
					"service",
					false,
				},
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeTarget, "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
		{
			"push ok, all target",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("target", "instance"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								time.Second,
								true,
							),
						),
					),
					expectPush(
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("request", "instance"),
							[]*execution.Target{
								{domain.ExecutionTargetTypeTarget, "target"},
							},
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"",
					"",
					true,
				},
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeTarget, "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
		{
			"push not found, method include",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				grpcMethodExists: existsMock(true),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"method",
					"",
					false,
				},
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeInclude, "request.include"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push ok, method include",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("request.include", "instance"),
								[]*execution.Target{
									{domain.ExecutionTargetTypeTarget, "target"},
								},
							),
						),
					),
					expectPush(
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("request.method", "instance"),
							[]*execution.Target{
								{domain.ExecutionTargetTypeInclude, "request.include"},
							},
						),
					),
				),
				grpcMethodExists: existsMock(true),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"method",
					"",
					false,
				},
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeInclude, "request.include"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
		{
			"push not found, service include",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				grpcServiceExists: existsMock(true),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"",
					"service",
					false,
				},
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeInclude, "request.include"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push ok, service include",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("request.include", "instance"),
								[]*execution.Target{
									{domain.ExecutionTargetTypeTarget, "target"},
								},
							),
						),
					),
					expectPush(
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("request.service", "instance"),
							[]*execution.Target{
								{domain.ExecutionTargetTypeInclude, "request.include"},
							},
						),
					),
				),
				grpcServiceExists: existsMock(true),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"",
					"service",
					false,
				},
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeInclude, "request.include"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
		{
			"push not found, all include",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"",
					"",
					true,
				},
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeInclude, "request.include"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push ok, all include",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("request.include", "instance"),
								[]*execution.Target{
									{domain.ExecutionTargetTypeTarget, "target"},
								},
							),
						),
					),
					expectPush(
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("request", "instance"),
							[]*execution.Target{
								{domain.ExecutionTargetTypeInclude, "request.include"},
							},
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"",
					"",
					true,
				},
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeInclude, "request.include"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				GrpcMethodExisting:  tt.fields.grpcMethodExists,
				GrpcServiceExisting: tt.fields.grpcServiceExists,
			}
			details, err := c.SetExecutionRequest(tt.args.ctx, tt.args.cond, tt.args.set, tt.args.resourceOwner)
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

func TestCommands_SetExecutionResponse(t *testing.T) {
	type fields struct {
		eventstore        func(t *testing.T) *eventstore.Eventstore
		grpcMethodExists  func(method string) bool
		grpcServiceExists func(method string) bool
	}
	type args struct {
		ctx           context.Context
		cond          *ExecutionAPICondition
		set           *SetExecution
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
			"no resourceowner, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				cond:          &ExecutionAPICondition{},
				set:           &SetExecution{},
				resourceOwner: "",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"no cond, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				cond:          &ExecutionAPICondition{},
				set:           &SetExecution{},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"no valid cond, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"notvalid",
					"notvalid",
					false,
				},
				set:           &SetExecution{},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"empty executionType, error",
			fields{
				eventstore:       expectEventstore(),
				grpcMethodExists: existsMock(true),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"notvalid",
					"",
					false,
				},
				set:           &SetExecution{},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"empty target, error",
			fields{
				eventstore:       expectEventstore(),
				grpcMethodExists: existsMock(true),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"notvalid",
					"",
					false,
				},
				set:           &SetExecution{},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"push failed, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						target.NewAddedEvent(context.Background(),
							target.NewAggregate("target", "instance"),
							"name",
							domain.TargetTypeWebhook,
							"https://example.com",
							time.Second,
							true,
						),
					),
					expectPushFailed(
						zerrors.ThrowPreconditionFailed(nil, "id", "name already exists"),
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("response.valid", "instance"),
							[]*execution.Target{
								{domain.ExecutionTargetTypeTarget, "target"},
							},
						),
					),
				),
				grpcMethodExists: existsMock(true),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"valid",
					"",
					false,
				},
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeTarget, "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			"method not found, error",
			fields{
				eventstore:       expectEventstore(),
				grpcMethodExists: existsMock(false),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"method",
					"",
					false,
				},
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeTarget, "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"service not found, error",
			fields{
				eventstore:        expectEventstore(),
				grpcServiceExists: existsMock(false),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"",
					"service",
					false,
				},
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeTarget, "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push ok, method target",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("target", "instance"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								time.Second,
								true,
							),
						),
					),
					expectPush(
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("response.method", "instance"),
							[]*execution.Target{
								{domain.ExecutionTargetTypeTarget, "target"},
							},
						),
					),
				),
				grpcMethodExists: existsMock(true),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"method",
					"",
					false,
				},
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeTarget, "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
		{
			"push ok, service target",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						targetAddEvent("target", "instance"),
					),
					expectPush(
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("response.service", "instance"),
							[]*execution.Target{
								{domain.ExecutionTargetTypeTarget, "target"},
							},
						),
					),
				),
				grpcServiceExists: existsMock(true),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"",
					"service",
					false,
				},
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeTarget, "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
		{
			"push ok, all target",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						targetAddEvent("target", "instance"),
					),
					expectPush(
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("response", "instance"),
							[]*execution.Target{
								{domain.ExecutionTargetTypeTarget, "target"},
							},
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"",
					"",
					true,
				},
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeTarget, "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore(t),
				GrpcMethodExisting:  tt.fields.grpcMethodExists,
				GrpcServiceExisting: tt.fields.grpcServiceExists,
			}
			details, err := c.SetExecutionResponse(tt.args.ctx, tt.args.cond, tt.args.set, tt.args.resourceOwner)
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

func TestCommands_SetExecutionEvent(t *testing.T) {
	type fields struct {
		eventstore       func(t *testing.T) *eventstore.Eventstore
		eventExists      func(string) bool
		eventGroupExists func(string) bool
	}
	type args struct {
		ctx           context.Context
		cond          *ExecutionEventCondition
		set           *SetExecution
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
			"no resourceowner, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				cond:          &ExecutionEventCondition{},
				set:           &SetExecution{},
				resourceOwner: "",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"no cond, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				cond:          &ExecutionEventCondition{},
				set:           &SetExecution{},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"no valid cond, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionEventCondition{
					"notvalid",
					"notvalid",
					false,
				},
				set:           &SetExecution{},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"empty executionType, error",
			fields{
				eventstore:  expectEventstore(),
				eventExists: existsMock(true),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionEventCondition{
					"notvalid",
					"",
					false,
				},
				set:           &SetExecution{},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"empty target, error",
			fields{
				eventstore:  expectEventstore(),
				eventExists: existsMock(true),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionEventCondition{
					"notvalid",
					"",
					false,
				},
				set:           &SetExecution{},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"push failed, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						targetAddEvent("target", "instance"),
					),
					expectPushFailed(
						zerrors.ThrowPreconditionFailed(nil, "id", "name already exists"),
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("event.valid", "instance"),
							[]*execution.Target{
								{domain.ExecutionTargetTypeTarget, "target"},
							},
						),
					),
				),
				eventExists: existsMock(true),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionEventCondition{
					"valid",
					"",
					false,
				},
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeTarget, "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			"event not found, error",
			fields{
				eventstore:  expectEventstore(),
				eventExists: existsMock(false),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionEventCondition{
					"event",
					"",
					false,
				},
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeTarget, "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"group not found, error",
			fields{
				eventstore:       expectEventstore(),
				eventGroupExists: existsMock(false),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionEventCondition{
					"",
					"group",
					false,
				},
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeTarget, "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push ok, event target",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						targetAddEvent("target", "instance"),
					),
					expectPush(
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("event.event", "instance"),
							[]*execution.Target{
								{domain.ExecutionTargetTypeTarget, "target"},
							},
						),
					),
				),
				eventExists: existsMock(true),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionEventCondition{
					"event",
					"",
					false,
				},
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeTarget, "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
		{
			"push ok, group target",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						targetAddEvent("target", "instance"),
					),
					expectPush(
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("event.group", "instance"),
							[]*execution.Target{
								{domain.ExecutionTargetTypeTarget, "target"},
							},
						),
					),
				),
				eventGroupExists: existsMock(true),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionEventCondition{
					"",
					"group",
					false,
				},
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeTarget, "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
		{
			"push ok, all target",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						targetAddEvent("target", "instance"),
					),
					expectPush(
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("event", "instance"),
							[]*execution.Target{
								{domain.ExecutionTargetTypeTarget, "target"},
							},
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionEventCondition{
					"",
					"",
					true,
				},
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeTarget, "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:         tt.fields.eventstore(t),
				EventExisting:      tt.fields.eventExists,
				EventGroupExisting: tt.fields.eventGroupExists,
			}
			details, err := c.SetExecutionEvent(tt.args.ctx, tt.args.cond, tt.args.set, tt.args.resourceOwner)
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

func TestCommands_SetExecutionFunction(t *testing.T) {
	type fields struct {
		eventstore           func(t *testing.T) *eventstore.Eventstore
		actionFunctionExists func(string) bool
	}
	type args struct {
		ctx           context.Context
		cond          ExecutionFunctionCondition
		set           *SetExecution
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
			"no resourceowner, error",
			fields{
				eventstore:           expectEventstore(),
				actionFunctionExists: existsMock(true),
			},
			args{
				ctx:           context.Background(),
				cond:          "",
				set:           &SetExecution{},
				resourceOwner: "",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"no cond, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				cond:          "",
				set:           &SetExecution{},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"empty executionType, error",
			fields{
				eventstore:           expectEventstore(),
				actionFunctionExists: existsMock(true),
			},
			args{
				ctx:           context.Background(),
				cond:          "function",
				set:           &SetExecution{},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"empty target, error",
			fields{
				eventstore:           expectEventstore(),
				actionFunctionExists: existsMock(true),
			},
			args{
				ctx:           context.Background(),
				cond:          "function",
				set:           &SetExecution{},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"push failed, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						targetAddEvent("target", "instance"),
					),
					expectPushFailed(
						zerrors.ThrowPreconditionFailed(nil, "id", "name already exists"),
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("function.function", "instance"),
							[]*execution.Target{
								{domain.ExecutionTargetTypeTarget, "target"},
							},
						),
					),
				),
				actionFunctionExists: existsMock(true),
			},
			args{
				ctx:  context.Background(),
				cond: "function",
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeTarget, "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		}, {
			"push error, function target",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
				actionFunctionExists: existsMock(true),
			},
			args{
				ctx:  context.Background(),
				cond: "function",
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeTarget, "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push error, function not existing",
			fields{
				eventstore:           expectEventstore(),
				actionFunctionExists: existsMock(false),
			},
			args{
				ctx:  context.Background(),
				cond: "function",
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeTarget, "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push ok, function target",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						targetAddEvent("target", "instance"),
					),
					expectPush(
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("function.function", "instance"),
							[]*execution.Target{
								{domain.ExecutionTargetTypeTarget, "target"},
							},
						),
					),
				),
				actionFunctionExists: existsMock(true),
			},
			args{
				ctx:  context.Background(),
				cond: "function",
				set: &SetExecution{
					Targets: []*execution.Target{
						{domain.ExecutionTargetTypeTarget, "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:             tt.fields.eventstore(t),
				ActionFunctionExisting: tt.fields.actionFunctionExists,
			}
			details, err := c.SetExecutionFunction(tt.args.ctx, tt.args.cond, tt.args.set, tt.args.resourceOwner)
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

func TestCommands_DeleteExecutionRequest(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		cond          *ExecutionAPICondition
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
			"no resourceowner, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				cond:          &ExecutionAPICondition{},
				resourceOwner: "",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"no cond, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				cond:          &ExecutionAPICondition{},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"no valid cond, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"notvalid",
					"notvalid",
					false,
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"push failed, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("request.valid", "instance"),
								[]*execution.Target{
									{domain.ExecutionTargetTypeTarget, "target"},
								},
							),
						),
					),
					expectPushFailed(
						zerrors.ThrowPreconditionFailed(nil, "id", "name already exists"),
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("request.valid", "instance"),
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"valid",
					"",
					false,
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			"not found, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"method",
					"",
					false,
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push ok, method target",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("request.method", "instance"),
								[]*execution.Target{
									{domain.ExecutionTargetTypeTarget, "target"},
								},
							),
						),
					),
					expectPush(
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("request.method", "instance"),
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"method",
					"",
					false,
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
		{
			"push ok, service target",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("request.service", "instance"),
								[]*execution.Target{
									{domain.ExecutionTargetTypeTarget, "target"},
								},
							),
						),
					),
					expectPush(
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("request.service", "instance"),
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"",
					"service",
					false,
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
		{
			"push ok, all target",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("request", "instance"),
								[]*execution.Target{
									{domain.ExecutionTargetTypeTarget, "target"},
								},
							),
						),
					),
					expectPush(
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("request", "instance"),
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"",
					"",
					true,
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			details, err := c.DeleteExecutionRequest(tt.args.ctx, tt.args.cond, tt.args.resourceOwner)
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

func TestCommands_DeleteExecutionResponse(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		cond          *ExecutionAPICondition
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
			"no resourceowner, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				cond:          &ExecutionAPICondition{},
				resourceOwner: "",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"no cond, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				cond:          &ExecutionAPICondition{},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"no valid cond, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"notvalid",
					"notvalid",
					false,
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"push failed, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("response.valid", "instance"),
								[]*execution.Target{
									{domain.ExecutionTargetTypeTarget, "target"},
								},
							),
						),
					),
					expectPushFailed(
						zerrors.ThrowPreconditionFailed(nil, "id", "name already exists"),
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("response.valid", "instance"),
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"valid",
					"",
					false,
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			"not found, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"method",
					"",
					false,
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push ok, method target",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("response.method", "instance"),
								[]*execution.Target{
									{domain.ExecutionTargetTypeTarget, "target"},
								},
							),
						),
					),
					expectPush(
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("response.method", "instance"),
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"method",
					"",
					false,
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
		{
			"push ok, service target",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("response.service", "instance"),
								[]*execution.Target{
									{domain.ExecutionTargetTypeTarget, "target"},
								},
							),
						),
					),
					expectPush(
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("response.service", "instance"),
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"",
					"service",
					false,
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
		{
			"push ok, all target",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("response", "instance"),
								[]*execution.Target{
									{domain.ExecutionTargetTypeTarget, "target"},
								},
							),
						),
					),
					expectPush(
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("response", "instance"),
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"",
					"",
					true,
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			details, err := c.DeleteExecutionResponse(tt.args.ctx, tt.args.cond, tt.args.resourceOwner)
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

func TestCommands_DeleteExecutionEvent(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		cond          *ExecutionEventCondition
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
			"no resourceowner, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				cond:          &ExecutionEventCondition{},
				resourceOwner: "",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"no cond, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				cond:          &ExecutionEventCondition{},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"push failed, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("event.valid", "instance"),
								[]*execution.Target{
									{domain.ExecutionTargetTypeTarget, "target"},
								},
							),
						),
					),
					expectPushFailed(
						zerrors.ThrowPreconditionFailed(nil, "id", "name already exists"),
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("event.valid", "instance"),
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionEventCondition{
					"valid",
					"",
					false,
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			"push error, not existing",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionEventCondition{
					"valid",
					"",
					false,
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push error, event",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionEventCondition{
					"valid",
					"",
					false,
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push ok, event",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("event.valid", "instance"),
								[]*execution.Target{
									{domain.ExecutionTargetTypeTarget, "target"},
								},
							),
						),
					),
					expectPush(
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("event.valid", "instance"),
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionEventCondition{
					"valid",
					"",
					false,
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
		{
			"push error, group",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionEventCondition{
					"",
					"valid",
					false,
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push ok, group",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("event.group", "instance"),
								[]*execution.Target{
									{domain.ExecutionTargetTypeTarget, "target"},
								},
							),
						),
					),
					expectPush(
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("event.group", "instance"),
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionEventCondition{
					"",
					"group",
					false,
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
		{
			"push error, all",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionEventCondition{
					"",
					"",
					true,
				},
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push ok, all",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("event", "instance"),
								[]*execution.Target{
									{domain.ExecutionTargetTypeTarget, "target"},
								},
							),
						),
					),
					expectPush(
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("event", "instance"),
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionEventCondition{
					"",
					"",
					true,
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			details, err := c.DeleteExecutionEvent(tt.args.ctx, tt.args.cond, tt.args.resourceOwner)
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

func TestCommands_DeleteExecutionFunction(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		cond          ExecutionFunctionCondition
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
			"no resourceowner, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				cond:          "",
				resourceOwner: "",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"no cond, error",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				cond:          "",
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"push failed, error",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("function.function", "instance"),
								[]*execution.Target{
									{domain.ExecutionTargetTypeTarget, "target"},
								},
							),
						),
					),
					expectPushFailed(
						zerrors.ThrowPreconditionFailed(nil, "id", "name already exists"),
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("function.function", "instance"),
						),
					),
				),
			},
			args{
				ctx:           context.Background(),
				cond:          "function",
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			"push error, not existing",
			fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args{
				ctx:           context.Background(),
				cond:          "function",
				resourceOwner: "instance",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push ok, function",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("function.function", "instance"),
								[]*execution.Target{
									{domain.ExecutionTargetTypeTarget, "target"},
								},
							),
						),
					),
					expectPush(
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("function.function", "instance"),
						),
					),
				),
			},
			args{
				ctx:           context.Background(),
				cond:          "function",
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			details, err := c.DeleteExecutionFunction(tt.args.ctx, tt.args.cond, tt.args.resourceOwner)
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

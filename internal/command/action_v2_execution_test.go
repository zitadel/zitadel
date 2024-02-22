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
		eventstore        *eventstore.Eventstore
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
				eventstore: eventstoreExpect(t),
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
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				cond:          &ExecutionAPICondition{},
				set:           &SetExecution{},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"no valid cond, error",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"notvalid",
					"notvalid",
					false,
				},
				set:           &SetExecution{},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"empty executionType, error",
			fields{
				eventstore:       eventstoreExpect(t),
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
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"empty target, error",
			fields{
				eventstore:       eventstoreExpect(t),
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
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"target and include, error",
			fields{
				eventstore:       eventstoreExpect(t),
				grpcMethodExists: existsMock(true),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"notvalid",
					"",
					false,
				},
				set: &SetExecution{
					Targets:  []string{"invalid"},
					Includes: []string{"invalid"},
				},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"push failed, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("target", "org1"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								time.Second,
								true,
								true,
							),
						),
					),
					expectPushFailed(
						zerrors.ThrowPreconditionFailed(nil, "id", "name already exists"),
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("grpc.valid", "org1"),
							domain.ExecutionTypeRequest,
							[]string{"target"},
							nil,
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
					Targets: []string{"target"},
				},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			"method not found, error",
			fields{
				eventstore:       eventstoreExpect(t),
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
					Targets: []string{"target"},
				},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"service not found, error",
			fields{
				eventstore:        eventstoreExpect(t),
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
					Targets: []string{"target"},
				},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push ok, method target",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("target", "org1"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								time.Second,
								true,
								true,
							),
						),
					),
					expectPush(
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("grpc.method", "org1"),
							domain.ExecutionTypeRequest,
							[]string{"target"},
							nil,
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
					Targets: []string{"target"},
				},
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"push ok, service target",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("target", "org1"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								time.Second,
								true,
								true,
							),
						),
					),
					expectPush(
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("grpc.service", "org1"),
							domain.ExecutionTypeRequest,
							[]string{"target"},
							nil,
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
					Targets: []string{"target"},
				},
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"push ok, all target",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("target", "org1"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								time.Second,
								true,
								true,
							),
						),
					),
					expectPush(
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("grpc", "org1"),
							domain.ExecutionTypeRequest,
							[]string{"target"},
							nil,
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
					Targets: []string{"target"},
				},
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		/*
			// TODO: include checks
			{
				"push ok, method include",
				fields{
					eventstore: eventstoreExpect(t,
						expectFilter(),
						expectPush(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("grpc.method", "org1"),
								domain.ExecutionTypeRequest,
								nil,
								[]string{"include"},
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
						Includes: []string{"include"},
					},
					resourceOwner: "org1",
				},
				res{
					details: &domain.ObjectDetails{
						ResourceOwner: "org1",
					},
				},
			},
			{
				"push ok, service include",
				fields{
					eventstore: eventstoreExpect(t,
						expectFilter(),
						expectPush(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("grpc.service", "org1"),
								domain.ExecutionTypeRequest,
								nil,
								[]string{"include"},
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
						Includes: []string{"include"},
					},
					resourceOwner: "org1",
				},
				res{
					details: &domain.ObjectDetails{
						ResourceOwner: "org1",
					},
				},
			},
			{
				"push ok, all include",
				fields{
					eventstore: eventstoreExpect(t,
						expectFilter(),
						expectPush(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("grpc", "org1"),
								domain.ExecutionTypeRequest,
								nil,
								[]string{"include"},
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
						Includes: []string{"include"},
					},
					resourceOwner: "org1",
				},
				res{
					details: &domain.ObjectDetails{
						ResourceOwner: "org1",
					},
				},
			},
			{
				"push ok, all include reset",
				fields{
					eventstore: eventstoreExpect(t,
						expectFilter(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("grpc", "org1"),
								domain.ExecutionTypeRequest,
								[]string{"target"},
								nil,
							),
						),
						expectPush(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("grpc", "org1"),
								domain.ExecutionTypeRequest,
								nil,
								[]string{"include"},
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
						Includes: []string{"include"},
					},
					resourceOwner: "org1",
				},
				res{
					details: &domain.ObjectDetails{
						ResourceOwner: "org1",
					},
				},
			},
			{
				"push ok, all include remove set",
				fields{
					eventstore: eventstoreExpect(t,
						expectFilter(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("grpc", "org1"),
								domain.ExecutionTypeRequest,
								[]string{"target"},
								nil,
							),
							execution.NewRemovedEvent(context.Background(),
								execution.NewAggregate("grpc", "org1"),
							),
						),
						expectPush(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("grpc", "org1"),
								domain.ExecutionTypeRequest,
								nil,
								[]string{"include"},
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
						Includes: []string{"include"},
					},
					resourceOwner: "org1",
				},
				res{
					details: &domain.ObjectDetails{
						ResourceOwner: "org1",
					},
				},
			},
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
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
		eventstore        *eventstore.Eventstore
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
				eventstore: eventstoreExpect(t),
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
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				cond:          &ExecutionAPICondition{},
				set:           &SetExecution{},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"no valid cond, error",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"notvalid",
					"notvalid",
					false,
				},
				set:           &SetExecution{},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"empty executionType, error",
			fields{
				eventstore:       eventstoreExpect(t),
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
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"empty target, error",
			fields{
				eventstore:       eventstoreExpect(t),
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
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"target and include, error",
			fields{
				eventstore:       eventstoreExpect(t),
				grpcMethodExists: existsMock(true),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"notvalid",
					"",
					false,
				},
				set: &SetExecution{
					Targets:  []string{"invalid"},
					Includes: []string{"invalid"},
				},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"push failed, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectFilter(
						target.NewAddedEvent(context.Background(),
							target.NewAggregate("target", "org1"),
							"name",
							domain.TargetTypeWebhook,
							"https://example.com",
							time.Second,
							true,
							true,
						),
					),
					expectPushFailed(
						zerrors.ThrowPreconditionFailed(nil, "id", "name already exists"),
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("grpc.valid", "org1"),
							domain.ExecutionTypeResponse,
							[]string{"target"},
							nil,
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
					Targets: []string{"target"},
				},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			"method not found, error",
			fields{
				eventstore:       eventstoreExpect(t),
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
					Targets: []string{"target"},
				},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"service not found, error",
			fields{
				eventstore:        eventstoreExpect(t),
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
					Targets: []string{"target"},
				},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push ok, method target",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("target", "org1"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								time.Second,
								true,
								true,
							),
						),
					),
					expectPush(
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("grpc.method", "org1"),
							domain.ExecutionTypeResponse,
							[]string{"target"},
							nil,
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
					Targets: []string{"target"},
				},
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"push ok, service target",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("target", "org1"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								time.Second,
								true,
								true,
							),
						),
					),
					expectPush(
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("grpc.service", "org1"),
							domain.ExecutionTypeResponse,
							[]string{"target"},
							nil,
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
					Targets: []string{"target"},
				},
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"push ok, all target",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("target", "org1"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								time.Second,
								true,
								true,
							),
						),
					),
					expectPush(
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("grpc", "org1"),
							domain.ExecutionTypeResponse,
							[]string{"target"},
							nil,
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
					Targets: []string{"target"},
				},
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		/*
			{
				"push ok, method include",
				fields{
					eventstore: eventstoreExpect(t,
						expectFilter(),
						expectFilter(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("target", "org1"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								time.Second,
								true,
								true,
							),
						),
						expectPush(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("grpc.method", "org1"),
								domain.ExecutionTypeResponse,
								nil,
								[]string{"include"},
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
						Includes: []string{"include"},
					},
					resourceOwner: "org1",
				},
				res{
					details: &domain.ObjectDetails{
						ResourceOwner: "org1",
					},
				},
			},
			{
				"push ok, service include",
				fields{
					eventstore: eventstoreExpect(t,
						expectFilter(),

						expectFilter(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("target", "org1"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								time.Second,
								true,
								true,
							),
						),
						expectPush(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("grpc.service", "org1"),
								domain.ExecutionTypeResponse,
								nil,
								[]string{"include"},
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
						Includes: []string{"include"},
					},
					resourceOwner: "org1",
				},
				res{
					details: &domain.ObjectDetails{
						ResourceOwner: "org1",
					},
				},
			},
			{
				"push ok, all include",
				fields{
					eventstore: eventstoreExpect(t,
						expectFilter(),
						expectFilter(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("target", "org1"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								time.Second,
								true,
								true,
							),
						),
						expectPush(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("grpc", "org1"),
								domain.ExecutionTypeResponse,
								nil,
								[]string{"include"},
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
						Includes: []string{"include"},
					},
					resourceOwner: "org1",
				},
				res{
					details: &domain.ObjectDetails{
						ResourceOwner: "org1",
					},
				},
			},
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
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
		eventstore       *eventstore.Eventstore
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
				eventstore: eventstoreExpect(t),
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
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				cond:          &ExecutionEventCondition{},
				set:           &SetExecution{},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"no valid cond, error",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionEventCondition{
					"notvalid",
					"notvalid",
					false,
				},
				set:           &SetExecution{},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"empty executionType, error",
			fields{
				eventstore:  eventstoreExpect(t),
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
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"empty target, error",
			fields{
				eventstore:  eventstoreExpect(t),
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
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"target and include, error",
			fields{
				eventstore:  eventstoreExpect(t),
				eventExists: existsMock(true),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionEventCondition{
					"notvalid",
					"",
					false,
				},
				set: &SetExecution{
					Targets:  []string{"invalid"},
					Includes: []string{"invalid"},
				},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"push failed, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("target", "org1"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								time.Second,
								true,
								true,
							),
						),
					),
					expectPushFailed(
						zerrors.ThrowPreconditionFailed(nil, "id", "name already exists"),
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("event.valid", "org1"),
							domain.ExecutionTypeEvent,
							[]string{"target"},
							nil,
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
					Targets: []string{"target"},
				},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			"event not found, error",
			fields{
				eventstore:  eventstoreExpect(t),
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
					Targets: []string{"target"},
				},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"group not found, error",
			fields{
				eventstore:       eventstoreExpect(t),
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
					Targets: []string{"target"},
				},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push ok, event target",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("target", "org1"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								time.Second,
								true,
								true,
							),
						),
					),
					expectPush(
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("event.event", "org1"),
							domain.ExecutionTypeEvent,
							[]string{"target"},
							nil,
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
					Targets: []string{"target"},
				},
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"push ok, group target",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("target", "org1"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								time.Second,
								true,
								true,
							),
						),
					),
					expectPush(
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("event.group", "org1"),
							domain.ExecutionTypeEvent,
							[]string{"target"},
							nil,
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
					Targets: []string{"target"},
				},
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"push ok, all target",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("target", "org1"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								time.Second,
								true,
								true,
							),
						),
					),
					expectPush(
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("event", "org1"),
							domain.ExecutionTypeEvent,
							[]string{"target"},
							nil,
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
					Targets: []string{"target"},
				},
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		/*
			{
				"push ok, event include",
				fields{
					eventstore: eventstoreExpect(t,
						expectFilter(),
						expectPush(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("event.event", "org1"),
								domain.ExecutionTypeEvent,
								nil,
								[]string{"include"},
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
						Includes: []string{"include"},
					},
					resourceOwner: "org1",
				},
				res{
					details: &domain.ObjectDetails{
						ResourceOwner: "org1",
					},
				},
			},
			{
				"push ok, group include",
				fields{
					eventstore: eventstoreExpect(t,
						expectFilter(),
						expectPush(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("event.group", "org1"),
								domain.ExecutionTypeEvent,
								nil,
								[]string{"include"},
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
						Includes: []string{"include"},
					},
					resourceOwner: "org1",
				},
				res{
					details: &domain.ObjectDetails{
						ResourceOwner: "org1",
					},
				},
			},
			{
				"push ok, all include",
				fields{
					eventstore: eventstoreExpect(t,
						expectFilter(),
						expectPush(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("event", "org1"),
								domain.ExecutionTypeEvent,
								nil,
								[]string{"include"},
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
						Includes: []string{"include"},
					},
					resourceOwner: "org1",
				},
				res{
					details: &domain.ObjectDetails{
						ResourceOwner: "org1",
					},
				},
			},
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:         tt.fields.eventstore,
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
		eventstore           *eventstore.Eventstore
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
				eventstore:           eventstoreExpect(t),
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
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				cond:          "",
				set:           &SetExecution{},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"empty executionType, error",
			fields{
				eventstore:           eventstoreExpect(t),
				actionFunctionExists: existsMock(true),
			},
			args{
				ctx:           context.Background(),
				cond:          "function",
				set:           &SetExecution{},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"empty target, error",
			fields{
				eventstore:           eventstoreExpect(t),
				actionFunctionExists: existsMock(true),
			},
			args{
				ctx:           context.Background(),
				cond:          "function",
				set:           &SetExecution{},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"target and include, error",
			fields{
				eventstore:           eventstoreExpect(t),
				actionFunctionExists: existsMock(true),
			},
			args{
				ctx:  context.Background(),
				cond: "function",
				set: &SetExecution{
					Targets:  []string{"invalid"},
					Includes: []string{"invalid"},
				},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"push failed, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("target", "org1"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								time.Second,
								true,
								true,
							),
						),
					),
					expectPushFailed(
						zerrors.ThrowPreconditionFailed(nil, "id", "name already exists"),
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("func.function", "org1"),
							domain.ExecutionTypeFunction,
							[]string{"target"},
							nil,
						),
					),
				),
				actionFunctionExists: existsMock(true),
			},
			args{
				ctx:  context.Background(),
				cond: "function",
				set: &SetExecution{
					Targets: []string{"target"},
				},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		}, {
			"push error, function target",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectFilter(),
				),
				actionFunctionExists: existsMock(true),
			},
			args{
				ctx:  context.Background(),
				cond: "function",
				set: &SetExecution{
					Targets: []string{"target"},
				},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push error, function not existing",
			fields{
				eventstore:           eventstoreExpect(t),
				actionFunctionExists: existsMock(false),
			},
			args{
				ctx:  context.Background(),
				cond: "function",
				set: &SetExecution{
					Targets: []string{"target"},
				},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push ok, function target",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("target", "org1"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								time.Second,
								true,
								true,
							),
						),
					),
					expectPush(
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("func.function", "org1"),
							domain.ExecutionTypeFunction,
							[]string{"target"},
							nil,
						),
					),
				),
				actionFunctionExists: existsMock(true),
			},
			args{
				ctx:  context.Background(),
				cond: "function",
				set: &SetExecution{
					Targets: []string{"target"},
				},
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		/*
			{
					"push error, function include",
					fields{
						eventstore: eventstoreExpect(t,
							expectFilter(),
							expectFilter(),
							expectPush(
								execution.NewSetEvent(context.Background(),
									execution.NewAggregate("func.function1", "org1"),
									domain.ExecutionTypeFunction,
									nil,
									[]string{"include"},
								),
							),
						),
					},
					args{
						ctx:  context.Background(),
						cond: "function1",
						set: &SetExecution{
							Includes: []string{"function2"},
						},
						resourceOwner: "org1",
					},
					res{
						err: zerrors.IsNotFound,
					},
				},
			{
				"push ok, function include",
				fields{
					eventstore: eventstoreExpect(t,
						expectFilter(),
						expectFilter(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("func.function2", "org1"),
								domain.ExecutionTypeFunction,
								nil,
								[]string{"include"},
							),
						),
						expectPush(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("func.function1", "org1"),
								domain.ExecutionTypeFunction,
								nil,
								[]string{"include"},
							),
						),
					),
				},
				args{
					ctx:  context.Background(),
					cond: "function1",
					set: &SetExecution{
						Includes: []string{"function2"},
					},
					resourceOwner: "org1",
				},
				res{
					details: &domain.ObjectDetails{
						ResourceOwner: "org1",
					},
				},
			},
		*/
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:             tt.fields.eventstore,
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
		eventstore *eventstore.Eventstore
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
				eventstore: eventstoreExpect(t),
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
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				cond:          &ExecutionAPICondition{},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"no valid cond, error",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"notvalid",
					"notvalid",
					false,
				},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"push failed, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("grpc.valid", "org1"),
								domain.ExecutionTypeRequest,
								[]string{"target"},
								nil,
							),
						),
					),
					expectPushFailed(
						zerrors.ThrowPreconditionFailed(nil, "id", "name already exists"),
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("grpc.valid", "org1"),
							domain.ExecutionTypeRequest,
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
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			"not found, error",
			fields{
				eventstore: eventstoreExpect(t,
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
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push ok, method target",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("grpc.method", "org1"),
								domain.ExecutionTypeRequest,
								[]string{"target"},
								nil,
							),
						),
					),
					expectPush(
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("grpc.method", "org1"),
							domain.ExecutionTypeRequest,
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
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"push ok, service target",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("grpc.service", "org1"),
								domain.ExecutionTypeRequest,
								[]string{"target"},
								nil,
							),
						),
					),
					expectPush(
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("grpc.service", "org1"),
							domain.ExecutionTypeRequest,
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
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"push ok, all target",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("grpc", "org1"),
								domain.ExecutionTypeRequest,
								[]string{"target"},
								nil,
							),
						),
					),
					expectPush(
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("grpc", "org1"),
							domain.ExecutionTypeRequest,
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
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
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
		eventstore *eventstore.Eventstore
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
				eventstore: eventstoreExpect(t),
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
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				cond:          &ExecutionAPICondition{},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"no valid cond, error",
			fields{
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"notvalid",
					"notvalid",
					false,
				},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"push failed, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("grpc.valid", "org1"),
								domain.ExecutionTypeResponse,
								[]string{"target"},
								nil,
							),
						),
					),
					expectPushFailed(
						zerrors.ThrowPreconditionFailed(nil, "id", "name already exists"),
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("grpc.valid", "org1"),
							domain.ExecutionTypeResponse,
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
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			"not found, error",
			fields{
				eventstore: eventstoreExpect(t,
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
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push ok, method target",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("grpc.method", "org1"),
								domain.ExecutionTypeResponse,
								[]string{"target"},
								nil,
							),
						),
					),
					expectPush(
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("grpc.method", "org1"),
							domain.ExecutionTypeResponse,
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
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"push ok, service target",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("grpc.service", "org1"),
								domain.ExecutionTypeResponse,
								[]string{"target"},
								nil,
							),
						),
					),
					expectPush(
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("grpc.service", "org1"),
							domain.ExecutionTypeResponse,
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
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"push ok, all target",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("grpc", "org1"),
								domain.ExecutionTypeResponse,
								[]string{"target"},
								nil,
							),
						),
					),
					expectPush(
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("grpc", "org1"),
							domain.ExecutionTypeResponse,
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
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
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
		eventstore *eventstore.Eventstore
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
				eventstore: eventstoreExpect(t),
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
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				cond:          &ExecutionEventCondition{},
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"push failed, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("event.valid", "org1"),
								domain.ExecutionTypeEvent,
								[]string{"target"},
								nil,
							),
						),
					),
					expectPushFailed(
						zerrors.ThrowPreconditionFailed(nil, "id", "name already exists"),
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("event.valid", "org1"),
							domain.ExecutionTypeEvent,
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
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			"push error, not existing",
			fields{
				eventstore: eventstoreExpect(t,
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
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push error, event",
			fields{
				eventstore: eventstoreExpect(t,
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
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push ok, event",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("event.valid", "org1"),
								domain.ExecutionTypeEvent,
								[]string{"target"},
								nil,
							),
						),
					),
					expectPush(
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("event.valid", "org1"),
							domain.ExecutionTypeEvent,
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
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"push error, group",
			fields{
				eventstore: eventstoreExpect(t,
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
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push ok, group",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("event.group", "org1"),
								domain.ExecutionTypeEvent,
								[]string{"target"},
								nil,
							),
						),
					),
					expectPush(
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("event.group", "org1"),
							domain.ExecutionTypeEvent,
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
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
		{
			"push error, all",
			fields{
				eventstore: eventstoreExpect(t,
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
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push ok, all",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("event", "org1"),
								domain.ExecutionTypeEvent,
								[]string{"target"},
								nil,
							),
						),
					),
					expectPush(
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("event", "org1"),
							domain.ExecutionTypeEvent,
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
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
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
		eventstore *eventstore.Eventstore
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
				eventstore: eventstoreExpect(t),
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
				eventstore: eventstoreExpect(t),
			},
			args{
				ctx:           context.Background(),
				cond:          "",
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			"push failed, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("func.function", "org1"),
								domain.ExecutionTypeFunction,
								[]string{"target"},
								nil,
							),
						),
					),
					expectPushFailed(
						zerrors.ThrowPreconditionFailed(nil, "id", "name already exists"),
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("func.function", "org1"),
							domain.ExecutionTypeFunction,
						),
					),
				),
			},
			args{
				ctx:           context.Background(),
				cond:          "function",
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			"push error, not existing",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
				),
			},
			args{
				ctx:           context.Background(),
				cond:          "function",
				resourceOwner: "org1",
			},
			res{
				err: zerrors.IsNotFound,
			},
		},
		{
			"push ok, function",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEvent(context.Background(),
								execution.NewAggregate("func.function", "org1"),
								domain.ExecutionTypeFunction,
								[]string{"target"},
								nil,
							),
						),
					),
					expectPush(
						execution.NewRemovedEvent(context.Background(),
							execution.NewAggregate("func.function", "org1"),
							domain.ExecutionTypeFunction,
						),
					),
				),
			},
			args{
				ctx:           context.Background(),
				cond:          "function",
				resourceOwner: "org1",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore,
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

package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/execution"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func grpcMethodExistsMock(exists bool) func(method string) bool {
	return func(method string) bool {
		return exists
	}
}

func grpcServiceExistsMock(exists bool) func(method string) bool {
	return func(service string) bool {
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
				grpcMethodExists: grpcMethodExistsMock(true),
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
				grpcMethodExists: grpcMethodExistsMock(true),
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
				grpcMethodExists: grpcMethodExistsMock(true),
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
			"unique constraint failed, error",
			fields{
				eventstore: eventstoreExpect(t,
					expectFilter(),
					expectPushFailed(
						zerrors.ThrowPreconditionFailed(nil, "id", "name already exists"),
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("grpc.valid", "org1"),
							domain.ExecutionTypeRequest,
							[]string{"valid"},
							nil,
						),
					),
				),
				grpcMethodExists: grpcMethodExistsMock(true),
			},
			args{
				ctx: context.Background(),
				cond: &ExecutionAPICondition{
					"valid",
					"",
					false,
				},
				set: &SetExecution{
					Targets: []string{"valid"},
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
				grpcMethodExists: grpcMethodExistsMock(false),
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
				grpcServiceExists: grpcServiceExistsMock(false),
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
					expectPush(
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("grpc.method", "org1"),
							domain.ExecutionTypeRequest,
							[]string{"target"},
							nil,
						),
					),
				),
				grpcMethodExists: grpcMethodExistsMock(true),
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
					expectPush(
						execution.NewSetEvent(context.Background(),
							execution.NewAggregate("grpc.service", "org1"),
							domain.ExecutionTypeRequest,
							[]string{"target"},
							nil,
						),
					),
				),
				grpcServiceExists: grpcServiceExistsMock(true),
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
				grpcMethodExists: grpcMethodExistsMock(true),
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
				grpcServiceExists: grpcServiceExistsMock(true),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore:          tt.fields.eventstore,
				grpcMethodExisting:  tt.fields.grpcMethodExists,
				grpcServiceExisting: tt.fields.grpcServiceExists,
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

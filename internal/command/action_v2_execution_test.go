package command

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/crypto"
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
				set:           &SetExecution{Targets: []*execution.Target{{}}},
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
						{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
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
						{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
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
					expectFilter(), // execution doesn't exist yet
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("target", "instance"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								time.Second,
								true,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("12345678"),
								},
							),
						),
					),
					expectPush(
						execution.NewSetEventV2(context.Background(),
							execution.NewAggregate("request/method", "instance"),
							[]*execution.Target{
								{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
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
						{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
					ID:            "request/method",
				},
			},
		},
		{
			"push ok, service target",
			fields{
				eventstore: expectEventstore(
					expectFilter(), // execution doesn't exist yet
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("target", "instance"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								time.Second,
								true,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("12345678"),
								},
							),
						),
					),
					expectPush(
						execution.NewSetEventV2(context.Background(),
							execution.NewAggregate("request/service", "instance"),
							[]*execution.Target{
								{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
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
						{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
					ID:            "request/service",
				},
			},
		},
		{
			"push ok, all target",
			fields{
				eventstore: expectEventstore(
					expectFilter(), // execution doesn't exist yet
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("target", "instance"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								time.Second,
								true,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("12345678"),
								},
							),
						),
					),
					expectPush(
						execution.NewSetEventV2(context.Background(),
							execution.NewAggregate("request", "instance"),
							[]*execution.Target{
								{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
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
						{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
					ID:            "request",
				},
			},
		},
		{
			"push not found, method include",
			fields{
				eventstore: expectEventstore(
					expectFilter(), // execution doesn't exist yet
					expectFilter(), // target doesn't exist
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
						{Type: domain.ExecutionTargetTypeInclude, Target: "request/include"},
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
					expectFilter(), // execution doesn't exist yet
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("request/include", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
								},
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("request/include", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
								},
							),
						),
					),
					expectPush(
						execution.NewSetEventV2(context.Background(),
							execution.NewAggregate("request/method", "instance"),
							[]*execution.Target{
								{Type: domain.ExecutionTargetTypeInclude, Target: "request/include"},
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
						{Type: domain.ExecutionTargetTypeInclude, Target: "request/include"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
					ID:            "request/method",
				},
			},
		},
		{
			"push not found, service include",
			fields{
				eventstore: expectEventstore(
					expectFilter(), // execution doesn't exist yet
					expectFilter(), // target doesn't exist
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
						{Type: domain.ExecutionTargetTypeInclude, Target: "request/include"},
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
					expectFilter(), // execution doesn't exist yet
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("request/include", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
								},
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("request/include", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
								},
							),
						),
					),
					expectPush(
						execution.NewSetEventV2(context.Background(),
							execution.NewAggregate("request/service", "instance"),
							[]*execution.Target{
								{Type: domain.ExecutionTargetTypeInclude, Target: "request/include"},
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
						{Type: domain.ExecutionTargetTypeInclude, Target: "request/include"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
					ID:            "request/service",
				},
			},
		},
		{
			"push not found, all include",
			fields{
				eventstore: expectEventstore(
					expectFilter(), // execution doesn't exist yet
					expectFilter(), // target doesn't exist
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
						{Type: domain.ExecutionTargetTypeInclude, Target: "request/include"},
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
					expectFilter(), // execution doesn't exist yet
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("request/include", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
								},
							),
						),
					),
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("request/include", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
								},
							),
						),
					),
					expectPush(
						execution.NewSetEventV2(context.Background(),
							execution.NewAggregate("request", "instance"),
							[]*execution.Target{
								{Type: domain.ExecutionTargetTypeInclude, Target: "request/include"},
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
						{Type: domain.ExecutionTargetTypeInclude, Target: "request/include"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
					ID:            "request",
				},
			},
		},
		{
			"push ok, remove all targets",
			fields{
				eventstore: expectEventstore(
					expectFilter( // execution has targets
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("request", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
								},
							),
						),
					),
					expectPush(
						execution.NewSetEventV2(context.Background(),
							execution.NewAggregate("request", "instance"),
							[]*execution.Target{},
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
					Targets: []*execution.Target{},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
					ID:            "request",
				},
			},
		},
		{
			"push ok, unchanged execution",
			fields{
				eventstore: expectEventstore(
					expectFilter( // execution has targets
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("request", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
								},
							),
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
					Targets: []*execution.Target{{
						Type:   domain.ExecutionTargetTypeTarget,
						Target: "target",
					}},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
					ID:            "request",
				},
			},
		},
		{
			"push ok, remove all targets",
			fields{
				eventstore: expectEventstore(
					expectFilter( // execution has targets
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("request", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
								},
							),
						),
					),
					expectPush(
						execution.NewSetEventV2(context.Background(),
							execution.NewAggregate("request", "instance"),
							[]*execution.Target{},
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
					Targets: []*execution.Target{},
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
			"push ok, unchanged execution",
			fields{
				eventstore: expectEventstore(
					expectFilter( // execution has targets
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("request", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
								},
							),
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
					Targets: []*execution.Target{{
						Type:   domain.ExecutionTargetTypeTarget,
						Target: "target",
					}},
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
				assertObjectDetails(t, tt.res.details, details)
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
				set:           &SetExecution{Targets: []*execution.Target{{}}},
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
					expectFilter(), // execution doesn't exist yet
					expectFilter(
						target.NewAddedEvent(context.Background(),
							target.NewAggregate("target", "instance"),
							"name",
							domain.TargetTypeWebhook,
							"https://example.com",
							time.Second,
							true,
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "enc",
								KeyID:      "id",
								Crypted:    []byte("12345678"),
							},
						),
					),
					expectPushFailed(
						zerrors.ThrowPreconditionFailed(nil, "id", "name already exists"),
						execution.NewSetEventV2(context.Background(),
							execution.NewAggregate("response/valid", "instance"),
							[]*execution.Target{
								{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
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
						{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
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
						{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
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
						{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
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
					expectFilter(), // execution doesn't exist yet
					expectFilter(
						eventFromEventPusher(
							target.NewAddedEvent(context.Background(),
								target.NewAggregate("target", "instance"),
								"name",
								domain.TargetTypeWebhook,
								"https://example.com",
								time.Second,
								true,
								&crypto.CryptoValue{
									CryptoType: crypto.TypeEncryption,
									Algorithm:  "enc",
									KeyID:      "id",
									Crypted:    []byte("12345678"),
								},
							),
						),
					),
					expectPush(
						execution.NewSetEventV2(context.Background(),
							execution.NewAggregate("response/method", "instance"),
							[]*execution.Target{
								{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
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
						{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
					ID:            "response/method",
				},
			},
		},
		{
			"push ok, service target",
			fields{
				eventstore: expectEventstore(
					expectFilter(), // execution doesn't exist yet
					expectFilter(
						targetAddEvent("target", "instance"),
					),
					expectPush(
						execution.NewSetEventV2(context.Background(),
							execution.NewAggregate("response/service", "instance"),
							[]*execution.Target{
								{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
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
						{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
					ID:            "response/service",
				},
			},
		},
		{
			"push ok, all target",
			fields{
				eventstore: expectEventstore(
					expectFilter(), // execution doesn't exist yet
					expectFilter(
						targetAddEvent("target", "instance"),
					),
					expectPush(
						execution.NewSetEventV2(context.Background(),
							execution.NewAggregate("response", "instance"),
							[]*execution.Target{
								{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
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
						{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
					ID:            "response",
				},
			},
		},
		{
			"push ok, remove all targets",
			fields{
				eventstore: expectEventstore(
					expectFilter( // execution has targets
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("response", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
								},
							),
						),
					),
					expectPush(
						execution.NewSetEventV2(context.Background(),
							execution.NewAggregate("response", "instance"),
							[]*execution.Target{},
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
					Targets: []*execution.Target{},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
					ID:            "response",
				},
			},
		},
		{
			"push ok, unchanged execution",
			fields{
				eventstore: expectEventstore(
					expectFilter( // execution has targets
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("response", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
								},
							),
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
					Targets: []*execution.Target{{
						Type:   domain.ExecutionTargetTypeTarget,
						Target: "target",
					}},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
					ID:            "response",
				},
			},
		},
		{
			"push ok, remove all targets",
			fields{
				eventstore: expectEventstore(
					expectFilter( // execution has targets
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("response", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
								},
							),
						),
					),
					expectPush(
						execution.NewSetEventV2(context.Background(),
							execution.NewAggregate("response", "instance"),
							[]*execution.Target{},
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
					Targets: []*execution.Target{},
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
			"push ok, unchanged execution",
			fields{
				eventstore: expectEventstore(
					expectFilter( // execution has targets
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("response", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
								},
							),
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
					Targets: []*execution.Target{{
						Type:   domain.ExecutionTargetTypeTarget,
						Target: "target",
					}},
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
				assertObjectDetails(t, tt.res.details, details)
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
				set:           &SetExecution{Targets: []*execution.Target{{Target: "target"}}},
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
				set:           &SetExecution{Targets: []*execution.Target{{}}},
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
					expectFilter(), // execution doesn't exist yet
					expectFilter(
						targetAddEvent("target", "instance"),
					),
					expectPushFailed(
						zerrors.ThrowPreconditionFailed(nil, "id", "name already exists"),
						execution.NewSetEventV2(context.Background(),
							execution.NewAggregate("event/valid", "instance"),
							[]*execution.Target{
								{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
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
						{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
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
						{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
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
						{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
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
					expectFilter(), // execution doesn't exist yet
					expectFilter(
						targetAddEvent("target", "instance"),
					),
					expectPush(
						execution.NewSetEventV2(context.Background(),
							execution.NewAggregate("event/event", "instance"),
							[]*execution.Target{
								{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
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
						{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
					ID:            "event/event",
				},
			},
		},
		{
			"push ok, group target",
			fields{
				eventstore: expectEventstore(
					expectFilter(), // execution doesn't exist yet
					expectFilter(
						targetAddEvent("target", "instance"),
					),
					expectPush(
						execution.NewSetEventV2(context.Background(),
							execution.NewAggregate("event/group.*", "instance"),
							[]*execution.Target{
								{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
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
						{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
					ID:            "event/group.*",
				},
			},
		},
		{
			"push ok, all target",
			fields{
				eventstore: expectEventstore(
					expectFilter(), // execution doesn't exist yet
					expectFilter(
						targetAddEvent("target", "instance"),
					),
					expectPush(
						execution.NewSetEventV2(context.Background(),
							execution.NewAggregate("event", "instance"),
							[]*execution.Target{
								{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
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
						{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
					ID:            "event",
				},
			},
		},
		{
			"push ok, remove all targets",
			fields{
				eventstore: expectEventstore(
					expectFilter( // execution has targets
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("event", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
								},
							),
						),
					),
					expectPush(
						execution.NewSetEventV2(context.Background(),
							execution.NewAggregate("event", "instance"),
							[]*execution.Target{},
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
					Targets: []*execution.Target{},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
					ID:            "event",
				},
			},
		},
		{
			"push ok, unchanged execution",
			fields{
				eventstore: expectEventstore(
					expectFilter( // execution has targets
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("event", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
								},
							),
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
					Targets: []*execution.Target{{
						Type:   domain.ExecutionTargetTypeTarget,
						Target: "target",
					}},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
					ID:            "event",
				},
			},
		},
		{
			"push ok, remove all targets",
			fields{
				eventstore: expectEventstore(
					expectFilter( // execution has targets
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("event", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
								},
							),
						),
					),
					expectPush(
						execution.NewSetEventV2(context.Background(),
							execution.NewAggregate("event", "instance"),
							[]*execution.Target{},
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
					Targets: []*execution.Target{},
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
			"push ok, unchanged execution",
			fields{
				eventstore: expectEventstore(
					expectFilter( // execution has targets
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("event", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
								},
							),
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
					Targets: []*execution.Target{{
						Type:   domain.ExecutionTargetTypeTarget,
						Target: "target",
					}},
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
				assertObjectDetails(t, tt.res.details, details)
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
			"empty target, error",
			fields{
				eventstore:           expectEventstore(),
				actionFunctionExists: existsMock(true),
			},
			args{
				ctx:           context.Background(),
				cond:          "function",
				set:           &SetExecution{Targets: []*execution.Target{{}}},
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
					expectFilter(), // execution doesn't exist yet
					expectFilter(
						targetAddEvent("target", "instance"),
					),
					expectPushFailed(
						zerrors.ThrowPreconditionFailed(nil, "id", "name already exists"),
						execution.NewSetEventV2(context.Background(),
							execution.NewAggregate("function/function", "instance"),
							[]*execution.Target{
								{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
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
						{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
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
					expectFilter(), // execution doesn't exist yet
					expectFilter(), // target doesn't exist
				),
				actionFunctionExists: existsMock(true),
			},
			args{
				ctx:  context.Background(),
				cond: "function",
				set: &SetExecution{
					Targets: []*execution.Target{
						{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
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
						{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
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
					expectFilter(), // execution doesn't exist yet
					expectFilter(
						targetAddEvent("target", "instance"),
					),
					expectPush(
						execution.NewSetEventV2(context.Background(),
							execution.NewAggregate("function/function", "instance"),
							[]*execution.Target{
								{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
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
						{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
					},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
					ID:            "function/function",
				},
			},
		},
		{
			"push ok, remove all targets",
			fields{
				actionFunctionExists: existsMock(true),
				eventstore: expectEventstore(
					expectFilter( // execution has targets
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("function/function", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
								},
							),
						),
					),
					expectPush(
						execution.NewSetEventV2(context.Background(),
							execution.NewAggregate("function/function", "instance"),
							[]*execution.Target{},
						),
					),
				),
			},
			args{
				ctx:  context.Background(),
				cond: "function",
				set: &SetExecution{
					Targets: []*execution.Target{},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
					ID:            "function/function",
				},
			},
		},
		{
			"push ok, unchanged execution",
			fields{
				actionFunctionExists: existsMock(true),
				eventstore: expectEventstore(
					expectFilter( // execution has targets
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("function/function", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
								},
							),
						),
					),
				),
			},
			args{
				ctx:  context.Background(),
				cond: "function",
				set: &SetExecution{
					Targets: []*execution.Target{{
						Type:   domain.ExecutionTargetTypeTarget,
						Target: "target",
					}},
				},
				resourceOwner: "instance",
			},
			res{
				details: &domain.ObjectDetails{
					ResourceOwner: "instance",
					ID:            "function/function",
				},
			},
		},
		{
			"push ok, remove all targets",
			fields{
				actionFunctionExists: existsMock(true),
				eventstore: expectEventstore(
					expectFilter( // execution has targets
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("function/function", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
								},
							),
						),
					),
					expectPush(
						execution.NewSetEventV2(context.Background(),
							execution.NewAggregate("function/function", "instance"),
							[]*execution.Target{},
						),
					),
				),
			},
			args{
				ctx:  context.Background(),
				cond: "function",
				set: &SetExecution{
					Targets: []*execution.Target{},
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
			"push ok, unchanged execution",
			fields{
				actionFunctionExists: existsMock(true),
				eventstore: expectEventstore(
					expectFilter( // execution has targets
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("function/function", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
								},
							),
						),
					),
				),
			},
			args{
				ctx:  context.Background(),
				cond: "function",
				set: &SetExecution{
					Targets: []*execution.Target{{
						Type:   domain.ExecutionTargetTypeTarget,
						Target: "target",
					}},
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
				assertObjectDetails(t, tt.res.details, details)
			}
		})
	}
}

func mockExecutionIncludesCache(cache map[string][]string) includeCacheFunc {
	return func(ctx context.Context, id string, resourceOwner string) ([]string, error) {
		included, ok := cache[id]
		if !ok {
			return nil, zerrors.ThrowPreconditionFailed(nil, "", "cache failed")
		}
		return included, nil
	}
}

func TestCommands_checkForIncludeCircular(t *testing.T) {
	type args struct {
		ctx           context.Context
		id            string
		resourceOwner string
		includes      []string
		cache         map[string][]string
	}
	type res struct {
		err func(error) bool
	}
	tests := []struct {
		name string
		args args
		res  res
	}{
		{
			"not found, error",
			args{
				ctx:           context.Background(),
				id:            "id",
				resourceOwner: "",
				includes:      []string{"notexistent"},
				cache:         map[string][]string{},
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			"single, ok",
			args{
				ctx:           context.Background(),
				id:            "id1",
				resourceOwner: "",
				includes:      []string{"id2"},
				cache: map[string][]string{
					"id2": {},
				},
			},
			res{},
		},
		{
			"single, circular",
			args{
				ctx:           context.Background(),
				id:            "id1",
				resourceOwner: "",
				includes:      []string{"id1"},
				cache:         map[string][]string{},
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			"multi 3, ok",
			args{
				ctx:           context.Background(),
				id:            "id1",
				resourceOwner: "",
				includes:      []string{"id2"},
				cache: map[string][]string{
					"id2": {"id3"},
					"id3": {},
				},
			},
			res{},
		},
		{
			"multi 3, circular",
			args{
				ctx:           context.Background(),
				id:            "id1",
				resourceOwner: "",
				includes:      []string{"id2"},
				cache: map[string][]string{
					"id2": {"id3"},
					"id3": {"id1"},
				},
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			"multi 5, ok",
			args{
				ctx:           context.Background(),
				id:            "id1",
				resourceOwner: "",
				includes:      []string{"id11", "id12"},
				cache: map[string][]string{
					"id11": {"id21", "id23"},
					"id12": {"id22"},
					"id21": {},
					"id22": {},
					"id23": {},
				},
			},
			res{},
		},
		{
			"multi 5, circular",
			args{
				ctx:           context.Background(),
				id:            "id1",
				resourceOwner: "",
				includes:      []string{"id11", "id12"},
				cache: map[string][]string{
					"id11": {"id21", "id23"},
					"id12": {"id22"},
					"id21": {},
					"id22": {},
					"id23": {"id1"},
				},
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			"multi 5, circular",
			args{
				ctx:           context.Background(),
				id:            "id1",
				resourceOwner: "",
				includes:      []string{"id11", "id12"},
				cache: map[string][]string{
					"id11": {"id21", "id23"},
					"id12": {"id22"},
					"id21": {},
					"id22": {},
					"id23": {"id11"},
				},
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			"multi 5, circular",
			args{
				ctx:           context.Background(),
				id:            "id1",
				resourceOwner: "",
				includes:      []string{"id11", "id12"},
				cache: map[string][]string{
					"id11": {"id21", "id23"},
					"id12": {"id22"},
					"id21": {"id11"},
					"id22": {},
					"id23": {},
				},
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			"multi 5, circular",
			args{
				ctx:           context.Background(),
				id:            "id1",
				resourceOwner: "",
				includes:      []string{"id11", "id12"},
				cache: map[string][]string{
					"id11": {"id21", "id23"},
					"id12": {"id22"},
					"id21": {},
					"id22": {"id12"},
					"id23": {},
				},
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			"multi 3, maxlevel",
			args{
				ctx:           context.Background(),
				id:            "id1",
				resourceOwner: "",
				includes:      []string{"id2"},
				cache: map[string][]string{
					"id2": {"id3"},
					"id3": {},
				},
			},
			res{},
		},
		{
			"multi 4, over maxlevel",
			args{
				ctx:           context.Background(),
				id:            "id1",
				resourceOwner: "",
				includes:      []string{"id2"},
				cache: map[string][]string{
					"id2": {"id3"},
					"id3": {"id4"},
					"id4": {},
				},
			},
			res{
				err: zerrors.IsPreconditionFailed,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := mockExecutionIncludesCache(tt.args.cache)
			err := checkForIncludeCircular(tt.args.ctx, tt.args.id, tt.args.resourceOwner, tt.args.includes, f, 3)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func mockExecutionIncludesCacheFuncs(cache map[string][]string) (func(string) ([]string, bool), func(string, []string)) {
	return func(s string) ([]string, bool) {
			includes, ok := cache[s]
			return includes, ok
		}, func(s string, strings []string) {
			cache[s] = strings
		}
}

func TestCommands_getExecutionIncludes(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		cache         map[string][]string
		id            string
		resourceOwner string
	}
	type res struct {
		includes []string
		cache    map[string][]string
		err      func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			"new empty, ok",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("request/include", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
								},
							),
						),
					),
				),
			},
			args{
				ctx:           context.Background(),
				cache:         map[string][]string{},
				id:            "id",
				resourceOwner: "instance",
			},
			res{
				includes: []string{},
				cache:    map[string][]string{"id": {}},
			},
		},
		{
			"new includes, ok",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("request/include", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
					),
				),
			},
			args{
				ctx:           context.Background(),
				cache:         map[string][]string{},
				id:            "id",
				resourceOwner: "instance",
			},
			res{
				includes: []string{"include"},
				cache:    map[string][]string{"id": {"include"}},
			},
		},
		{
			"found, ok",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				cache:         map[string][]string{"id": nil},
				id:            "id",
				resourceOwner: "instance",
			},
			res{
				includes: nil,
				cache:    map[string][]string{"id": nil},
			},
		},
		{
			"found includes, ok",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx:           context.Background(),
				cache:         map[string][]string{"id": {"include1", "include2", "include3"}},
				id:            "id",
				resourceOwner: "instance",
			},
			res{
				includes: []string{"include1", "include2", "include3"},
				cache:    map[string][]string{"id": {"include1", "include2", "include3"}},
			},
		},
		{
			"found multiple, ok",
			fields{
				eventstore: expectEventstore(),
			},
			args{
				ctx: context.Background(),
				cache: map[string][]string{
					"id1": {"include1", "include2", "include3"},
					"id2": {"include1", "include2", "include3"},
					"id3": {"include1", "include2", "include3"},
				},
				id:            "id2",
				resourceOwner: "instance",
			},
			res{
				includes: []string{"include1", "include2", "include3"},
				cache: map[string][]string{
					"id1": {"include1", "include2", "include3"},
					"id2": {"include1", "include2", "include3"},
					"id3": {"include1", "include2", "include3"},
				},
			},
		},
		{
			"new multiple, ok",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("request/include", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
								},
							),
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				cache: map[string][]string{
					"id1": {"include1", "include2", "include3"},
					"id2": {"include1", "include2", "include3"},
					"id3": {"include1", "include2", "include3"},
				},
				id:            "id",
				resourceOwner: "instance",
			},
			res{
				includes: []string{},
				cache: map[string][]string{
					"id1": {"include1", "include2", "include3"},
					"id2": {"include1", "include2", "include3"},
					"id3": {"include1", "include2", "include3"},
					"id":  {},
				},
			},
		},
		{
			"new multiple includes, ok",
			fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("request/include", "instance"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
					),
				),
			},
			args{
				ctx: context.Background(),
				cache: map[string][]string{
					"id1": {"include1", "include2", "include3"},
					"id2": {"include1", "include2", "include3"},
					"id3": {"include1", "include2", "include3"},
				},
				id:            "id",
				resourceOwner: "instance",
			},
			res{
				includes: []string{"include"},
				cache: map[string][]string{
					"id1": {"include1", "include2", "include3"},
					"id2": {"include1", "include2", "include3"},
					"id3": {"include1", "include2", "include3"},
					"id":  {"include"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			includes, err := c.getExecutionIncludes(mockExecutionIncludesCacheFuncs(tt.args.cache))(tt.args.ctx, tt.args.id, tt.args.resourceOwner)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.cache, tt.args.cache)
				assert.Equal(t, tt.res.includes, includes)
			}
		})
	}
}

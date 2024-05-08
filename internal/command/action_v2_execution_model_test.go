package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/execution"
)

func TestCommandSide_executionsExistsWriteModel(t *testing.T) {
	type fields struct {
		eventstore func(t *testing.T) *eventstore.Eventstore
	}
	type args struct {
		ctx           context.Context
		ids           []string
		resourceOwner string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		res    bool
	}{
		{
			name: "execution, single",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				ids: []string{"execution"},
			},
			res: true,
		},
		{
			name: "execution, single reset",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
						eventFromEventPusher(
							execution.NewRemovedEvent(context.Background(),
								execution.NewAggregate("execution", "org1"),
							),
						),
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				ids: []string{"execution"},
			},
			res: true,
		},
		{
			name: "execution, single before removed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewRemovedEvent(context.Background(),
								execution.NewAggregate("execution", "org1"),
							),
						),
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				ids: []string{"execution"},
			},
			res: true,
		},
		{
			name: "execution, single removed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
						eventFromEventPusher(
							execution.NewRemovedEvent(context.Background(),
								execution.NewAggregate("execution", "org1"),
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				ids: []string{"execution"},
			},
			res: false,
		},
		{
			name: "execution, multiple",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution1", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution2", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution3", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				ids: []string{"execution1", "execution2", "execution3"},
			},
			res: true,
		},

		{
			name: "execution, multiple, first removed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution1", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution2", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
						eventFromEventPusher(
							execution.NewRemovedEvent(context.Background(),
								execution.NewAggregate("execution1", "org1"),
							),
						),
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution3", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				ids: []string{"execution1", "execution2", "execution3"},
			},
			res: false,
		},
		{
			name: "execution, multiple, second removed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution1", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution2", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
						eventFromEventPusher(
							execution.NewRemovedEvent(context.Background(),
								execution.NewAggregate("execution2", "org1"),
							),
						),
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution3", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				ids: []string{"execution1", "execution2", "execution3"},
			},
			res: false,
		},
		{
			name: "execution, multiple, third removed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution1", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution2", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution3", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
						eventFromEventPusher(
							execution.NewRemovedEvent(context.Background(),
								execution.NewAggregate("execution3", "org1"),
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				ids: []string{"execution1", "execution2", "execution3"},
			},
			res: false,
		},
		{
			name: "execution, multiple, before removed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewRemovedEvent(context.Background(),
								execution.NewAggregate("execution1", "org1"),
							),
						),
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution1", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution2", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution3", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				ids: []string{"execution1", "execution2", "execution3"},
			},
			res: true,
		},
		{
			name: "execution, multiple, all removed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution1", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
						eventFromEventPusher(
							execution.NewRemovedEvent(context.Background(),
								execution.NewAggregate("execution1", "org1"),
							),
						),
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution2", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution3", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
						eventFromEventPusher(
							execution.NewRemovedEvent(context.Background(),
								execution.NewAggregate("execution2", "org1"),
							),
						),
						eventFromEventPusher(
							execution.NewRemovedEvent(context.Background(),
								execution.NewAggregate("execution3", "org1"),
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				ids: []string{"execution1", "execution2", "execution3"},
			},
			res: false,
		},

		{
			name: "execution, multiple, two removed",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution1", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
						eventFromEventPusher(
							execution.NewRemovedEvent(context.Background(),
								execution.NewAggregate("execution1", "org1"),
							),
						),
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution2", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
						eventFromEventPusher(
							execution.NewSetEventV2(context.Background(),
								execution.NewAggregate("execution3", "org1"),
								[]*execution.Target{
									{Type: domain.ExecutionTargetTypeTarget, Target: "target"},
									{Type: domain.ExecutionTargetTypeInclude, Target: "include"},
								},
							),
						),
						eventFromEventPusher(
							execution.NewRemovedEvent(context.Background(),
								execution.NewAggregate("execution2", "org1"),
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				ids: []string{"execution1", "execution2", "execution3"},
			},
			res: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Commands{
				eventstore: tt.fields.eventstore(t),
			}
			assert.Equal(t, tt.res, c.existsExecutionsByIDs(tt.args.ctx, tt.args.ids, tt.args.resourceOwner))

		})
	}
}

package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

func TestCommandSide_UpdateInstance(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx  context.Context
		name string
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "empty name, invalid error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:  authz.WithInstanceID(context.Background(), "INSTANCE"),
				name: "",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "instance not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:  authz.WithInstanceID(context.Background(), "INSTANCE"),
				name: "INSTANCE_CHANGED",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "instance removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewInstanceAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"INSTANCE",
							),
						),
						eventFromEventPusher(
							instance.NewInstanceRemovedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"INSTANCE",
								nil,
							),
						),
					),
				),
			},
			args: args{
				ctx:  authz.WithInstanceID(context.Background(), "INSTANCE"),
				name: "INSTANCE_CHANGED",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "no changes, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewInstanceAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"INSTANCE",
							),
						),
					),
				),
			},
			args: args{
				ctx:  authz.WithInstanceID(context.Background(), "INSTANCE"),
				name: "INSTANCE",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "instance change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							instance.NewInstanceAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"INSTANCE",
							),
						),
					),
					expectPush(
						instance.NewInstanceChangedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"INSTANCE_CHANGED",
						),
					),
				),
			},
			args: args{
				ctx:  authz.WithInstanceID(context.Background(), "INSTANCE"),
				name: "INSTANCE_CHANGED",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.UpdateInstance(tt.args.ctx, tt.args.name)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_RemoveInstance(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		instanceID string
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "instance not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:        authz.WithInstanceID(context.Background(), "INSTANCE"),
				instanceID: "INSTANCE",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "instance removed, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewInstanceAddedEvent(
								context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"INSTANCE",
							),
						),
						eventFromEventPusher(
							instance.NewInstanceRemovedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"INSTANCE",
								nil,
							),
						),
					),
				),
			},
			args: args{
				ctx:        authz.WithInstanceID(context.Background(), "INSTANCE"),
				instanceID: "INSTANCE",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "instance remove, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							instance.NewInstanceAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"INSTANCE",
							),
						),
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							instance.NewDomainAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"instance.domain",
								true,
							),
						),
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							instance.NewDomainAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"custom.domain",
								false,
							),
						),
					),
					expectPush(
						instance.NewInstanceRemovedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"INSTANCE",
							[]string{
								"instance.domain",
								"custom.domain",
							},
						),
					),
				),
			},
			args: args{
				ctx:        authz.WithInstanceID(context.Background(), "INSTANCE"),
				instanceID: "INSTANCE",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.RemoveInstance(tt.args.ctx, tt.args.instanceID)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}

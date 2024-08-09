package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/notification/channels/fs"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/settings"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddDefaultDebugNotificationProviderFile(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx      context.Context
		provider *fs.Config
	}
	type res struct {
		want *domain.ObjectDetails
		err  func(error) bool
	}
	ctx := authz.WithInstanceID(context.Background(), "INSTANCE")
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "provider already existing, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewDebugNotificationProviderFileAddedEvent(ctx,
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx: ctx,
				provider: &fs.Config{
					Compact: true,
					Enabled: true,
				},
			},
			res: res{
				err: zerrors.IsErrorAlreadyExists,
			},
		},
		{
			name: "add provider,ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						instance.NewDebugNotificationProviderFileAddedEvent(ctx,
							&instance.NewAggregate("INSTANCE").Aggregate,
							true,
						),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(ctx, "INSTANCE"),
				provider: &fs.Config{
					Compact: true,
				},
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
			got, err := r.AddDebugNotificationProviderFile(tt.args.ctx, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_ChangeDebugNotificationProviderFile(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx      context.Context
		provider *fs.Config
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
			name: "provider not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				provider: &fs.Config{
					Compact: true,
					Enabled: true,
				},
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "no changes, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewDebugNotificationProviderFileAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				provider: &fs.Config{
					Compact: true,
					Enabled: false,
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "no changes enabled, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewDebugNotificationProviderFileAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				provider: &fs.Config{
					Compact: true,
					Enabled: true,
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							instance.NewDebugNotificationProviderFileAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
							),
						),
					),
					expectPush(
						newDefaultDebugNotificationFileChangedEvent(context.Background(),
							false),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				provider: &fs.Config{
					Compact: false,
					Enabled: false,
				},
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
			got, err := r.ChangeDefaultNotificationFile(tt.args.ctx, tt.args.provider)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func TestCommandSide_RemoveDebugNotificationProviderFile(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx context.Context
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
			name: "provider not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: zerrors.IsNotFound,
			},
		},
		{
			name: "remove, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusherWithInstanceID(
							"INSTANCE",
							instance.NewDebugNotificationProviderFileAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
							),
						),
					),
					expectPush(
						instance.NewDebugNotificationProviderFileRemovedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
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
			got, err := r.RemoveDefaultNotificationFile(tt.args.ctx)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
			if tt.res.err == nil {
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}
func newDefaultDebugNotificationFileChangedEvent(ctx context.Context, compact bool) *instance.DebugNotificationProviderFileChangedEvent {
	event, _ := instance.NewDebugNotificationProviderFileChangedEvent(ctx,
		&instance.NewAggregate("INSTANCE").Aggregate,
		[]settings.DebugNotificationProviderChanges{
			settings.ChangeCompact(compact),
		},
	)
	return event
}

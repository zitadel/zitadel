package command

import (
	"context"
	"testing"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/notification/channels/fs"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/settings"

	"github.com/stretchr/testify/assert"
)

func TestCommandSide_AddDefaultDebugNotificationProviderFile(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx      context.Context
		provider *fs.FSConfig
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
			name: "provider already existing, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewDebugNotificationProviderFileAddedEvent(context.Background(),
								&instance.NewAggregate().Aggregate,
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				provider: &fs.FSConfig{
					Compact: true,
					Enabled: true,
				},
			},
			res: res{
				err: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "add provider,ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								instance.NewDebugNotificationProviderFileAddedEvent(context.Background(),
									&instance.NewAggregate().Aggregate,
									true,
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				provider: &fs.FSConfig{
					Compact: true,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: domain.IAMID,
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
				assert.Equal(t, tt.res.want, got)
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
		provider *fs.FSConfig
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
				provider: &fs.FSConfig{
					Compact: true,
					Enabled: true,
				},
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
							instance.NewDebugNotificationProviderFileAddedEvent(context.Background(),
								&instance.NewAggregate().Aggregate,
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				provider: &fs.FSConfig{
					Compact: true,
					Enabled: false,
				},
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
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
								&instance.NewAggregate().Aggregate,
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				provider: &fs.FSConfig{
					Compact: true,
					Enabled: true,
				},
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewDebugNotificationProviderFileAddedEvent(context.Background(),
								&instance.NewAggregate().Aggregate,
								true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newDefaultDebugNotificationFileChangedEvent(context.Background(),
									false),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				provider: &fs.FSConfig{
					Compact: false,
					Enabled: false,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
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
				assert.Equal(t, tt.res.want, got)
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
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "remove, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewDebugNotificationProviderFileAddedEvent(context.Background(),
								&instance.NewAggregate().Aggregate,
								true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								instance.NewDebugNotificationProviderFileRemovedEvent(context.Background(),
									&instance.NewAggregate().Aggregate),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
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
				assert.Equal(t, tt.res.want, got)
			}
		})
	}
}
func newDefaultDebugNotificationFileChangedEvent(ctx context.Context, compact bool) *instance.DebugNotificationProviderFileChangedEvent {
	event, _ := instance.NewDebugNotificationProviderFileChangedEvent(ctx,
		&instance.NewAggregate().Aggregate,
		[]settings.DebugNotificationProviderChanges{
			settings.ChangeCompact(compact),
		},
	)
	return event
}

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

func TestCommandSide_AddDefaultDebugNotificationProviderLog(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		instanceID string
		provider   *fs.FSConfig
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
							instance.NewDebugNotificationProviderLogAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
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
			name: "add disabled provider,ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								instance.NewDebugNotificationProviderLogAddedEvent(context.Background(),
									&instance.NewAggregate("INSTANCE").Aggregate,
									true,
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				provider: &fs.FSConfig{
					Compact: true,
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
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
								instance.NewDebugNotificationProviderLogAddedEvent(context.Background(),
									&instance.NewAggregate("INSTANCE").Aggregate,
									true,
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				provider: &fs.FSConfig{
					Compact: true,
					Enabled: true,
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
			got, err := r.AddDebugNotificationProviderLog(tt.args.ctx, tt.args.instanceID, tt.args.provider)
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

func TestCommandSide_ChangeDebugNotificationProviderLog(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx        context.Context
		instanceID string
		provider   *fs.FSConfig
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
				ctx:        context.Background(),
				instanceID: "INSTANCE",
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
							instance.NewDebugNotificationProviderLogAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
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
			name: "change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewDebugNotificationProviderLogAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newDefaultDebugNotificationLogChangedEvent(context.Background(),
									false),
							),
						},
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				provider: &fs.FSConfig{
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
		{
			name: "change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewDebugNotificationProviderLogAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newDefaultDebugNotificationLogChangedEvent(context.Background(),
									false),
							),
						},
					),
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
				provider: &fs.FSConfig{
					Compact: false,
					Enabled: true,
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
			got, err := r.ChangeDefaultNotificationLog(tt.args.ctx, tt.args.instanceID, tt.args.provider)
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

func TestCommandSide_RemoveDebugNotificationProviderLog(t *testing.T) {
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
			name: "provider not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:        context.Background(),
				instanceID: "INSTANCE",
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
							instance.NewDebugNotificationProviderLogAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								instance.NewDebugNotificationProviderLogRemovedEvent(context.Background(),
									&instance.NewAggregate("INSTANCE").Aggregate),
							),
						},
					),
				),
			},
			args: args{
				ctx:        context.Background(),
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
			got, err := r.RemoveDefaultNotificationLog(tt.args.ctx, tt.args.instanceID)
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
func newDefaultDebugNotificationLogChangedEvent(ctx context.Context, compact bool) *instance.DebugNotificationProviderLogChangedEvent {
	event, _ := instance.NewDebugNotificationProviderLogChangedEvent(ctx,
		&instance.NewAggregate("INSTANCE").Aggregate,
		[]settings.DebugNotificationProviderChanges{
			settings.ChangeCompact(compact),
		},
	)
	return event
}

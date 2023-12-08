package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/policy"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddDefaultNotificationPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx            context.Context
		resourceOwner  string
		passwordChange bool
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
			name: "notification policy already existing, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewNotificationPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx:            context.Background(),
				resourceOwner:  "INSTANCE",
				passwordChange: true,
			},
			res: res{
				err: zerrors.IsErrorAlreadyExists,
			},
		},
		{
			name: "add policy,ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						instance.NewNotificationPolicyAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							true,
						),
					),
				),
			},
			args: args{
				ctx:            context.Background(),
				resourceOwner:  "INSTANCE",
				passwordChange: true,
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "add empty policy,ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						instance.NewNotificationPolicyAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							true,
						),
					),
				),
			},
			args: args{
				ctx:            context.Background(),
				resourceOwner:  "INSTANCE",
				passwordChange: true,
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
			got, err := r.AddDefaultNotificationPolicy(tt.args.ctx, tt.args.resourceOwner, tt.args.passwordChange)
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

func TestCommandSide_ChangeDefaultNotificationPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx            context.Context
		resourceOwner  string
		passwordChange bool
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
			name: "privacy policy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:            context.Background(),
				resourceOwner:  "INSTANCE",
				passwordChange: true,
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
							instance.NewNotificationPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx:            context.Background(),
				resourceOwner:  "INSTANCE",
				passwordChange: true,
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
						eventFromEventPusher(
							instance.NewNotificationPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								false,
							),
						),
					),
					expectPush(
						newDefaultNotificationPolicyChangedEvent(context.Background(),
							true,
						),
					),
				),
			},
			args: args{
				ctx:            context.Background(),
				resourceOwner:  "INSTANCE",
				passwordChange: true,
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
			got, err := r.ChangeDefaultNotificationPolicy(tt.args.ctx, tt.args.resourceOwner, tt.args.passwordChange)
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

func newDefaultNotificationPolicyChangedEvent(ctx context.Context, passwordChange bool) *instance.NotificationPolicyChangedEvent {
	event, _ := instance.NewNotificationPolicyChangedEvent(ctx,
		&instance.NewAggregate("INSTANCE").Aggregate,
		[]policy.NotificationPolicyChanges{
			policy.ChangePasswordChange(passwordChange),
		},
	)
	return event
}

package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/repository"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

func TestCommandSide_AddDefaultLockoutPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		policy *domain.LockoutPolicy
	}
	type res struct {
		want *domain.LockoutPolicy
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "lockout policy already existing, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewLockoutPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								10,
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.LockoutPolicy{
					MaxPasswordAttempts: 10,
					ShowLockOutFailures: true,
				},
			},
			res: res{
				err: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "add policy,ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusherWithInstanceID(
								"INSTANCE",
								instance.NewLockoutPolicyAddedEvent(context.Background(),
									&instance.NewAggregate("INSTANCE").Aggregate,
									10,
									true,
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				policy: &domain.LockoutPolicy{
					MaxPasswordAttempts: 10,
					ShowLockOutFailures: true,
				},
			},
			res: res{
				want: &domain.LockoutPolicy{
					ObjectRoot: models.ObjectRoot{
						InstanceID:    "INSTANCE",
						AggregateID:   "INSTANCE",
						ResourceOwner: "INSTANCE",
					},
					MaxPasswordAttempts: 10,
					ShowLockOutFailures: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.AddDefaultLockoutPolicy(tt.args.ctx, tt.args.policy)
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

func TestCommandSide_ChangeDefaultLockoutPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		policy *domain.LockoutPolicy
	}
	type res struct {
		want *domain.LockoutPolicy
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "lockout policy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.LockoutPolicy{
					MaxPasswordAttempts: 10,
					ShowLockOutFailures: true,
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
							instance.NewLockoutPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								10,
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.LockoutPolicy{
					MaxPasswordAttempts: 10,
					ShowLockOutFailures: true,
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
							instance.NewLockoutPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								10,
								true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newDefaultLockoutPolicyChangedEvent(context.Background(), 20, false),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.LockoutPolicy{
					MaxPasswordAttempts: 20,
					ShowLockOutFailures: false,
				},
			},
			res: res{
				want: &domain.LockoutPolicy{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "INSTANCE",
						ResourceOwner: "INSTANCE",
					},
					MaxPasswordAttempts: 20,
					ShowLockOutFailures: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeDefaultLockoutPolicy(tt.args.ctx, tt.args.policy)
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

func newDefaultLockoutPolicyChangedEvent(ctx context.Context, maxAttempts uint64, showLockoutFailure bool) *instance.LockoutPolicyChangedEvent {
	event, _ := instance.NewLockoutPolicyChangedEvent(ctx,
		&instance.NewAggregate("INSTANCE").Aggregate,
		[]policy.LockoutPolicyChanges{
			policy.ChangeMaxAttempts(maxAttempts),
			policy.ChangeShowLockOutFailures(showLockoutFailure),
		},
	)
	return event
}

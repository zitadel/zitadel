package command

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/policy"
	"github.com/stretchr/testify/assert"
	"testing"
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
							iam.NewLockoutPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
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
							eventFromEventPusher(
								iam.NewLockoutPolicyAddedEvent(context.Background(),
									&iam.NewAggregate().Aggregate,
									10,
									true,
								),
							),
						},
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
				want: &domain.LockoutPolicy{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "IAM",
						ResourceOwner: "IAM",
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
							iam.NewLockoutPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
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
							iam.NewLockoutPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
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
						AggregateID:   "IAM",
						ResourceOwner: "IAM",
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

func newDefaultLockoutPolicyChangedEvent(ctx context.Context, maxAttempts uint64, showLockoutFailure bool) *iam.LockoutPolicyChangedEvent {
	event, _ := iam.NewLockoutPolicyChangedEvent(ctx,
		&iam.NewAggregate().Aggregate,
		[]policy.LockoutPolicyChanges{
			policy.ChangeMaxAttempts(maxAttempts),
			policy.ChangeShowLockOutFailures(showLockoutFailure),
		},
	)
	return event
}

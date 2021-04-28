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

func TestCommandSide_AddDefaultPasswordLockoutPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		policy *domain.PasswordLockoutPolicy
	}
	type res struct {
		want *domain.PasswordLockoutPolicy
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "password lockout policy already existing, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewPasswordLockoutPolicyAddedEvent(context.Background(),
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
				policy: &domain.PasswordLockoutPolicy{
					MaxAttempts:         10,
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
								iam.NewPasswordLockoutPolicyAddedEvent(context.Background(),
									&iam.NewAggregate().Aggregate,
									10,
									true,
								),
							),
						},
						nil,
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.PasswordLockoutPolicy{
					MaxAttempts:         10,
					ShowLockOutFailures: true,
				},
			},
			res: res{
				want: &domain.PasswordLockoutPolicy{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "IAM",
						ResourceOwner: "IAM",
					},
					MaxAttempts:         10,
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
			got, err := r.AddDefaultPasswordLockoutPolicy(tt.args.ctx, tt.args.policy)
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

func TestCommandSide_ChangeDefaultPasswordLockoutPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		policy *domain.PasswordLockoutPolicy
	}
	type res struct {
		want *domain.PasswordLockoutPolicy
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "password lockout policy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.PasswordLockoutPolicy{
					MaxAttempts:         10,
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
							iam.NewPasswordLockoutPolicyAddedEvent(context.Background(),
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
				policy: &domain.PasswordLockoutPolicy{
					MaxAttempts:         10,
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
							iam.NewPasswordLockoutPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								10,
								true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newDefaultPasswordLockoutPolicyChangedEvent(context.Background(), 20, false),
							),
						},
						nil,
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.PasswordLockoutPolicy{
					MaxAttempts:         20,
					ShowLockOutFailures: false,
				},
			},
			res: res{
				want: &domain.PasswordLockoutPolicy{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "IAM",
						ResourceOwner: "IAM",
					},
					MaxAttempts:         20,
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
			got, err := r.ChangeDefaultPasswordLockoutPolicy(tt.args.ctx, tt.args.policy)
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

func newDefaultPasswordLockoutPolicyChangedEvent(ctx context.Context, maxAttempts uint64, showLockoutFailure bool) *iam.PasswordLockoutPolicyChangedEvent {
	event, _ := iam.NewPasswordLockoutPolicyChangedEvent(ctx,
		&iam.NewAggregate().Aggregate,
		[]policy.PasswordLockoutPolicyChanges{
			policy.ChangeMaxAttempts(maxAttempts),
			policy.ChangeShowLockOutFailures(showLockoutFailure),
		},
	)
	return event
}

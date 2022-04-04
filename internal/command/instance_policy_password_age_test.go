package command

import (
	"context"
	"testing"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/policy"
	"github.com/stretchr/testify/assert"
)

func TestCommandSide_AddDefaultPasswordAgePolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		policy *domain.PasswordAgePolicy
	}
	type res struct {
		want *domain.PasswordAgePolicy
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "password age policy already existing, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewPasswordAgePolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								365,
								10,
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.PasswordAgePolicy{
					MaxAgeDays:     365,
					ExpireWarnDays: 10,
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
								instance.NewPasswordAgePolicyAddedEvent(context.Background(),
									&instance.NewAggregate("INSTANCE").Aggregate,
									365,
									10,
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				policy: &domain.PasswordAgePolicy{
					ExpireWarnDays: 365,
					MaxAgeDays:     10,
				},
			},
			res: res{
				want: &domain.PasswordAgePolicy{
					ObjectRoot: models.ObjectRoot{
						InstanceID:    "INSTANCE",
						AggregateID:   "INSTANCE",
						ResourceOwner: "INSTANCE",
					},
					ExpireWarnDays: 365,
					MaxAgeDays:     10,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.AddDefaultPasswordAgePolicy(tt.args.ctx, tt.args.policy)
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

func TestCommandSide_ChangeDefaultPasswordAgePolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		policy *domain.PasswordAgePolicy
	}
	type res struct {
		want *domain.PasswordAgePolicy
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "password age policy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.PasswordAgePolicy{
					MaxAgeDays:     365,
					ExpireWarnDays: 10,
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
							instance.NewPasswordAgePolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								365,
								10,
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.PasswordAgePolicy{
					ExpireWarnDays: 365,
					MaxAgeDays:     10,
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
							instance.NewPasswordAgePolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								365,
								10,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newDefaultPasswordAgePolicyChangedEvent(context.Background(), 125, 5),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.PasswordAgePolicy{
					MaxAgeDays:     125,
					ExpireWarnDays: 5,
				},
			},
			res: res{
				want: &domain.PasswordAgePolicy{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "INSTANCE",
						ResourceOwner: "INSTANCE",
					},
					MaxAgeDays:     125,
					ExpireWarnDays: 5,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeDefaultPasswordAgePolicy(tt.args.ctx, tt.args.policy)
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

func newDefaultPasswordAgePolicyChangedEvent(ctx context.Context, maxAgeDays, expiryWarnDays uint64) *instance.PasswordAgePolicyChangedEvent {
	event, _ := instance.NewPasswordAgePolicyChangedEvent(ctx,
		&instance.NewAggregate("INSTANCE").Aggregate,
		[]policy.PasswordAgePolicyChanges{
			policy.ChangeExpireWarnDays(expiryWarnDays),
			policy.ChangeMaxAgeDays(maxAgeDays),
		},
	)
	return event
}

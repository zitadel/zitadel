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

func TestCommandSide_AddDefaultOrgIAMPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		policy *domain.OrgIAMPolicy
	}
	type res struct {
		want *domain.OrgIAMPolicy
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "orgiam policy already existing, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewOrgIAMPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.OrgIAMPolicy{
					UserLoginMustBeDomain: true,
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
								iam.NewOrgIAMPolicyAddedEvent(context.Background(),
									&iam.NewAggregate().Aggregate,
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
				policy: &domain.OrgIAMPolicy{
					UserLoginMustBeDomain: true,
				},
			},
			res: res{
				want: &domain.OrgIAMPolicy{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "IAM",
						ResourceOwner: "IAM",
					},
					UserLoginMustBeDomain: true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.AddDefaultOrgIAMPolicy(tt.args.ctx, tt.args.policy)
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

func TestCommandSide_ChangeDefaultOrgIAMPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		policy *domain.OrgIAMPolicy
	}
	type res struct {
		want *domain.OrgIAMPolicy
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "orgiampolicy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.OrgIAMPolicy{
					UserLoginMustBeDomain: true,
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
							iam.NewOrgIAMPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								true,
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.OrgIAMPolicy{
					UserLoginMustBeDomain: true,
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
							iam.NewOrgIAMPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								true,
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newDefaultOrgIAMPolicyChangedEvent(context.Background(), false),
							),
						},
						nil,
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.OrgIAMPolicy{
					UserLoginMustBeDomain: false,
				},
			},
			res: res{
				want: &domain.OrgIAMPolicy{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "IAM",
						ResourceOwner: "IAM",
					},
					UserLoginMustBeDomain: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeDefaultOrgIAMPolicy(tt.args.ctx, tt.args.policy)
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

func newDefaultOrgIAMPolicyChangedEvent(ctx context.Context, userLoginMustBeDomain bool) *iam.OrgIAMPolicyChangedEvent {
	event, _ := iam.NewOrgIAMPolicyChangedEvent(ctx,
		&iam.NewAggregate().Aggregate,
		[]policy.OrgIAMPolicyChanges{
			policy.ChangeUserLoginMustBeDomain(userLoginMustBeDomain),
		},
	)
	return event
}

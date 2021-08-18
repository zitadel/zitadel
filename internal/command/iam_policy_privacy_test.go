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

func TestCommandSide_AddDefaultPrivacyPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		policy *domain.PrivacyPolicy
	}
	type res struct {
		want *domain.PrivacyPolicy
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "privacy policy already existing, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewPrivacyPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"TOSLink",
								"PrivacyLink",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.PrivacyPolicy{
					TOSLink:     "TOSLink",
					PrivacyLink: "PrivacyLink",
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
								iam.NewPrivacyPolicyAddedEvent(context.Background(),
									&iam.NewAggregate().Aggregate,
									"TOSLink",
									"PrivacyLink",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.PrivacyPolicy{
					TOSLink:     "TOSLink",
					PrivacyLink: "PrivacyLink",
				},
			},
			res: res{
				want: &domain.PrivacyPolicy{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "IAM",
						ResourceOwner: "IAM",
					},
					TOSLink:     "TOSLink",
					PrivacyLink: "PrivacyLink",
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
						[]*repository.Event{
							eventFromEventPusher(
								iam.NewPrivacyPolicyAddedEvent(context.Background(),
									&iam.NewAggregate().Aggregate,
									"",
									"",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.PrivacyPolicy{
					TOSLink:     "",
					PrivacyLink: "",
				},
			},
			res: res{
				want: &domain.PrivacyPolicy{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "IAM",
						ResourceOwner: "IAM",
					},
					TOSLink:     "",
					PrivacyLink: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.AddDefaultPrivacyPolicy(tt.args.ctx, tt.args.policy)
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

func TestCommandSide_ChangeDefaultPrivacyPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		policy *domain.PrivacyPolicy
	}
	type res struct {
		want *domain.PrivacyPolicy
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
				ctx: context.Background(),
				policy: &domain.PrivacyPolicy{
					TOSLink:     "TOSLink",
					PrivacyLink: "PrivacyLink",
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
							iam.NewPrivacyPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"TOSLink",
								"PrivacyLink",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.PrivacyPolicy{
					TOSLink:     "TOSLink",
					PrivacyLink: "PrivacyLink",
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
							iam.NewPrivacyPolicyAddedEvent(context.Background(),
								&iam.NewAggregate().Aggregate,
								"TOSLink",
								"PrivacyLink",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newDefaultPrivacyPolicyChangedEvent(context.Background(),
									"TOSLinkChanged",
									"PrivacyLinkChanged",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.PrivacyPolicy{
					TOSLink:     "TOSLinkChanged",
					PrivacyLink: "PrivacyLinkChanged",
				},
			},
			res: res{
				want: &domain.PrivacyPolicy{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "IAM",
						ResourceOwner: "IAM",
					},
					TOSLink:     "TOSLinkChanged",
					PrivacyLink: "PrivacyLinkChanged",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeDefaultPrivacyPolicy(tt.args.ctx, tt.args.policy)
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

func newDefaultPrivacyPolicyChangedEvent(ctx context.Context, tosLink, privacyLink string) *iam.PrivacyPolicyChangedEvent {
	event, _ := iam.NewPrivacyPolicyChangedEvent(ctx,
		&iam.NewAggregate().Aggregate,
		[]policy.PrivacyPolicyChanges{
			policy.ChangeTOSLink(tosLink),
			policy.ChangePrivacyLink(privacyLink),
		},
	)
	return event
}

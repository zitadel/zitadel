package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/repository/org"
	"github.com/caos/zitadel/internal/repository/policy"
)

func TestCommandSide_AddPrivacyPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		orgID  string
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
			name: "org id missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "policy already existing, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewPrivacyPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"TOSLink",
								"PrivacyLink",
							),
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
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
								org.NewPrivacyPolicyAddedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"TOSLink",
									"PrivacyLink",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.PrivacyPolicy{
					TOSLink:     "TOSLink",
					PrivacyLink: "PrivacyLink",
				},
			},
			res: res{
				want: &domain.PrivacyPolicy{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "org1",
						ResourceOwner: "org1",
					},
					TOSLink:     "TOSLink",
					PrivacyLink: "PrivacyLink",
				},
			},
		},
		{
			name: "add policy empty links, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewPrivacyPolicyAddedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"",
									"",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.PrivacyPolicy{
					TOSLink:     "",
					PrivacyLink: "",
				},
			},
			res: res{
				want: &domain.PrivacyPolicy{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "org1",
						ResourceOwner: "org1",
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
			got, err := r.AddPrivacyPolicy(tt.args.ctx, tt.args.orgID, tt.args.policy)
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

func TestCommandSide_ChangePrivacyPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		orgID  string
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
			name: "org id missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
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
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "policy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
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
							org.NewPrivacyPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"TOSLink",
								"PrivacyLink",
							),
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
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
							org.NewPrivacyPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"TOSLink",
								"PrivacyLink",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newPrivacyPolicyChangedEvent(context.Background(), "org1", "TOSLinkChange", "PrivacyLinkChange"),
							),
						},
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.PrivacyPolicy{
					TOSLink:     "TOSLinkChange",
					PrivacyLink: "PrivacyLinkChange",
				},
			},
			res: res{
				want: &domain.PrivacyPolicy{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "org1",
						ResourceOwner: "org1",
					},
					TOSLink:     "TOSLinkChange",
					PrivacyLink: "PrivacyLinkChange",
				},
			},
		},
		{
			name: "change to empty links, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewPrivacyPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"TOSLink",
								"PrivacyLink",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newPrivacyPolicyChangedEvent(context.Background(), "org1", "", ""),
							),
						},
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.PrivacyPolicy{
					TOSLink:     "",
					PrivacyLink: "",
				},
			},
			res: res{
				want: &domain.PrivacyPolicy{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "org1",
						ResourceOwner: "org1",
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
			got, err := r.ChangePrivacyPolicy(tt.args.ctx, tt.args.orgID, tt.args.policy)
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

func TestCommandSide_RemovePrivacyPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx   context.Context
		orgID string
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
			name: "org id missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "policy not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
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
							org.NewPrivacyPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"TOSLink",
								"PrivacyLink",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewPrivacyPolicyRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate),
							),
						},
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "org1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.RemovePrivacyPolicy(tt.args.ctx, tt.args.orgID)
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

func newPrivacyPolicyChangedEvent(ctx context.Context, orgID string, tosLink, privacyLink string) *org.PrivacyPolicyChangedEvent {
	event, _ := org.NewPrivacyPolicyChangedEvent(ctx,
		&org.NewAggregate(orgID, orgID).Aggregate,
		[]policy.PrivacyPolicyChanges{
			policy.ChangeTOSLink(tosLink),
			policy.ChangePrivacyLink(privacyLink),
		},
	)
	return event
}

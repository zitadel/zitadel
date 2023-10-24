package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/policy"
)

func TestCommandSide_AddDefaultPrivacyPolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx          context.Context
		tosLink      string
		privacyLink  string
		helpLink     string
		supportEmail domain.EmailAddress
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
			name: "privacy policy already existing, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewPrivacyPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"TOSLink",
								"PrivacyLink",
								"HelpLink",
								"support@example.com",
							),
						),
					),
				),
			},
			args: args{
				ctx:          context.Background(),
				tosLink:      "TOSLink",
				privacyLink:  "PrivacyLink",
				helpLink:     "HelpLink",
				supportEmail: "support@example.com",
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
						instance.NewPrivacyPolicyAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"TOSLink",
							"PrivacyLink",
							"HelpLink",
							"support@example.com",
						),
					),
				),
			},
			args: args{
				ctx:          authz.WithInstanceID(context.Background(), "INSTANCE"),
				tosLink:      "TOSLink",
				privacyLink:  "PrivacyLink",
				helpLink:     "HelpLink",
				supportEmail: "support@example.com",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "INSTANCE",
				},
			},
		},
		{
			name: "wrong email, can't add policy",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:          authz.WithInstanceID(context.Background(), "INSTANCE"),
				tosLink:      "TOSLink",
				privacyLink:  "PrivacyLink",
				helpLink:     "HelpLink",
				supportEmail: "wrong email",
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "add empty policy,ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						instance.NewPrivacyPolicyAddedEvent(context.Background(),
							&instance.NewAggregate("INSTANCE").Aggregate,
							"",
							"",
							"",
							"",
						),
					),
				),
			},
			args: args{
				ctx:          authz.WithInstanceID(context.Background(), "INSTANCE"),
				tosLink:      "",
				privacyLink:  "",
				helpLink:     "",
				supportEmail: "",
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
			got, err := r.AddDefaultPrivacyPolicy(tt.args.ctx, tt.args.tosLink, tt.args.privacyLink, tt.args.helpLink, tt.args.supportEmail)
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
					TOSLink:      "TOSLink",
					PrivacyLink:  "PrivacyLink",
					HelpLink:     "HelpLink",
					SupportEmail: "support@example.com",
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
							instance.NewPrivacyPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"TOSLink",
								"PrivacyLink",
								"HelpLink",
								"support@example.com",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.PrivacyPolicy{
					TOSLink:      "TOSLink",
					PrivacyLink:  "PrivacyLink",
					HelpLink:     "HelpLink",
					SupportEmail: "support@example.com",
				},
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "wrong email, can't change policy",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.PrivacyPolicy{
					TOSLink:      "TOSLink",
					PrivacyLink:  "PrivacyLink",
					HelpLink:     "HelpLink",
					SupportEmail: "wrong email",
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewPrivacyPolicyAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								"TOSLink",
								"PrivacyLink",
								"HelpLink",
								"support@example.com",
							),
						),
					),
					expectPush(
						newDefaultPrivacyPolicyChangedEvent(context.Background(),
							"TOSLinkChanged",
							"PrivacyLinkChanged",
							"HelpLinkChanged",
							"support2@example.com",
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.PrivacyPolicy{
					TOSLink:      "TOSLinkChanged",
					PrivacyLink:  "PrivacyLinkChanged",
					HelpLink:     "HelpLinkChanged",
					SupportEmail: "support2@example.com",
				},
			},
			res: res{
				want: &domain.PrivacyPolicy{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "INSTANCE",
						ResourceOwner: "INSTANCE",
						InstanceID:    "INSTANCE",
					},
					TOSLink:      "TOSLinkChanged",
					PrivacyLink:  "PrivacyLinkChanged",
					HelpLink:     "HelpLinkChanged",
					SupportEmail: "support2@example.com",
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

func newDefaultPrivacyPolicyChangedEvent(ctx context.Context, tosLink, privacyLink, helpLink, supportEmail string) *instance.PrivacyPolicyChangedEvent {
	event, _ := instance.NewPrivacyPolicyChangedEvent(ctx,
		&instance.NewAggregate("INSTANCE").Aggregate,
		[]policy.PrivacyPolicyChanges{
			policy.ChangeTOSLink(tosLink),
			policy.ChangePrivacyLink(privacyLink),
			policy.ChangeHelpLink(helpLink),
			policy.ChangeSupportEmail(domain.EmailAddress(supportEmail)),
		},
	)
	return event
}

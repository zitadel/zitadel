package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/policy"
	"github.com/zitadel/zitadel/internal/zerrors"
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
					TOSLink:        "TOSLink",
					PrivacyLink:    "PrivacyLink",
					HelpLink:       "HelpLink",
					SupportEmail:   "support@example.com",
					DocsLink:       "DocsLink",
					CustomLink:     "CustomLink",
					CustomLinkText: "CustomLinkText",
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
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
								&org.NewAggregate("org1").Aggregate,
								"TOSLink",
								"PrivacyLink",
								"HelpLink",
								"support@example.com",
								"DocsLink",
								"CustomLink",
								"CustomLinkText"),
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.PrivacyPolicy{
					TOSLink:        "TOSLink",
					PrivacyLink:    "PrivacyLink",
					HelpLink:       "HelpLink",
					SupportEmail:   "support@example.com",
					DocsLink:       "DocsLink",
					CustomLink:     "CustomLink",
					CustomLinkText: "CustomLinkText",
				},
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
						org.NewPrivacyPolicyAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"TOSLink",
							"PrivacyLink",
							"HelpLink",
							"support@example.com",
							"DocsLink",
							"CustomLink",
							"CustomLinkText",
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.PrivacyPolicy{
					TOSLink:        "TOSLink",
					PrivacyLink:    "PrivacyLink",
					HelpLink:       "HelpLink",
					SupportEmail:   "support@example.com",
					DocsLink:       "DocsLink",
					CustomLink:     "CustomLink",
					CustomLinkText: "CustomLinkText",
				},
			},
			res: res{
				want: &domain.PrivacyPolicy{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "org1",
						ResourceOwner: "org1",
					},
					TOSLink:        "TOSLink",
					PrivacyLink:    "PrivacyLink",
					HelpLink:       "HelpLink",
					SupportEmail:   "support@example.com",
					DocsLink:       "DocsLink",
					CustomLink:     "CustomLink",
					CustomLinkText: "CustomLinkText",
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
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.PrivacyPolicy{
					TOSLink:        "TOSLink",
					PrivacyLink:    "PrivacyLink",
					HelpLink:       "HelpLink",
					SupportEmail:   "wrong email",
					DocsLink:       "DocsLink",
					CustomLink:     "CustomLink",
					CustomLinkText: "CustomLinkText",
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "add policy empty links and empty support email, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						org.NewPrivacyPolicyAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							"",
							"",
							"",
							"",
							"",
							"",
							"",
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.PrivacyPolicy{
					TOSLink:        "",
					PrivacyLink:    "",
					HelpLink:       "",
					SupportEmail:   "",
					DocsLink:       "",
					CustomLink:     "",
					CustomLinkText: "",
				},
			},
			res: res{
				want: &domain.PrivacyPolicy{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "org1",
						ResourceOwner: "org1",
					},
					TOSLink:        "",
					PrivacyLink:    "",
					HelpLink:       "",
					SupportEmail:   "",
					DocsLink:       "",
					CustomLink:     "",
					CustomLinkText: "",
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
					TOSLink:        "TOSLink",
					PrivacyLink:    "PrivacyLink",
					HelpLink:       "HelpLink",
					SupportEmail:   "support@example.com",
					DocsLink:       "DocsLink",
					CustomLink:     "CustomLink",
					CustomLinkText: "CustomLinkText",
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
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
					TOSLink:        "TOSLink",
					PrivacyLink:    "PrivacyLink",
					HelpLink:       "HelpLink",
					SupportEmail:   "support@example.com",
					DocsLink:       "DocsLink",
					CustomLink:     "CustomLink",
					CustomLinkText: "CustomLinkText",
				},
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
							org.NewPrivacyPolicyAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								"TOSLink",
								"PrivacyLink",
								"HelpLink",
								"support@example.com",
								"DocsLink",
								"CustomLink",
								"CustomLinkText",
							),
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.PrivacyPolicy{
					TOSLink:        "TOSLink",
					PrivacyLink:    "PrivacyLink",
					HelpLink:       "HelpLink",
					SupportEmail:   "support@example.com",
					DocsLink:       "DocsLink",
					CustomLink:     "CustomLink",
					CustomLinkText: "CustomLinkText",
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "wrong email, can't change policy",
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.PrivacyPolicy{
					TOSLink:        "TOSLinkChange",
					PrivacyLink:    "PrivacyLinkChange",
					HelpLink:       "HelpLinkChange",
					SupportEmail:   "wrong email",
					DocsLink:       "DocsLink",
					CustomLink:     "CustomLink",
					CustomLinkText: "CustomLinkText",
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
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
								&org.NewAggregate("org1").Aggregate,
								"TOSLink",
								"PrivacyLink",
								"HelpLink",
								"support@example.com",
								"DocsLink",
								"CustomLink",
								"CustomLinkText",
							),
						),
					),
					expectPush(
						newPrivacyPolicyChangedEvent(context.Background(), "org1", "TOSLinkChange", "PrivacyLinkChange", "HelpLinkChange", "support2@example.com", "DocsLinkChange", "CustomLinkChange", "CustomLinkTextChange"),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.PrivacyPolicy{
					TOSLink:        "TOSLinkChange",
					PrivacyLink:    "PrivacyLinkChange",
					HelpLink:       "HelpLinkChange",
					SupportEmail:   "support2@example.com",
					DocsLink:       "DocsLinkChange",
					CustomLink:     "CustomLinkChange",
					CustomLinkText: "CustomLinkTextChange",
				},
			},
			res: res{
				want: &domain.PrivacyPolicy{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "org1",
						ResourceOwner: "org1",
					},
					TOSLink:        "TOSLinkChange",
					PrivacyLink:    "PrivacyLinkChange",
					HelpLink:       "HelpLinkChange",
					SupportEmail:   "support2@example.com",
					DocsLink:       "DocsLinkChange",
					CustomLink:     "CustomLinkChange",
					CustomLinkText: "CustomLinkTextChange",
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
								&org.NewAggregate("org1").Aggregate,
								"TOSLink",
								"PrivacyLink",
								"HelpLink",
								"support@example.com",
								"DocsLink",
								"CustomLink",
								"CustomLinkText",
							),
						),
					),
					expectPush(
						newPrivacyPolicyChangedEvent(context.Background(), "org1", "", "", "", "", "", "", ""),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.PrivacyPolicy{
					TOSLink:        "",
					PrivacyLink:    "",
					HelpLink:       "",
					SupportEmail:   "",
					DocsLink:       "",
					CustomLink:     "",
					CustomLinkText: "",
				},
			},
			res: res{
				want: &domain.PrivacyPolicy{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "org1",
						ResourceOwner: "org1",
					},
					TOSLink:        "",
					PrivacyLink:    "",
					HelpLink:       "",
					SupportEmail:   "",
					DocsLink:       "",
					CustomLink:     "",
					CustomLinkText: "",
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
				err: zerrors.IsErrorInvalidArgument,
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
				err: zerrors.IsNotFound,
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
								&org.NewAggregate("org1").Aggregate,
								"TOSLink",
								"PrivacyLink",
								"HelpLink",
								"support@example.com",
								"DocsLink",
								"CustomLink",
								"CustomLinkText",
							),
						),
					),
					expectPush(
						org.NewPrivacyPolicyRemovedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate),
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
				assertObjectDetails(t, tt.res.want, got)
			}
		})
	}
}

func newPrivacyPolicyChangedEvent(ctx context.Context, orgID string, tosLink, privacyLink, helpLink, supportEmail, docsLink, customLink, customLinkText string) *org.PrivacyPolicyChangedEvent {
	event, _ := org.NewPrivacyPolicyChangedEvent(ctx,
		&org.NewAggregate(orgID).Aggregate,
		[]policy.PrivacyPolicyChanges{
			policy.ChangeTOSLink(tosLink),
			policy.ChangePrivacyLink(privacyLink),
			policy.ChangeHelpLink(helpLink),
			policy.ChangeSupportEmail(domain.EmailAddress(supportEmail)),
			policy.ChangeDocsLink(docsLink),
			policy.ChangeCustomLink(customLink),
			policy.ChangeCustomLinkText(customLinkText),
		},
	)
	return event
}

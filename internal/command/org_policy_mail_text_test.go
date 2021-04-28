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

func TestCommandSide_AddMailText(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		orgID  string
		policy *domain.MailText
	}
	type res struct {
		want *domain.MailText
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
				policy: &domain.MailText{
					MailTextType: "mail-text-type",
					Language:     "de",
					Title:        "title",
					PreHeader:    "pre-header",
					Subject:      "subject",
					Greeting:     "greeting",
					Text:         "text",
					ButtonText:   "button-text",
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "mail text already existing, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewMailTextAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"mail-text-type",
								"de",
								"title",
								"pre-header",
								"subject",
								"greeting",
								"text",
								"button-text",
							),
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.MailText{
					MailTextType: "mail-text-type",
					Language:     "de",
					Title:        "title",
					PreHeader:    "pre-header",
					Subject:      "subject",
					Greeting:     "greeting",
					Text:         "text",
					ButtonText:   "button-text",
				},
			},
			res: res{
				err: caos_errs.IsErrorAlreadyExists,
			},
		},
		{
			name: "mail text already existing, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewMailTextAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"mail-text-type",
								"de",
								"title",
								"pre-header",
								"subject",
								"greeting",
								"text",
								"button-text",
							),
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.MailText{
					MailTextType: "mail-text-type",
					Language:     "de",
					Title:        "title",
					PreHeader:    "pre-header",
					Subject:      "subject",
					Greeting:     "greeting",
					Text:         "text",
					ButtonText:   "button-text",
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
								org.NewMailTextAddedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"mail-text-type",
									"de",
									"title",
									"pre-header",
									"subject",
									"greeting",
									"text",
									"button-text",
								),
							),
						},
						nil,
						uniqueConstraintsFromEventConstraint(policy.NewAddMailTextUniqueConstraint("org1", "mail-text-type", "de")),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.MailText{
					MailTextType: "mail-text-type",
					Language:     "de",
					Title:        "title",
					PreHeader:    "pre-header",
					Subject:      "subject",
					Greeting:     "greeting",
					Text:         "text",
					ButtonText:   "button-text",
				},
			},
			res: res{
				want: &domain.MailText{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "org1",
						ResourceOwner: "org1",
					},
					MailTextType: "mail-text-type",
					Language:     "de",
					Title:        "title",
					PreHeader:    "pre-header",
					Subject:      "subject",
					Greeting:     "greeting",
					Text:         "text",
					ButtonText:   "button-text",
					State:        domain.PolicyStateActive,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.AddMailText(tt.args.ctx, tt.args.orgID, tt.args.policy)
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

func TestCommandSide_ChangeMailText(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		orgID  string
		policy *domain.MailText
	}
	type res struct {
		want *domain.MailText
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
				policy: &domain.MailText{
					MailTextType: "mail-text-type",
					Language:     "de",
					Title:        "title",
					PreHeader:    "pre-header",
					Subject:      "subject",
					Greeting:     "greeting",
					Text:         "text",
					ButtonText:   "button-text",
				},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "mailtext invalid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				policy: &domain.MailText{},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "mail template not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.MailText{
					MailTextType: "mail-text-type",
					Language:     "de",
					Title:        "title",
					PreHeader:    "pre-header",
					Subject:      "subject",
					Greeting:     "greeting",
					Text:         "text",
					ButtonText:   "button-text",
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
							org.NewMailTextAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"mail-text-type",
								"de",
								"title",
								"pre-header",
								"subject",
								"greeting",
								"text",
								"button-text",
							),
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.MailText{
					MailTextType: "mail-text-type",
					Language:     "de",
					Title:        "title",
					PreHeader:    "pre-header",
					Subject:      "subject",
					Greeting:     "greeting",
					Text:         "text",
					ButtonText:   "button-text",
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
							org.NewMailTextAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"mail-text-type",
								"de",
								"title",
								"pre-header",
								"subject",
								"greeting",
								"text",
								"button-text",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newMailTextChangedEvent(
									context.Background(),
									"org1",
									"mail-text-type",
									"de",
									"title-change",
									"pre-header-change",
									"subject-change",
									"greeting-change",
									"text-change",
									"button-text-change"),
							),
						},
						nil,
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.MailText{
					MailTextType: "mail-text-type",
					Language:     "de",
					Title:        "title-change",
					PreHeader:    "pre-header-change",
					Subject:      "subject-change",
					Greeting:     "greeting-change",
					Text:         "text-change",
					ButtonText:   "button-text-change",
				},
			},
			res: res{
				want: &domain.MailText{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "org1",
						ResourceOwner: "org1",
					},
					MailTextType: "mail-text-type",
					Language:     "de",
					Title:        "title-change",
					PreHeader:    "pre-header-change",
					Subject:      "subject-change",
					Greeting:     "greeting-change",
					Text:         "text-change",
					ButtonText:   "button-text-change",
					State:        domain.PolicyStateActive,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeMailText(tt.args.ctx, tt.args.orgID, tt.args.policy)
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

func TestCommandSide_RemoveMailText(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx          context.Context
		orgID        string
		mailTextType string
		language     string
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
				ctx:          context.Background(),
				orgID:        "org1",
				mailTextType: "mail-text-type",
				language:     "de",
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
							org.NewMailTextAddedEvent(context.Background(),
								&org.NewAggregate("org1", "org1").Aggregate,
								"mail-text-type",
								"de",
								"title",
								"pre-header",
								"subject",
								"greeting",
								"text",
								"button-text",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								org.NewMailTextRemovedEvent(context.Background(),
									&org.NewAggregate("org1", "org1").Aggregate,
									"mail-text-type",
									"de"),
							),
						},
						nil,
						uniqueConstraintsFromEventConstraint(policy.NewRemoveMailTextUniqueConstraint("org1", "mail-text-type", "de")),
					),
				),
			},
			args: args{
				ctx:          context.Background(),
				orgID:        "org1",
				mailTextType: "mail-text-type",
				language:     "de",
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
			err := r.RemoveMailText(tt.args.ctx, tt.args.orgID, tt.args.mailTextType, tt.args.language)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func newMailTextChangedEvent(ctx context.Context, orgID, mailTextType, language, title, preHeader, subject, greeting, text, buttonText string) *org.MailTextChangedEvent {
	event, _ := org.NewMailTextChangedEvent(ctx,
		&org.NewAggregate(orgID, orgID).Aggregate,
		mailTextType,
		language,
		[]policy.MailTextChanges{
			policy.ChangeTitle(title),
			policy.ChangePreHeader(preHeader),
			policy.ChangeSubject(subject),
			policy.ChangeGreeting(greeting),
			policy.ChangeText(text),
			policy.ChangeButtonText(buttonText),
		},
	)
	return event
}

package command

//
//import (
//	"context"
//	"github.com/caos/zitadel/internal/domain"
//	caos_errs "github.com/caos/zitadel/internal/errors"
//	"github.com/caos/zitadel/internal/eventstore"
//	"github.com/caos/zitadel/internal/eventstore/repository"
//	"github.com/caos/zitadel/internal/eventstore/v1/models"
//	"github.com/caos/zitadel/internal/repository/iam"
//	"github.com/caos/zitadel/internal/repository/policy"
//	"github.com/stretchr/testify/assert"
//	"testing"
//)
//
//func TestCommandSide_AddDefaultMessageTextPolicy(t *testing.T) {
//	type fields struct {
//		eventstore *eventstore.Eventstore
//	}
//	type args struct {
//		ctx    context.Context
//		policy *domain.CustomMessageText
//	}
//	type res struct {
//		want *domain.CustomMessageText
//		err  func(error) bool
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		res    res
//	}{
//		{
//			name: "mail text invalid, invalid argument error",
//			fields: fields{
//				eventstore: eventstoreExpect(
//					t,
//				),
//			},
//			args: args{
//				ctx:    context.Background(),
//				policy: &domain.CustomMessageText{},
//			},
//			res: res{
//				err: caos_errs.IsErrorInvalidArgument,
//			},
//		},
//		{
//			name: "mail text already existing, already exists error",
//			fields: fields{
//				eventstore: eventstoreExpect(
//					t,
//					expectFilter(
//						eventFromEventPusher(
//							iam.NewMessageTextAddedEvent(context.Background(),
//								&iam.NewAggregate().Aggregate,
//								"mail-text-type",
//								"de",
//								"title",
//								"pre-header",
//								"subject",
//								"greeting",
//								"text",
//								"button-text",
//							),
//						),
//					),
//				),
//			},
//			args: args{
//				ctx: context.Background(),
//				policy: &domain.MessageText{
//					Template: "mail-text-type",
//					Language:     "de",
//					Title:        "title",
//					PreHeader:    "pre-header",
//					Subject:      "subject",
//					Greeting:     "greeting",
//					Text:         "text",
//					ButtonText:   "button-text",
//				},
//			},
//			res: res{
//				err: caos_errs.IsErrorAlreadyExists,
//			},
//		},
//		{
//			name: "add mail template,ok",
//			fields: fields{
//				eventstore: eventstoreExpect(
//					t,
//					expectFilter(),
//					expectPush(
//						[]*repository.Event{
//							eventFromEventPusher(
//								iam.NewMessageTextAddedEvent(context.Background(),
//									&iam.NewAggregate().Aggregate,
//									"mail-text-type",
//									"de",
//									"title",
//									"pre-header",
//									"subject",
//									"greeting",
//									"text",
//									"button-text",
//								),
//							),
//						},
//						uniqueConstraintsFromEventConstraint(policy.NewAddMessageTextUniqueConstraint("IAM", "mail-text-type", "de")),
//					),
//				),
//			},
//			args: args{
//				ctx: context.Background(),
//				policy: &domain.MessageText{
//					Template: "mail-text-type",
//					Language:     "de",
//					Title:        "title",
//					PreHeader:    "pre-header",
//					Subject:      "subject",
//					Greeting:     "greeting",
//					Text:         "text",
//					ButtonText:   "button-text",
//				},
//			},
//			res: res{
//				want: &domain.MessageText{
//					ObjectRoot: models.ObjectRoot{
//						AggregateID:   "IAM",
//						ResourceOwner: "IAM",
//					},
//					Template: "mail-text-type",
//					Language:     "de",
//					Title:        "title",
//					PreHeader:    "pre-header",
//					Subject:      "subject",
//					Greeting:     "greeting",
//					Text:         "text",
//					ButtonText:   "button-text",
//					State:        domain.PolicyStateActive,
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			r := &Commands{
//				eventstore: tt.fields.eventstore,
//			}
//			got, err := r.AddDefaultMessageText(tt.args.ctx, tt.args.policy)
//			if tt.res.err == nil {
//				assert.NoError(t, err)
//			}
//			if tt.res.err != nil && !tt.res.err(err) {
//				t.Errorf("got wrong err: %v ", err)
//			}
//			if tt.res.err == nil {
//				assert.Equal(t, tt.res.want, got)
//			}
//		})
//	}
//}
//
//func TestCommandSide_ChangeDefaultMessageTextPolicy(t *testing.T) {
//	type fields struct {
//		eventstore *eventstore.Eventstore
//	}
//	type args struct {
//		ctx    context.Context
//		policy *domain.MessageText
//	}
//	type res struct {
//		want *domain.MessageText
//		err  func(error) bool
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		res    res
//	}{
//		{
//			name: "mailtext invalid, invalid argument error",
//			fields: fields{
//				eventstore: eventstoreExpect(
//					t,
//				),
//			},
//			args: args{
//				ctx:    context.Background(),
//				policy: &domain.MessageText{},
//			},
//			res: res{
//				err: caos_errs.IsErrorInvalidArgument,
//			},
//		},
//		{
//			name: "mail text not existing, not found error",
//			fields: fields{
//				eventstore: eventstoreExpect(
//					t,
//					expectFilter(),
//				),
//			},
//			args: args{
//				ctx: context.Background(),
//				policy: &domain.MessageText{
//					Template: "mail-text-type",
//					Language:     "de",
//					Title:        "title",
//					PreHeader:    "pre-header",
//					Subject:      "subject",
//					Greeting:     "greeting",
//					Text:         "text",
//					ButtonText:   "button-text",
//				},
//			},
//			res: res{
//				err: caos_errs.IsNotFound,
//			},
//		},
//		{
//			name: "no changes, precondition error",
//			fields: fields{
//				eventstore: eventstoreExpect(
//					t,
//					expectFilter(
//						eventFromEventPusher(
//							iam.NewMessageTextAddedEvent(context.Background(),
//								&iam.NewAggregate().Aggregate,
//								"mail-text-type",
//								"de",
//								"title",
//								"pre-header",
//								"subject",
//								"greeting",
//								"text",
//								"button-text",
//							),
//						),
//					),
//				),
//			},
//			args: args{
//				ctx: context.Background(),
//				policy: &domain.MessageText{
//					Template: "mail-text-type",
//					Language:     "de",
//					Title:        "title",
//					PreHeader:    "pre-header",
//					Subject:      "subject",
//					Greeting:     "greeting",
//					Text:         "text",
//					ButtonText:   "button-text",
//				},
//			},
//			res: res{
//				err: caos_errs.IsPreconditionFailed,
//			},
//		},
//		{
//			name: "change, ok",
//			fields: fields{
//				eventstore: eventstoreExpect(
//					t,
//					expectFilter(
//						eventFromEventPusher(
//							iam.NewMessageTextAddedEvent(context.Background(),
//								&iam.NewAggregate().Aggregate,
//								"mail-text-type",
//								"de",
//								"title",
//								"pre-header",
//								"subject",
//								"greeting",
//								"text",
//								"button-text",
//							),
//						),
//					),
//					expectPush(
//						[]*repository.Event{
//							eventFromEventPusher(
//								newDefaultMessageTextPolicyChangedEvent(
//									context.Background(),
//									"mail-text-type",
//									"de",
//									"title-change",
//									"pre-header-change",
//									"subject-change",
//									"greeting-change",
//									"text-change",
//									"button-text-change"),
//							),
//						},
//					),
//				),
//			},
//			args: args{
//				ctx: context.Background(),
//				policy: &domain.MessageText{
//					Template: "mail-text-type",
//					Language:     "de",
//					Title:        "title-change",
//					PreHeader:    "pre-header-change",
//					Subject:      "subject-change",
//					Greeting:     "greeting-change",
//					Text:         "text-change",
//					ButtonText:   "button-text-change",
//				},
//			},
//			res: res{
//				want: &domain.MessageText{
//					ObjectRoot: models.ObjectRoot{
//						AggregateID:   "IAM",
//						ResourceOwner: "IAM",
//					},
//					Template: "mail-text-type",
//					Language:     "de",
//					Title:        "title-change",
//					PreHeader:    "pre-header-change",
//					Subject:      "subject-change",
//					Greeting:     "greeting-change",
//					Text:         "text-change",
//					ButtonText:   "button-text-change",
//					State:        domain.PolicyStateActive,
//				},
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			r := &Commands{
//				eventstore: tt.fields.eventstore,
//			}
//			got, err := r.ChangeDefaultMessageText(tt.args.ctx, tt.args.policy)
//			if tt.res.err == nil {
//				assert.NoError(t, err)
//			}
//			if tt.res.err != nil && !tt.res.err(err) {
//				t.Errorf("got wrong err: %v ", err)
//			}
//			if tt.res.err == nil {
//				assert.Equal(t, tt.res.want, got)
//			}
//		})
//	}
//}
//
//func newDefaultMessageTextPolicyChangedEvent(ctx context.Context, mailTextType, language, title, preHeader, subject, greeting, text, buttonText string) *iam.MessageTextChangedEvent {
//	event, _ := iam.NewMessageTextChangedEvent(ctx,
//		&iam.NewAggregate().Aggregate,
//		mailTextType,
//		language,
//		[]policy.MessageTextChanges{
//			policy.ChangeTitle(title),
//			policy.ChangePreHeader(preHeader),
//			policy.ChangeSubject(subject),
//			policy.ChangeGreeting(greeting),
//			policy.ChangeText(text),
//			policy.ChangeButtonText(buttonText),
//		},
//	)
//	return event
//}

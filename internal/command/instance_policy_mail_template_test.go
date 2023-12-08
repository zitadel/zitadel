package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/policy"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestCommandSide_AddDefaultMailTemplatePolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		policy *domain.MailTemplate
	}
	type res struct {
		want *domain.MailTemplate
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "mailtemplate invalid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				policy: &domain.MailTemplate{},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "mailtemplate already existing, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewMailTemplateAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								[]byte("template"),
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.MailTemplate{
					Template: []byte("template"),
				},
			},
			res: res{
				err: zerrors.IsErrorAlreadyExists,
			},
		},
		{
			name: "add mail template,ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						instance.NewMailTemplateAddedEvent(
							authz.WithInstanceID(context.Background(), "INSTANCE"),
							&instance.NewAggregate("INSTANCE").Aggregate,
							[]byte("template"),
						),
					),
				),
			},
			args: args{
				ctx: authz.WithInstanceID(context.Background(), "INSTANCE"),
				policy: &domain.MailTemplate{
					Template: []byte("template"),
				},
			},
			res: res{
				want: &domain.MailTemplate{
					ObjectRoot: models.ObjectRoot{
						InstanceID:    "INSTANCE",
						AggregateID:   "INSTANCE",
						ResourceOwner: "INSTANCE",
					},
					Template: []byte("template"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.AddDefaultMailTemplate(tt.args.ctx, tt.args.policy)
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

func TestCommandSide_ChangeDefaultMailTemplatePolicy(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		policy *domain.MailTemplate
	}
	type res struct {
		want *domain.MailTemplate
		err  func(error) bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		res    res
	}{
		{
			name: "mailtemplate invalid, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx:    context.Background(),
				policy: &domain.MailTemplate{},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "mailtempalte not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.MailTemplate{
					Template: []byte("template-change"),
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
							instance.NewMailTemplateAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								[]byte("template"),
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.MailTemplate{
					Template: []byte("template"),
				},
			},
			res: res{
				err: zerrors.IsPreconditionFailed,
			},
		},
		{
			name: "change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							instance.NewMailTemplateAddedEvent(context.Background(),
								&instance.NewAggregate("INSTANCE").Aggregate,
								[]byte("template"),
							),
						),
					),
					expectPush(
						newDefaultMailTemplatePolicyChangedEvent(context.Background(), []byte("template-change")),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.MailTemplate{
					Template: []byte("template-change"),
				},
			},
			res: res{
				want: &domain.MailTemplate{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "INSTANCE",
						ResourceOwner: "INSTANCE",
						InstanceID:    "INSTANCE",
					},
					Template: []byte("template-change"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeDefaultMailTemplate(tt.args.ctx, tt.args.policy)
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

func newDefaultMailTemplatePolicyChangedEvent(ctx context.Context, template []byte) *instance.MailTemplateChangedEvent {
	event, _ := instance.NewMailTemplateChangedEvent(ctx,
		&instance.NewAggregate("INSTANCE").Aggregate,
		[]policy.MailTemplateChanges{
			policy.ChangeTemplate(template),
		},
	)
	return event
}

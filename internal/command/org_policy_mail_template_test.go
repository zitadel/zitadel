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

func TestCommandSide_AddMailTemplate(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		orgID  string
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
			name: "org id missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.MailTemplate{
					Template: []byte("template"),
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
			},
		},
		{
			name: "mail template already existing, already exists error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							org.NewMailTemplateAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								[]byte("template"),
							),
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.MailTemplate{
					Template: []byte("template"),
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
						org.NewMailTemplateAddedEvent(context.Background(),
							&org.NewAggregate("org1").Aggregate,
							[]byte("template"),
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.MailTemplate{
					Template: []byte("template"),
				},
			},
			res: res{
				want: &domain.MailTemplate{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "org1",
						ResourceOwner: "org1",
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
			got, err := r.AddMailTemplate(tt.args.ctx, tt.args.orgID, tt.args.policy)
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

func TestCommandSide_ChangeMailTemplate(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx    context.Context
		orgID  string
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
			name: "org id missing, invalid argument error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				policy: &domain.MailTemplate{
					Template: []byte("template"),
				},
			},
			res: res{
				err: zerrors.IsErrorInvalidArgument,
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
				policy: &domain.MailTemplate{
					Template: []byte("template"),
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
							org.NewMailTemplateAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								[]byte("template"),
							),
						),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
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
							org.NewMailTemplateAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								[]byte("template"),
							),
						),
					),
					expectPush(
						newMailTemplateChangedEvent(context.Background(), "org1", "template2"),
					),
				),
			},
			args: args{
				ctx:   context.Background(),
				orgID: "org1",
				policy: &domain.MailTemplate{
					Template: []byte("template2"),
				},
			},
			res: res{
				want: &domain.MailTemplate{
					ObjectRoot: models.ObjectRoot{
						AggregateID:   "org1",
						ResourceOwner: "org1",
					},
					Template: []byte("template2"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeMailTemplate(tt.args.ctx, tt.args.orgID, tt.args.policy)
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

func TestCommandSide_RemoveMailTemplate(t *testing.T) {
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
							org.NewMailTemplateAddedEvent(context.Background(),
								&org.NewAggregate("org1").Aggregate,
								[]byte("template"),
							),
						),
					),
					expectPush(
						org.NewMailTemplateRemovedEvent(context.Background(),
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
			err := r.RemoveMailTemplate(tt.args.ctx, tt.args.orgID)
			if tt.res.err == nil {
				assert.NoError(t, err)
			}
			if tt.res.err != nil && !tt.res.err(err) {
				t.Errorf("got wrong err: %v ", err)
			}
		})
	}
}

func newMailTemplateChangedEvent(ctx context.Context, orgID string, template string) *org.MailTemplateChangedEvent {
	event, _ := org.NewMailTemplateChangedEvent(ctx,
		&org.NewAggregate(orgID).Aggregate,
		[]policy.MailTemplateChanges{
			policy.ChangeTemplate([]byte(template)),
		},
	)
	return event
}

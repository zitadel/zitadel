package command

import (
	"context"
	"testing"

	"github.com/caos/zitadel/internal/domain"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/repository"
	"github.com/caos/zitadel/internal/id"
	id_mock "github.com/caos/zitadel/internal/id/mock"
	"github.com/caos/zitadel/internal/notification/channels/twilio"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/stretchr/testify/assert"
)

func TestCommandSide_AddSMSConfigTwilio(t *testing.T) {
	type fields struct {
		eventstore  *eventstore.Eventstore
		idGenerator id.Generator
	}
	type args struct {
		ctx context.Context
		sms *twilio.TwilioConfig
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
			name: "add sms config twilio, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(iam.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								"providerid",
								"sid",
								"token",
								"from",
							),
							),
						},
					),
				),
				idGenerator: id_mock.NewIDGeneratorExpectIDs(t, "providerid"),
			},
			args: args{
				ctx: context.Background(),
				sms: &twilio.TwilioConfig{
					SID:   "sid",
					Token: "token",
					From:  "from",
				},
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore:  tt.fields.eventstore,
				idGenerator: tt.fields.idGenerator,
			}
			_, got, err := r.AddSMSConfigTwilio(tt.args.ctx, tt.args.sms)
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

func TestCommandSide_ChangeSMSConfigTwilio(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx context.Context
		id  string
		sms *twilio.TwilioConfig
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
			name: "id empty, precondition error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
				),
			},
			args: args{
				ctx: context.Background(),
				sms: &twilio.TwilioConfig{},
			},
			res: res{
				err: caos_errs.IsErrorInvalidArgument,
			},
		},
		{
			name: "sms not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				sms: &twilio.TwilioConfig{},
				id:  "id",
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
							iam.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								"providerid",
								"sid",
								"token",
								"from",
							),
						),
					),
				),
			},
			args: args{
				ctx: context.Background(),
				sms: &twilio.TwilioConfig{
					SID:   "sid",
					Token: "token",
					From:  "from",
				},
				id: "providerid",
			},
			res: res{
				err: caos_errs.IsPreconditionFailed,
			},
		},
		{
			name: "sms config twilio change, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								"providerid",
								"sid",
								"token",
								"from",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								newSMSConfigTwilioChangedEvent(
									context.Background(),
									"providerid",
									"sid2",
									"token2",
									"from2",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				sms: &twilio.TwilioConfig{
					SID:   "sid2",
					Token: "token2",
					From:  "from2",
				},
				id: "providerid",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ChangeSMSConfigTwilio(tt.args.ctx, tt.args.id, tt.args.sms)
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

func TestCommandSide_ActivateSMSConfigTwilio(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx context.Context
		id  string
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
			name: "id empty, invalid error",
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
			name: "sms not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				id:  "id",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "sms config twilio activate, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								"providerid",
								"sid",
								"token",
								"from",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								iam.NewSMSConfigTwilioActivatedEvent(
									context.Background(),
									&iam.NewAggregate().Aggregate,
									"providerid",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				id:  "providerid",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.ActivateSMSConfigTwilio(tt.args.ctx, tt.args.id)
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

func TestCommandSide_DeactivateSMSConfigTwilio(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx context.Context
		id  string
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
			name: "id empty, invalid error",
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
			name: "sms not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				id:  "id",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "sms config twilio deactivate, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								"providerid",
								"sid",
								"token",
								"from",
							),
						),
						eventFromEventPusher(
							iam.NewSMSConfigTwilioActivatedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								"providerid",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								iam.NewSMSConfigTwilioDeactivatedEvent(
									context.Background(),
									&iam.NewAggregate().Aggregate,
									"providerid",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				id:  "providerid",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.DeactivateSMSConfigTwilio(tt.args.ctx, tt.args.id)
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

func TestCommandSide_RemoveSMSConfigTwilio(t *testing.T) {
	type fields struct {
		eventstore *eventstore.Eventstore
	}
	type args struct {
		ctx context.Context
		id  string
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
			name: "id empty, invalid error",
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
			name: "sms not existing, not found error",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(),
				),
			},
			args: args{
				ctx: context.Background(),
				id:  "id",
			},
			res: res{
				err: caos_errs.IsNotFound,
			},
		},
		{
			name: "sms config twilio remove, ok",
			fields: fields{
				eventstore: eventstoreExpect(
					t,
					expectFilter(
						eventFromEventPusher(
							iam.NewSMSConfigTwilioAddedEvent(
								context.Background(),
								&iam.NewAggregate().Aggregate,
								"providerid",
								"sid",
								"token",
								"from",
							),
						),
					),
					expectPush(
						[]*repository.Event{
							eventFromEventPusher(
								iam.NewSMSConfigTwilioRemovedEvent(
									context.Background(),
									&iam.NewAggregate().Aggregate,
									"providerid",
								),
							),
						},
					),
				),
			},
			args: args{
				ctx: context.Background(),
				id:  "providerid",
			},
			res: res{
				want: &domain.ObjectDetails{
					ResourceOwner: "IAM",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Commands{
				eventstore: tt.fields.eventstore,
			}
			got, err := r.RemoveSMSConfigTwilio(tt.args.ctx, tt.args.id)
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

func newSMSConfigTwilioChangedEvent(ctx context.Context, id, sid, token, from string) *iam.SMSConfigTwilioChangedEvent {
	changes := []iam.SMSConfigTwilioChanges{
		iam.ChangeSMSConfigTwilioSID(sid),
		iam.ChangeSMSConfigTwilioToken(token),
		iam.ChangeSMSConfigTwilioFrom(from),
	}
	event, _ := iam.NewSMSConfigTwilioChangedEvent(ctx,
		&iam.NewAggregate().Aggregate,
		id,
		changes,
	)
	return event
}
